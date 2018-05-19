package liborchid

import (
	"os"
	"path/filepath"

	"github.com/dhowden/tag"
	"github.com/faiface/beep/mp3"
)

func FindSongs(dir string, recursive bool) (songs []*Song) {
	songs = []*Song{}
	if recursive {
		_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && filepath.Ext(info.Name()) == ".mp3" {
				songs = append(songs, NewSong(path))
			}
			return err
		})
		return
	}
	f, err := os.Open(dir)
	if err != nil {
		return
	}
	defer f.Close()
	if names, err := f.Readdirnames(-1); err == nil {
		for _, name := range names {
			if filepath.Ext(name) == ".mp3" {
				songs = append(songs, NewSong(filepath.Join(dir, name)))
			}
		}
	}
	return
}

type Song struct {
	path string
}

func NewSong(path string) *Song {
	return &Song{path: path}
}

func (s *Song) Name() string {
	u := filepath.Base(s.path)
	ext := filepath.Ext(u)
	return u[:len(u)-len(ext)]
}

func (s *Song) Stream() (*Stream, error) {
	f, err := os.Open(s.path)
	if err != nil {
		return nil, err
	}
	stream, format, err := mp3.Decode(f)
	if err != nil {
		return nil, err
	}
	return NewStream(stream, format), nil
}

func (s *Song) Metadata() tag.Metadata {
	f, err := os.Open(s.path)
	if err != nil {
		return nil
	}
	defer f.Close()
	metadata, err := tag.ReadFrom(f)
	if err != nil {
		return nil
	}
	return metadata
}
