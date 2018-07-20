package multiplex

// FanInNaive receives from 2 channels, writes to an output channel
// This is a naive implementation, using 2 goroutines, instead of select
func FanInNaive(input1, input2 <-chan string) <-chan string {
	out := make(chan string)
	// Receive infinitely from input1, write out
	go func() {
		for {
			out <- <-input1
		}
	}()
	// Receive infinitely from input2, write out
	go func() {
		for {
			out <- <-input2
		}
	}()
	return out
}

// FanIn receives from 2 channels, writes to an output channel
// This implementation uses select
//
// The select statement provides another way to handle multiple channels.
// It's like a switch statement, but each case is a communication:
//  * all channels are evaluated
//  * selection blocks until one communication can proceed, which then does
//  * if multiple can proceed, select chooses one pseudo-randomly
//  * a default clause, if present, executes immediately if not channel is ready
func FanIn(input1, input2 <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for {
			select {
			case s := <-input1:
				out <- s
			case s := <-input2:
				out <- s
			}
		}
	}()
	return out
}
