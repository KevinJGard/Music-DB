package test

import (
	"testing"
	"os"
	"path/filepath"
	_ "github.com/mattn/go-sqlite3"
	"github.com/KevinJGard/MusicDB/src/model"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *model.DataBase {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	db := model.NewDataBase()
	return db
}

func assertTableExists(t *testing.T, db *model.DataBase, tableName string) {
	var count int
	query := `SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name = ?`
	err := db.Db.QueryRow(query, tableName).Scan(&count)
	assert.NoError(t, err, "Expected no error querying for '%s' table.", tableName)
	assert.Equal(t, 1, count, "Expected '%s' table to exist.", tableName)
}

func assertSongInserted(t *testing.T, db *model.DataBase, song *model.Song) {
	err := db.InsertSong(song)
	assert.NoError(t, err, "Failed inserting song.")
	id, err := db.GetSongID(song.PerformerID, song.AlbumID, song.Path, song.Title, song.Genre, song.Track, song.Year)
	assert.NoError(t, err, "Expected no error while counting songs.")
	assert.Equal(t, int64(1), id, "Expected one song to be inserted.")
}

func assertPerformerInserted(t *testing.T, db *model.DataBase, performer *model.Performer) {
	err := db.InsertPerformer(performer)
	assert.NoError(t, err, "Failed inserting performer.")
	id, err := db.GetPerformerID(performer.Name)
	assert.NoError(t, err, "Expected no error while counting performers.")
	assert.Equal(t, int64(1), id, "Expected one performer to be inserted.")
}

func assertAlbumInserted(t *testing.T, db *model.DataBase, album *model.Album) {
	err := db.InsertAlbum(album)
	assert.NoError(t, err, "Failed inserting album.")
	id, err := db.GetAlbumID(album.Name, album.Year)
	assert.NoError(t, err, "Expected no error while counting albums.")
	assert.Equal(t, int64(1), id, "Expected one album to be inserted.")
}

func TestNewDataBase(t *testing.T) {
	db := setupTestDB(t)
	assert.NotNil(t, db.Db, "Database should not be nil.")

	tables := []string{"types", "performers", "persons", "groups", "albums", "rolas", "in_group"}
	for _, table := range tables {
		assertTableExists(t, db, table)
	}
    
    dbFilePath := filepath.Join(os.Getenv("HOME"), ".local", "share", "MusicDB", "music.db")
    _, err := os.Stat(dbFilePath)
    assert.NoError(t, err, "Database file should exist.")
    defer db.Db.Close()
}

func TestInsertSong(t *testing.T) {
	db := setupTestDB(t)
	song := &model.Song{
		PerformerID: 1,
		AlbumID: 1,
		Path: "/path/test/song1.mp3",
		Title: "Song 1",
		Track: 86,
		Year: 1521,
		Genre: "Rap",
	}

	assertSongInserted(t, db, song)
}

func TestInsertPerformer(t *testing.T) {
	db := setupTestDB(t)
	performer := &model.Performer{
		Type: 1,
		Name: "Test Group",
	}

	assertPerformerInserted(t, db, performer)
}

func TestInsertAlbum(t *testing.T) {
	db := setupTestDB(t)
	album := &model.Album{
		Path: "/path/test/song1.mp3",
		Name: "Test Album",
		Year: 1945,
	}

	assertAlbumInserted(t, db, album)
}

func TestGetSongID(t *testing.T) {
	db := setupTestDB(t)
	song := &model.Song{
		PerformerID: 1,
        AlbumID: 1,
        Path: "/path/test/song1.mp3", 
        Title: "song1.mp3", 
        Track: 34, 
        Year: 1901, 
        Genre: "Pop", 
	}

	assertSongInserted(t, db, song)

	noID, err := db.GetSongID(5, 3, "test/music/song.mp3", "song.mp3", "Rap", 5, 2014)
	assert.NoError(t, err, "Expected no error while getting no-song ID.")
	assert.Zero(t, noID, "Expected no-song ID to be 0.")
}

func TestGetPerformerID(t *testing.T) {
	db := setupTestDB(t)
	performer := &model.Performer{
		Type: 1,
		Name: "Test Group",
	}

	assertPerformerInserted(t, db, performer)

	noID, err := db.GetPerformerID("No Performer")
	assert.NoError(t, err, "Expected no error while getting no-performer ID.")
	assert.Zero(t, noID, "Expected no-performer ID to be 0.")
}

func TestGetAlbumID(t *testing.T) {
	db := setupTestDB(t)
	album := &model.Album{
		Path: "/path/test/song1.mp3",
		Name: "Test Album",
		Year: 1945,
	}

	assertAlbumInserted(t, db, album)

	noID, err := db.GetAlbumID("No album",  1)
	assert.NoError(t, err, "Expected no error while getting no-album ID.")
	assert.Zero(t, noID, "Expected no-album ID to be 0.")
}

func TestInsertSongIfNotExists(t *testing.T) {
	db := setupTestDB(t)
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
	db := setupTestDB(t)
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
	db := setupTestDB(t)
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
	db := setupTestDB(t)
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
	db := setupTestDB(t)
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
	db := setupTestDB(t)
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
	db := setupTestDB(t)
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

	sameName, err := db.GetAlbumName(1)
	assert.NoError(t, err, "Expected no error while getting album name")
	assert.Equal(t, "Test name album", sameName, "Expected album name to match.")
}

func TestUpdatePerformer(t *testing.T) {
	db := setupTestDB(t)
	performer := &model.Performer{
		Type: 1,
		Name: "Existing performer",
	}

	err := db.InsertPerformer(performer)
	assert.NoError(t, err, "Failed inserting performer.")
	err = db.UpdatePerformer(1, 2, "New name")
	assert.NoError(t, err, "Expected no error while editing performer.")

	id, err := db.GetPerformerID("New name")
	assert.NoError(t, err, "Expected no error while getting performer ID.")
	assert.Equal(t, int64(1), id, "Expected same performer id.")

	sameName, err := db.GetPerformerName(1)
	assert.NoError(t, err, "Expected no error while getting performer name")
	assert.Equal(t, "New name", sameName, "Expected performer name to match.")
}

func TestUpdateNamePerformer(t *testing.T) {
	db := setupTestDB(t)
	performer := &model.Performer{
		Type: 1,
		Name: "Existing performer",
	}

	err := db.InsertPerformer(performer)
	assert.NoError(t, err, "Failed inserting performer.")
	err = db.UpdateNamePerformer(1, "New name")
	assert.NoError(t, err, "Expected no error while editing performer.")

	id, err := db.GetPerformerID("New name")
	assert.NoError(t, err, "Expected no error while getting performer ID.")
	assert.Equal(t, int64(1), id, "Expected same performer id.")

	sameName, err := db.GetPerformerName(1)
	assert.NoError(t, err, "Expected no error while getting performer name")
	assert.Equal(t, "New name", sameName, "Expected performer name to match.")
}

func assertDefinePerson(t *testing.T, db *model.DataBase, stageName, realName, birthDate, deathDate string) {
	err := db.DefinePerson(stageName, realName, birthDate, deathDate)
	assert.NoError(t, err, "Failed inserting person.")
	id, err := db.GetPersonID(stageName, realName, birthDate, deathDate)
	assert.NoError(t, err, "Expected no error while counting persons.")
	assert.Equal(t, int64(1), id, "Expected one person to be inserted.")
}

func TestDefinePerson(t *testing.T) {
	db := setupTestDB(t)
	const stageName = "Stage Name"
	const realName = "Real Name"
	const birthDate = "1945"
	const deathDate = "2011"

	assertDefinePerson(t, db, stageName, realName, birthDate, deathDate)
}

func TestGetPersonID(t *testing.T) {
	db := setupTestDB(t)
	const stageName = "Stage Name"
	const realName = "Real Name"
	const birthDate = "1945"
	const deathDate = "2011"

	assertDefinePerson(t, db, stageName, realName, birthDate, deathDate)


	noID, err := db.GetPersonID("No Person", "No Person Name", "2000", "0")
	assert.NoError(t, err, "Expected no error while getting no-person ID.")
	assert.Zero(t, noID, "Expected no-person ID to be 0.")
}

func TestInsertPersonIfNotExists(t *testing.T) {
	db := setupTestDB(t)
	const stageName = "Stage Name"
	const realName = "Real Name"
	const birthDate = "1945"
	const deathDate = "2011"

	personID, err := db.InsertPersonIfNotExists(stageName, realName, birthDate, deathDate)
	assert.NoError(t, err, "Failed inserting new person.")
	assert.NotZero(t, personID, "Expected person ID to be greater than 0.")

	sameID, err := db.InsertPersonIfNotExists(stageName, realName, birthDate, deathDate)
	assert.NoError(t, err, "Failed inserting person.")
	assert.Equal(t, personID, sameID, "Expected same ID for existing person.")
}

func assertDefineGroup(t *testing.T, db *model.DataBase, name, startDate, endDate string) {
	err := db.DefineGroup(name, startDate, endDate)
	assert.NoError(t, err, "Failed inserting group.")
	id, err := db.GetGroupID(name, startDate, endDate)
	assert.NoError(t, err, "Expected no error while counting groups.")
	assert.Equal(t, int64(1), id, "Expected one group to be inserted.")
}

func TestDefineGroup(t *testing.T) {
	db := setupTestDB(t)
	const name = "Name Group"
	const startDate = "2006"
	const endDate = "2019"

	assertDefineGroup(t, db, name, startDate, endDate)
}

func TestGetGroupID(t *testing.T) {
	db := setupTestDB(t)
	const name = "Name Group"
	const startDate = "2006"
	const endDate = "2019"

	assertDefineGroup(t, db, name, startDate, endDate)

	noID, err := db.GetGroupID("No Group", "2000", "0")
	assert.NoError(t, err, "Expected no error while getting no-group ID.")
	assert.Zero(t, noID, "Expected no-group ID to be 0.")
}

func TestInsertGroupIfNotExists(t *testing.T) {
	db := setupTestDB(t)
	const name = "Name Group"
	const startDate = "2006"
	const endDate = "2019"

	groupID, err := db.InsertGroupIfNotExists(name, startDate, endDate)
	assert.NoError(t, err, "Failed inserting new group.")
	assert.NotZero(t, groupID, "Expected group ID to be greater than 0.")

	sameID, err := db.InsertGroupIfNotExists(name, startDate, endDate)
	assert.NoError(t, err, "Failed inserting group.")
	assert.Equal(t, groupID, sameID, "Expected same ID for existing group.")
}

func TestGetGroupIDByName(t *testing.T) {
	db := setupTestDB(t)
	const name = "Name Group"
	const startDate = "2006"
	const endDate = "2019"

	err := db.DefineGroup(name, startDate, endDate)
	assert.NoError(t, err, "Failed inserting group.")
	id, err := db.GetGroupIDByName(name)
	assert.NoError(t, err, "Expected no error while getting group ID.")
	assert.NotZero(t, id, "Expected group ID to be greater than 0.")

	noID, err := db.GetGroupIDByName("No Group")
	assert.NoError(t, err, "Expected no error while getting no-group ID.")
	assert.Zero(t, noID, "Expected no-group ID to be 0.")
}

func TestInsertPersonInGroup(t *testing.T) {
	db := setupTestDB(t)

	err := db.InsertPersonInGroup(1, 1)
	assert.NoError(t, err, "Failed inserting in_group.")
}

func assertInsert(t *testing.T, db *model.DataBase) {
	performer := &model.Performer{Name: "Test Performer", Type: 1}
    performerID, err := db.InsertPerformerIfNotExists(performer.Name, performer.Type)
    assert.NoError(t, err, "Failed inserting performer.")
    album := &model.Album{Name: "Test Album", Year: 1901, Path: "/path/test"}
    albumID, err := db.InsertAlbumIfNotExists(album.Name, album.Year, album.Path)
    assert.NoError(t, err, "Failed inserting album.")
	song := &model.Song{
		PerformerID: performerID,
        AlbumID: albumID,
        Path: "/path/test/song1.mp3", 
        Title: "song1", 
        Track: 34, 
        Year: 1901, 
        Genre: "Pop", 
	}
	err = db.InsertSong(song)
	assert.NoError(t, err, "Failed inserting song.")
}

func TestSearchByTitle(t *testing.T) {
	db := setupTestDB(t)
	assertInsert(t, db)
	
	songs, err := db.SearchByTitle("song1")
	assert.NoError(t, err, "Expected no error searching by title.")
	assert.Len(t, songs, 1, "Expected one song returned.")
	assert.Equal(t, "song1", songs[0].Title, "Expected song title to match.")
}

func TestSearchByPerformer(t *testing.T) {
	db := setupTestDB(t)
	assertInsert(t, db)
	
	songs, err := db.SearchByPerformer("Test Performer")
	assert.NoError(t, err, "Expected no error searching by performer.")
	assert.Len(t, songs, 1, "Expected one song returned.")
	assert.Equal(t, "song1", songs[0].Title, "Expected song title to match.")
}

func TestSearchSearchByAlbum(t *testing.T) {
	db := setupTestDB(t)
	assertInsert(t, db)
	
	songs, err := db.SearchByAlbum("Test Album")
	assert.NoError(t, err, "Expected no error searching by album.")
	assert.Len(t, songs, 1, "Expected one song returned.")
	assert.Equal(t, "song1", songs[0].Title, "Expected song title to match.")
}

func TestSearchSearchByYear(t *testing.T) {
	db := setupTestDB(t)
	assertInsert(t, db)
	
	songs, err := db.SearchByYear(1901)
	assert.NoError(t, err, "Expected no error searching by year.")
	assert.Len(t, songs, 1, "Expected one song returned.")
	assert.Equal(t, "song1", songs[0].Title, "Expected song title to match.")
}

func TestSearchByGenre(t *testing.T) {
	db := setupTestDB(t)
	assertInsert(t, db)
	
	songs, err := db.SearchByGenre("Pop")
	assert.NoError(t, err, "Expected no error searching by year.")
	assert.Len(t, songs, 1, "Expected one song returned.")
	assert.Equal(t, "song1", songs[0].Title, "Expected song title to match.")
}