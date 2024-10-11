package model

// Performer represents an artist or group that performs the song.
type Performer struct {
	ID int64
	Type int
	Name string
}