package main

import (
	"github.com/KevinJGard/MusicDB/src/model"
	"fmt"
	"log"
)

func main() {
	miner := model.NewMiner()

	directory := "/home/kevingardhp/Música"
	files, err := miner.FindMP3Files(directory)
	if err != nil {
		log.Fatalf("Error traversing directory: %v", err)
	}
	fmt.Println("MP3 files found:")
	for _, file := range files {
		fmt.Println(file)
	}
}