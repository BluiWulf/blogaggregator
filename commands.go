package main

import (
	"errors"

	"github.com/bluiwulf/blogaggregator/internal/config"
)

type state struct {
	config		*config.Config
}

type command struct {
	name		string
	args		[]string
}

type commands struct {
	handlers	map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if s == nil || *s == (state{}) {
		return errors.New("cannot run command with nonexistent state")
	}
	err := c.handlers[cmd.name](s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}
