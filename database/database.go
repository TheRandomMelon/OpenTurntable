package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// DB wrapper for SQLite operations
type DB struct {
	conn *sql.DB
}

// Song represents a song in the database
type Song struct {
	ID        int
	Path      string
	Title     string
	Artist_ID int
	Album_ID  int
	Composer  string
	Comment   string
	Genre     string
	Year      string
}

// Gather where the database should be
func getDatabasePath() (string, error) {
	// Uses configuration directory. This is stored depending on OS:
	// Windows: %appdata%
	// macOS:   $HOME/Library/Application Support
	// Linux:   $XDG_CONFIG_HOME (or $HOME/.config)
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("could not get config directory: %w", err)
	}

	appDir := filepath.Join(configDir, "OpenTurntable")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", fmt.Errorf("could not create app directory: %w", err)
	}

	return filepath.Join(appDir, "app.db"), nil
}

// NewDB creates a new database connection and initializes tables
func NewDB() (*DB, error) {
	dbPath, err := getDatabasePath()
	if err != nil {
		return nil, err
	}

	// Check if database file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		// Create the database file
		file, err := os.Create(dbPath)
		if err != nil {
			return nil, fmt.Errorf("could not create database file: %w", err)
		}
		file.Close()
	}

	// Begin connection
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create tables if they don't exist
	createTables := `
    CREATE TABLE IF NOT EXISTS songs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		path TEXT NOT NULL,
		title TEXT,
		artist_id INTEGER,
		album_id INTEGER,
		composer TEXT,
		comment TEXT,
		genre TEXT,
		year TEXT,
		FOREIGN KEY (artist_id) REFERENCES artists(id),
		FOREIGN KEY (album_id) REFERENCES albums(id)
	);

	CREATE TABLE IF NOT EXISTS artists (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		pfp TEXT
	);

	CREATE TABLE IF NOT EXISTS albums (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		art TEXT,
		artist_id INTEGER,
		FOREIGN KEY (artist_id) REFERENCES artists(id)
	);

	CREATE TABLE IF NOT EXISTS tags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);

	CREATE TABLE IF NOT EXISTS tag_values (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		tag_id INTEGER,
		value TEXT,
		song_id INTEGER,
		FOREIGN KEY (tag_id) REFERENCES tags(id),
		FOREIGN KEY (song_id) REFERENCES songs(id)
	);
	`

	// If table creation fails
	if _, err = conn.Exec(createTables); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &DB{conn: conn}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// CreateSong inserts a new song into the database
func (db *DB) CreateSong(song Song) (int64, error) {
	result, err := db.conn.Exec(
		"INSERT INTO songs (path, title, artist_id, album_id, composer, comment, genre, year) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		song.Path, song.Title, song.Artist_ID, song.Album_ID, song.Composer, song.Comment, song.Genre, song.Year,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create song: %w", err)
	}

	return result.LastInsertId()
}

// GetSongs retrieves all songs from the database
func (db *DB) GetSongs() ([]Song, error) {
	rows, err := db.conn.Query("SELECT * FROM songs")
	if err != nil {
		return nil, fmt.Errorf("failed to get songs: %w", err)
	}
	defer rows.Close()

	var songs []Song
	for rows.Next() {
		var s Song
		if err := rows.Scan(&s.ID, &s.Path, &s.Title, &s.Artist_ID, &s.Album_ID, &s.Composer, &s.Comment, &s.Genre, &s.Year); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		songs = append(songs, s)
	}

	return songs, nil
}

// DeleteSong removes a song by ID
func (db *DB) DeleteSong(id int) error {
	_, err := db.conn.Exec("DELETE FROM songs WHERE id = ?", id)
	return err
}
