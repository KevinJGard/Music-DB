package model

import (
	"log"
	"os"
	"path/filepath"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type DataBase struct {
	Db *sql.DB
}

func NewDataBase() *DataBase {
	dbDir := filepath.Join(os.Getenv("HOME"), ".local", "share", "MusicDB")
	os.MkdirAll(dbDir, os.ModePerm)
	dbFile := filepath.Join(dbDir, "music.sqlite")

	database, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if err := createTables(database); err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}

	return &DataBase{Db: database}
}

func createTables(db *sql.DB) error {
	tableCreationQueries := []string {
		`CREATE TABLE IF NOT EXISTS types (
			id_type INTEGER PRIMARY KEY,
			description TEXT
		);`,
		`INSERT OR IGNORE INTO types (id_type, description) VALUES (0, 'Person');`,
		`INSERT OR IGNORE INTO types (id_type, description) VALUES (1, 'Group');`,
		`INSERT OR IGNORE INTO types (id_type, description) VALUES (2, 'Unknown');`,
		`CREATE TABLE IF NOT EXISTS performers (
			id_performer INTEGER PRIMARY KEY,
			id_type INTEGER,
			name TEXT,
			FOREIGN KEY (id_type) REFERENCES types(id_type)
		);`,
		`CREATE TABLE IF NOT EXISTS persons (
			id_person INTEGER PRIMARY KEY,
			stage_name TEXT,
			real_name TEXT,
			birth_date TEXT,
			death_date TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS groups (
			id_group INTEGER PRIMARY KEY,
			name TEXT,
			start_date TEXT,
			end_date TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS in_group (
			id_person INTEGER,
			id_group INTEGER,
			PRIMARY KEY (id_person, id_group),
			FOREIGN KEY (id_person) REFERENCES persons(id_person),
			FOREIGN KEY (id_group) REFERENCES groups(id_group)
		);`,
		`CREATE TABLE IF NOT EXISTS albums (
			id_album INTEGER PRIMARY KEY,
			path TEXT,
			name TEXT,
			year INTEGER
		);`,
		`CREATE TABLE IF NOT EXISTS rolas (
			id_rola INTEGER PRIMARY KEY,
			id_performer INTEGER,
			id_album INTEGER,
			path TEXT,
			title TEXT,
			track INTEGER,
			year INTEGER,
			genre TEXT,
			FOREIGN KEY (id_performer) REFERENCES performers(id_performer),
			FOREIGN KEY (id_album) REFERENCES albums(id_album)
		);`,
	}

	for _, query := range tableCreationQueries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

func (db *DataBase) InsertSong(song *Song) error {
	query := `INSERT INTO rolas (id_performer, id_album, path, title, track, year, genre) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Db.Exec(query, song.PerformerID, song.AlbumID, song.Path, song.Title, song.Track, song.Year, song.Genre)
	return err
}

func (db *DataBase) InsertPerformer(performer *Performer) error {
	query := `INSERT INTO performers (id_type, name) 
              VALUES (?, ?)`
	_, err := db.Db.Exec(query, performer.Type, performer.Name)
	return err
}

func (db *DataBase) InsertAlbum(album *Album) error {
	query := `INSERT INTO albums (path, name, year) 
              VALUES (?, ?, ?)`
	_, err := db.Db.Exec(query, album.Path, album.Name, album.Year)
	return err
}

func (db *DataBase) GetSongID(performer, album int64, path, title, genre  string, track , year int) (int64, error) {
	var id int64
	query := `SELECT id_rola FROM rolas WHERE id_performer = ? AND id_album = ? AND path = ? AND title = ? AND track = ? AND year = ? AND genre = ?`
    err := db.Db.QueryRow(query, performer, album, path, title, track, year, genre).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return id, err
}

func (db *DataBase) GetPerformerID(name string) (int64, error) {
	var id int64
	query := `SELECT id_performer FROM performers WHERE name = ?`
	err := db.Db.QueryRow(query, name).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return id, err
}

func (db *DataBase) GetAlbumID(album string, year int) (int64, error) {
	var id int64
	query := `SELECT id_album FROM albums WHERE name = ? AND year = ?`
	err := db.Db.QueryRow(query, album, year).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return id, err
}


func (db *DataBase) InsertSongIfNotExists(performer, album int64, path, title, genre  string, track , year int) (int64, error) {
    id, err := db.GetSongID(performer, album, path, title, genre, track, year)
    if err == nil && id != 0 {
        return id, nil
    }

    song := Song{
        PerformerID: performer,
        AlbumID: album,
        Path: path, 
        Title: title, 
        Track: track, 
        Year: year, 
        Genre: genre, 
    }
    err = db.InsertSong(&song)
    if err != nil {
        return 0, err
    }
    song.ID, err = db.GetSongID(song.PerformerID, song.AlbumID, song.Path, song.Title, song.Genre, song.Track, song.Year)
    return song.ID, nil
}

func (db *DataBase) InsertPerformerIfNotExists(name string, performerType int) (int64, error) {
	id, err := db.GetPerformerID(name)
	if err == nil && id != 0 {
		return id, nil
	}

	performer := Performer{
		Type: performerType,
		Name: name,
	}
	err = db.InsertPerformer(&performer)
	if err != nil {
		return 0, err
	}
	performer.ID, err = db.GetPerformerID(performer.Name)
	return performer.ID, nil
}

func (db *DataBase) InsertAlbumIfNotExists(name string, year int, path string) (int64, error) {
	id, err := db.GetAlbumID(name, year)
	if err == nil && id != 0 {
		return id, nil
	}

	album := Album{
		Path: path,
		Name: name,
		Year: year,
	}
	err = db.InsertAlbum(&album)
	if err != nil {
		return 0, err
	}
	album.ID, err = db.GetAlbumID(album.Name, album.Year)
	return album.ID, nil
}

func (db *DataBase) GetPerformerName(performerID int64) (string, error) {
	var name string
	err := db.Db.QueryRow("SELECT name FROM performers WHERE id_performer = ?", performerID).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (db *DataBase) GetAlbumName(albumID int64) (string, error) {
	var name string
	err := db.Db.QueryRow("SELECT name FROM albums WHERE id_album = ?", albumID).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (db *DataBase) UpdateSong(idRola int64, newTitle, newGenre string, newTrack, newYear int) error {
	query := `UPDATE rolas SET title = ?, track = ?, year = ?, genre = ? WHERE id_rola = ?`
	_, err := db.Db.Exec(query, newTitle, newTrack, newYear, newGenre, idRola)
	return err
}

func (db *DataBase) UpdateAlbum(idAlbum int64, newName string, newYear int) error {
	query := `UPDATE albums SET name = ?, year = ? WHERE id_album = ?`
	_, err := db.Db.Exec(query, newName,newYear, idAlbum)
	return err
}

func (db *DataBase) UpdatePerformer(idPerformer int64, typePerf int, newName string) error {
	query := `UPDATE performers SET id_type = ?, name = ? WHERE id_performer = ?`
	_, err := db.Db.Exec(query, typePerf, newName, idPerformer)
	return err
}

func (db *DataBase) UpdateNamePerformer(idPerformer int64, newName string) error {
	query := `UPDATE performers SET name = ? WHERE id_performer = ?`
	_, err := db.Db.Exec(query, newName, idPerformer)
	return err
}

func (db *DataBase) DefinePerson(stageName, realName, birthDate, deathDate string) error {
	query := `INSERT INTO persons (stage_name, real_name, birth_date, death_date) 
              VALUES (?, ?, ?, ?)`
	_, err := db.Db.Exec(query, stageName, realName, birthDate, deathDate)
	return err
}

func (db *DataBase) GetPersonID(stageName, realName, birthDate, deathDate string) (int64, error) {
	var id int64
	query := `SELECT id_person FROM persons WHERE stage_name = ? AND real_name = ? AND birth_date = ? AND death_date = ?`
	err := db.Db.QueryRow(query, stageName, realName, birthDate, deathDate).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return id, err
}

func (db *DataBase) InsertPersonIfNotExists(stageName, realName, birthDate, deathDate string) (int64, error) {
	id, err := db.GetPersonID(stageName, realName, birthDate, deathDate)
	if err == nil && id != 0 {
		return id, nil
	}

	err = db.DefinePerson(stageName, realName, birthDate, deathDate)
	if err != nil {
		return 0, err
	}
	personID, err := db.GetPersonID(stageName, realName, birthDate, deathDate)
	return personID, nil
}

func (db *DataBase) DefineGroup(name, startDate, endDate string) error {
	query := `INSERT INTO groups (name, start_date, end_date) 
              VALUES (?, ?, ?)`
	_, err := db.Db.Exec(query, name, startDate, endDate)
	return err
}

func (db *DataBase) GetGroupID(name, startDate, endDate string) (int64, error) {
	var id int64
	query := `SELECT id_group FROM groups WHERE name = ? AND start_date = ? AND end_date = ?`
	err := db.Db.QueryRow(query, name, startDate, endDate).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return id, err
}

func (db *DataBase) InsertGroupIfNotExists(name, startDate, endDate string) (int64, error) {
	id, err := db.GetGroupID(name, startDate, endDate)
	if err == nil && id != 0 {
		return id, nil
	}

	err = db.DefineGroup(name, startDate, endDate)
	if err != nil {
		return 0, err
	}
	groupID, err := db.GetGroupID(name, startDate, endDate)
	return groupID, nil
}