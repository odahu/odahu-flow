package workersmanager

import (
	"context"
	"errors"
	"fmt"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sync"
	"time"
)

var (
	log = logf.Log.WithName("workers-manager")
)

type Runnable interface {
	Run(ctx context.Context) error
	String() string
}

type WorkersManager struct {
	runners          []Runnable
	running          uint
	mu               *sync.Mutex
	cancelRunnersCtx context.CancelFunc
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
	r.runners = append(r.runners, runnable)
}

func (r *WorkersManager) Shutdown(ctx context.Context) error {
	t := time.NewTicker(time.Millisecond * 200)
	defer t.Stop()

	r.cancelRunnersCtx()

	for {
		select {
		case <-t.C:
			log.Info("Runners manager is trying to shutdown gracefully")
			r.mu.Lock()
			allStopped := r.running == 0
			r.mu.Unlock()
			if allStopped {
				return nil
			}

		case <-ctx.Done():
			log.Error(ctx.Err(), "Unable to shutdown gracefully")
			return ctx.Err()
		}
	}
}

type runnerResponse struct {
	err      error
	runnerID string
}

func (r *WorkersManager) Run() error {

	runnersTotal := len(r.runners)
	respCh := make(chan runnerResponse, runnersTotal)

	ctx, cancel := context.WithCancel(context.Background())
	r.cancelRunnersCtx = cancel

	for _, runnable := range r.runners {
		runnable := runnable

		r.mu.Lock()
		r.running++
		r.mu.Unlock()

		go func() {
			err := runnable.Run(ctx)
			respCh <- runnerResponse{
				err:      err,
				runnerID: runnable.String(),
			}
		}()
	}
	log.Info(fmt.Sprintf("%v runners will be launched", runnersTotal))


	// Finishing code
	var err error
	for range r.runners {
		resp := <-respCh

		if resp.err != nil {
			err = errors.New("one or more runners return error")
			log.Error(resp.err, fmt.Sprintf("Runner return error. Id: %v", resp.runnerID))
			cancel()
		} else {
			log.Info(fmt.Sprintf("Runner was stopped. Id: %v", resp.runnerID))
		}

		r.mu.Lock()
		r.running--
		r.mu.Unlock()

	}

	return err
}
