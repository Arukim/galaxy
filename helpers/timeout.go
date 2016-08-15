package helpers

import (
	"time"
)

// Timeout is way to wait for some event, use NewTimeout to create
type Timeout struct {
	Alarm chan bool
}

// NewTimeout creation
func NewTimeout(d time.Duration) *Timeout {
	timeout := new(Timeout)
	timeout.Alarm = make(chan bool, 1)
	go func() {
		time.Sleep(d)
		timeout.Alarm <- true
	}()
	return timeout
}
