package main

import "io"
import "bytes"
import "os"
import "strings"
import "path/filepath"
import "github.com/faiface/beep"
import "github.com/faiface/beep/speaker"
import "github.com/faiface/beep/mp3"
import "github.com/dhowden/tag"

//var log = make(chan string, 100)

type song struct {
	path string
}

func (s *song) abs() string {
	f, err := filepath.Abs(s.path)
	if err != nil {
		panic(err)
	}
	return f
}

func (s *song) RelPath(root string) string {
	return strings.TrimPrefix(s.abs(), root+"/")
}

func (s *song) Mp3() (beep.StreamSeekCloser, beep.Format, error) {
	f, err := os.Open(s.path)
	if err != nil {
		return nil, beep.Format{}, err
	}
	return mp3.Decode(f)
}

func (s *song) SongStream(done func()) (ss *songStream, err error) {
	stream, format, err := s.Mp3()
	if err != nil {
		return
	}
	ss = newSongStream(stream, format, done)
	return
}

func (s *song) Picture() (io.Reader, bool) {
	f, err := os.Open(s.path)
	if err != nil {
		return nil, false
	}
	m, err := tag.ReadFrom(f)
	if err != nil {
		return nil, false
	}
	if m.Picture() == nil {
		return nil, false
	}
	return bytes.NewReader(m.Picture().Data), true
}

type songStream struct {
	stream   beep.StreamCloser
	format   beep.Format
	paused   bool
	finished bool
	done     func()
}

func newSongStream(stream beep.StreamCloser, format beep.Format, done func()) *songStream {
	return &songStream{
		stream:   stream,
		format:   format,
		paused:   true,
		finished: false,
		done:     done,
	}
}

func (s *songStream) initSpeaker() error {
	return nil
}

func (s *songStream) Teardown(d bool) {
	if !s.finished {
		speaker.Lock()
		s.stream.Close()
		speaker.Unlock()
		if d {
			s.done()
		}
		s.finished = true
	}
}

func (s *songStream) Toggle() {
	s.paused = !s.paused
}

func (s *songStream) Play() (err error) {
	//log <- "Play()"
	s.paused = false
	//log <- "InitSpeaker()"
	err = s.initSpeaker()
	if err != nil {
		return
	}
	//log <- "speaker.Play(...)"
	speaker.Play(beep.Seq(
		beep.StreamerFunc(func(sample [][2]float64) (int, bool) {
			if !s.paused {
				return s.stream.Stream(sample)
			}
			for i := 0; i < len(sample); i++ {
				sample[i] = [2]float64{0, 0}
			}
			return len(sample), true
		}),
		beep.Callback(func() {
			s.Teardown(true)
		}),
	))
	//log <- "Play() end"
	return
}
