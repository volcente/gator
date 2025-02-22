package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/volcente/gator/internal/config"
	"github.com/volcente/gator/internal/database"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	dbQueries := database.New(db)

	programState := state{
		db:     dbQueries,
		config: &cfg,
	}

	cmds := commands{
		commandList: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handleReset)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	commandName := os.Args[1]
	commandArgs := os.Args[2:]

	if err = cmds.run(&programState, command{Name: commandName, Args: commandArgs}); err != nil {
		log.Fatal(err)
	}
}
