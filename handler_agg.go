package main

import (
	"context"
	"fmt"
	"time"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expects request interval (1s, 1m, 1h, etc) as argument, usage: %s <time_between_reqs>", cmd.name)
	}

	// feedURL := "https://www.wagslane.dev/index.xml"
	// feed, err := fetchFeed(context.Background(), feedURL)
	// if err != nil {
	// 	return fmt.Errorf("couldn't fetch feed: %w", err)
	// }
	// fmt.Printf("Feed: %+v\n", feed)

	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't parse time between reqs: %w", err)
	}
	fmt.Printf("Collecting feeds every %v\n", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			fmt.Println(err)
		}
	}

}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get next feed to fetch: %w", err)
	}

	if err = s.db.MarkFeedFetched(context.Background(), feed.ID); err != nil {
		return fmt.Errorf("couldn't mark feed as fetched: %w", err)
	}

	fetchedFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}

	fmt.Printf("Items in %s feed:\n", fetchedFeed.Channel.Title)
	for _, item := range fetchedFeed.Channel.Item {
		fmt.Printf(" * %s\n", item.Title)
	}
	fmt.Println("***********************")
	return nil
}
