// Package scheduler provides manager to run long running function.
//
// The scheduler package should be used inside another goroutine.
package scheduler

import (
	"context"
	"fmt"
	"sync"
)

// Routine is an interface for process that want to be scheduled.
type Routine interface {
	// Stream return a channel and receive [context.Context] for cancellation.
	// The method must be implemented by waiting on [context.Context.Done]
	// before passing data to channel for scheduler cancellation.
	//
	// Example implementations:
	//  func (s *scheduledRoutine) Stream(ctx context.Context) <-chan any {
	//  	ch := make(chan any, 1)
	//  	go func() {
	//          arr := fetchItems()
	//  		for _, v := range arr {
	//  			select {
	//  			case <-ctx.Done():
	//  				close(ch)
	//  				return
	//  			default:
	//  				ch <- v
	//  			}
	//  		}
	//  		close(ch)
	//  	}()
	//
	//  	return ch
	//  }
	Stream(context.Context) <-chan any
	// Consume read from the channel returned by Stream()
	//  func (s *scheduledRoutine) Consume(ch <-chan any) {
	//  	for v := range ch {
	//  		n := v.(string)
	//  		fmt.Println("Hello,", n)
	//  		time.Sleep(3 * time.Second)
	//  	}
	//  }
	Consume(<-chan any)
}

type scheduledRoutine struct {
	timer   Timer
	routine Routine
}

type scheduleManager struct {
	wg       *sync.WaitGroup
	routines []scheduledRoutine
	doneFunc context.CancelFunc
	closeCh  chan struct{}
}

// New returns new schedule manager instance
func New() *scheduleManager {
	return &scheduleManager{
		wg:      new(sync.WaitGroup),
		closeCh: make(chan struct{}),
	}
}

// Register add Timer and Routine to manager for running
func (m *scheduleManager) Register(timer Timer, r Routine) {
	sr := scheduledRoutine{
		timer,
		r,
	}
	m.routines = append(m.routines, sr)
}

// Run start scheduler (will block on call)
func (m *scheduleManager) Run() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	m.doneFunc = cancelFunc
	for _, r := range m.routines {
		m.wg.Add(1)
		go worker(ctx, m.wg, r)
	}

	m.wg.Wait()
	m.closeCh <- struct{}{}
}

// Stop will notify worker to stop processing scheduler (will block until finish)
func (m *scheduleManager) Stop() {
	m.doneFunc()
	<-m.closeCh
}

func worker(ctx context.Context, wg *sync.WaitGroup, sr scheduledRoutine) {
	defer wg.Done()
	for {
		fmt.Println("Start Update Name scheduler")
		<-sr.timer()

		ch := sr.routine.Stream(ctx)
		sr.routine.Consume(ch)

		if ctx.Err() != nil {
			fmt.Println("Exit from worker")
			return
		}
	}
}
