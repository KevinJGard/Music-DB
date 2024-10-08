package test

import (
	"testing"
	"os"
	"path/filepath"
	_ "github.com/mattn/go-sqlite3"
	"github.com/KevinJGard/MusicDB/src/model"
	"github.com/stretchr/testify/assert"
)

func TestNewDataBase(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	db := model.NewDataBase()

	assert.NotNil(t, db.Db, "Database should not be nil.")

	var count int
	query := `SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name = 'types'`
	err := db.Db.QueryRow(query).Scan(&count)
	assert.NoError(t, err, "Expected no error querying for 'types' table.")
    assert.Equal(t, 1, count, "Expected 'types' table to exist.")

    query = `SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name = 'performers'`
	err = db.Db.QueryRow(query).Scan(&count)
	assert.NoError(t, err, "Expected no error querying for 'performers' table.")
    assert.Equal(t, 1, count, "Expected 'performers' table to exist.")

    query = `SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name= 'persons'`
	err = db.Db.QueryRow(query).Scan(&count)
	assert.NoError(t, err, "Expected no error querying for 'persons' table.")
    assert.Equal(t, 1, count, "Expected 'persons' table to exist.")

    query = `SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name = 'groups'`
	err = db.Db.QueryRow(query).Scan(&count)
	assert.NoError(t, err, "Expected no error querying for 'groups' table.")
    assert.Equal(t, 1, count, "Expected 'groups' table to exist.")

    query = `SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name = 'albums'`
	err = db.Db.QueryRow(query).Scan(&count)
	assert.NoError(t, err, "Expected no error querying for 'albums' table.")
    assert.Equal(t, 1, count, "Expected 'albums' table to exist.")

    query = `SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name = 'rolas'`
	err = db.Db.QueryRow(query).Scan(&count)
	assert.NoError(t, err, "Expected no error querying for 'rolas' table.")
    assert.Equal(t, 1, count, "Expected 'rolas' table to exist.")

    query = `SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name = 'in_group'`
	err = db.Db.QueryRow(query).Scan(&count)
	assert.NoError(t, err, "Expected no error querying for 'in_group' table.")
    assert.Equal(t, 1, count, "Expected 'in_group' table to exist.")
    
    dbFilePath := filepath.Join(tempDir, ".local", "share", "MusicDB", "music.sqlite")
    _, err = os.Stat(dbFilePath)
    assert.NoError(t, err, "Database file should exist.")
    defer db.Db.Close()
}

func TestInsertSong(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	db := model.NewDataBase()
	song := &model.Song{
		PerformerID: 1,
		AlbumID: 1,
		Path: "/path/test/song1.mp3",
		Title: "Song 1",
		Track: 86,
		Year: 1521,
		Genre: "Rap",
	}

	err := db.InsertSong(song)
	assert.NoError(t, err, "Failed inserting song.")
	var id int
	query := `SELECT id_rola FROM rolas WHERE title = ?`
	err = db.Db.QueryRow(query, song.Title).Scan(&id)
	assert.NoError(t, err, "Expected no error while counting songs.")
	assert.Equal(t, 1, id, "Expected one song to be inserted.")
}

func TestInsertPerformer(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	db := model.NewDataBase()
	performer := &model.Performer{
		Type: 1,
		Name: "Test Group",
	}

	err := db.InsertPerformer(performer)
	assert.NoError(t, err, "Failed inserting performer.")
	id, err := db.GetPerformerID(performer.Name)
	assert.NoError(t, err, "Expected no error while counting performers.")
	assert.Equal(t, int64(1), id, "Expected one performer to be inserted.")
}

func TestInsertAlbum(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	db := model.NewDataBase()
	album := &model.Album{
		Path: "/path/test/song1.mp3",
		Name: "Test Album",
		Year: 1945,
	}

	err := db.InsertAlbum(album)
	assert.NoError(t, err, "Failed inserting album.")
	id, err := db.GetAlbumID(album.Name, album.Year)
	assert.NoError(t, err, "Expected no error while counting albums.")
	assert.Equal(t, int64(1), id, "Expected one album to be inserted.")
}

func TestGetSongID(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	db := model.NewDataBase()
	song := &model.Song{
		PerformerID: 1,
        AlbumID: 1,
        Path: "/path/test/song1.mp3", 
        Title: "song1.mp3", 
        Track: 34, 
        Year: 1901, 
        Genre: "Pop", 
	}

	err := db.InsertSong(song)
	assert.NoError(t, err, "Failed inserting song.")
	id, err := db.GetSongID(song.PerformerID, song.AlbumID, song.Path, song.Title, song.Genre, song.Track, song.Year)
	assert.NoError(t, err, "Expected no error while getting song ID.")
	assert.NotZero(t, id, "Expected song ID to be greater than 0.")

	noID, err := db.GetSongID(5, 3, "test/music/song.mp3", "song.mp3", "Rap", 5, 2014)
	assert.NoError(t, err, "Expected no error while getting no-song ID.")
	assert.Zero(t, noID, "Expected no-song ID to be 0.")
}

func TestGetPerformerID(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	db := model.NewDataBase()
	performer := &model.Performer{
		Type: 1,
		Name: "Test Group",
	}

	err := db.InsertPerformer(performer)
	assert.NoError(t, err, "Failed inserting performer.")
	id, err := db.GetPerformerID(performer.Name)
	assert.NoError(t, err, "Expected no error while getting performer ID.")
	assert.NotZero(t, id, "Expected performer ID to be greater than 0.")

	noID, err := db.GetPerformerID("No Performer")
	assert.NoError(t, err, "Expected no error while getting no-performer ID.")
	assert.Zero(t, noID, "Expected no-performer ID to be 0.")
}

