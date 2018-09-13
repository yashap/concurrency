package pingpong

import (
	"fmt"
	"time"
)

// ball represents a ping-pong ball, with a count of how many times its been hit
type ball struct {
	hits int
}

// receives the ball, hits it, sending it back
func player(name string, table chan *ball) {
	for {
		b := <-table // blocks until a ball is sent on table
		b.hits++
		fmt.Println(name, b.hits)
		time.Sleep(100 * time.Millisecond)
		table <- b // blocks until someone receives on table
	}
}

// PlayGame plays a game of ping-pong got a given duration, printing moves to stdout
func PlayGame(duration time.Duration) {
	table := make(chan *ball)
	go player("ping", table)
	go player("pong", table)
	table <- new(ball)   // send the ball to a (basically random) player
	time.Sleep(duration) // let the game proceed for awhile
	fmt.Printf("Game over, I grabbed the ball! It was hit %d times\n", (<-table).hits)
}
