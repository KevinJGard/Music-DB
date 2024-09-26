package model

import (
	"os"
	"path/filepath"
)

type Miner struct{}

func NewMiner() *Miner {
	return &Miner{}
}

func (miner *Miner) FindMP3Files(directory string) ([]string, error) {
	const mp3 = ".mp3"
	var files []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == mp3 {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return files, nil
}