package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/bluiwulf/blogaggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	s := state{ config: &cfg }
	cmds := commands{ handlers: make(map[string]func(*state, command) error) }

	cliArgs := os.Args
	if len(cliArgs) < 2 {
		log.Fatal(errors.New("no command has been provided"))
	}
	cmd := command{
		name: cliArgs[1],
		args: cliArgs[2:],
	}
	cmds.register(cmd.name, handlerLogin)
	err = cmds.run(&s, cmd)
	if err != nil {
		log.Fatal(err)
	}
}