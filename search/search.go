package search

import (
	"fmt"
	"math/rand"
	"time"
)

// Result represents a search result
type Result string

// Search is a function that, given a query, can be used to execute a search
type Search func(query string) Result

func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q\n", kind, query))
	}
}

var (
	web1   = fakeSearch("web1")
	web2   = fakeSearch("web2")
	web3   = fakeSearch("web3")
	image1 = fakeSearch("image1")
	image2 = fakeSearch("image2")
	image3 = fakeSearch("image3")
	video1 = fakeSearch("video1")
	video2 = fakeSearch("video2")
	video3 = fakeSearch("video3")
)

// GoogleSynchronous is a function that, given a query, pretends to search for matching
// websites, images and videos. It performs the searches synchronously
func GoogleSynchronous(query string) (results []Result) {
	return []Result{web1(query), image1(query), video1(query)}
}

// Google is a function that, given a query, pretends to search for matching websites,
// images and videos. It performs the searches concurrently
func Google(query string) (results []Result) {
	res := make(chan Result)
	go func() { res <- web1(query) }()
	go func() { res <- image1(query) }()
	go func() { res <- video1(query) }()

	for i := 0; i < 3; i++ {
		results = append(results, <-res)
	}
	return
}

// GoogleWithTimeout is like Google, but with a timeout
func GoogleWithTimeout(query string) (results []Result) {
	res := make(chan Result)
	go func() { res <- web1(query) }()
	go func() { res <- image1(query) }()
	go func() { res <- video1(query) }()

	timeout := time.After(80 * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case result := <-res:
			results = append(results, result)
		case <-timeout:
			fmt.Println("timed out")
			return
		}
	}
	return
}

// First takes a query and a set of Search services, returns the first one that
// comes back
func First(query string, replicas ...Search) Result {
	c := make(chan Result)
	searchReplica := func(i int) { c <- replicas[i](query) }
	for i := range replicas {
		go searchReplica(i)
	}
	return <-c // return only the first to come back
}

// GoogleWithReplicas is like Google, but has a replica of each search service, takes
// first result for each (for improved performance)
func GoogleWithReplicas(query string) (results []Result) {
	res := make(chan Result)
	go func() { res <- First(query, web1, web2, web3) }()
	go func() { res <- First(query, image1, image2, image3) }()
	go func() { res <- First(query, video1, video2, video3) }()

	for i := 0; i < 3; i++ {
		results = append(results, <-res)
	}
	return
}

// GoogleWithReplicasAndTimeout is a combo of GoogleWithReplicas and GoogleWithTimeout
func GoogleWithReplicasAndTimeout(query string) (results []Result) {
	res := make(chan Result)
	go func() { res <- First(query, web1, web2, web3) }()
	go func() { res <- First(query, image1, image2, image3) }()
	go func() { res <- First(query, video1, video2, video3) }()

	timeout := time.After(80 * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case result := <-res:
			results = append(results, result)
		case <-timeout:
			fmt.Println("timed out")
			return
		}
	}
	return
}
