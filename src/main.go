package main

import (
	"github.com/KevinJGard/MusicDB/src/model"
	"fmt"
	"os"
	"log"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <directory>", os.Args[0])
	}
	miner := model.NewMiner()

	directory := os.Args[1]
	files, err := miner.FindMP3Files(directory)
	if err != nil {
		log.Fatalf("Error traversing directory: %v", err)
	}
	fmt.Println("MP3 files found:")
	for _, file := range files {
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
		track:= metadata["Track"].(map[string]int)
		fmt.Printf("Track: %d of %d \n", track["Number"], track["Total"])
		fmt.Printf("Composer: %s \n", metadata["Composer"])
	}
}