package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"openturntable/database"
	"openturntable/playback"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx    context.Context
	player *playback.Player
	db     *database.DB
}

func NewApp() *App {
	return &App{
		player: playback.NewPlayer(),
	}
}

// Called at application startup. Currently initializes the DB
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	db, err := database.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	a.db = db

	// Get all songs
	songs, err := a.db.GetSongs()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Songs:")
	for _, song := range songs {
		fmt.Printf("%d: %s (%s)\n", song.ID, song.Path, song.Title)
	}
}

// Called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// Called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Close database
	a.db.Close()
}

/// =================
///  PLAYER BINDINGS
/// =================

// Select file and tell player to begin playing the file
func (a *App) SelectAndPlayFile() error {
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Audio Files (*.mp3, *.wav, *.flac, *.ogg)",
				Pattern:     "*.mp3;*.wav;*.flac;*.ogg",
			},
		},
	})
	if err != nil {
		return err
	}
	return a.player.Play(filePath)
}

func (a *App) PlayFile(filePath string) error {
	return a.player.Play(filePath)
}

// Binding to call  pause function in player
func (a *App) PauseMusic() {
	a.player.Pause()
}

// Binding to call IsPlaying in player
func (a *App) IsPlaying() bool {
	return a.player.IsPlaying()
}

// Binding to call SetVolume in player
func (a *App) SetVolume(vol float64) {
	a.player.SetVolume(vol)
}

// Binding to call Seek in player
func (a *App) Seek(seconds float64) error {
	return a.player.Seek(seconds)
}

// Binding to call GetPosition in player, alongside event handling
func (a *App) GetPosition() (float64, error) {
	position, err := a.player.GetPosition()
	if err != nil {
		return 0, err
	}

	duration, err2 := a.GetDuration()
	if err2 != nil {
		return 0, err
	}

	if position >= duration {
		runtime.EventsEmit(a.ctx, "playbackComplete")
	}

	return position, err
}

// Binding to call GetDuration in player
func (a *App) GetDuration() (float64, error) {
	return a.player.GetDuration()
}

// Binding to call GetFilePath in player
func (a *App) GetFilePath() string {
	return a.player.GetFilePath()
}

// Binding to call GetMetadata in player
func (a *App) GetMetadata() map[string]string {
	return a.player.GetMetadata()
}

/// =================
///    DB BINDINGS
/// =================

// Binding to call GetSongs in db
func (a *App) GetSongs() ([]database.Song, error) {
	return a.db.GetSongs()
}

// Binding to call GetSongsWithDetails in db
func (a *App) GetSongsWithDetails() ([]database.SongWithDetails, error) {
	return a.db.GetSongsWithDetails()
}

// Inserts all songs in directory (recursive) from user provided directory
func (a *App) ImportSongsFromDirectory() (string, error) {
	// Have user choose directory
	dirPath, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Audio Files (*.mp3, *.wav, *.flac, *.ogg)",
				Pattern:     "*.mp3;*.wav;*.flac;*.ogg",
			},
		},
	})
	if err != nil {
		return "", err
	}

	runtime.EventsEmit(a.ctx, "toggleImporting")

	allowed := map[string]bool{
		".mp3":  true,
		".flac": true,
		".wav":  true,
		".ogg":  true,
	}

	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing %s: %v\n", path, err)
			return nil
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		runtime.EventsEmit(a.ctx, "currentImportFileWorking", path)

		// Check extension
		ext := strings.ToLower(filepath.Ext(path))
		if allowed[ext] {
			a.CreateSongFromFilePath(path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		return "", nil
	}

	runtime.EventsEmit(a.ctx, "toggleImporting")

	return dirPath, nil
}

// Inserts a new song into the database from file provided
func (a *App) CreateSongFromFilePath(filePath string) (int64, error) {
	// Check if file exists in the DB, and delete the old record if so
	existingSong, err := a.db.GetSongByPath(filePath)
	if err != nil {
		log.Println("error finding existing song: ", err)
	} else {
		err = a.db.DeleteSong(existingSong.ID)
		if err != nil {
			log.Println("failed to delete old song record: ", err)
		}
	}

	// Open file for reading
	f, err := os.Open(filePath)
	if err != nil {
		return -1, err
	}

	// Read metadata
	metadata := playback.ReadMetadata(f)

	// Initialize variables
	var artist database.Artist
	var album database.Album
	var song database.Song

	// Check for artist
	if artistName, ok := metadata["artist"]; ok && artistName != "" {
		// Check if artist already exists
		artist, err = a.db.GetArtistByName(artistName)
		if err != nil {
			// Try creating artist upon error (assuming artist is not found)
			log.Println("error finding existing artist:", err)

			artist = database.Artist{
				Name: metadata["artist"],
				PFP:  "",
			}

			createArtist, err := a.db.CreateArtist(artist)
			if err != nil {
				log.Println("error creating artist: ", err)
			}

			artist.ID = createArtist
		}
	}

	// Check for album
	if albumName, ok := metadata["album"]; ok && albumName != "" {
		album, err = a.db.GetAlbumByNameAndArtistId(albumName, artist.ID)
		if err != nil {
			log.Println("error finding existing album:", err)

			// Try creating album upon error (assuming album is not found)
			album = database.Album{
				Name:      metadata["album"],
				Art:       metadata["albumArt"],
				Artist_ID: artist.ID,
			}

			createAlbum, err := a.db.CreateAlbum(album)
			if err != nil {
				log.Println("error creating album: ", err)
			}

			album.ID = createAlbum
		}
	}

	// Create song
	song = database.Song{
		Path:      filePath,
		Title:     metadata["title"],
		Artist_ID: sql.NullInt64{Int64: 0, Valid: false},
		Album_ID:  sql.NullInt64{Int64: 0, Valid: false},
		Composer:  sql.NullString{String: metadata["composer"], Valid: metadata["composer"] != ""},
		Comment:   sql.NullString{String: metadata["comment"], Valid: metadata["comment"] != ""},
		Genre:     sql.NullString{String: metadata["genre"], Valid: metadata["genre"] != ""},
		Year:      sql.NullString{String: metadata["year"], Valid: metadata["year"] != ""},
	}

	// Set artist ID if valid
	if artist.ID != 0 {
		song.Artist_ID = sql.NullInt64{
			Int64: artist.ID,
			Valid: true,
		}
	}

	// Set album ID if valid
	if album.ID != 0 {
		song.Album_ID = sql.NullInt64{
			Int64: album.ID,
			Valid: true,
		}
	}

	createSong, err := a.db.CreateSong(song)
	if err != nil {
		log.Println("error creating song: ", err)
	}

	return createSong, nil
}

// Inserts a new song into the database with file selection
func (a *App) ChooseAndCreateSong() (int64, error) {
	// Have user choose file
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Audio Files (*.mp3, *.wav, *.flac, *.ogg)",
				Pattern:     "*.mp3;*.wav;*.flac;*.ogg",
			},
		},
	})
	if err != nil {
		return -1, err
	}

	runtime.EventsEmit(a.ctx, "toggleImporting")

	song, err := a.CreateSongFromFilePath(filePath)
	if err != nil {
		return -1, err
	}

	runtime.EventsEmit(a.ctx, "toggleImporting")

	return song, nil
}
