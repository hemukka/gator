package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hemukka/gator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expects url as arguments, usage: %s <feed_url>", cmd.name)
	}

	return followFeed(s, cmd.args[0], user)
}

func followFeed(s *state, url string, user database.User) error {
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("couldn't find feed by url: %w", err)
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return fmt.Errorf("couldn't follow the feed: %w", err)
	}

	fmt.Println("Feed followed successfully:")
	fmt.Printf(" * Feed: %s\n", feedFollow.FeedName)
	fmt.Printf(" * User: %s\n", feedFollow.UserName)

	return nil
}

func handlerListFeedFollows(s *state, cmd command, user database.User) error {
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get feed follows: %w", err)
	}

	if len(feedFollows) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
	}

	fmt.Printf("Feed follows for user %s:\n", feedFollows[0].UserName)
	for _, follow := range feedFollows {
		fmt.Printf(" * %s\n", follow.FeedName)
	}
	return nil
}
