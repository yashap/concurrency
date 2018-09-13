package rss

import (
	"time"

	"github.com/mmcdole/gofeed"
)

// Fetcher fetches items from an RSS feed
type Fetcher interface {
	Fetch() FetchResult
}

// NewFetcher creates a Fetcher for a domain
func NewFetcher(url string) Fetcher {
	return &fetcher{parser: gofeed.NewParser(), url: url}
}

// fetcher implementes the Fetcher interface
type fetcher struct {
	parser *gofeed.Parser // fetches from feeds, parses content
	url    string         // url to fetch content from
}

// Fetch fetches items from an RSS feed
func (f *fetcher) Fetch() FetchResult {
	feed, err := f.parser.ParseURL(f.url)
	next := time.Now().Add(10 * time.Second)
	var items []Item
	if err == nil {
		for _, item := range feed.Items {
			items = append(
				items,
				Item{
					Title:     item.Title,
					Link:      item.Link,
					Channel:   feed.Link,
					Author:    item.Author.Name,
					Published: item.Published,
					GUID:      item.GUID,
				})
		}
	}
	return FetchResult{items, next, err}
}

// FetchResult is a struct to hold all the results of Fetcher.Fetch()
type FetchResult struct {
	Fetched []Item
	Next    time.Time
	Err     error
}

// Item is an item in an RSS feed
type Item struct {
	Title, Link, Channel, Author, Published, GUID string
}
