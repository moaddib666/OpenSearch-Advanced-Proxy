package ports

import "time"

type TimeTracker interface {
	Start()
	Stop()
	GetDuration() time.Duration
}
