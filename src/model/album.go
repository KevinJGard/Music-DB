package model

// Album represents an album associated with the song.
type Album struct {
	ID int64
	Path string
	Name string
	Year int
}