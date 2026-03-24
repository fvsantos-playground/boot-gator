package commands

import "fmt"

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("")
	}

	err := s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	return nil
}
