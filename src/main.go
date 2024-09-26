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
		fmt.Println(file)
	}
}