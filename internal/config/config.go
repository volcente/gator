package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	config_filename = ".gatorconfig.json"
)

type Config struct {
	DBUrl           string `db:"db_url" json:"db_url"`
	CurrentUsername string `db:"current_username" json:"current_username"`
}

func Read() (Config, error) {
	filepath, err := getConfigFilepath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(filepath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	config := Config{}
	if err = json.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, err
	}

	return config, err
}

func (c *Config) SetUser(userName string) error {
	c.CurrentUsername = userName
	return write(*c)
}

func getConfigFilepath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filepath := fmt.Sprintf("%s/%s", homeDir, config_filename)
	return filepath, nil
}

func write(cfg Config) error {
	filepath, err := getConfigFilepath()
	if err != nil {
		return err
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = json.NewEncoder(file).Encode(cfg); err != nil {
		return err
	}

	return nil
}
