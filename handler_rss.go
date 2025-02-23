package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/volcente/gator/internal/database"
)

func handlerAggregate(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	feed, err := fetchFeed(context.Background(), feedURL)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s feed_name feed_url", cmd.Name)
	}

	user, err := s.db.GetUser(context.Background(), s.config.CurrentUsername)
	if err != nil {
		return fmt.Errorf("could not get the current user from databse: %w", err)
	}

	feedName, feedURL := cmd.Args[0], cmd.Args[1]
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("could not add the feed: %w", err)
	}

	fmt.Printf("%+v\n", feed)
	return nil
}

func handlerShowFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("could not get feeds: %w", err)
	}

	fmt.Printf("---------\n")
	for _, feed := range feeds {
		fmt.Printf("* Feed: %s\n", feed.Name)
		fmt.Printf("* Feed URL: %s\n", feed.Url)
		fmt.Printf("* Created by: %s\n", feed.Username)
		fmt.Printf("---------\n")
	}

	return nil
}
