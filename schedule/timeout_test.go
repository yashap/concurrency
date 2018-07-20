package schedule

import (
	"testing"
	"time"
)

func TestReceiveWithTimeoutSuccess(t *testing.T) {
	c := make(chan string)
	expected := "test"
	go func() {
		time.Sleep(10 * time.Millisecond)
		c <- expected
	}()

	actual, err := ReceiveWithTimeout(c, 1*time.Second)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	} else if actual != expected {
		t.Errorf("Expected %q, but got %q", expected, actual)
	}
}

func TestReceiveWithTimeoutFailure(t *testing.T) {
	c := make(chan string)
	go func() {
		time.Sleep(1 * time.Second)
		c <- "test"
	}()

	_, err := ReceiveWithTimeout(c, 10*time.Millisecond)
	if err == nil || err.Error() != "Timeout" {
		t.Errorf("Unexpected error %v", err)
	}
}

func TestReceiveMultiWithTimeout(t *testing.T) {
	frequency := 5 * time.Millisecond
	message := "test"
	timeout := 100 * time.Millisecond
	c := ReceiveMultiWithTimeout(Every(frequency, message), timeout)
	counter := 0

	for s := range c {
		counter++
		if s != message {
			t.Errorf("Expected %q, got %q", message, s)
		}
	}

	maxCount := int(timeout/frequency) + 1 // 1 is an arbitrary buffer, time is imprecise
	if counter == 0 {
		t.Errorf("Expected to receive values from channel, but didn't")
	} else if counter > maxCount {
		t.Errorf("Expected count to be no higher than %d, was %d", maxCount, counter)
	}
}
