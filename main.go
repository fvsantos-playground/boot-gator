package main

import (
	"fmt"
	"log"

	"github.com/fvsantos-playground/boot-gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("Error reading config file", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	err = cfg.SetUser("fvsantos")
	if err != nil {
		log.Fatalf("couldn't set current user: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config again: %+v\n", cfg)
}
