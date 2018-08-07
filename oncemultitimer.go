// Package oncemultitimer is a timer that can be scheduled multiple times but
// will fire only once when the first timer expires.
package oncemultitimer

import (
	"fmt"
	"sync"
	"time"
)

// test only variable used to make the reader goroutine start after some
// seconds. Used to test the condition where timers already fired before the
// reader is started
var startDelay = 0 * time.Second

// OnceMultiTimer is a timer that can be scheduled multiple times but will
// fire only once when the first timer expires.
type OnceMultiTimer struct {
	// using the same time.Timer semantic we cannot define a receive only
	// channel for users since we also have to write to it. To do this we
	// should define a C() function
	C chan time.Time

	stopCh      chan struct{}
	schedulerCh chan time.Time

	m      sync.Mutex
	done   bool
	timers []*time.Timer

	wg sync.WaitGroup

	once sync.Once
}

// NewTimer create a new OnceMultiTimer. No timer is scheduled (use AddTimer to schedule new timers).
func NewTimer() *OnceMultiTimer {
	s := &OnceMultiTimer{
		C:      make(chan time.Time, 1),
		stopCh: make(chan struct{}),
		// buffered channel of size 1 so the it'll be enqueued without blocking
		// since the reader goroutine could be already exited (closed schedulerStopCh)
		// or not yet started when the message is enqueued
		schedulerCh: make(chan time.Time, 1),
	}

	s.wg.Add(1)
	go func() {
		if startDelay > 0 {
			time.Sleep(startDelay)
		}
		select {
		case <-s.stopCh:
		case t := <-s.schedulerCh:
			s.C <- t
		}
		s.wg.Done()

		s.stop()
	}()

	return s
}

// AddTimer schedules a new timer after duration d. Returns an error when the
// timer already fired or has been stopped.
func (s *OnceMultiTimer) AddTimer(d time.Duration) error {
	timer := time.NewTimer(d)
	s.m.Lock()
	if s.done {
		s.m.Unlock()
		return fmt.Errorf("scheduler already stopped/fired")
	}
	s.timers = append(s.timers, timer)
	s.m.Unlock()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		select {
		case <-s.stopCh:
		case t := <-timer.C:
			// fire only one time
			s.once.Do(func() { s.schedulerCh <- t })
		}
	}()

	return nil
}

func (s *OnceMultiTimer) stop() bool {
	s.m.Lock()
	// return if already stopping/stopped
	if s.done {
		s.m.Unlock()
		return false
	}

	// stop timers, just to clean them so they won't fire in the future
	// since we don't need them anymore (note that multiple timers may have
	// already fired)
	for _, timer := range s.timers {
		timer.Stop()
	}

	close(s.stopCh)

	s.done = true
	s.m.Unlock()

	// wait for all goroutines to exit
	s.wg.Wait()

	return true
}

// Stop prevents the Timer from firing.
// It follows the same stop semantic and caveats of time.Timer Stop() function.
func (s *OnceMultiTimer) Stop() bool {
	return s.stop()
}
