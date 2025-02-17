package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/volcente/gator/internal/config"
)

type state struct {
	config *config.Config
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Login command expects username as an argument!")
	}

	username := cmd.args[0]
	err := s.config.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("User has been set successfully!")
	return nil
}

type commands struct {
	commandList map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandList[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	command, exists := c.commandList[cmd.name]
	if !exists {
		return errors.New("Command doesn't exist!")
	}

	err := command(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		os.Exit(1)
	}

	s := state{config: &cfg}
	commands := commands{
		commandList: make(map[string]func(*state, command) error),
	}

	commands.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Printf("Not enough arguments provided, expected at least 2.\n")
		os.Exit(1)
	}

	args := os.Args[1:]
	commandName := args[0]
	commandArgs := args[1:]

	if err = commands.run(&s, command{name: commandName, args: commandArgs}); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
