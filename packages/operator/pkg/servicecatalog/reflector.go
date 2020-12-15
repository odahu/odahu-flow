package servicecatalog

import (
	"context"
	"go.uber.org/multierr"
	"go.uber.org/zap"
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
	HandleTimeout time.Duration
	JobCapacity int
	WorkersCount int

}


type versionedEvent struct {
	version int
	event   GenericEvent
}


// RunProcessor process events in jobChan. Can re-schedule events in case of temporary errors.
func (r Reflector) RunProcessor(ctx context.Context,
	jobChan <-chan versionedEvent, schedChan chan<- versionedEvent) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case event := <-jobChan:
			if err := r.H.Handle(event.event.Embedded); err != nil {
				if !IsTemporary(err){
					return err
				}
				r.Log.Warnw("Temporary error during event handling. Reschedule versionedEvent",
					"event", event, zap.Error(err))
				schedChan <- event
			}
		}
	}
}


// RunFetcher fetch events from API server and schedule them to processing.
func (r Reflector) RunFetcher(ctx context.Context, schedChan chan<- versionedEvent) error {
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
				schedChan <- versionedEvent{
					version: cursor,
					event:   event,
				}
			}
		}

	}
}


// RunScheduler fills jobChan according some rules.
// Event skipped if more fresh version of event was scheduled previously.
func (r Reflector) RunScheduler(ctx context.Context,
	jobChan chan<- versionedEvent, schedChan <-chan versionedEvent) {

	lastCursorForEntityID := make(map[string]int)

	var log = r.Log

	for {
		select {
		case <-ctx.Done():
			log.Debugw("Cancel event was received")
			return
		case item := <-schedChan:

			log = log.With("eventVersion", item.version, "entityID", item.event.EntityID)

			log.Debugw("New event is scheduled")

			lastScheduledVersion, ok := lastCursorForEntityID[item.event.EntityID]
			if !ok || lastScheduledVersion <= item.version {
				lastCursorForEntityID[item.event.EntityID] = item.version
				jobChan <- item
				log.Debugw("Event was enqueued for jobChan")
			} else {
				log.Debugw("Event version is stale. Skip enqueueing",
					"lastScheduledVersion", lastScheduledVersion)
			}

		}
	}

}


func (r Reflector) Run(ctx context.Context) error {



	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobChan := make(chan versionedEvent, r.JobCapacity)
	schedChan := make(chan versionedEvent, r.JobCapacity)

	// Grab all goroutines error into err
	var err error
	errLock := sync.Mutex{}
	wg := &sync.WaitGroup{}

	// Goroutine fetch events from API server and schedule them
	wg.Add(1)
	go func() {
		defer wg.Done()
		if gErr := r.RunFetcher(ctx, schedChan); gErr != nil {
			errLock.Lock()
			err = multierr.Append(err, gErr)
			errLock.Unlock()
		}
	}()

	// Scheduler
	wg.Add(1)
	go func() {
		defer wg.Done()
		r.RunScheduler(ctx, jobChan, schedChan)
	}()

	// Handler. Can re-schedule event if temporary error happens
	for i := 0; i < r.WorkersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if gErr := r.RunProcessor(ctx, jobChan, schedChan); gErr != nil {
				errLock.Lock()
				err = multierr.Append(err, gErr)
				errLock.Unlock()
			}
		}()
	}

	wg.Wait()

	return err
}