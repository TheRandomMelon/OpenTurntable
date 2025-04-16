package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// Wrapper for SQLite operations
type DB struct {
	conn *sql.DB
}

// Represents a song in the database
type Song struct {
	ID        int64
	Path      string
	Title     string
	Artist_ID sql.NullInt64
	Album_ID  sql.NullInt64
	Composer  sql.NullString
	Comment   sql.NullString
	Genre     sql.NullString
	Year      sql.NullString
}

// Represents an artist in the database
type Artist struct {
	ID   int64
	Name string
	PFP  string
}

// Represents an album in the database
type Album struct {
	ID        int64
	Name      string
	Art       string
	Artist_ID int64
}

// Represents a song with album/artist name and details
type SongWithDetails struct {
	Song
	ArtistName sql.NullString
	ArtistPFP  sql.NullString
	AlbumName  sql.NullString
	AlbumArt   sql.NullString
}

// Represents a playlist in the database
type Playlist struct {
	ID          int64
	Name        string
	Description string
	Picture     string
}

// Represents a playlist entry
type PlaylistEntry struct {
	ID          int64
	Playlist_ID int64
	Song_ID     int64
	ListOrder   int64
}

// Represents a playlist and its entries
type PlaylistWithEntries struct {
	Playlist Playlist
	Entries  []PlaylistEntry
}

// Represents a playlist entry with its song data included
type PlaylistEntryWithSong struct {
	ID          int64
	Playlist_ID int64
	ListOrder   int64
	Song        SongWithDetails
}

