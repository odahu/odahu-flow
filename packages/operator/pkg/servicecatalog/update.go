package servicecatalog

import (
	"container/list"
	"go.uber.org/zap"
	"sync"
	"time"
)

type LatestGenericEvents struct {
	Cursor int
	Events []interface{}
}

type EventFetcher interface {
	GetLastEvents(cursor int) (LatestGenericEvents, error)
}

type EventHandler interface {
	Handle(event interface{}) (err error)
}

type Updater struct {
	Log            *zap.SugaredLogger

	H             EventHandler
	C             EventFetcher
	FetchTimeout  time.Duration
	HandleTimeout time.Duration

}

func (u Updater) Run(stopCh <-chan struct{}) {

	// We need two goroutines that communicate with each other through queue.
	// First goroutine fill queue by events,
	// Second goroutine retrieves events and pass them to handlers.
	// If handler return error for some event then second goroutine push this event back to queue to try handle it later

	queue := list.New()
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	wg.Add(2)

	// Fill queue from API Server
	go func() {
		defer wg.Done()

		var cursor int
		t := time.NewTicker(u.FetchTimeout * time.Second)
		for {
			select {
			case  <-stopCh:
				return
			case  <-t.C:
				lastEvents, err := u.C.GetLastEvents(cursor)
				if err != nil {
					u.Log.Errorw("Unable to get last events", zap.Error(err))
					continue
				}
				cursor = lastEvents.Cursor

				mu.Lock()
				for _, event := range lastEvents.Events {
					// Write error handling
					queue.PushBack(event)
				}
				mu.Unlock()
			}
		}

	}()

	// Process queue
	go func() {
		defer wg.Done()
		t := time.NewTicker(u.HandleTimeout * time.Second)
		for {
			select {
			case <-stopCh:
				return
			case <-t.C:
				for {
					mu.Lock()
					el := queue.Front()
					mu.Unlock()
					if el == nil {
						break
					}

					event := el.Value
					if err := u.H.Handle(event); err != nil {
						u.Log.Warnw("Error during event handling. Push event back to queue",
							"event", event)
						mu.Lock()
						queue.PushBack(event)
						mu.Unlock()
					}


				}
			}
		}
	}()

	wg.Wait()
	u.Log.Info("Stop signal was received")
}