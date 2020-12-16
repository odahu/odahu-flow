package servicecatalog

import (
	"context"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"k8s.io/client-go/util/workqueue"
	"sync"
	"time"
)


type eventFetcher interface {
	GetLastEvents(cursor int) (LatestGenericEvents, error)
}

type eventHandler interface {
	Handle(event interface{}) (err error)
}

// Reflector subscribe on changes in external system and handle these events to reflect them
// in Local state (eg. App storage)
type Reflector struct {
	Log            *zap.SugaredLogger

	H             eventHandler
	C             eventFetcher
	FetchTimeout  time.Duration
	WorkersCount int
	q workqueue.RateLimitingInterface
	eventCache map[interface{}]interface{}
	eventCacheMu *sync.RWMutex
}

type ReflectorOpts struct {
	WorkersCount int
	FetchTimeout time.Duration
}

func NewReflector(log *zap.SugaredLogger, handler eventHandler, fetcher eventFetcher,
	q workqueue.RateLimitingInterface, opts ReflectorOpts) Reflector {

	r := Reflector{
		Log:          log,
		H:            handler,
		C:            fetcher,
		q:            q,
		eventCache:   make(map[interface{}]interface{}),
		eventCacheMu: &sync.RWMutex{},
		WorkersCount: 1,
		FetchTimeout: 10 * time.Second,
	}

	if opts.WorkersCount != 0 {
		r.WorkersCount = opts.WorkersCount
	}
	if opts.FetchTimeout != 0 {
		r.FetchTimeout = opts.FetchTimeout
	}

	return r

}


func (r Reflector) RunProcessor(ctx context.Context) error {

	for {
		EntityID, shutdown := r.q.Get()
		r.eventCacheMu.RLock()
		event := r.eventCache[EntityID]
		r.eventCacheMu.RUnlock()

		if shutdown {
			return nil
		}

		err := r.H.Handle(event)
		r.q.Done(EntityID)
		if err != nil {
			r.q.AddRateLimited(EntityID)
			continue
		}
		r.q.Forget(EntityID)

	}

}


func (r Reflector) RunFetcher(ctx context.Context) error {
	var cursor int
	t := time.NewTicker(r.FetchTimeout)

	for {

		select {
		case  <-ctx.Done():
			return nil
		case  <-t.C:
			lastEvents, err := r.C.GetLastEvents(cursor)
			if err != nil {
				r.Log.Errorw("Unable to get last events", zap.Error(err))

				if !IsTemporary(err) {
					return err
				}

				r.Log.Warnw("Temporary error during event fetching.", zap.Error(err))
				continue
			}
			cursor = lastEvents.Cursor

			for _, event := range lastEvents.Events {
				r.q.Add(event.VersionedID.EntityID)
				r.eventCacheMu.Lock()
				r.eventCache[event.VersionedID.EntityID] = event.Embedded
				r.eventCacheMu.Unlock()
			}
		}

	}
}


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
		if gErr := r.RunFetcher(ctx); gErr != nil {
			errLock.Lock()
			err = multierr.Append(err, gErr)
			errLock.Unlock()
			r.Log.Errorw("Error during fetching events. Send signal to stop", zap.Error(err))
			cancel()
		}
	}()


	// Handler. Can re-schedule event if temporary error happens
	for i := 0; i < r.WorkersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if gErr := r.RunProcessor(ctx); gErr != nil {
				errLock.Lock()
				err = multierr.Append(err, gErr)
				errLock.Unlock()
				r.Log.Errorw("Error during processing events. Send signal to stop", zap.Error(err))
				cancel()
			}
		}()
	}

	go func() {
		<-ctx.Done()
		r.Log.Debug("Get cancel signal. Shutdown queue")
		r.q.ShutDown()
	}()

	wg.Wait()

	return err
}