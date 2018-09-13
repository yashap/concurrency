package rss

// Merge merges several RSS subscriptions into a single subscription
func Merge(subs ...Subscription) Subscription {
	var stoppableSubs []stoppableSubscription
	for _, sub := range subs {
		stoppableSubs = append(stoppableSubs, stoppableSubscription{sub, make(chan bool)})
	}
	m := &mergedSubscription{
		stoppableSubscriptions: stoppableSubs,
		updates:                make(chan Item),
	}
	m.loop()
	return m
}

// mergedSubscription is a subscription made a zero-to-many subscriptions
type mergedSubscription struct {
	stoppableSubscriptions []stoppableSubscription
	updates                chan Item
}

// stoppableSubscription is a subscription, and a channel that lets you signal to stop reading from it
type stoppableSubscription struct {
	sub  Subscription
	stop chan bool
}

func (m *mergedSubscription) Updates() <-chan Item {
	return m.updates
}

func (m *mergedSubscription) Close() (err error) {
	for _, s := range m.stoppableSubscriptions {
		s.stop <- true
		closeResult := s.sub.Close()
		if err != nil {
			err = closeResult // we'll just report the "last" error if there are more than one
		}
	}
	close(m.updates)
	return
}

func (m *mergedSubscription) loop() {
	for _, closeableSub := range m.stoppableSubscriptions {

		// interact with each underlying subscription on a different goroutine
		go func(s Subscription, stop chan bool) {
			var pending []Item

			// endlessly poll for updates and put them on a single channel, until we receive a stop signal
			for {
				var updates chan Item
				var first Item
				if len(pending) > 0 {
					first = pending[0]
					updates = m.updates
				}

				select {
				// This will be blocked, and thus skipped, if the underlying sub has no updates
				// If there ARE updates in the underlying subscription, we put them in the merged sub's
				// pending queue. We only remove them from said queue when someone is listening for updates
				// (see below)
				case item := <-s.Updates():
					pending = append(pending, item)

				// This will be blocked, and thus skipped, if there's nothing listening for merged updates
				case updates <- first:
					updates = nil
					pending = pending[1:] // after sending, remove the item from pending

				case <-stop:
					close(stop)
					return
				}
			}
		}(closeableSub.sub, closeableSub.stop)
	}
}
