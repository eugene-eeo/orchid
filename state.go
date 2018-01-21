package main

import "os"
import "strings"
import "path/filepath"
import "math/rand"

func remove(s *song, xs []*song) []*song {
	j := 0
	f := false
	for i, x := range xs {
		if x == s {
			f = true
			j = i
			break
		}
	}
	if f {
		return append(xs[:j], xs[j+1:]...)
	}
	return xs
}

func findSongs(dir string) (songs []*song, err error) {
	f, err := os.Open(dir)
	if err != nil {
		return
	}
	files, err := f.Readdirnames(-1)
	if err != nil {
		return
	}
	songs = []*song{}
	for _, name := range files {
		if strings.HasSuffix(name, ".mp3") {
			songs = append(songs, &song{name})
		}
	}
	return
}

func mod(r int, m int) int {
	t := r % m
	if t < 0 {
		t += m
	}
	return t
}

type state struct {
	cursor     int
	repeat     bool
	NowPlaying chan *song
	directory  string
	songs      []*song
	songsQueue chan *song
	toggle     chan bool
	stop       chan bool
	next       chan bool
}

func newState(dir string) (s *state, err error) {
	dir, err = filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	songs, err := findSongs(dir)
	s = &state{
		directory:  dir,
		songs:      songs,
		NowPlaying: make(chan *song),
		songsQueue: make(chan *song),
		toggle:     make(chan bool),
		stop:       make(chan bool),
		next:       make(chan bool),
	}
	return
}

func (s *state) currentSong() *song {
	return s.songs[s.cursor]
}

func (s *state) NameOf(so *song) string {
	return so.RelPath(s.directory)
}

func (s *state) TogglePlay() {
	s.toggle <- true
}

func (s *state) Shuffle() {
	n := len(s.songs)
	u := s.currentSong()
	for i := 0; i < n; i++ {
		j := rand.Intn(n)
		s.songs[i], s.songs[j] = s.songs[j], s.songs[i]
		if s.songs[i] == u {
			s.cursor = i
		}
		if s.songs[j] == u {
			s.cursor = j
		}
	}
}

func (s *state) Loop() {
	var stream *songStream = nil
	var err error = nil
	for {
		select {
		case <-s.toggle:
			if stream != nil {
				stream.Toggle()
			}
		case <-s.stop:
			if stream != nil {
				stream.Teardown(false)
			}
		case u := <-s.songsQueue:
			// when we are done let it naturally go to next stream
			stream, err = u.SongStream(func() {
				s.Next(1, false)
			})
			// forcefully move to next song if there are errors
			if err != nil {
				s.songs = remove(u, s.songs)
				go s.Next(1, true)
				continue
			}
			if err = stream.Play(); err != nil {
				s.songs = remove(u, s.songs)
				go s.Next(1, true)
			}
		}
	}
}

func (s *state) ToggleRepeat() {
	s.repeat = !s.repeat
}

func (s *state) Next(i int, force bool) {
	if len(s.songs) == 0 {
		return
	}
	if !s.repeat || force {
		s.cursor = mod(s.cursor+i, len(s.songs))
	}
	s.stop <- true
	s.songsQueue <- s.songs[s.cursor]
	s.NowPlaying <- s.songs[s.cursor]
}

func (s *state) Peek(i int) *song {
	if len(s.songs) == 0 {
		return nil
	}
	return s.songs[mod(s.cursor+i, len(s.songs))]
}
