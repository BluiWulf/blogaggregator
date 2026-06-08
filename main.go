package main

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/bluiwulf/blogaggregator/internal/config"
	"github.com/bluiwulf/blogaggregator/internal/database"
)

import _ "github.com/lib/pq"

type state struct {
	db 		*database.Queries
	cfg 	*config.Config
}

type errorCode struct {
	code 	int
	msg		string
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbQs := database.New(db)
	s := state{
		db:	 dbQs,
		cfg: &cfg,
	}
	cmds := commands{ cmds: make(map[string]func(*state, command) error) }

	cliArgs := os.Args
	if len(cliArgs) < 2 {
		log.Fatal(errors.New("no command has been provided"))
	}
	cmd := command{
		name: cliArgs[1],
		args: cliArgs[2:],
	}

	cmds.register("login", commandLogin)
	cmds.register("register", commandRegister)
	cmds.register("reset", commandReset)
	cmds.register("users", commandUsers)
	err = cmds.run(&s, cmd)
	if err != nil {
		log.Fatal(err)
	}
}