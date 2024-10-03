package model

type Song struct {
	ID int64
	PerformerID int64
	AlbumID    int64
	Path       string
	Title      string
	Track      int
	Year       int
	Genre      string
}