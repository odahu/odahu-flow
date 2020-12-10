package servicecatalog

import (
	"context"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"sync"
	"time"
)



type eventFetcher interface {
	GetLastEvents(cursor int) (latestGenericEvents, error)
}

type eventHandler interface {
	Handle(event interface{}) (err error)
}

// Reflector subscribe on changes in external system and handle this events to reflect them
// in Local state (eg. App storage)
type Reflector struct {
	Log            *zap.SugaredLogger

	H             eventHandler
	C             eventFetcher
	FetchTimeout  time.Duration
	HandleTimeout time.Duration

}

type temporaryError interface {
	Temporary() bool
}

func IsTemporary(err error) bool {
	tErr, ok := err.(temporaryError)
	return ok && tErr.Temporary()
}

func (u Reflector) Run(ctx context.Context) error {

	// We need two goroutines that communicate with each other through queue.
	// First goroutine fill queue by events,
	// Second goroutine retrieves events and pass them to handlers.
	// If handler return error for some event then second goroutine push this event back to queue to try handle it later

	queue := NewQueue()
	wg := &sync.WaitGroup{}
	wg.Add(2)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var fetchEventsPermanentErr error
	var handleEventPermanentErr error

	// Fill queue from API Server
	go func() {

		defer wg.Done()
		defer cancel()

		var cursor int
		t := time.NewTicker(u.FetchTimeout * time.Second)

		for {

			select {
			case  <-ctx.Done():
				return
			case  <-t.C:
				lastEvents, err := u.C.GetLastEvents(cursor)
				if err != nil {
					u.Log.Errorw("Unable to get last events", zap.Error(err))

					if !IsTemporary(err) {
						fetchEventsPermanentErr = err
						return
					}

					u.Log.Warnw("Transient error during event fetching. Will try fetch later", zap.Error(err))
					continue
				}
				cursor = lastEvents.Cursor

				for _, event := range lastEvents.Events {
					queue.Push(ctx, &Item{Value: event})
				}
			}

		}

	}()

	// Process queue
	go func() {

		defer wg.Done()
		defer cancel()

		t := time.NewTicker(u.HandleTimeout * time.Second)
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				for {
					item := queue.Pop()
					if item == nil {
						// Queue is empty. Sleep again
						break
					}

					event := item.Value
					if err := u.H.Handle(event); err != nil {

						if !IsTemporary(err){
							handleEventPermanentErr = err
							return
						}
						u.Log.Warnw("Transient error during event handling. Push event back to queue",
							"event", event, zap.Error(err))
						queue.Push(ctx, item)
					}


				}
			}

		}
	}()

	wg.Wait()
	var err error
	if fetchEventsPermanentErr != nil {
		err = multierr.Append(err, fetchEventsPermanentErr)
	}
	if handleEventPermanentErr != nil {
		err = multierr.Append(err, fetchEventsPermanentErr)
	}
	return err
}