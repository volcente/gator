package main

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

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

func handlerBrowsePosts(s *state, cmd command) error {
	pagination := pagination{page: defaultPage, limit: defaultLimit}
	sorting := sorting{sortBy: defaultSortBy, order: defaultOrder}
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

// command should work this way
// browse: --page n --limit n --sort title --order asc --filter title=foo
const (
	pageFlag   = "--page"
	limitFlag  = "--limit"
	sortByFlag = "--sort-by"
	orderFlag  = "--sort-order"
)

func validateUserInput(cmd command, pagination *pagination, sorting *sorting) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("usage: %s %s <page> %s <limit> %s <sortBy> %s <asc | desc>", cmd.Name, pageFlag, limitFlag, sortByFlag, orderFlag)
	}

	subArgs := slices.Chunk(cmd.Args, 2)
	subArgsMap := make(map[string]string)
	for subArg := range subArgs {
		argName, argValue := strings.ToLower(subArg[0]), strings.ToLower(subArg[1])
		subArgsMap[argName] = argValue
	}

	if rawPage, exists := subArgsMap[pageFlag]; exists {
		if page, err := strconv.Atoi(rawPage); err == nil {
			pagination.page = page
		} else {
			return fmt.Errorf("invalid page number: %s", rawPage)
		}
	}

	if rawLimit, exists := subArgsMap[limitFlag]; exists {
		if limit, err := strconv.Atoi(rawLimit); err == nil {
			pagination.limit = limit
		} else {
			return fmt.Errorf("invalid limit number: %s", rawLimit)
		}
	}

	acceptedSortByFields := []string{"title", "url", "publishedAt"}
	if rawSortBy, exists := subArgsMap[sortByFlag]; exists {
		if !slices.Contains(acceptedSortByFields, rawSortBy) {
			return fmt.Errorf("invalid sortBy: %s. accepted sortBy fields: %s", rawSortBy, strings.Join(acceptedSortByFields, "|"))
		}
		sorting.sortBy = rawSortBy
	}

	acceptedOrderFields := []string{"asc", "desc"}
	if rawOrder, exists := subArgsMap[orderFlag]; exists {
		if rawOrder != "asc" && rawOrder != "desc" {
			return fmt.Errorf("invalid order: %s. accepted order values: %s", rawOrder, strings.Join(acceptedOrderFields, "|"))
		}
		sorting.order = rawOrder
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
