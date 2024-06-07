package scheduler

import "time"

var waitFor = time.After
var timeNow = time.Now

// Timer represent function that return channel of [time.Time].
// It works like [time.After].
type Timer func() <-chan time.Time

// Every returns a function that implements a timer, which triggers after the specified duration has passed.
// It is useful for routines that need to run continuously at regular intervals.
//
// The input is a time.Duration value specifying the interval.
//
// Example usage:
//
//	everyFunc := Every(2 * time.Hour)
//	for {
//	    <-everyFunc()  // This will wait for 2 hours and then unblock, and repeat continuously
//	    // Your periodic task here
//	}
func Every(d time.Duration) func() <-chan time.Time {
	return func() <-chan time.Time {
		return waitFor(d)
	}
}

// At returns a function that implements a timer, which triggers at the next occurrence
// of the specified hour and minute. It is useful for routines that need to run at an exact time.
//
// The input time string should be in the format "HH:MM".
// The returned function creates a time.Time channel that will send the current time when the specified time is reached.
//
// Example usage:
//
//	atFunc := At("14:30")
//	<-atFunc()  // This will wait until the next 14:30 and then unblock
func At(s string) func() <-chan time.Time {
	return func() <-chan time.Time {
		t, err := time.Parse("15:04", s)
		if err != nil {
			panic("Can't parse string time")
		}
		now := timeNow()

		t = time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
		if now.After(t) {
			t = t.AddDate(0, 0, 1)
		}

		d := t.Sub(now)
		return waitFor(d)
	}
}
