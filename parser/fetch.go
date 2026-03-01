package parser

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog/log"
)

type Feed struct {
	URL string
}

func (f *Feed) Fetch(ctx context.Context) (*gofeed.Feed, error) {
	log.Debug().Msgf("Fetching feed from URL: %s", f.URL)

	req, err := http.NewRequestWithContext(ctx, "GET", f.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch feed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	parser := gofeed.NewParser()
	feed, err := parser.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse feed: %w", err)
	}

	log.Debug().Msgf("Successfully fetched and parsed feed: %s", f.URL)
	return feed, nil
}

func FetchAll(ctx context.Context, feeds *[]string) ([]*gofeed.Feed, error) {
	var allFeeds []*gofeed.Feed
	for _, feedURL := range *feeds {
		feed := &Feed{URL: feedURL}
		parsedFeed, err := feed.Fetch(ctx)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to fetch feed from URL: %s", feedURL)
			continue // Skip this feed and continue with the next one
		}
		allFeeds = append(allFeeds, parsedFeed)
	}
	return allFeeds, nil
}
