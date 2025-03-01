package main

import (
	"context"
	"fmt"
	"github.com/volcente/gator/internal/database"
	"strconv"
)

const defaultLimit = 2

func handlerBrowsePosts(s *state, cmd command) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %s <optional: limit>", cmd.Name)
	}

	limit := defaultLimit
	if len(cmd.Args) == 1 {
		parsedLimit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return err
		}
		limit = parsedLimit
	}

	posts, err := s.db.GetPostsForUser(context.Background(), int32(limit))
	if err != nil {
		return fmt.Errorf("could not get posts: %w", err)
	}

	for _, post := range posts {
		printPost(post)
	}

	return nil
}

func printPost(post database.Post) {
	fmt.Printf("\n---\n")
	fmt.Printf("Title: %s\n", post.Title)
	fmt.Printf("Published: %s\n", post.PublishedAt)
	fmt.Printf("Description: %s\n", post.Description)
	fmt.Printf("Source: %s\n", post.Url)
	fmt.Printf("---\n")
}
