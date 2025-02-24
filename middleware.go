package main

import (
	"context"

	"github.com/volcente/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.config.CurrentUsername)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}
