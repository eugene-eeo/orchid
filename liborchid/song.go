package liborchid

import "os"
import "path/filepath"
import "github.com/faiface/beep/mp3"
import "github.com/dhowden/tag"

func FindSongs(dir string) (songs []*Song, err error) {
	f, err := os.Open(dir)
	if err != nil {
		return
	}
	files, err := f.Readdirnames(-1)
	if err != nil {
		return
	}
	songs = []*Song{}
	for _, name := range files {
		if filepath.Ext(name) == ".mp3" {
			songs = append(songs, NewSong(name))
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

func (s *Song) file() (*os.File, error) {
	return os.Open(s.path)
}

func (s *Song) Stream() (*Stream, error) {
	f, err := s.file()
	if err != nil {
		return nil, err
	}
	stream, format, _ := mp3.Decode(f)
	return NewStream(stream, format), nil
}

func (s *Song) Tags() (tag.Metadata, error) {
	f, err := s.file()
	if err != nil {
		return nil, err
	}
	return tag.ReadFrom(f)
}