func TestGetAlbumID(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	db := model.NewDataBase()
	album := &model.Album{
		Path: "/path/test/song1.mp3",
		Name: "Test Album",
		Year: 1945,
	}

	err := db.InsertAlbum(album)
	assert.NoError(t, err, "Failed inserting album.")
	id, err := db.GetAlbumID(album.Name, album.Year)
	assert.NoError(t, err, "Expected no error while getting album ID.")
	assert.NotZero(t, id, "Expected album ID to be greater than 0.")

	noID, err := db.GetAlbumID("No album",  1)
	assert.NoError(t, err, "Expected no error while getting no-album ID.")
	assert.Zero(t, noID, "Expected no-album ID to be 0.")
}

func TestInsertSongIfNotExists(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	db := model.NewDataBase()
	song := &model.Song{
		PerformerID: 1,
        AlbumID: 1,
        Path: "/path/test/song1.mp3", 
        Title: "song1.mp3", 
        Track: 34, 
        Year: 1901, 
        Genre: "Pop", 
	}
	songID, err := db.InsertSongIfNotExists(song.PerformerID, song.AlbumID, song.Path, song.Title, song.Genre, song.Track, song.Year)
	assert.NoError(t, err, "Failed inserting new song.")
	assert.NotZero(t, songID, "Expected song ID to be greater than 0.")

	sameID, err := db.InsertSongIfNotExists(song.PerformerID, song.AlbumID, song.Path, song.Title, song.Genre, song.Track, song.Year)
	assert.NoError(t, err, "Failed inserting new song.")
	assert.Equal(t, songID, sameID, "Expected same ID for existing song.")
}

func TestInsertPerformerIfNotExists(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	db := model.NewDataBase()
	performer := &model.Performer{
		Type: 0,
		Name: "Test Performer",
	}
	performerID, err := db.InsertPerformerIfNotExists(performer.Name, performer.Type)
	assert.NoError(t, err, "Failed inserting new performer.")
	assert.NotZero(t, performerID, "Expected performer ID to be greater than 0.")

	sameID, err := db.InsertPerformerIfNotExists(performer.Name, performer.Type)
	assert.NoError(t, err, "Failed inserting performer.")
	assert.Equal(t, performerID, sameID, "Expected same ID for existing performer.")
}

func TestInsertAlbumIfNotExists(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	db := model.NewDataBase()
	album := &model.Album{
		Path: "/path/test/song1.mp3",
		Name: "Test Album",
		Year: 2005,
	}
	albumID, err := db.InsertAlbumIfNotExists(album.Name, album.Year, album.Path)
	assert.NoError(t, err, "Failed inserting new album.")
	assert.NotZero(t, albumID, "Expected album ID to be greater than 0.")

	sameID, err := db.InsertAlbumIfNotExists(album.Name, album.Year, album.Path)
	assert.NoError(t, err, "Failed inserting album.")
	assert.Equal(t, albumID, sameID, "Expected same ID for existing album.")
}

func TestGetPerformerName(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	db := model.NewDataBase()
	performer := &model.Performer{
		Type: 1,
		Name: "Test Group",
	}

	err := db.InsertPerformer(performer)
	assert.NoError(t, err, "Failed inserting performer.")
	name, err := db.GetPerformerName(1)
	assert.NoError(t, err, "Expected no error while getting performer name.")
	assert.Equal(t, performer.Name, name, "Expected performer name to match.")

	noName, err := db.GetPerformerName(999)
	assert.Error(t, err, "Expected no error while getting no-performer name.")
	assert.Empty(t, noName, "Expected no-performer name to be empty.")
}

func TestGetAlbumName(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	db := model.NewDataBase()
	album := &model.Album{
		Path: "/path/test/song1.mp3",
		Name: "Test Album",
		Year: 1945,
	}

	err := db.InsertAlbum(album)
	assert.NoError(t, err, "Failed inserting album.")
	name, err := db.GetAlbumName(1)
	assert.NoError(t, err, "Expected no error while getting album name.")
	assert.Equal(t, album.Name, name, "Expected album name to match.")

	noName, err := db.GetAlbumName(999)
	assert.Error(t, err, "Expected no error while getting no-album name.")
	assert.Empty(t, noName, "Expected no-album name to be empty.")
}

func TestUpdateSong(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	db := model.NewDataBase()
	song := &model.Song{
		PerformerID: 1,
        AlbumID: 1,
        Path: "/path/test/song1.mp3", 
        Title: "song1.mp3", 
        Track: 34, 
        Year: 1901, 
        Genre: "Pop", 
	}

	err := db.InsertSong(song)
	assert.NoError(t, err, "Failed inserting song.")
	err = db.UpdateSong(1, "Test Title", "Jazz", 25, 2005)
	assert.NoError(t, err, "Expected no error while editing song.")

	id, err := db.GetSongID(song.PerformerID, song.AlbumID, song.Path, "Test Title", "Jazz", 25, 2005)
	assert.NoError(t, err, "Expected no error while getting song ID.")
	assert.Equal(t, int64(1), id, "Expected same song id.")
}

func TestUpdateAlbum(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	db := model.NewDataBase()
	album := &model.Album{
		Path: "/path/test/song1.mp3",
		Name: "Test Album",
		Year: 1945,
	}

	err := db.InsertAlbum(album)
	assert.NoError(t, err, "Failed inserting album.")
	err = db.UpdateAlbum(1, "Test name album", 2003)
	assert.NoError(t, err, "Expected no error while editing album.")

	id, err := db.GetAlbumID("Test name album", 2003)
	assert.NoError(t, err, "Expected no error while getting album ID.")
	assert.Equal(t, int64(1), id, "Expected same album id.")
}