package main

import (
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expects username as argument, usage: %s <name>", cmd.name)
	}

	if err := s.cfg.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Printf("Username set to: %v\n", s.cfg.CurrentUserName)
	return nil
}
