package playback

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dhowden/tag"
)

func (p *Player) ReadMetadata(file *os.File) map[string]string {
	metadata := make(map[string]string)

	tags, err := tag.ReadFrom(file)
	if err != nil {
		metadata["title"] = filepath.Base(file.Name())
		metadata["artist"] = ""
		metadata["album"] = ""
		metadata["albumartist"] = ""
		metadata["composer"] = ""
		metadata["comment"] = ""
		metadata["genre"] = ""
		metadata["year"] = "0"
		return metadata
	}

	title := tags.Title()
	if title == "" {
		title = filepath.Base(file.Name())
	}
	metadata["title"] = title
	metadata["artist"] = tags.Artist()
	metadata["album"] = tags.Album()
	metadata["albumartist"] = tags.AlbumArtist()
	metadata["composer"] = tags.Composer()
	metadata["comment"] = tags.Comment()
	metadata["genre"] = tags.Genre()
	metadata["year"] = fmt.Sprintf("%d", tags.Year())

	if pic := tags.Picture(); pic != nil {
		metadata["albumArt"] = fmt.Sprintf(
			"data:%s;base64,%s",
			pic.MIMEType,
			base64.StdEncoding.EncodeToString(pic.Data),
		)
	}

	return metadata
}
