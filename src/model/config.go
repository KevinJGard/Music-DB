package model

import (
	"os"
	"path/filepath"
	"encoding/json"
	"log"
	"strings"
)

// Config holds the music directory path where the MP3 files are located.
type Config struct {
	MusicDirectory string `json:"music_directory"`
}

// NewConfig creates a new Config instance.
// It checks for the existence of the configuration file, and if it doesn't exist,
// it creates a default configuration. If the file exists, it loads the configuration.
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

// SaveConfig saves the configuration to a JSON file.
// Converts the Config structure to a JSON string and writes it to the specified file.
func SaveConfig(file string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0644)
}

// LoadConfig loads the configuration from the JSON file.
// It reads the file and unmarshals the JSON data into the Config struct.
func LoadConfig(file string, config *Config) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, config)
}

// SetDirectory updates the music directory in the configuration file.
func (config *Config) SetDirectory(newDir string) error {
	config.MusicDirectory = newDir
	configFile := filepath.Join(os.Getenv("HOME"), ".config", "MusicDB", "config.json")
	return SaveConfig(configFile, config)
}

// GetDefaultDir returns the default music directory based on the user's language setting.
func GetDefaultDir() string {
	lang := os.Getenv("LANG")
	if strings.HasPrefix(lang, "es") {
		return filepath.Join(os.Getenv("HOME"), "Música")
	}
	return filepath.Join(os.Getenv("HOME"), "Music")
}