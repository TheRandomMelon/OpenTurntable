package playback

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/flac"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/vorbis"
	"github.com/gopxl/beep/wav"
)

type Player struct {
	ctrl     *beep.Ctrl
	volume   *effects.Volume
	filePath string
	format   beep.Format
	metadata map[string]string
	streamer beep.StreamSeekCloser
}

func NewPlayer() *Player {
	return &Player{}
}

func (p *Player) Play(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	// Close previous streamer if exists
	if p.streamer != nil {
		p.streamer.Close()
		p.streamer = nil
	}

	// Stop any current playback
	speaker.Clear()

	// Reset structures
	p.ctrl = nil
	p.volume = nil
	p.filePath = ""
	p.metadata = nil

	// Read metadata
	metadata := p.ReadMetadata(f)
	p.metadata = metadata

	// Reset file pointer
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to reset file pointer: %w", err)
	}

	// Determine file type
	ext := strings.ToLower(filepath.Ext(filePath))

	var streamer beep.StreamSeekCloser
	var format beep.Format

	switch ext {
	case ".mp3":
		streamer, format, err = mp3.Decode(f)
	case ".flac":
		streamer, format, err = flac.Decode(f)
	case ".wav":
		streamer, format, err = wav.Decode(f)
	case ".ogg":
		streamer, format, err = vorbis.Decode(f)
	default:
		return errors.New("unsupported_file_type")
	}

	if err != nil {
		return err
	}
	p.streamer = streamer

	p.ctrl = &beep.Ctrl{Streamer: streamer, Paused: false}
	p.format = format
	p.filePath = filePath

	p.volume = &effects.Volume{
		Streamer: p.ctrl,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(p.volume)

	return nil
}

func (p *Player) Pause() {
	if p.ctrl != nil {
		p.ctrl.Paused = !p.ctrl.Paused
	}
}

func (p *Player) IsPlaying() bool {
	return p.ctrl != nil && !p.ctrl.Paused
}

func (p *Player) SetVolume(vol float64) {
	if p.volume == nil {
		return
	}
	if vol > -5 {
		p.volume.Silent = false
		p.volume.Volume = vol
	} else {
		p.volume.Silent = true
	}
}

func (p *Player) Seek(seconds float64) error {
	if p.streamer == nil {
		return errors.New("seeking not supported")
	}

	targetSample := int(seconds * float64(p.format.SampleRate))
	if targetSample < 0 || targetSample > p.streamer.Len() {
		return errors.New("seek position out of bounds")
	}

	speaker.Lock()
	defer speaker.Unlock()
	return p.streamer.Seek(targetSample)
}

func (p *Player) GetPosition() (float64, error) {
	if p.streamer == nil {
		return 0, errors.New("no active stream")
	}

	speaker.Lock()
	defer speaker.Unlock()
	return float64(p.streamer.Position()) / float64(p.format.SampleRate), nil
}

func (p *Player) GetDuration() (float64, error) {
	if p.streamer == nil {
		return 0, errors.New("no active stream")
	}

	speaker.Lock()
	defer speaker.Unlock()
	return float64(p.streamer.Len()) / float64(p.format.SampleRate), nil
}

func (p *Player) GetFilePath() string {
	return p.filePath
}

func (p *Player) GetMetadata() map[string]string {
	return p.metadata
}
