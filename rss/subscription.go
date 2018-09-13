package rss

import (
	"time"
)

// Subscription is a subscription to an RSS feed
type Subscription interface {
	Updates() <-chan Item // stream of Items
	Close() error         // shuts down the stream
}

// Subscribe uses a Fetcher to create a Subscription. It will immediately start
// fetching items from the feed, and sending them to the updates channel
func Subscribe(fetcher Fetcher) Subscription {
	s := &sub{
		fetcher: fetcher,
		updates: make(chan Item),
		closed:  false,
		err:     nil,
		closing: make(chan chan error),
	}
	go s.loop()
	return s
}

// sub implementes the Subscription interface. We make it private, so that you can only
// construct it with Subscribe(fetcher Fetcher), as the initialization is a bit complex
type sub struct {
	fetcher Fetcher   // fetches items
	updates chan Item // delivers items to the consumer of the Subscription
	closed  bool
	err     error
	// This `chan chan` enables a request/response style of communication.
	//  * The service (loop) listens for requests on its channel, s.closing
	//  * The client (Close) sends a request on closing: exit and reply with
	//		the error
	closing chan chan error
}

func (s *sub) Updates() <-chan Item {
	return s.updates
}

func (s *sub) Close() error {
	// close expects to receive an error on this channel, if there is one when closing
	errchan := make(chan error)
	s.closing <- errchan // this is like a request to a service, to close and return any errors
	return <-errchan     // wait for the response, then return it
}

// loop fetches items using s.fetcher, and sends them on s.updates (which is a channel
// returned by s.Updates()). Exits when it receives on s.closing (this is triggered by
// s.Close())
func (s *sub) loop() {
	var pending []Item               // fetches write here; reading updates consumes from here
	var seen = make(map[string]bool) // set of item ids, so we don't double-deliver
	var next time.Time               // zero value is epoch
	var err error                    // set when Fetch fails
	var fetchDone chan FetchResult   // if non-nil, fetcher.Fetch() is running
	const maxPending = 10            // max number of items we'll keep in our queue before we pause fetching

	for {
		var fetchDelay time.Duration // initially 0 (no delay)
		if now := time.Now(); next.After(now) {
			fetchDelay = next.Sub(now)
		}
		// nil channel trick, so that we don't fetch if we already have too many pending items
		// see below for details, but basically we, by default, don't have any fetches scheduled,
		// only schedule a fetch if we don't have too many pending items. Also, we only start a
		// fetch if there isn't one currently running
		var startFetch <-chan time.Time
		if fetchDone == nil && len(pending) < maxPending {
			startFetch = time.After(fetchDelay) // schedule a fetch
		}

		// this is a nil channel. Reading from a nil channel blocks forever, and select
		// will skip blocked channels. So selecting from updates will only proceed if
		// we set updates to be non-nil
		var updates chan Item
		var first Item
		if len(pending) > 0 {
			first = pending[0]
			updates = s.updates
		}

		select {
		// If it's time, we'll fetch from the feed asynchronously
		// The reason we do it asynchronously (i.e. fetch on a separate goroutine, put result
		// in a channel) is that we want this select block to keep churning, so we can react
		// immediately to things like closing
		case <-startFetch:
			fetchDone = make(chan FetchResult, 1) // "fetching in progress"
			go func() {
				fetchDone <- s.fetcher.Fetch()
			}()

		// When the fetch is done, we'll grab the results
		case result := <-fetchDone:
			fetchDone = nil // "fetching not in progress"
			err = result.Err
			next = result.Next
			if err != nil {
				next = time.Now().Add(10 * time.Second)
				break
			}
			for _, item := range result.Fetched {
				if !seen[item.GUID] {
					// We can't just send each `item`` into `s.updates`, could block forever.
					// Our use of `pending` helps with that
					pending = append(pending, item)
					seen[item.GUID] = true
				}
			}

		// See above notes about enabling/disabling updates channel. But basically, this
		// tries to send an item into the channel, only when there's something to send.
		// Will only actually send if there's something receiving at the other end
		case updates <- first:
			pending = pending[1:] // after sending, remove the item from pending

		// Close() has asked us to close, and return any errors
		case errchan := <-s.closing:
			close(s.updates) // tells receiver we're done
			errchan <- err   // send errors back to Close() via the channel it provided
			return
		}
	}
}