// Represents a playlist with entries (with songs) included
type PlaylistWithSongs struct {
	Playlist Playlist
	Entries  []PlaylistEntryWithSong
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

// Creates a new database connection and initializes tables
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

	CREATE TABLE IF NOT EXISTS playlists (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		description TEXT,
		picture TEXT
	);

	CREATE TABLE IF NOT EXISTS playlist_entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		playlist_id INTEGER,
		song_id INTEGER,
		list_order INTEGER,
		FOREIGN KEY (playlist_id) REFERENCES playlists(id),
		FOREIGN KEY (song_id) REFERENCES songs(id)
	);
	`

	// If table creation fails
	if _, err = conn.Exec(createTables); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &DB{conn: conn}, nil
}

// Closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// Inserts a new song into the database
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

// Inserts a new artist into the database
func (db *DB) CreateArtist(artist Artist) (int64, error) {
	result, err := db.conn.Exec(
		"INSERT INTO artists (name, pfp) VALUES (?, ?)",
		artist.Name, artist.PFP,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create artist: %w", err)
	}

	return result.LastInsertId()
}

// Retrieves an artist from the database by ID
func (db *DB) GetArtistById(id int64) (Artist, error) {
	var artist Artist
	err := db.conn.QueryRow("SELECT * FROM artists WHERE id = ?", id).Scan(&artist.ID, &artist.Name, &artist.PFP)
	if err != nil {
		if err == sql.ErrNoRows {
			return Artist{}, fmt.Errorf("artist with ID %d not found", id)
		}
		return Artist{}, err
	}

	return artist, nil
}

// Retrieves an album from the database by ID
func (db *DB) GetAlbumById(id int64) (Album, error) {
	var album Album
	err := db.conn.QueryRow("SELECT * FROM albums WHERE id = ?", id).Scan(&album.ID, &album.Name, &album.Art, &album.Artist_ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return Album{}, fmt.Errorf("album with ID %d not found", id)
		}
		return Album{}, err
	}

	return album, nil
}

// Retrieves an artist from the database by name
func (db *DB) GetArtistByName(name string) (Artist, error) {
	var artist Artist
	err := db.conn.QueryRow("SELECT * FROM artists WHERE name = ?", name).Scan(&artist.ID, &artist.Name, &artist.PFP)
	if err != nil {
		if err == sql.ErrNoRows {
			return Artist{}, fmt.Errorf("artist with name %s not found", name)
		}
		return Artist{}, err
	}

	return artist, nil
}

// Retrieves an album from the database by name and artist ID
func (db *DB) GetAlbumByNameAndArtistId(name string, artist_id int64) (Album, error) {
	var album Album
	err := db.conn.QueryRow("SELECT * FROM albums WHERE name = ? AND artist_id = ?", name, artist_id).Scan(&album.ID, &album.Name, &album.Art, &album.Artist_ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return Album{}, fmt.Errorf("album with name %s and artist_id %d not found", name, artist_id)
		}
		return Album{}, err
	}

	return album, nil
}

// Inserts a new album into the database
func (db *DB) CreateAlbum(album Album) (int64, error) {
	result, err := db.conn.Exec(
		"INSERT INTO albums (name, art, artist_id) VALUES (?, ?, ?)",
		album.Name, album.Art, album.Artist_ID,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create album: %w", err)
	}

	return result.LastInsertId()
}

// Retrieves song by ID
func (db *DB) GetSongById(id int64) (Song, error) {
	var song Song
	err := db.conn.QueryRow("SELECT * FROM songs WHERE id = ?", id).Scan(&song.ID, &song.Path, &song.Title, &song.Album_ID, &song.Artist_ID, &song.Composer, &song.Comment, &song.Genre, &song.Year)
	if err != nil {
		if err == sql.ErrNoRows {
			return Song{}, fmt.Errorf("song with ID %d not found", id)
		}
		return Song{}, err
	}

	return song, nil
}

// Retrieves a song by file path
func (db *DB) GetSongByPath(path string) (Song, error) {
	var song Song
	err := db.conn.QueryRow("SELECT * FROM songs WHERE path = ?", path).Scan(&song.ID, &song.Path, &song.Title, &song.Album_ID, &song.Artist_ID, &song.Composer, &song.Comment, &song.Genre, &song.Year)
	if err != nil {
		if err == sql.ErrNoRows {
			return Song{}, fmt.Errorf("song with path %s not found", path)
		}
		return Song{}, err
	}

	return song, nil
}

// Retrieves all songs from the database
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
			return nil, fmt.Errorf("failed to scan song: %w", err)
		}
		songs = append(songs, s)
	}

	return songs, nil
}

// Removes a song by ID
func (db *DB) DeleteSong(id int64) error {
	_, err := db.conn.Exec("DELETE FROM songs WHERE id = ?", id)
	return err
}

// Retrieves a song with details (album name/art, artist name/pfp)
func (db *DB) GetSongsWithDetails() ([]SongWithDetails, error) {
	query := `
        SELECT 
            songs.id,
            songs.path,
            songs.title,
            songs.artist_id,
            songs.album_id,
            songs.composer,
            songs.comment,
            songs.genre,
            songs.year,
            COALESCE(artists.name, '') AS artist_name,
            COALESCE(artists.pfp, '') AS artist_pfp,
            COALESCE(albums.name, '') AS album_name,
            COALESCE(albums.art, '') AS album_art
        FROM songs
        LEFT JOIN artists ON songs.artist_id = artists.id
        LEFT JOIN albums ON songs.album_id = albums.id
    `
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get songs (w/ details): %w", err)
	}
	defer rows.Close()

	var songs []SongWithDetails
	for rows.Next() {
		var s SongWithDetails
		err := rows.Scan(
			&s.ID,
			&s.Path,
			&s.Title,
			&s.Artist_ID,
			&s.Album_ID,
			&s.Composer,
			&s.Comment,
			&s.Genre,
			&s.Year,
			&s.ArtistName,
			&s.ArtistPFP,
			&s.AlbumName,
			&s.AlbumArt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan song details: %w", err)
		}
		songs = append(songs, s)
	}
	return songs, nil
}

/// ===========
///  PLAYLISTS
/// ===========

// Inserts a new playlist into the database
func (db *DB) CreatePlaylist(playlist Playlist) (int64, error) {
	result, err := db.conn.Exec(
		"INSERT INTO playlists (name, description, picture) VALUES (?, ?, ?)",
		playlist.Name, playlist.Description, playlist.Picture,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create playlist: %w", err)
	}

	return result.LastInsertId()
}

// Gets all playlists
func (db *DB) GetPlaylists() ([]Playlist, error) {
	rows, err := db.conn.Query("SELECT * FROM playlists")
	if err != nil {
		return nil, fmt.Errorf("failed to get playlists: %w", err)
	}
	defer rows.Close()

	var playlists []Playlist
	for rows.Next() {
		var p Playlist
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Picture); err != nil {
			return nil, fmt.Errorf("failed to scan playlist: %w", err)
		}
		playlists = append(playlists, p)
	}

	return playlists, nil
}

// Inserts a new playlist entry into the database
func (db *DB) CreatePlaylistEntry(pe PlaylistEntry) (int64, error) {
	result, err := db.conn.Exec(
		"INSERT INTO playlist_entries (playlist_id, song_id, list_order) VALUES (?, ?, ?)",
		pe.Playlist_ID, pe.Song_ID, pe.ListOrder,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create playlist: %w", err)
	}

	return result.LastInsertId()
}

// Gets a playlist with all of its song entries
func (db *DB) GetPlaylistWithSongs(playlistID int64) (*PlaylistWithSongs, error) {
	var playlist Playlist

	// Fetch playlist metadata
	err := db.conn.QueryRow(`
		SELECT ID, Name, Description, Picture
		FROM Playlists
		WHERE ID = ?`, playlistID).Scan(
		&playlist.ID,
		&playlist.Name,
		&playlist.Description,
		&playlist.Picture,
	)
	if err != nil {
		return nil, err
	}

	// Fetch playlist entries with detailed song info
	rows, err := db.conn.Query(`
		SELECT 
			pe.ID, pe.Playlist_ID, pe.ListOrder,

			s.ID, s.Path, s.Title, s.Artist_ID, s.Album_ID,
			s.Composer, s.Comment, s.Genre, s.Year,

			ar.Name AS ArtistName, ar.PFP AS ArtistPFP,
			al.Name AS AlbumName, al.Art AS AlbumArt

		FROM PlaylistEntries pe
		JOIN Songs s ON pe.Song_ID = s.ID
		LEFT JOIN Artists ar ON s.Artist_ID = ar.ID
		LEFT JOIN Albums al ON s.Album_ID = al.ID
		WHERE pe.Playlist_ID = ?
		ORDER BY pe.ListOrder
	`, playlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []PlaylistEntryWithSong
	for rows.Next() {
		var entry PlaylistEntryWithSong
		var song SongWithDetails

		err := rows.Scan(
			&entry.ID,
			&entry.Playlist_ID,
			&entry.ListOrder,

			&song.ID,
			&song.Path,
			&song.Title,
			&song.Artist_ID,
			&song.Album_ID,
			&song.Composer,
			&song.Comment,
			&song.Genre,
			&song.Year,

			&song.ArtistName,
			&song.ArtistPFP,
			&song.AlbumName,
			&song.AlbumArt,
		)
		if err != nil {
			return nil, err
		}

		entry.Song = song
		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &PlaylistWithSongs{
		Playlist: playlist,
		Entries:  entries,
	}, nil
}
