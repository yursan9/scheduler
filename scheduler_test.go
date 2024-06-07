package scheduler

import (
	"context"
	"testing"
	"time"
)

type TestScheduler struct{}

func (s *TestScheduler) Stream(ctx context.Context) <-chan any {
	ch := make(chan any, 1)
	return ch
}

func (s *TestScheduler) Consume(ch <-chan any) {}

func TestManager(t *testing.T) {
	m := New()
	m.Register(Every(1*time.Second), &TestScheduler{})

	if len(m.routines) != 1 {
		t.Error("Expect 1 routine to be registered")
	}

	go m.Run()
	<-time.After(10 * time.Millisecond)
	m.Stop()
}
