package controller

import (
	"context"
	"fmt"
	"sync"
)

type Runnable interface {
	Run(ctx context.Context) error
	String() string
}

type WorkersManager struct {
	runners          []Runnable
	running          uint
	mu               *sync.Mutex
}

func NewWorkersManager() *WorkersManager {

	tm := WorkersManager{
		runners: make([]Runnable, 0),
		running: 0,
		mu:      &sync.Mutex{},
	}

	return &tm
}

func (r *WorkersManager) AddRunnable(runnable Runnable) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.runners = append(r.runners, runnable)
}

type runnerResponse struct {
	err      error
	runnerID string
}

func (r *WorkersManager) Run(ctx context.Context) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	respCh := make(chan runnerResponse, len(r.runners))
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, runnable := range r.runners {
		runnable := runnable
		go func() {
			err := runnable.Run(ctx)
			respCh <- runnerResponse{
				err:      err,
				runnerID: runnable.String(),
			}
		}()
	}
	log.Info(fmt.Sprintf("%v runners are launched", len(r.runners)))


	for range r.runners {
		resp := <-respCh
		if resp.err != nil {
			log.Error(resp.err, fmt.Sprintf("Runner return error. Id: %v", resp.runnerID))
			cancel()
		} else {
			log.Info(fmt.Sprintf("Runner was stopped. Id: %v", resp.runnerID))
		}
	}

	// If runners == 0 we anyway should wait cancellation before return from function
	<-ctx.Done()

	return nil
}
