package multiplex

import (
	"testing"
	"time"

	"github.com/yashap/concurrency/schedule"
)

func TestFanIn(t *testing.T) {
	delay := 20 * time.Millisecond
	message1 := "test1"
	message2 := "test2"
	n := 10

	c := FanIn(schedule.Every(delay, message1), schedule.Every(delay, message2))
	messages := []string{}
	start := time.Now()

	// receive n messages from channel
	for i := 0; i < n; i++ {
		actualMessage := <-c
		if !(actualMessage == message1 || actualMessage == message2) {
			t.Errorf(
				"Expected value at index %d to be %q or %q, was %q",
				i, message1, message2, actualMessage,
			)
		}
		messages = append(messages, <-c)
	}

	// Make sure we received n messages
	if len(messages) != n {
		t.Errorf(
			"Expected to receive %d messages out of the channel, instead received %d",
			n, len(messages),
		)
	}

	// Make sure the delays were respected
	runtime := time.Since(start)
	minRuntime := time.Duration(n) * delay / 2 // divide by 2 because we execute with concurrency of 2
	maxRuntime := 3 * minRuntime               // in a perfect world runtime is exactly minRuntime, but give a buffer
	if runtime < minRuntime {
		t.Errorf("Expected runtime to be at least %v, was %v", minRuntime, runtime)
	}
	if runtime > maxRuntime {
		t.Errorf("Expected runtime to be at most %v, was %v", maxRuntime, runtime)
	}
}

func TestFanInNaive(t *testing.T) {
	delay := 20 * time.Millisecond
	message1 := "test1"
	message2 := "test2"
	n := 10

	c := FanInNaive(schedule.Every(delay, message1), schedule.Every(delay, message2))
	messages := []string{}
	start := time.Now()

	// receive n messages from channel
	for i := 0; i < n; i++ {
		actualMessage := <-c
		if !(actualMessage == message1 || actualMessage == message2) {
			t.Errorf(
				"Expected value at index %d to be %q or %q, was %q",
				i, message1, message2, actualMessage,
			)
		}
		messages = append(messages, <-c)
	}

	// Make sure we received n messages
	if len(messages) != n {
		t.Errorf(
			"Expected to receive %d messages out of the channel, instead received %d",
			n, len(messages),
		)
	}

	// Make sure the delays were respected
	runtime := time.Since(start)
	minRuntime := time.Duration(n) * delay / 2 // divide by 2 because we execute with concurrency of 2
	maxRuntime := 3 * minRuntime               // in a perfect world runtime is exactly minRuntime, but give a buffer
	if runtime < minRuntime {
		t.Errorf("Expected runtime to be at least %v, was %v", minRuntime, runtime)
	}
	if runtime > maxRuntime {
		t.Errorf("Expected runtime to be at most %v, was %v", maxRuntime, runtime)
	}
}
