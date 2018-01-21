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

type playerRequest func(*playerState) *playerState

type playerState struct {
	stream *songStream
	cursor int
	Repeat bool
	Queue  []*song
}

func (s *playerState) Song() *song {
	return s.Queue[s.cursor]
}

func (s *playerState) Peek(i int) *song {
	if len(s.Queue) == 0 {
		return nil
	}
	return s.Queue[mod(s.cursor+i, len(s.Queue))]
}

func (s *playerState) Paused() bool {
	return s.stream.ctrl.Paused
}

type state struct {
	directory string
	State     chan *playerState
	Requests  chan playerRequest
	songs     []*song
}

func newState(dir string) (s *state, err error) {
	dir, err = filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	songs, err := findSongs(dir)
	s = &state{
		directory: dir,
		songs:     songs,
		State:     make(chan *playerState),
		Requests:  make(chan playerRequest),
	}
	return
}

func (s *state) TogglePlay() {
	s.Requests <- func(st *playerState) *playerState {
		st.stream.Toggle()
		return st
	}
}

func (s *state) Shuffle() {
	s.Requests <- func(st *playerState) *playerState {
		n := len(st.Queue)
		z := st.Song()
		for i := 0; i < n; i++ {
			j := rand.Intn(n)
			st.Queue[i], st.Queue[j] = st.Queue[j], st.Queue[i]
			if st.Queue[i] == z {
				st.cursor = i
			}
			if st.Queue[j] == z {
				st.cursor = j
			}
		}
		return st
	}
}

func (s *state) Loop() {
	state := &playerState{
		Repeat: false,
		Queue:  s.songs,
		stream: nil,
		cursor: 0,
	}
	for {
		req := <-s.Requests
		state = req(state)
		s.State <- state
	}
}

func (s *state) ToggleRepeat() {
	s.Requests <- func(st *playerState) *playerState {
		st.Repeat = !st.Repeat
		return st
	}
}

func (s *state) Next(i int, force bool) {
	s.Requests <- func(st *playerState) *playerState {
		if len(st.Queue) == 0 {
			return st
		}
		if !st.Repeat || force {
			st.cursor = mod(st.cursor+i, len(st.Queue))
		}
		sng := st.Song()
		stream, err := sng.SongStream(func() {
			// when we are done let it naturally go to next stream
			s.Next(1, false)
		})
		if err != nil {
			return st
		}
		if err = stream.Play(); err != nil {
			return st
		}
		st.stream = stream
		return st
	}
}
