package controller

import (
	"log"
	"github.com/KevinJGard/MusicDB/src/model"
)

type Controller struct {
	DB *model.DataBase
	Miner *model.Miner
	Config *model.Config
}

func NewController() *Controller {
	config := model.NewConfig()
	db := model.NewDataBase()
	miner := model.NewMiner()
	return &Controller{DB: db, Miner: miner, Config: config}
}

func (c *Controller) SetMusicDirectory(newDir string) error {
	return c.Config.SetDirectory(newDir)
}


func (c *Controller) MineMetadata(updateProgress func(int), complete func()) error {
	directory := c.Config.MusicDirectory
	files, err := c.Miner.FindMP3Files(directory)
	if err != nil {
		return err
	}

	totalFiles := len(files)
	for i, file := range files {
		if err := c.Miner.ProcessFile(c.DB, file); err != nil {
			log.Printf("Error procesing file %s: %v", file, err)
			continue
		}
		updateProgress((i + 1) * 100 / totalFiles)
	}
	complete()
	return nil
}