package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dhowden/tag"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/flac"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/vorbis"
	"github.com/gopxl/beep/wav"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx      context.Context       // App context for Wails
	ctrl     *beep.Ctrl            // Base Ctrl struct
	volume   *effects.Volume       // Volume control
	filePath string                // Current file path being played
	format   beep.Format           // To store sample rate
	metadata map[string]string     // Currently loaded metadata
	streamer beep.StreamSeekCloser // Tracks streamer
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Perform playback with a given file selected in a dialog box.
func (a *App) PlayMusic() error {
	// Grab file path with native dialog
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

	// Set filepath variable
	a.filePath = filePath

	// Attempt to open file for reading
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	// Close previous streamer if it exists
	if a.streamer != nil {
		a.streamer.Close()
		a.streamer = nil
	}

	// Stop any current playback from the Speaker
	speaker.Clear()

	// Properly reset structures to avoid old references
	a.ctrl = nil
	a.volume = nil
	a.streamer = nil
	a.filePath = ""
	a.metadata = nil

	// Begin reading metadata
	metadata := a.ReadMetadata(f)
	a.metadata = metadata

	// Ensure file pointer is reset to beginning before decoding
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to reset file pointer: %w", err)
	}

	// Detect type of file and decode
	ext := strings.ToLower(filepath.Ext(filePath))

	var streamer beep.StreamSeekCloser
	var format beep.Format

	// Check extension of file and decode it
	switch ext {
	case ".mp3":
		streamer, format, err = mp3.Decode(f)
		if err != nil {
			return err
		}
		a.streamer = streamer
	case ".flac":
		streamer, format, err = flac.Decode(f)
		if err != nil {
			return err
		}
		a.streamer = streamer
	case ".wav":
		streamer, format, err = wav.Decode(f)
		if err != nil {
			return err
		}
		a.streamer = streamer
	case ".ogg":
		streamer, format, err = vorbis.Decode(f)
		if err != nil {
			return err
		}
		a.streamer = streamer
	default:
		return errors.New("unsupported_file_type")
	}

	// Set up base ctrl
	a.ctrl = &beep.Ctrl{Streamer: streamer, Paused: false}

	// Store format and file path for future use
	a.format = format

	// Wraps the controller in a volume effect for dynamic volume control.
	a.volume = &effects.Volume{
		Streamer: a.ctrl,
		Base:     2,
		Volume:   0, // 0dB = default volume
		Silent:   false,
	}

	// Initialize playback via the speaker
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(a.volume)

	return nil
}

// Pause/unpause currently playing music
func (a *App) PauseMusic() {
	if a.ctrl != nil {
		a.ctrl.Paused = !a.ctrl.Paused
	}
}

// Check if music is currently playing
func (a *App) IsPlaying() bool {
	if a.ctrl != nil {
		return !a.ctrl.Paused
	} else {
		return false
	}
}

// Change volume on the speaker
func (a *App) SetVolume(vol float64) {
	if a.volume != nil && vol > -5 {
		a.volume.Silent = false
		a.volume.Volume = vol
	} else if a.volume != nil {
		a.volume.Silent = true
	}
}

// Seek to a different position in a track
func (a *App) Seek(seconds float64) error {
	if a.streamer == nil {
		return errors.New("seeking not supported for this file type")
	}

	sampleRate := a.format.SampleRate
	targetSample := int(seconds * float64(sampleRate))

	if targetSample < 0 || targetSample > a.streamer.Len() {
		return errors.New("seek position out of bounds")
	}

	speaker.Lock()
	defer speaker.Unlock()
	return a.streamer.Seek(targetSample)
}

// Get current position in track
func (a *App) GetPosition() (float64, error) {
	if a.streamer == nil {
		return 0, errors.New("seeking not supported for this file type")
	}

	speaker.Lock()
	defer speaker.Unlock()
	pos := a.streamer.Position()

	finalPos := float64(pos) / float64(a.format.SampleRate)

	return finalPos, nil
}

// Get total duration of the currently selected track
func (a *App) GetDuration() (float64, error) {
	if a.streamer == nil {
		return 0, errors.New("seeking not supported for this file type")
	}

	speaker.Lock()
	defer speaker.Unlock()
	dur := a.streamer.Len()
	return float64(dur) / float64(a.format.SampleRate), nil
}

// Get the currently chosen file path
func (a *App) GetFilePath() string {
	return a.filePath
}

// Read tag metadata from supported files
func (a *App) ReadMetadata(fileBase *os.File) map[string]string {
	metadata := make(map[string]string)

	// Read tags from file
	file, err := tag.ReadFrom(fileBase)
	if err != nil {
		fmt.Println("Error reading metadata:", err)

		// Fallback to default values on error
		metadata["title"] = filepath.Base(fileBase.Name())
		metadata["artist"] = ""
		metadata["album"] = ""
		metadata["genre"] = ""
		metadata["year"] = "0"

		return metadata
	}

	// If no error, populate metadata from the file
	metadata["title"] = file.Title()
	metadata["artist"] = file.Artist()
	metadata["album"] = file.Album()
	metadata["genre"] = file.Genre()
	metadata["year"] = fmt.Sprintf("%d", file.Year())

	// Handle album art
	if pic := file.Picture(); pic != nil {
		dataURL := fmt.Sprintf("data:%s;base64,%s",
			pic.MIMEType,
			base64.StdEncoding.EncodeToString(pic.Data),
		)
		metadata["albumArt"] = dataURL
	}

	return metadata
}

// Get currently loaded metadata
func (a *App) GetMetadata() map[string]string {
	return a.metadata
}
