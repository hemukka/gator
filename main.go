package main

import (
	"log"
	"os"

	"github.com/hemukka/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	appState := &state{
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
	if err = cmds.run(appState, cmd); err != nil {
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
