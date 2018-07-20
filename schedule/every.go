package schedule

import (
	"time"
)

// Every generates a message every d duration, and writes it to a channel. Follows
// the "generator pattern", where a function creates and returns a channel, and
// writes values to it
func Every(d time.Duration, msg string) <-chan string {
	c := make(chan string)
	go func() { // Launch an infinite loop producer
		for {
			time.Sleep(d)
			c <- msg
		}
	}()
	return c // Return channel to caller
}

// EveryStoppable generates a message every d duration, and writes it to a
// channel. It returns a struct containing the main read-only output channel,
// as well as a write-only channel that consumers can write to, telling every
// to stop writing
func EveryStoppable(d time.Duration, msg string) Channels {
	output := make(chan string)
	stop := make(chan bool)
	go func() { // Launch an infinite loop producer, that also listens for stops
		time.Sleep(d)
		for {
			select {
			case output <- msg:
				time.Sleep(d)
			case <-stop:
				close(output)
				return
			}
		}
	}()
	return Channels{output, stop}
}

// Channels represents the two channels returned by EveryStoppable:
//  * output is the main channel, that you can read values from
//  * stop will stop the output channel when a value is written to it
type Channels struct {
	output <-chan string
	stop   chan<- bool
}
