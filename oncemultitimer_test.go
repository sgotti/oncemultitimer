package oncemultitimer

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/uber/tchannel-go/testutils/goroutines"
)

func TestMultiTimer(t *testing.T) {
	// make the reader goroutine start after some timers already fired
	startDelay = 1 * time.Second

	s := NewTimer()

	// Some timers  twon't fire
	s.AddTimer(2 * time.Second)
	s.AddTimer(1 * time.Second)

	// Add a lot of timers firing ("quite") at the same time to check that only one triggers the
	// function and that all the others don't block their goroutines, they could also fire before the main reader goroutine has started
	s.AddTimer(0 * time.Second)
	s.AddTimer(0 * time.Second)
	s.AddTimer(0 * time.Second)
	s.AddTimer(0 * time.Second)
	s.AddTimer(0 * time.Second)
	s.AddTimer(0 * time.Second)
	s.AddTimer(0 * time.Second)
	s.AddTimer(0 * time.Second)

	time.Sleep(3 * time.Second)

	<-s.C

	// check that no other timers are enqueued
	time.Sleep(1 * time.Second)

	select {
	case <-s.C:
		t.Fatalf("another message was enqueued, this shouldn't happen!")
	default:
	}

	s.Stop()
}

func TestMain(m *testing.M) {
	exitCode := m.Run()
	if err := goroutines.IdentifyLeaks(&goroutines.VerifyOpts{}); err != nil && exitCode == 0 {
		fmt.Fprintf(os.Stderr, "Found goroutine leaks on successful test run: %v", err)
		exitCode = 1
	}
	os.Exit(exitCode)
}
