package rssprovider

import (
	"net/url"
	"time"

	"github.com/chackett/zignews/pkg/storage"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
)

// RSSProvider provides functionality for retrieving RSS feeds
type RSSProvider struct {
	Label         string
	FeedURL       string
	PollFrequency time.Duration
}

// NewRSSProvider returns a new instance of RSSProvider
func NewRSSProvider(label, feedURL string, pollFrequency time.Duration) (*RSSProvider, error) {
	if label == "" {
		return nil, errors.New("label is required")
	}
	if _, err := url.ParseRequestURI(feedURL); err != nil {
		return nil, errors.Wrap(err, "validate feed URL")
	}

	result := &RSSProvider{
		Label:         label,
		FeedURL:       feedURL,
		PollFrequency: pollFrequency,
	}

	return result, nil
}

// Latest returns the latest items in the feed
func (r *RSSProvider) Latest() ([]storage.Article, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(r.FeedURL)
	if err != nil {
		return nil, errors.Wrap(err, "rss parse URL")
	}

	var result []storage.Article
	for _, item := range feed.Items {
		// Nullable values checked
		var imgURL string
		if item.Image != nil {
			imgURL = item.Image.URL
		}
		var published time.Time
		if item.PublishedParsed != nil {
			published = *item.PublishedParsed
		}

		a := storage.Article{
			Title:       item.Title,
			Categories:  item.Categories,
			Description: item.Description,
			GUID:        item.GUID,
			Link:        item.Link,
			Published:   published,
			Thumbnail:   imgURL,
			Provider:    r.Label,
		}
		result = append(result, a)
	}

	return result, nil
}

// PollingFrequency returns the amount of time to wait between calls to `Latest()`
func (r *RSSProvider) PollingFrequency() time.Duration {
	return r.PollFrequency
}
