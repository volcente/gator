package main

import "errors"

type command struct {
	Name string
	Args []string
}

type commands struct {
	commandList map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandList[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	command, exists := c.commandList[cmd.Name]
	if !exists {
		return errors.New("Command not found!")
	}
	return command(s, cmd)
}
