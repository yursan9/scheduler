package scheduler

import (
	"testing"
	"time"
)

func TestEvery(t *testing.T) {
	expectedDuration := 1 * time.Hour
	waitFor = func(d time.Duration) <-chan time.Time {
		if d != expectedDuration {
			t.Errorf("Expect duration to be %s got %s", expectedDuration, d)
		}
		return nil
	}

	Every(1 * time.Hour)()
}

func TestAt(t *testing.T) {
	expectedDuration := 1 * time.Hour
	input := "13:00"

	timeNow = func() time.Time {
		return time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC)
	}
	waitFor = func(d time.Duration) <-chan time.Time {
		if d != expectedDuration {
			t.Errorf("Expect duration to be %s got %s", expectedDuration, d)
		}
		return nil
	}

	At(input)()

	expectedDuration = 23 * time.Hour
	input = "11:00"

	timeNow = func() time.Time {
		return time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC)
	}
	waitFor = func(d time.Duration) <-chan time.Time {
		if d != expectedDuration {
			t.Errorf("Expect duration to be %s got %s", expectedDuration, d)
		}
		return nil
	}

	At(input)()
}
