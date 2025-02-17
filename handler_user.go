package main

import (
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("useage: %s <name>", cmd.Name)
	}

	username := cmd.Args[0]
	err := s.config.SetUser(username)
	if err != nil {
		return fmt.Errorf("Could not set current user: %w", err)
	}

	fmt.Printf("User switched successfully!")
	return nil
}
