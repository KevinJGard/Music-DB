package test

import (
	"testing"
	"os"
	"path/filepath"
	"encoding/json"
	"github.com/KevinJGard/MusicDB/src/model"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	config := model.NewConfig()

	configDir := filepath.Join(tempDir, ".config", "MusicDB")
	assert.DirExists(t, configDir, "Expected config directory to exist.")

	configFile := filepath.Join(configDir, "config.json")
	assert.FileExists(t, configFile, "Expected config file to exist.")

	assert.Empty(t, config.MusicDirectory, "Expected MusicDirectory to be empty.")
}

func TestSaveConfig(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	configDir := filepath.Join(tempDir, ".config", "MusicDB")

	err := os.MkdirAll(configDir, os.ModePerm)
	assert.NoError(t, err, "Failed to create config directory.")

	config := &model.Config{MusicDirectory: "/test/music/dir"}
	configFile := filepath.Join(configDir, "config.json")
	err = model.SaveConfig(configFile, config)
	assert.NoError(t, err, "Failed to save config.")

	data, err := os.ReadFile(configFile)
	assert.NoError(t, err, "Failed to read config file.")

	var loadedConfig model.Config
	err = json.Unmarshal(data, &loadedConfig)
	assert.NoError(t, err, "Failed to unmarshal config data.")
	assert.Equal(t, config.MusicDirectory, loadedConfig.MusicDirectory, "MusicDirectory does not match.")
}

func TestLoadConfig(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	configDir := filepath.Join(tempDir, ".config", "MusicDB")

	err := os.MkdirAll(configDir, os.ModePerm)
	assert.NoError(t, err, "Failed to create config directory.")

	config := &model.Config{MusicDirectory: "/test/music/dir"}
	configFile := filepath.Join(configDir, "config.json")
	err = model.SaveConfig(configFile, config)
	assert.NoError(t, err, "Failed to save config.")

	newConfig := &model.Config{}
	err = model.LoadConfig(configFile, newConfig)
	assert.NoError(t, err, "Failed to load config.")
	assert.Equal(t, config.MusicDirectory, newConfig.MusicDirectory, "MusicDirectory does not match.")

}

func TestSetDirectory(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	config := model.NewConfig()
	newDir := filepath.Join(tempDir, "TestMusicDir")

	err := config.SetDirectory(newDir)
	assert.NoError(t, err, "Failed to set new directory.")
	assert.Equal(t, newDir, config.MusicDirectory, "MusicDirectory does not match.")
	
	configDir := filepath.Join(tempDir, ".config", "MusicDB")
	err = os.MkdirAll(configDir, os.ModePerm)
	assert.NoError(t, err, "Failed to create config directory.")

	configFile := filepath.Join(configDir, "config.json")
	assert.FileExists(t, configFile, "Expected config file to exist after SetDirectory.")
}