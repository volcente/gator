package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/volcente/gator/internal/database"
)

const (
	defaultPage  = 1
	defaultLimit = 2
)

func handlerBrowsePosts(s *state, cmd command) error {
	limit, page := defaultLimit, defaultPage
	if err := validateUserInput(cmd, &limit, &page); err != nil {
		return err
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		Limit:  int32(limit),
		Offset: int32((page - 1) * limit),
	})
	if err != nil {
		return fmt.Errorf("could not get posts: %w", err)
	}

	for _, post := range posts {
		printPost(post)
	}

	return nil
}

func validateUserInput(cmd command, limit *int, page *int) error {
	argsCount := len(cmd.Args)
	if argsCount == 1 {
		parsedPage, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return err
		}
		*page = parsedPage
	} else if argsCount == 2 {
		parsedPage, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return err
		}
		*page = parsedPage

		parsedLimit, err := strconv.Atoi(cmd.Args[1])
		if err != nil {
			return err
		}
		*limit = parsedLimit
	} else {
		return fmt.Errorf("usage: %s <optional: page> <optional: limit>", cmd.Name)
	}
	return nil
}

func printPost(post database.GetPostsForUserRow) {
	fmt.Printf("\n---\n")
	fmt.Printf("Title: %s\n", post.Title)
	fmt.Printf("Published: %s\n", post.PublishedAt)
	fmt.Printf("Description: %s\n", post.Description)
	fmt.Printf("Source: %s\n", post.Url)
	fmt.Printf("---\n")
}
