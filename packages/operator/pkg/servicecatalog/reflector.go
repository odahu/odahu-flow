package servicecatalog

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"k8s.io/client-go/util/workqueue"
	"sync"
	"time"
)


// EventFetcher fetch events from event source
type EventFetcher interface {
	// GetLastEvents get last events from event source using cursor to determine offset of events
	// If returned error implements interface with func Temporary() and Temporary() == true
	// then attempt to fetch events will be retried
	GetLastEvents(cursor int) (LatestGenericEvents, error)
}

// EventHandler process event somehow
type EventHandler interface {
	// Handle function will be called for each new event that was fetched by EventFetcher
	// Type assertion can be applied to event to get access for underlying concrete event value.
	//
	// if error is returned that event will be rescheduled for further processing using
	// exponential backoff algorithm
	Handle(event interface{}, logger *zap.SugaredLogger) (err error)
}

// Reflector periodically fetches new events using EventFetcher and planning them for handling by
// EventHandler. The main goal of reflector is to provide ability to reflect changes in external system.
// The same idea as in Kubernetes Reconciliation (level-based triggering)
// https://hackernoon.com/level-triggering-and-reconciliation-in-kubernetes-1f17fe30333d
//
//
// Reflector assumes that only last event with concrete GenericEvent.EntityID makes sense to reflect
// external system state changes.
//
// For example
// If more than one event with the same GenericEvent.EntityID was fetched before worker
// is able to process them then only last one will be processed using EventHandler.Handle.
//
// Another example (with retry).
// If some event with GenericEvent.EntityID was re-enqueued because of failing (EventHandler.Handle returned error)
// but new event was fetched after that then only the last one event will be processed in next time.
//
type Reflector struct {
	log *zap.SugaredLogger

	handler      EventHandler
	fetcher      EventFetcher
	fetchTimeout time.Duration
	workersCount int
	queue        workqueue.RateLimitingInterface
	eventCache   *eventCache
}

// ReflectorOpts are options for Reflector
type ReflectorOpts struct {
	// How many concurrent workers will process event using EventHandler.Handle
	WorkersCount int
	// How often EventFetcher will fetch new events from source
	FetchTimeout time.Duration
}

func NewReflector(log *zap.SugaredLogger, handler EventHandler, fetcher EventFetcher,
	opts ReflectorOpts) Reflector {

	r := Reflector{
		log:     log,
		handler: handler,
		fetcher: fetcher,
		queue:   workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		eventCache:   &eventCache{
			store: make(map[interface{}]interface{}),
			mu:    &sync.RWMutex{},
		},
		workersCount: 1,
		fetchTimeout: 10 * time.Second,
	}

	if opts.WorkersCount != 0 {
		r.workersCount = opts.WorkersCount
	}
	if opts.FetchTimeout != 0 {
		r.fetchTimeout = opts.FetchTimeout
	}

	return r

}

// Run is blocking function that start reflector
// ctx can be used to stop execution
func (r Reflector) Run(ctx context.Context) error {


	ctx, cancel := context.WithCancel(ctx)
	defer cancel()


	// Grab all goroutines error into err
	var err error
	errLock := sync.Mutex{}
	wg := &sync.WaitGroup{}

	// Goroutine fetch events from API server and schedule them
	wg.Add(1)
	go func() {
		defer wg.Done()
		if gErr := r.runFetcher(ctx); gErr != nil {
			errLock.Lock()
			err = multierr.Append(err, gErr)
			errLock.Unlock()
			r.log.Errorw("Error during fetching events. Send signal to stop", zap.Error(err))
			cancel()
		} else {
			r.log.Info("Event fetcher was stopped")
		}
	}()


	// Handler. Can re-schedule event if temporary error happens
	for i := 0; i < r.workersCount; i++ {
		wg.Add(1)
		log := r.log.With("WorkerID", i)
		go func() {
			defer wg.Done()
			r.runProcessor(log)
			log.Info("Worker was stopped")
		}()
	}
	r.log.Infof("Start %d workers", r.workersCount)

	go func() {
		<-ctx.Done()
		r.log.Debug("Get cancel signal. Shutdown queue")
		r.queue.ShutDown()
	}()

	wg.Wait()

	return err
}


func (r Reflector) runProcessor(log *zap.SugaredLogger) {

	for {

		EntityID, shutdown := r.queue.Get()
		processingTraceID := uuid.New().String()
		log := log.With("EntityID", EntityID, "TraceID", processingTraceID)
		if shutdown {
			return
		}
		event := r.eventCache.get(EntityID)

		log.Info("Start processing event")
		err := r.handler.Handle(event, log)
		r.queue.Done(EntityID)
		if err != nil {
			log.Errorf("Error during processing event: %s. Retry", err.Error())
			r.queue.AddRateLimited(EntityID)
			continue
		}
		log.Info("Event was processed successfully")
		r.queue.Forget(EntityID)
	}
}

func (r Reflector) runFetcher(ctx context.Context) error {
	var cursor int
	t := time.NewTicker(r.fetchTimeout)

	for {

		select {
		case  <-ctx.Done():
			return nil
		case  <-t.C:
			lastEvents, err := r.fetcher.GetLastEvents(cursor)
			if err != nil {
				r.log.Errorw("Unable to get last events", zap.Error(err))

				if !IsTemporary(err) {
					return err
				}

				r.log.Warnw("Temporary error during event fetching.", zap.Error(err))
				continue
			}
			if lastEvents.Cursor > cursor {
				cursor = lastEvents.Cursor
			}

			for _, event := range lastEvents.Events {
				r.queue.Add(event.EntityID)
				r.eventCache.put(event.EntityID, event.Event)
			}
		}

	}
}

type eventCache struct {
	store map[interface{}]interface{}
	mu *sync.RWMutex
}

func (e *eventCache) get(key interface{}) interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.store[key]
}

func (e *eventCache) put(key interface{}, obj interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.store[key] = obj
}

