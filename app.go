package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
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
	ctx      context.Context
	ctrl     *beep.Ctrl
	volume   *effects.Volume
	filePath string
	seeker   beep.StreamSeeker // For seeking support
	format   beep.Format       // To store sample rate
	metadata map[string]string // Currently loaded metadata
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

func (a *App) PlayMusic() error {
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

	a.filePath = filePath

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	metadata, err := a.ReadMetadata(f)
	if err != nil {
		fmt.Println(err)
	} else {
		a.metadata = metadata
	}

	// Detect type of file and decode
	ext := strings.ToLower(filepath.Ext(filePath))

	var streamer beep.Streamer
	var format beep.Format

	switch ext {
	case ".mp3":
		streamer, format, err = mp3.Decode(f)
		if err != nil {
			return err
		}
	case ".flac":
		streamer, format, err = flac.Decode(f)
		if err != nil {
			return err
		}
	case ".wav":
		streamer, format, err = wav.Decode(f)
		if err != nil {
			return err
		}
	case ".ogg":
		streamer, format, err = vorbis.Decode(f)
		if err != nil {
			return err
		}
	default:
		return errors.New("unsupported_file_type")
	}

	a.ctrl = &beep.Ctrl{Streamer: streamer, Paused: false}

	// Check if streamer supports seeking
	if seekerStream, ok := streamer.(beep.StreamSeeker); ok {
		a.seeker = seekerStream
	} else {
		a.seeker = nil
	}

	// Store format and file path for future use
	a.format = format

	// Wraps the controller in a volume effect for dynamic volume control.
	a.volume = &effects.Volume{
		Streamer: a.ctrl,
		Base:     2,
		Volume:   0, // 0dB = default volume
		Silent:   false,
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(a.volume)

	return nil
}

func (a *App) PauseMusic() {
	if a.ctrl != nil {
		a.ctrl.Paused = !a.ctrl.Paused
	}
}

func (a *App) IsPlaying() bool {
	if a.ctrl != nil {
		return !a.ctrl.Paused
	} else {
		return false
	}
}

func (a *App) SetVolume(vol float64) {
	if a.volume != nil && vol > -5 {
		a.volume.Silent = false
		a.volume.Volume = vol
	} else if a.volume != nil {
		a.volume.Silent = true
	}
}

func (a *App) Seek(seconds float64) error {
	if a.seeker == nil {
		return errors.New("seeking not supported for this file type")
	}

	sampleRate := a.format.SampleRate
	targetSample := int(seconds * float64(sampleRate))

	if targetSample < 0 || targetSample > a.seeker.Len() {
		return errors.New("seek position out of bounds")
	}

	speaker.Lock()
	defer speaker.Unlock()
	return a.seeker.Seek(targetSample)
}

func (a *App) GetPosition() (float64, error) {
	if a.seeker == nil {
		return 0, errors.New("seeking not supported for this file type")
	}

	speaker.Lock()
	defer speaker.Unlock()
	pos := a.seeker.Position()

	finalPos := float64(pos) / float64(a.format.SampleRate)

	return finalPos, nil
}

func (a *App) GetDuration() (float64, error) {
	if a.seeker == nil {
		return 0, errors.New("seeking not supported for this file type")
	}

	speaker.Lock()
	defer speaker.Unlock()
	dur := a.seeker.Len()
	return float64(dur) / float64(a.format.SampleRate), nil
}

func (a *App) GetFilePath() string {
	return a.filePath
}

func (a *App) ReadMetadata(fileBase *os.File) (map[string]string, error) {
	metadata := make(map[string]string)

	// Read tags from file
	file, err := tag.ReadFrom(fileBase)
	if err != nil {
		return metadata, err
	}

	// Basic metadata
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

	return metadata, nil
}

func (a *App) GetMetadata() map[string]string {
	return a.metadata
}
