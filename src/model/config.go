package model

import (
	"os"
	"path/filepath"
	"encoding/json"
	"log"
	"strings"
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
		config.MusicDirectory = GetDefaultDir()
		if err := SaveConfig(configFile, config); err != nil {
			log.Fatalf("Error creating config file: %v", err)
		}
	} else {
		err := LoadConfig(configFile, config)
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}
		if config.MusicDirectory == "" {
			config.MusicDirectory = GetDefaultDir()
			if err := SaveConfig(configFile, config); err != nil {
				log.Fatalf("Error saving new information: %v", err)
			}
		}
	}
	return config
}

func SaveConfig(file string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0644)
}

func LoadConfig(file string, config *Config) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, config)
}

func (config *Config) SetDirectory(newDir string) error {
	config.MusicDirectory = newDir
	configFile := filepath.Join(os.Getenv("HOME"), ".config", "MusicDB", "config.json")
	return SaveConfig(configFile, config)
}

func GetDefaultDir() string {
	lang := os.Getenv("LANG")
	if strings.HasPrefix(lang, "es") {
		return filepath.Join(os.Getenv("HOME"), "Música")
	}
	return filepath.Join(os.Getenv("HOME"), "Music")
}