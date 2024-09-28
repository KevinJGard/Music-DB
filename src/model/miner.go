package model

import (
	"os"
	"path/filepath"
	"github.com/dhowden/tag"
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
	return miner.assignTag(metadata), nil
}

func (miner *Miner) assignTag(metadata tag.Metadata) map[string]interface{} {
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

func checkStringTag(tag string) string {
	if tag == "" {
		return "Unknown"
	}
	return tag
}

func checkYearTag(year int) int {
	if year == 0 {
		return 1
	}
	return year
}

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