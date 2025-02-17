package main

import (
	"log"
	"os"

	"github.com/volcente/gator/internal/config"
)

type state struct {
	config *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		os.Exit(1)
	}

	programState := state{config: &cfg}

	cmds := commands{
		commandList: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	commandName := os.Args[1]
	commandArgs := os.Args[2:]

	if err = cmds.run(&programState, command{Name: commandName, Args: commandArgs}); err != nil {
		log.Fatal(err)
	}
}
