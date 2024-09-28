package test

import (
	"testing"
	"io"
	"os"
	"path/filepath"
	"github.com/KevinJGard/MusicDB/src/model"
	"github.com/dhowden/tag"
)

func copyMp3(tempDir string, fileName string) error {
	file, err := os.Open(filepath.Join("..", "..", "testdata", "testdata_with_tags_sample.id3v24.mp3"))
    if err != nil {
        return err
    }
    defer file.Close()

    destFile, err := os.Create(filepath.Join(tempDir, fileName))
    if err != nil {
        return err
    }
    defer destFile.Close()

    _, err = io.Copy(destFile, file)
    return err
}

func createTempDirectorywithFiles(files []string) (string, error) {
	tempDir, err := os.MkdirTemp("", "MusicDB_Test")
    if err != nil {
        return "", err
    }

    for _, file := range files {
        if err := copyMp3(tempDir, file); err != nil {
	    	return "", err
	    }
    }

    return tempDir, nil
}

func TestFindMP3Files(t *testing.T) {
	files := []string{"test1.mp3", "test2.mp3", "test3.mp3"}
	tempDir, err := createTempDirectorywithFiles(files)
	if err != nil {
        t.Fatalf("failed to create temp dir with files: %v", err)
    }
    defer os.RemoveAll(tempDir)
	miner := model.NewMiner()
	foundFiles, err := miner.FindMP3Files(tempDir)
	if err != nil {
		t.Fatalf("Error traversing directory %s: %v", tempDir, err)
	}

	if len(foundFiles) != len(files) {
		t.Fatalf("Expected %d MP3 files, but found %d.", len(files), len(foundFiles))
	}

	for _, file := range files {
		filePath := filepath.Join(tempDir, file)
		found := false
		for _, foundFile := range foundFiles {
			if foundFile == filePath {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected file %s not found in results.", filePath)
		}
	}
}

func TestMineMetadata(t *testing.T) {
	files := []string{"test1.mp3", "test2.mp3", "test3.mp3"}
	tempDir, err := createTempDirectorywithFiles(files)
	if err != nil {
        t.Fatalf("failed to create temp dir with files: %v", err)
    }
    defer os.RemoveAll(tempDir)
	miner := model.NewMiner()

	for _, file := range files {
		filePath := filepath.Join(tempDir, file)
		metadata, err := miner.MineMetadata(filePath)
		
		if err != nil {
			t.Fatalf("Error reading metadata for %s: %v", file, err)
		}

		if metadata["Title"] == "" {
			t.Errorf("Tag \"Title\" not found in %s.", filePath)
		}
		if metadata["Artist"] == "" {
			t.Errorf("Tag \"Artist\" not found in %s.", filePath)
		}
		if metadata["Album"] == "" {
			t.Errorf("Tag \"Album\" not found in %s.", filePath)
		}
		if metadata["Genre"] == "" {
			t.Errorf("Tag \"Genre\" not found in %s.", filePath)
		}
		if metadata["Year"] == 0 {
			t.Errorf("Tag \"Year\" not found in %s.", filePath)
		}
		track := metadata["Track"].(map[string]int)
		if track["Number"] == 0 && track["Total"] == 0 {
			t.Errorf("Tag \"Track\" not found in %s.", filePath)
		}
	}
}

func TestAssignTag(t *testing.T) {
	file, err := os.Open(filepath.Join("..", "..", "testdata", "testdata_without_tags_sample.mp3"))
    if err != nil {
        t.Fatalf("Error opening test file: %v", err)
    }
    defer file.Close()

    metadata, err := tag.ReadFrom(file)
	if err != nil {
		t.Fatalf("Error reading metadata: %v", err)
	}

	miner := model.NewMiner()

	tags := miner.AssignTag(metadata)
	if tags["Title"] == "" {
		t.Error("Expected a valid Title or \"Unknown\"")
	}
	if tags["Artist"] == "" {
		t.Error("Expected a valid Artist or \"Unknown\"")
	}
	if tags["Album"] == "" {
		t.Error("Expected a valid Album or \"Unknown\"")
	}
	if tags["Genre"] == "" {
		t.Error("Expected a valid Genre or \"Unknown\"")
	}
	if tags["Year"] == 0 {
		t.Error("Expected a valid Year or \"1\"")
	}
	track := tags["Track"].(map[string]int)
	if track["Number"] == 0 && track["Total"] == 0 {
		t.Error("Expected valid Track number and total or \"1 and 1\"")
	}
}