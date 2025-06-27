package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hemukka/gator/internal/database"
	"github.com/lib/pq"
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
			log.Println(err)
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

	fmt.Printf("\033[2K\rScraping feed: %s (%s)... ", feed.Name, fetchedFeed.Channel.Title)

	alreadySaved := 0
	for _, item := range fetchedFeed.Channel.Item {
		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			fmt.Println()
			log.Println("coudn't parse pupDate: ", err)
		}

		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       newNullString(item.Title),
			Url:         item.Link,
			Description: newNullString(item.Description),
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			pqErr, ok := err.(*pq.Error)
			if ok && pqErr.Code == pq.ErrorCode("23505") {
				alreadySaved += 1
				// fmt.Printf("\033[2K\r * post already exists: %s", item.Title)
			} else {
				log.Println("\nerror saving post: ", err)
			}
			continue
		}
		// fmt.Printf("\033[2K\r * post saved: %s", post.Title.String)
	}
	fmt.Printf("%v new posts + %v already saved posts", len(fetchedFeed.Channel.Item)-alreadySaved, alreadySaved)
	return nil
}

func newNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{String: "", Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
