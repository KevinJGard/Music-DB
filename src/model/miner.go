package model

import (
	"os"
	"path/filepath"
	"time"
	"github.com/dhowden/tag"
)

// Miner is responsible for mining metadata from MP3 files.
type Miner struct{}

// NewMiner creates and returns a new Miner instance.
func NewMiner() *Miner {
	return &Miner{}
}

// FindMP3Files traverses the specified directory and returns a list of MP3 files.
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

// MineMetadata extracts metadata from a given MP3 file.
func (miner *Miner) MineMetadata(file string) (map[string]interface{}, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	metadata, err := tag.ReadFrom(f)
	if err != nil {
		return nil, err
	}
	return miner.AssignTag(metadata), nil
}

// AssignTag replaces empty extracted metadata by assigning default values such as 'Unknown',
// '1' and the current year.
func (miner *Miner) AssignTag(metadata tag.Metadata) map[string]interface{} {
	disc, totalDiscs := metadata.Disc()
	trackNumber, totalTracks := metadata.Track()
	disc, totalDiscs = checkTrackTag(disc, totalDiscs)
	trackNumber, totalTracks = checkTrackTag(trackNumber, totalTracks)
	return map[string]interface{} {
		"Title":        checkStringTag(metadata.Title()),
		"Artist":       checkStringTag(metadata.Artist()),
		"Album":        checkStringTag(metadata.Album()),
		"AlbumArtist":  checkStringTag(metadata.AlbumArtist()),
		"Genre":        checkStringTag(metadata.Genre()),
		"Year":         checkYearTag(metadata.Year()),
		"Disc":         map[string]int{"Number": disc, "Total": totalDiscs},
		"Comment":      checkStringTag(metadata.Comment()),
		"Track":        map[string]int{"Number": trackNumber, "Total": totalTracks},
		"Composer":     checkStringTag(metadata.Composer()),
	}
}

// checkStringTag returns "Unknown" if the tag is empty.
func checkStringTag(tag string) string {
	if tag == "" {
		return "Unknown"
	}
	return tag
}

// checkYearTag returns the current year if the year is 0.
func checkYearTag(year int) int {
	if year == 0 {
		return time.Now().Year()
	}
	return year
}

// checkTrackTag returns default values if any track tag is 0.
func checkTrackTag(trackNumber int, totalTracks int) (int, int) {
	if trackNumber == 0 && totalTracks == 0{
		return 1, 1
	} else if trackNumber == 0 {
		return 1, totalTracks
	} else if totalTracks == 0 {
		return trackNumber, 1
	}
	return trackNumber, totalTracks
}

// ProcessFile mines metadata from an MP3 file and inserts it into the database.
func (miner *Miner) ProcessFile(db *DataBase, file string) error {
	metadata, err := miner.MineMetadata(file)
	if err != nil {
		return err
	}

	var performerID int64
	if metadata["Artist"] == "Unknown" {
		performerID, err = db.InsertPerformerIfNotExists(metadata["Artist"].(string), 2)
	} else {
		performerID, err = db.InsertPerformerIfNotExists(metadata["Artist"].(string), 0)
	}
	if err != nil {
		return err
	}

	albumPath := filepath.Dir(file)
	albumID, err := db.InsertAlbumIfNotExists(metadata["Album"].(string), metadata["Year"].(int), albumPath)
	if err != nil {
		return err
	}

	song := Song{
		PerformerID: performerID,
		AlbumID: albumID,
		Path: file,
		Title: metadata["Title"].(string),
		Track: metadata["Track"].(map[string]int)["Number"],
		Year: metadata["Year"].(int),
		Genre: metadata["Genre"].(string),
	}
	_, err = db.InsertSongIfNotExists(song.PerformerID, song.AlbumID, song.Path, song.Title, song.Genre, song.Track, song.Year)

	return err
}