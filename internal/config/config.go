package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory")
		return "", err
	}
	return home + "/" + configFileName, nil
}

func write(cfg Config) error {
	content, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling config file")
		return err
	}
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Error creating config file")
		return err
	}
	defer file.Close()
	_, err = file.Write(content)
	if err != nil {
		fmt.Println("Error writing config file")
		return err
	}
	return nil
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		fmt.Println("Error getting user home directory")
		return Config{}, err
	}
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening config file")
		return Config{}, err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading config file")
		return Config{}, err
	}
	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		fmt.Println("Error unmarshalling config file")
		return Config{}, err
	}
	return config, nil
}

func SetUser(username string) error {
	config, err := Read()
	if err != nil {
		return err
	}
	config.CurrentUserName = username
	return write(config)
}
