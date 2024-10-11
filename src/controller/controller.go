package controller

import (
	"log"
	"fmt"
	"strconv"
	"github.com/KevinJGard/MusicDB/src/model"
)

// Controller manages the interaction between the model and the view.
type Controller struct {
	DB *model.DataBase
	Miner *model.Miner
	Config *model.Config
}

// NewController creates and returns a new Controller instance.
func NewController() *Controller {
	config := model.NewConfig()
	db := model.NewDataBase()
	miner := model.NewMiner()
	return &Controller{DB: db, Miner: miner, Config: config}
}

// SetMusicDirectory updates the music directory in the configuration.
func (c *Controller) SetMusicDirectory(newDir string) error {
	return c.Config.SetDirectory(newDir)
}

// MineMetadata finds MP3 files in the directory, extracts metadata from an MP3 file and 
// inserts it into the database.
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

// GetSongs retrieves all songs from the database.
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

// EditSong updates the details of a song.
func (c *Controller) EditSong(idRola int64, newTitle, newGenre string, newTrack, newYear int) error {
	err := c.DB.UpdateSong(idRola, newTitle, newGenre, newTrack, newYear)
	return err
}

// EditAlbum updates the details of an album.
func (c *Controller) EditAlbum(idAlbum int64, newName string, newYear int) error {
	err := c.DB.UpdateAlbum(idAlbum, newName, newYear)
	return err
}

// DefPerson defines a performer as a person and inserts their details into the database.
func (c *Controller) DefPerson(idPerf int64, stageName, realName, birthDate, deathDate string) error {
	err := c.DB.UpdatePerformer(idPerf, 0, stageName)
	if err != nil {
		return err
	}

	_, err = c.DB.InsertPersonIfNotExists(stageName, realName, birthDate, deathDate)
	return err
}

// DefGroup defines a performer as a group and inserts their details into the database.
func (c *Controller) DefGroup(idPerf int64, name, startDate, endDate string) error {
	err := c.DB.UpdatePerformer(idPerf, 1, name)
	if err != nil {
		return err
	}

	_, err = c.DB.InsertGroupIfNotExists(name, startDate, endDate)
	return err
}

// EditPerf updates the name of a performer
func (c *Controller) EditPerf(idPerf int64, newName string) error {
	err := c.DB.UpdateNamePerformer(idPerf, newName)
	return err
}

// AddPersonToGroup adds a person to a specified group in the database.
func (c *Controller) AddPersonToGroup(stageName, realName, birthDate, deathDate, nameGroup string) error {
	personID, err := c.DB.GetPersonID(stageName, realName, birthDate, deathDate)
	if err != nil {
		return err
	}
	groupID, err := c.DB.GetGroupIDByName(nameGroup)
	if err != nil {
		return err
	}

	if groupID == 0 {
		return fmt.Errorf("the group '%s' is not found in the database.", nameGroup)
	}
	
	return c.DB.InsertPersonInGroup(personID, groupID)
}

// GetSearchSongs searches for songs according to the request.
func (c *Controller) GetSearchSongs(search string) ([]model.Song, error) {
	results := splitString(search)

	var allSongs []model.Song
	if len(results["titles"]) > 0 {
		for _, title := range results["titles"] {
			songsByTitle, err := c.DB.SearchByTitle(title)
			if err != nil {
				return nil, err
			}
			allSongs = append(allSongs, songsByTitle...)
		}
	}
	if len(results["artists"]) > 0 {
		for _, artist := range results["artists"] {
			songsByPerformer, err := c.DB.SearchByPerformer(artist)
			if err != nil {
				return nil, err
			}
			allSongs = append(allSongs, songsByPerformer...)
		}
	}
	if len(results["albums"]) > 0 {
		for _, album := range results["albums"] {
			songsByAlbum, err := c.DB.SearchByAlbum(album)
			if err != nil {
				return nil, err
			}
			allSongs = append(allSongs, songsByAlbum...)
		}
	}
	if len(results["years"]) > 0 {
		for _, year := range results["years"] {
			yearNum, err := strconv.Atoi(year)
			if err != nil {
				return nil, err
			}
			songsByYear, err := c.DB.SearchByYear(yearNum)
			if err != nil {
				return nil, err
			}
			allSongs = append(allSongs, songsByYear...)
		}
	}
	if len(results["genres"]) > 0 {
		for _, genre := range results["genres"] {
			songsByGenre, err := c.DB.SearchByGenre(genre)
			if err != nil {
				return nil, err
			}
			allSongs = append(allSongs, songsByGenre...)
		}
	}
	return allSongs, nil
}