package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"time"

	"github.com/google/uuid"
	"github.com/volcente/gator/internal/database"
)

func handlerAggregate(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time_between_requests>", cmd.Name)
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("%s is not a valid time! %w", cmd.Args[0], err)
	}

	fmt.Printf("Collecting feeds every %s...\n", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s feed_name feed_url", cmd.Name)
	}

	feedName, feedURL := cmd.Args[0], cmd.Args[1]
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		AuthorID:  user.ID,
	})
	if err != nil {
		return fmt.Errorf("could not add the feed: %w", err)
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not create a feed follow: %w", err)
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

func handlerFollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	}

	feedUrl := cmd.Args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("could not retrieve feed based on passed url: %w", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("could not create a feed follow: %w", err)
	}

	fmt.Println("---------")
	fmt.Printf("Feed: %s\n", feedFollow.Name)
	fmt.Printf("Has been followed by: %s\n", feedFollow.Username)
	fmt.Println("---------")

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	followedFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("could not get followed feeds for current user: %w", err)
	}

	fmt.Println("---------")
	fmt.Printf("feeds for %s:\n", s.config.CurrentUsername)
	if len(followedFeeds) == 0 {
		fmt.Println("not feeds are followed.")
	} else {
		for _, followedFeed := range followedFeeds {
			fmt.Printf("* %s\n", followedFeed.FeedName)
		}
	}
	fmt.Println("---------")
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	}
	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("could not find feed: %w", err)
	}

	err = s.db.DeleteFeedFollower(context.Background(), database.DeleteFeedFollowerParams{
		FeedID: feed.ID,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("could not delete feed follow: %w", err)
	}

	fmt.Printf("%s unfollowed successfully!\n", feed.Name)
	return nil
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("could not get next feed to fetch: %w", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return fmt.Errorf("could not mark feed as fetched: %w", err)
	}

	updatedFeed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("could not get latest feed: %w", err)
	}

	for _, item := range updatedFeed.Channel.Item {
		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return fmt.Errorf("could not parse published date: %w", err)
		}

		err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			FeedID:      nextFeed.ID,
			PublishedAt: publishedAt,
		})
		if err != nil {
			if isDuplicatedKeyError(err) {
				continue
			}
			return fmt.Errorf("could not create post: %w", err)
		}
		fmt.Println("post has been successfully saved to the database")
	}
	return nil
}

func isDuplicatedKeyError(err error) bool {
	var pqError *pq.Error
	if isPqError := errors.As(err, &pqError); isPqError {
		return pqError.Code == "23505"
	}
	return false
}
