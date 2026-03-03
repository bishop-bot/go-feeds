package parser

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog/log"
)

type Feed struct {
	URL string
}

type FeedJob struct {
	ID  int
	URL string
}

type FeedResult struct {
	ID     int
	Status string
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
	const maxWorkers = 5
	results := []*gofeed.Feed{}
	feedJobs := make(chan FeedJob, len(*feeds))
	semaphore := make(chan struct{}, maxWorkers)
	wg := &sync.WaitGroup{}

	// Start worker goroutines
	for i := 0; i < len(*feeds); i++ {
		feedJobs <- FeedJob{ID: i, URL: (*feeds)[i]}
	}
	close(feedJobs)

	for job := range feedJobs {
		wg.Add(1)
		go func(job FeedJob) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire a slot
			defer func() { <-semaphore }() // Release the slot

			feed := &Feed{URL: job.URL}
			parsedFeed, err := feed.Fetch(ctx)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to fetch feed for job ID: %d", job.ID)
			} else {
				log.Info().Msgf("Successfully fetched feed for job ID: %d", job.ID)
				results = append(results, parsedFeed)
			}
		}(job)
	}

	wg.Wait()
	return results, nil
}
