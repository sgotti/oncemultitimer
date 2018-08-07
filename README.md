## OnceMultiTimer
[![GoDoc](https://godoc.org/github.com/sgotti/oncemultitimer?status.svg)](https://godoc.org/github.com/sgotti/oncemultitimer)
[![Build Status](https://travis-ci.org/sgotti/oncemultitimer.svg?branch=master)](https://travis-ci.org/sgotti/oncemultitimer)

OnceMultiTimer is a timer that can be scheduled multiple times but will fire only once when the first timer expires.

## Documentation

See the [godoc package documentation](http://godoc.org/github.com/sgotti/oncemultitimer).

## Getting Started

Install oncemultitimer in the usual way:

    go get github.com/sgotti/oncemultitimer

Example program:

``` go
package main

import (
	"log"
	"time"

	"github.com/sgotti/oncemultitimer"
)

func main() {
	log.Printf("start")
	t := oncemultitimer.NewTimer()

	// schedule multiple timers with different durations
	t.AddTimer(3 * time.Second)
	t.AddTimer(2 * time.Second)
	t.AddTimer(1 * time.Second)

	<-t.C

	// timer will expire after 1 second (the lowest timer duration)
	log.Printf("timer expired")
}
```

