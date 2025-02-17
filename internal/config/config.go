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

func getConfigFilepath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filepath := fmt.Sprintf("%s/%s", homeDir, config_filename)
	return filepath, nil
}

func GetConfig() (Config, error) {
	filepath, err := getConfigFilepath()
	if err != nil {
		return Config{}, err
	}

	content, err := os.ReadFile(filepath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err = json.Unmarshal(content, &config); err != nil {
		return Config{}, err
	}

	return config, err
}

func (c *Config) updateConfigFile() error {
	configFilepath, err := getConfigFilepath()
	if err != nil {
		return fmt.Errorf("Error[SetUser]: %v", err)
	}

	marshaledConfig, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return fmt.Errorf("Error[SetUser]: %v", err)
	}

	err = os.WriteFile(configFilepath, marshaledConfig, 0666)
	if err != nil {
		return fmt.Errorf("Error[SetUser]: %v", err)
	}
	return nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUsername = username

	if err := c.updateConfigFile(); err != nil {
		return err
	}
	return nil
}
