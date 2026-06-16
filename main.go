package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	// Reading configuration data
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	// Opening database
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbQs := database.New(db)
	s := state{
		db:	 dbQs,
		cfg: &cfg,
	}

	// Getting command from user
	cmds := commands{ cmds: make(map[string]func(*state, command) error) }
	cliArgs := os.Args
	if len(cliArgs) < 2 {
		log.Fatal(errors.New("no command has been provided"))
	}
	cmd := command{
		name: cliArgs[1],
		args: cliArgs[2:],
	}

	// Registering available commands
	cmds.register("login", commandLogin)
	cmds.register("register", commandRegister)
	cmds.register("reset", commandReset)
	cmds.register("users", commandUsers)
	cmds.register("agg", commandAgg)
	cmds.register("addfeed", middlewareLoggedIn(commandAddFeed))
	cmds.register("feeds", commandFeeds)
	cmds.register("follow", middlewareLoggedIn(commandFollow))
	cmds.register("following", middlewareLoggedIn(commandFollowing))
	cmds.register("unfollow", middlewareLoggedIn(commandUnfollow))
	cmds.register("browse", middlewareLoggedIn(commandBrowse))
	
	// Running command
	err = cmds.run(&s, cmd)
	if err != nil {
		log.Fatal(err)
	}
}

func middlewareLoggedIn(cmdHandler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUser)
		if err != nil {
			return fmt.Errorf("current user '%v' is not registered: %v", s.cfg.CurrentUser, err)
		}
		return cmdHandler(s, cmd, user)
	}
}
