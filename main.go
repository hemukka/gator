package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/hemukka/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("login expects username as argument")
	}

	if err := s.cfg.SetUser(cmd.args[0]); err != nil {
		return err
	}
	fmt.Printf("username set to: %v\n", s.cfg.CurrentUserName)
	return nil
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, exists := c.handlers[cmd.name]
	if !exists {
		return errors.New("unknown command")
	}
	if err := f(s, cmd); err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string,
	f func(s *state, cmd command) error) {
	c.handlers[name] = f
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	appState := state{
		cfg: &cfg,
	}

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("arguments missing")
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}
	if err = cmds.run(&appState, cmd); err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("Read config: %+v\n", cfg)

	// if err := cfg.SetUser("hemu"); err != nil {
	// 	log.Fatalf("couldn't set current user: %v", err)
	// }

	// cfg, err = config.Read()
	// if err != nil {
	// 	log.Fatalf("error reading config: %v", err)
	// }
	// fmt.Printf("Read config again: %v\n", cfg)

}
