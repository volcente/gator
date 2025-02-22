package main

import (
	"context"
	"fmt"
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
