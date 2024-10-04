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
    
    dbFilePath := filepath.Join(tempDir, ".config", "MusicDB", "music.sqlite")
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
