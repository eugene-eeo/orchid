package player

import "os"
import "path/filepath"
import "time"
import "strings"
import "sync"

import "github.com/faiface/beep"
import "github.com/faiface/beep/speaker"
import "github.com/faiface/beep/mp3"
import "github.com/dhowden/tag"

type Song string

func (s *Song) Name() string {
	return strings.TrimSuffix(
		filepath.Base(string(*s)),
		filepath.Ext(string(*s)),
	)
}

func (s *Song) Stream(done func(graceful bool)) (*Stream, error) {
	f, err := os.Open(string(*s))
	if err != nil {
		return nil, err
	}
	stream, format, err := mp3.Decode(f)
	if err != nil {
		return nil, err
	}
	return NewStream(stream, format, func(graceful bool) {
		defer f.Close()
		done(graceful)
	}), nil
}

func (s *Song) Picture() ([]byte, bool) {
	f, err := os.Open(string(*s))
	if err != nil {
		return nil, false
	}
	defer f.Close()
	m, err := tag.ReadFrom(f)
	if err != nil {
		return nil, false
	}
	if m.Picture() == nil {
		return nil, false
	}
	return m.Picture().Data, true
}

type Stream struct {
	ctrl   *beep.Ctrl
	stream beep.StreamCloser
	format beep.Format
	done   func(bool)
	once   sync.Once
}

func NewStream(stream beep.StreamCloser, format beep.Format, done func(bool)) *Stream {
	return &Stream{
		stream: stream,
		ctrl:   &beep.Ctrl{Streamer: stream},
		format: format,
		done:   done,
	}
}

func (s *Stream) initSpeaker() error {
	return speaker.Init(s.format.SampleRate, s.format.SampleRate.N(time.Second/10))
}

func (s *Stream) Teardown(d bool) {
	s.once.Do(func() {
		s.stream.Close()
		s.done(d)
	})
}

func (s *Stream) Paused() bool {
	return s.ctrl.Paused
}

func (s *Stream) Toggle() {
	speaker.Lock()
	s.ctrl.Paused = !s.ctrl.Paused
	speaker.Unlock()
}

func (s *Stream) Play() error {
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
