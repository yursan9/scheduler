package scheduler

import (
	"iter"
	"testing"
	"time"
)

type TestScheduler struct {
	v string
}

func (s *TestScheduler) Stream() iter.Seq[string] {
	return func(yield func(string) bool) {
		data := []string{"Hello World"}
		for _, v := range data {
			if !yield(v) {
				return
			}
		}
	}
}

func (s *TestScheduler) Consume(v string) error {
	s.v = v
	return ErrAbort
}

func TestManager(t *testing.T) {
	m := New()
	s := &TestScheduler{}
	Register(m, Every(1*time.Second), s)

	if len(m.routines) != 1 {
		t.Error("Expect 1 routine to be registered")
	}

	m.Run()
	m.Stop()

	if s.v != "Hello World" {
		t.Errorf("Expect consuming string: %s got %s", "Hello World", s.v)
	}
}
