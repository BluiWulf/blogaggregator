package main

import (
	"errors"
)

type command struct {
	name	string
	args	[]string
}

type commands struct {
	cmds	map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if s == nil || *s == (state{}) {
		return errors.New("cannot run command with nonexistent state")
	}
	err := c.cmds[cmd.name](s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}
