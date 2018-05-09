package liborchid

import "os"
import "path/filepath"
import "github.com/faiface/beep/mp3"
import "github.com/dhowden/tag"

func FindSongs(dir string) []*Song {
	songs := []*Song{}
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".mp3" {
			songs = append(songs, NewSong(path))
		}
		return nil
	})
	return songs
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

func (s *Song) file() (*os.File, error) {
	return os.Open(s.path)
}

func (s *Song) Stream() (*Stream, error) {
	f, err := s.file()
	if err != nil {
		return nil, err
	}
	stream, format, err := mp3.Decode(f)
	if err != nil {
		return nil, err
	}
	return NewStream(stream, format), nil
}

func (s *Song) Image() *tag.Picture {
	f, err := s.file()
	if err != nil {
		return nil
	}
	metadata, err := tag.ReadFrom(f)
	if err != nil {
		return nil
	}
	return metadata.Picture()
}
