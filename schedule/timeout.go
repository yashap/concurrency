package schedule

import (
	"errors"
	"time"
)

// ReceiveWithTimeout receives a single value from a channel, with a timeout.
// It returns either the value, or a timeout error
func ReceiveWithTimeout(c <-chan string, timeout time.Duration) (string, error) {
	for {
		select {
		case s := <-c:
			return s, nil
		case <-time.After(timeout):
			return "", errors.New("Timeout")
		}
	}
}

// ReceiveMultiWithTimeout receives 0-to-many values from a channel, and writes them
// to an output channel. Once a timeout is reached, it closes the output channel
func ReceiveMultiWithTimeout(c <-chan string, timeout time.Duration) <-chan string {
	out := make(chan string)
	timeoutChan := time.After(timeout)
	go func() {
		for {
			select {
			case s := <-c:
				out <- s
			case <-timeoutChan:
				close(out)
				return
			}
		}
	}()
	return out
}
