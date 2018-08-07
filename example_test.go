package oncemultitimer_test

import (
	"log"
	"time"

	"github.com/sgotti/oncemultitimer"
)

func Example() {
	log.Println("start: ", time.Now())

	t := oncemultitimer.NewTimer()

	// schedile multiple timer with different durations
	t.AddTimer(3 * time.Second)
	t.AddTimer(2 * time.Second)
	t.AddTimer(1 * time.Second)

	et := <-t.C

	// timer will expire after 1 second (the lowest timer)
	log.Println("timer expired: ", et)
}
