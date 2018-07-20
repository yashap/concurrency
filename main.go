package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/yashap/concurrency/search"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	results := search.GoogleWithReplicasAndTimeout("golang")
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}
