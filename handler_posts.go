package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/volcente/gator/internal/database"
)

const (
	defaultPage   = 1
	defaultLimit  = 2
	defaultOrder  = "asc"
	defaultSortBy = "title"
)

type pagination struct {
	page  int
	limit int
}

func (p pagination) getOffset() int {
	return p.limit * (p.page - 1)
}

type sorting struct {
	sortBy string
	order  string
}

// command should work this way
// browse: --page n --limit n --sort title --order asc --filter title=foo

func handlerBrowsePosts(s *state, cmd command) error {
	pagination := pagination{
		page:  defaultPage,
		limit: defaultLimit,
	}
	sorting := sorting{
		sortBy: defaultSortBy,
		order:  defaultOrder,
	}

	if err := validateUserInput(cmd, &pagination, &sorting); err != nil {
		return err
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		Limit:  int32(pagination.limit),
		Offset: int32(pagination.getOffset()),
	})
	if err != nil {
		return fmt.Errorf("could not get posts: %w", err)
	}

	for _, post := range posts {
		printPost(post)
	}

	return nil
}

func validateUserInput(cmd command, pagination *pagination, sorting *sorting) error {
	argsCount := len(cmd.Args)
	if argsCount == 1 {
		parsedPage, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return err
		}
		pagination.page = parsedPage
	} else if argsCount == 2 {
		parsedPage, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return err
		}
		pagination.page = parsedPage

		parsedLimit, err := strconv.Atoi(cmd.Args[1])
		if err != nil {
			return err
		}
		pagination.limit = parsedLimit

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
