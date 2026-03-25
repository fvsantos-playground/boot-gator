package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fvsantos-playground/boot-gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("Error reading config file", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	s := state{
		config: &cfg,
	}

	handlers := commands{
		commands: make(map[string]func(*state, command) error),
	}

	handlers.register("login", handlerLogin)

	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("No command provided")
		os.Exit(1)
	}

	cmd := command{
		name: args[0],
		args: args[1:],
	}

	if err := handlers.run(&s, cmd); err != nil {
		os.Exit(1)
	}
}
