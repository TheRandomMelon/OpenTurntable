package main

import (
	"context"
	"openturntable/playback"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx    context.Context
	player *playback.Player
}

func NewApp() *App {
	return &App{
		player: playback.NewPlayer(),
	}
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

// Binding to call GetPosition in player
func (a *App) GetPosition() (float64, error) {
	return a.player.GetPosition()
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
