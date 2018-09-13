package main

import (
	"fmt"
	"time"

	"github.com/yashap/concurrency/rss"
)

func main() {
	// Subscribe to some feeds, and created a merged update stream
	merged := rss.Merge(
		rss.Subscribe(rss.NewFetcher("https://dave.cheney.net/category/golang/feed")),
		rss.Subscribe(rss.NewFetcher("https://blog.learngoprogramming.com/feed")),
		rss.Subscribe(rss.NewFetcher("https://blog.golang.org/feed.atom?format=xml")),
	)

	// Close the subscriptions after a bit
	time.AfterFunc(3*time.Second, func() {
		fmt.Printf("Closing subscription: %v\n", merged)
		merged.Close()
	})

	// Receive items from stream, print them. Will continue until the subscription is closed
	for item := range merged.Updates() {
		fmt.Printf(
			"%s\n * Link: %s\n * Channel: %s\n * Author: %s\n * Published: %s\n * GUID: %s\n",
			item.Title, item.Link, item.Channel, item.Author, item.Published, item.GUID,
		)
	}

	fmt.Println()
	message := "This is just an easy way to see active goroutines, to ensure we're cleaning up " +
		"properly. There should only be the main.main() one."
	panic(message)
}
