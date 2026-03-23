package main

import (
	"fmt"

	"github.com/fvsantos-playground/boot-gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config file")
		return
	}
	fmt.Println(cfg)
	config.SetUser("fvsantos")
	cfg, err = config.Read()
	if err != nil {
		fmt.Println("Error reading config file")
		return
	}
	fmt.Println(cfg)
}
