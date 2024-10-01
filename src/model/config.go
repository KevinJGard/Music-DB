package model

import (
	"os"
	"path/filepath"
	"encoding/json"
	"log"
)

type Config struct {
	MusicDirectory string `json:"music_directory"`
}

func NewConfig() *Config {
	configDir := filepath.Join(os.Getenv("HOME"), ".config", "MusicDB")
	os.MkdirAll(configDir, os.ModePerm)

	configFile := filepath.Join(configDir, "config.json")
	config := &Config{}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		_, err = os.Create(configFile)
		if err != nil {
			log.Fatalf("Error creating config file: %v", err)
		}
	} else {
		err := loadConfig(configFile, config)
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}
	}
	return config
}

func saveConfig(file string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0644)
}

func loadConfig(file string, config *Config) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, config)
}

func (c *Config) SetDirectory(newDir string) error {
	c.MusicDirectory = newDir
	configFile := filepath.Join(os.Getenv("HOME"), ".config", "MusicDB", "config.json")
	return saveConfig(configFile, c)
}