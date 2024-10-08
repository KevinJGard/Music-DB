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

func (c *Controller) GetSongs() ([]model.Song, error) {
	var songs []model.Song
	rows, err := c.DB.Db.Query("SELECT * FROM rolas")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var song model.Song
		if err := rows.Scan(&song.ID, &song.PerformerID, &song.AlbumID, &song.Path, &song.Title, &song.Track, &song.Year, &song.Genre); err != nil {
			return nil, err
		}

		performerName, err := c.DB.GetPerformerName(song.PerformerID)
		if err != nil {
			return nil, err
		}
		albumName, err := c.DB.GetAlbumName(song.AlbumID)
		if err != nil {
			return nil, err
		}

		song.PerformerName = performerName
		song.AlbumName = albumName

		songs = append(songs, song)
	}

	return songs, nil
}

func (c *Controller) EditSong(idRola int64, newTitle, newGenre string, newTrack, newYear int) error {
	err := c.DB.UpdateSong(idRola, newTitle, newGenre, newTrack, newYear)
	return err
}

func (c *Controller) EditAlbum(idAlbum int64, newName string, newYear int) error {
	err := c.DB.UpdateAlbum(idAlbum, newName, newYear)
	return err
}

func (c *Controller) DefPerson(idPerf int64, stageName, realName, birthDate, deathDate string) error {
	err := c.DB.UpdatePerformer(idPerf, 0, stageName)
	if err != nil {
		return err
	}

	_, err = c.DB.InsertPersonIfNotExists(stageName, realName, birthDate, deathDate)
	return err
}

func (c *Controller) DefGroup(idPerf int64, name, startDate, endDate string) error {
	err := c.DB.UpdatePerformer(idPerf, 1, name)
	if err != nil {
		return err
	}

	_, err = c.DB.InsertGroupIfNotExists(name, startDate, endDate)
	return err
}

func (c *Controller) EditPerf(idPerf int64, newName string) error {
	err := c.DB.UpdateNamePerformer(idPerf, newName)
	return err
}