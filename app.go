package main

import (
	"context"
	"os"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx    context.Context
	ctrl   *beep.Ctrl
	volume *effects.Volume
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

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return err
	}
	a.ctrl = &beep.Ctrl{Streamer: streamer, Paused: false}
	// Wrap the controller in a volume effect for dynamic volume control.
	a.volume = &effects.Volume{
		Streamer: a.ctrl,
		Base:     2,
		Volume:   0, // 0dB default volume
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

func (a *App) SetVolume(vol float64) {
	if a.volume != nil && vol > -5 {
		a.volume.Silent = false
		a.volume.Volume = vol
	} else if a.volume != nil {
		a.volume.Silent = true
	}
}
