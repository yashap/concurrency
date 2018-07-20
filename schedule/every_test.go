package schedule

import (
	"testing"
	"time"
)

func TestEvery(t *testing.T) {
	delay := 20 * time.Millisecond
	message := "test"
	n := 5
	counter := 0
	c := Every(delay, message)
	start := time.Now()

	// receive n messages from channel
	for i := 0; i < n; i++ {
		actualMessage := <-c
		if actualMessage != message {
			t.Errorf("Expected value at index %d to be %q, was %q", i, message, actualMessage)
		}
		counter++
	}

	// ensure we received n messages
	if counter != n {
		t.Errorf("Expected to receive %d messages out of the channel, instead received %d", n, counter)
	}

	// ensure the delays were respected
	runtime := time.Since(start)
	minRuntime := time.Duration(n) * delay
	maxRuntime := 3 * minRuntime // in a perfect world runtime is exactly minRuntime, but give a buffer
	if runtime < minRuntime {
		t.Errorf("Expected runtime to be at least %v, was %v", minRuntime, runtime)
	}
	if runtime > maxRuntime {
		t.Errorf("Expected runtime to be at most %v, was %v", maxRuntime, runtime)
	}
}

func TestEveryStoppable(t *testing.T) {
	delay := 20 * time.Millisecond
	message := "test"
	n := 5
	counter := 0
	res := EveryStoppable(delay, message)
	start := time.Now()

	// receive n messages from channel
	for i := 0; i < n; i++ {
		actualMessage := <-res.output
		if actualMessage != message {
			t.Errorf("Expected value at index %d to be %q, was %q", i, message, actualMessage)
		}
		counter++
	}

	// ensure we received n messages
	if counter != n {
		t.Errorf("Expected to receive %d messages out of the channel, instead received %d", n, counter)
	}

	// ensure the delays were respected
	runtime := time.Since(start)
	minRuntime := time.Duration(n) * delay
	maxRuntime := 3 * minRuntime // in a perfect world runtime is exactly minRuntime, but give a buffer
	if runtime < minRuntime {
		t.Errorf("Expected runtime to be at least %v, was %v", minRuntime, runtime)
	}
	if runtime > maxRuntime {
		t.Errorf("Expected runtime to be at most %v, was %v", maxRuntime, runtime)
	}

	// Stopping should close the output channel
	res.stop <- true
	if value, ok := <-res.output; ok {
		t.Errorf("Expected output channel to be closed, but still delivered %q", value)
	}
}
