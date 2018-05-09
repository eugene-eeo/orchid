package liborchid

import "sync"
import "time"
import "github.com/faiface/beep"
import "github.com/faiface/beep/speaker"

type Stream struct {
	stream     beep.StreamCloser
	format     beep.Format
	ctrl       *beep.Ctrl
	done       chan bool
	finishOnce sync.Once
	playOnce   sync.Once
}

func NewStream(stream beep.StreamCloser, format beep.Format) *Stream {
	return &Stream{
		stream: stream,
		format: format,
		ctrl:   &beep.Ctrl{Streamer: stream},
		done:   make(chan bool),
	}
}

func (s *Stream) finish(completed bool) {
	s.finishOnce.Do((func() {
		s.stream.Close()
		s.done <- completed
	}))
}

func (s *Stream) Stop() {
	s.finish(false)
}

func (s *Stream) Play() {
	s.playOnce.Do(func() {
		speaker.Init(s.format.SampleRate, s.format.SampleRate.N(time.Second/10))
		speaker.Play(beep.Seq(
			s.ctrl,
			beep.Callback(func() {
				s.finish(true)
			}),
		))
	})
}

func (s *Stream) Toggle() bool {
	speaker.Lock()
	defer speaker.Unlock()
	s.ctrl.Paused = !s.ctrl.Paused
	return s.ctrl.Paused
}

func (s *Stream) Paused() bool {
	return s.ctrl.Paused
}

func (s *Stream) Complete() <-chan bool {
	return s.done
}
