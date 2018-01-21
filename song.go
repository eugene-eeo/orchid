package main

import "time"
import "io"
import "bytes"
import "os"
import "strings"
import "path/filepath"
import "github.com/faiface/beep"
import "github.com/faiface/beep/speaker"
import "github.com/faiface/beep/mp3"
import "github.com/dhowden/tag"

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
	done     func()
	ctrl     *beep.Ctrl
	finished bool
}

func newSongStream(stream beep.StreamCloser, format beep.Format, done func()) *songStream {
	return &songStream{
		stream:   stream,
		format:   format,
		finished: false,
		done:     done,
		ctrl:     &beep.Ctrl{Streamer: stream},
	}
}

func (s *songStream) initSpeaker() error {
	return speaker.Init(s.format.SampleRate, s.format.SampleRate.N(time.Second/10))
}

func (s *songStream) Teardown(d bool) {
	if !s.finished {
		_ = s.stream.Close()
		if d {
			s.done()
		}
		s.finished = true
	}
}

func (s *songStream) Toggle() {
	speaker.Lock()
	s.ctrl.Paused = !s.ctrl.Paused
	speaker.Unlock()
}

func (s *songStream) Play() error {
	err := s.initSpeaker()
	if err != nil {
		return err
	}
	speaker.Play(beep.Seq(
		s.ctrl,
		beep.Callback(func() {
			s.Teardown(true)
		}),
	))
	return nil
}
