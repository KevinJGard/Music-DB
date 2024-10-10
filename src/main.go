package main

import (
	"github.com/KevinJGard/MusicDB/src/model"
	"github.com/KevinJGard/MusicDB/src/controller"
	"fmt"
	"os"
	"log"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <directory> <search>", os.Args[0])
	}
	config := model.NewConfig()
	err := config.SetDirectory(os.Args[1])
	if err != nil {
		log.Fatalf("Error setting directory: %v", err)
	}
	directory := config.MusicDirectory
	miner := model.NewMiner()
	database := model.NewDataBase()
	controller := controller.NewController()
	songs, err := controller.GetSearchSongs(os.Args[2])
	if err != nil {
		log.Fatalf("Error searching songs: %v", err)
	}
	fmt.Println("Search results:")
	for _, song := range songs {
		fmt.Printf("Title: %s, Artist: %s, Album: %s, Year: %d\n", song.Title, song.PerformerName, song.AlbumName, song.Year)
	}

	files, err := miner.FindMP3Files(directory)
	if err != nil {
		log.Fatalf("Error traversing directory %s: %v", directory, err)
	}
	fmt.Println("MP3 files found:")
	for i, file := range files {
		metadata, err := miner.MineMetadata(file)
		if err != nil {
			log.Printf("Error reading metadata for %s: %v", file, err)
			continue
		}
		fmt.Printf("File: %s \n", file)
		fmt.Printf("Title: %s \n", metadata["Title"])
		fmt.Printf("Artist: %s \n", metadata["Artist"])
		fmt.Printf("Album: %s \n", metadata["Album"])
		fmt.Printf("AlbumArtist: %s \n", metadata["AlbumArtist"])
		fmt.Printf("Genre: %s \n", metadata["Genre"])
		fmt.Printf("Year: %d \n", metadata["Year"])
		disc := metadata["Disc"].(map[string]int)
		fmt.Printf("Disc Number: %d \n", disc["Number"])
		fmt.Printf("Comment: %s \n", metadata["Comment"])
		track := metadata["Track"].(map[string]int)
		fmt.Printf("Track: %d of %d \n", track["Number"], track["Total"])
		fmt.Printf("Composer: %s \n", metadata["Composer"])
		fmt.Printf("%d----------------------------------------------------------------------\n", i)

		if err := miner.ProcessFile(database, file); err != nil {
			log.Printf("Error procesing file %s: %v", file, err)
			continue
		}
	}

	fmt.Println("Everything was done correctly.")
}