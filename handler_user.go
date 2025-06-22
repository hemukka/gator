package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hemukka/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expects username as argument, usage: %s <name>", cmd.name)
	}

	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}

	if err := s.cfg.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Printf("Username set to: %v\n", s.cfg.CurrentUserName)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expects username as argument, usage: %s <name>", cmd.name)
	}

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}

	user, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return fmt.Errorf("couldn't create new user: %w", err)
	}

	fmt.Println("User created successfully:")
	printUser(user)

	if err := s.cfg.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}
