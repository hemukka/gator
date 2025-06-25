package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hemukka/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) != 1 {
		fmt.Printf("(You can browse specific number of posts, usage: %v <limit>)\n", cmd.name)
	} else {
		specifiedLimit, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("coudn't convert limit arg to int: %w", err)
		}
		limit = specifiedLimit
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("coudn't get posts: %w", err)
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Println("\n*****************************")
		fmt.Printf("%s from %s\n", post.PublishedAt.Format("Mon Jan 2 15:04"), post.FeedName)
		fmt.Println(post.Title.String)
		fmt.Println("")
		fmt.Println(post.Description.String)
		fmt.Printf("\nLink: %s\n", post.Url)
	}

	return nil
}
