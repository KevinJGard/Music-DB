package test

import (
	"testing"
	"github.com/KevinJGard/MusicDB/src/model"
)

func TestFindMP3Files(t *testing.T) {
	miner := model.NewMiner()

	directory := "/home/kevingardhp/Música/"
	_, err := miner.FindMP3Files(directory)
	if err != nil {
		t.Fatalf("Error traversing directory: %v", err)
	}
}

func TestMineMetadata(t *testing.T) {
	miner := model.NewMiner()
	file := "/home/kevingardhp/Música/full.mp3"
	metadata, err := miner.MineMetadata(file)
	
	if err != nil {
		t.Fatalf("Error reading metadata for %s: %v", file, err)
	}

	if metadata["Title"] == "" {
		t.Error("Tag not found.")
	}
	if metadata["Artist"] == "" {
		t.Error("Tag not found.")
	}
	if metadata["Album"] == "" {
		t.Error("Tag not found.")
	}
	if metadata["Genre"] == "" {
		t.Error("Tag not found.")
	}
	if metadata["Year"] == 0 {
		t.Error("Tag not found.")
	}
	track:= metadata["Track"].(map[string]int)
	if track["Number"] == 0 && track["Total"] == 0 {
		t.Error("Tag not found.")
	}
}