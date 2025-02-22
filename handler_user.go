package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/volcente/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("useage: %s <name>", cmd.Name)
	}

	username := cmd.Args[0]
	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		fmt.Println("user does not exist.")
		os.Exit(1)
	}

	err = s.config.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("could not set current user: %w", err)
	}

	fmt.Printf("user switched successfully!")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("useage: %s <name>", cmd.Name)
	}

	username := cmd.Args[0]
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})
	if err != nil {
		fmt.Println("user already exists.")
		os.Exit(1)
	}

	if err = s.config.SetUser(username); err != nil {
		return fmt.Errorf("could not set current user: %w", err)
	}

	fmt.Printf("user has been successfuly created!\nUser: %+v", user)
	return nil
}
