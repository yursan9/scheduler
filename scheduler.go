// Package scheduler provides manager to run long running function.
//
// The scheduler package should be used inside another goroutine.
package scheduler

import (
	"context"
	"errors"
	"iter"
	"log/slog"
	"sync"
)

var (
	ErrAbort = errors.New("consumer abort")
)

// Routine is an interface for process that want to be scheduled.
type Routine[T any] interface {
	// Stream return a channel and receive [context.Context] for cancellation.
	// The method must be implemented by waiting on [context.Context.Done]
	// before passing data to channel for scheduler cancellation.
	//
	// Example implementations:
	//  func (s *scheduledRoutine) Stream() iter.Seq[string] {
	//  	return func(yield func(E) bool) {
	// 			for _, v := range s.fetch() {
	//     			if !yield(v) {
	//         			return
	//     			}
	// 			}
	//  	}
	//  }
	Stream() iter.Seq[T]
	// Consume read from the channel returned by Stream()
	//  func (s *scheduledRoutine) Consume(s string) error {
	//  	n := v.(string)
	//  	fmt.Println("Hello,", n)
	//  	time.Sleep(3 * time.Second)
	//		return nil
	//  }
	Consume(T) error
}

type scheduledRoutine[T any] struct {
	timer   Timer
	routine Routine[T]
}

func (sr *scheduledRoutine[T]) Execute(ctx context.Context) error {
	for {
		<-sr.timer()
		for v := range sr.routine.Stream() {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				err := sr.routine.Consume(v)
				slog.WarnContext(ctx, "error from consumer", "error", err)
				if errors.Is(err, ErrAbort) {
					return err
				}
			}
		}
	}
}

type ScheduledRoutine interface {
	Execute(context.Context) error
}

type Manager struct {
	wg       *sync.WaitGroup
	routines []ScheduledRoutine
	doneFunc context.CancelFunc
	closeCh  chan struct{}
}

// New returns new schedule manager instance
func New() *Manager {
	return &Manager{
		wg:      new(sync.WaitGroup),
		closeCh: make(chan struct{}),
	}
}

// Register add Timer and Routine to manager for running
func Register[T any](m *Manager, timer Timer, r Routine[T]) {
	sr := &scheduledRoutine[T]{
		timer,
		r,
	}
	m.routines = append(m.routines, sr)
}

// Run start scheduler (will block on call)
func (m *Manager) Run() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	m.doneFunc = cancelFunc
	for _, r := range m.routines {
		m.wg.Add(1)
		go func() {
			r.Execute(ctx)
			m.wg.Done()
		}()
	}

	m.wg.Wait()
}

// Stop will notify worker to stop processing scheduler (will block until finish)
func (m *Manager) Stop() {
	m.doneFunc()
}
