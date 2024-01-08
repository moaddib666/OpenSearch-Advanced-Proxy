package tracker

import "time"

type DefaultTimeTracker struct {
	startTime time.Time
	endTime   time.Time
}

func (d *DefaultTimeTracker) Start() {
	d.startTime = time.Now()
}

func (d *DefaultTimeTracker) Stop() {
	d.endTime = time.Now()
}

func (d *DefaultTimeTracker) GetDuration() time.Duration {
	return d.endTime.Sub(d.startTime)
}

func NewDefaultTimeTracker() *DefaultTimeTracker {
	return &DefaultTimeTracker{}
}
