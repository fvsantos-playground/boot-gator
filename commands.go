package main

import "fmt"

type commands struct {
	commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	h, ok := c.commands[cmd.name]
	if !ok {
		return fmt.Errorf("non ecxiste")
	}

	err := h(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}
