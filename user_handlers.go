package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/bluiwulf/blogaggregator/internal/database"
	"github.com/google/uuid"
)

func commandLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("%v command expects a single argument, the username", cmd.name)
	}
	name := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("user '%v' is not registered: %v", name, err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("failed to set current user '%v': %v", name, err)
	}
	fmt.Println("User has been set")

	return nil
}

func commandRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("%v command expects a single argument, the username", cmd.name)
	}
	name := cmd.args[0]
	user, err := s.db.GetUser(context.Background(), name)
	if err == nil {
		fmt.Println(fmt.Errorf("user '%v' already exists", name))
		os.Exit(1)
	}

	userParams := database.CreateUserParams{
		ID: 		uuid.New(),
		CreatedAt: 	time.Now(),
		UpdatedAt: 	time.Now(),
		Name: 		name,
	}
	user, err = s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return fmt.Errorf("failed to register user '%v': %v", name, err)
	}
	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("failed to set current user '%v': %v", name, err)
	}

	fmt.Println("User has been registered")
	fmt.Println()
	fmt.Printf("ID: %v\n", user.ID)
	fmt.Printf("Name: %v\n", user.Name)
	fmt.Println()

	return nil
}

func commandReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("%v command expects no argument", cmd.name)
	}
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset users database: %v", err)
	}
	fmt.Println("Successfully reset users database")

	return nil
}

func commandUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("%v command expects no argument", cmd.name)
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get list of users: %v", err)
	}
	
	fmt.Println()
	for _, user := range users {
		if user.Name == s.cfg.CurrentUser {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Printf("* %v\n", user.Name)
		}
	}
	fmt.Println()

	return nil
}
