package cron

import (
	"math/rand"
	"time"
)

// ConstantDelaySchedule represents a simple recurring duty cycle, e.g. "Every 5 minutes".
// StartTime defines, when the first run shoud be executed.
// This allows to run the job immediatly (@every 5s,0s) or
// with a random delay (@every 5s,@rand) within the Delay.
// It does not support jobs more frequent than once a second.
type ConstantDelaySchedule struct {
	Delay     time.Duration
	StartTime time.Time
}

// Every returns a crontab Schedule that activates once every duration.
// Delays of less than a second are not supported (will round up to 1 second).
// Any fields less than a Second are truncated.
func Every(duration time.Duration) ConstantDelaySchedule {
	if duration < time.Second {
		duration = time.Second
	}
	return ConstantDelaySchedule{
		Delay:     duration - time.Duration(duration.Nanoseconds())%time.Second,
		StartTime: time.Unix(0, 0),
	}
}

// Every returns a crontab Schedule that activates once every duration,
// but with a explicit initial delay. This allows to run the job immediatly.
// Delays of less than a second are not supported (will round up to 1 second).
// Any fields less than a Second are truncated.
func EveryWithInitial(duration time.Duration, initial time.Duration) ConstantDelaySchedule {
	t := time.Now()
	cds := Every(duration)
	cds.StartTime = t.Add(initial - time.Duration(t.Nanosecond())%time.Second)
	return cds
}

// Every returns a crontab Schedule that activates once every duration,
// but with a random initial delay. This allows to distribut multiple jobs
// with the same interval.
// Delays of less than a second are not supported (will round up to 1 second).
// Any fields less than a Second are truncated.
func EveryWithRandInitial(duration time.Duration) ConstantDelaySchedule {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	t := time.Now()
	cds := Every(duration)
	cds.StartTime = t.Add(cds.Delay - time.Duration(t.Nanosecond())%time.Second - time.Duration(r.Int63()%int64(duration.Seconds()))*time.Second)
	return cds
}

// Next returns the next time this should be run.
// This rounds so that the next activation time will be on the second.
func (schedule ConstantDelaySchedule) Next(t time.Time) time.Time {
	if schedule.StartTime.Sub(t).Seconds() > 0 {
		// Initial run
		return schedule.StartTime
	} else {
		return t.Add(schedule.Delay - time.Duration(t.Nanosecond())*time.Nanosecond)
	}
}
