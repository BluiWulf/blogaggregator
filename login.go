package main

import (
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("the %v handler expects a single argument, the username", cmd.name)
	}
	user := cmd.args[0]
	err := s.config.SetUser(user)
	if err != nil {
		return err
	}
	fmt.Println("User has been set")

	return nil
}
