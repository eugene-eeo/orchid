package liborchid

import "sync"
import "time"
import "github.com/faiface/beep"
import "github.com/faiface/beep/effects"
import "github.com/faiface/beep/speaker"

type VolumeInfo struct {
	V   float64
	Min float64
	Max float64
}

func (v VolumeInfo) Volume() float64 {
	if v.V > v.Max {
		return v.Max
	}
	if v.V < v.Min {
		return v.Min
	}
	return v.V
}

func (v VolumeInfo) Silent() bool {
	return v.V <= v.Min
}

type Stream struct {
	stream     beep.StreamSeekCloser
	format     beep.Format
	volume     *effects.Volume
	ctrl       *beep.Ctrl
	done       chan bool
	finishOnce sync.Once
}

func NewStream(stream beep.StreamSeekCloser, format beep.Format) *Stream {
	volume := &effects.Volume{
		Streamer: stream,
		Volume:   0,
		Base:     2,
		Silent:   false,
	}
	return &Stream{
		stream: stream,
		format: format,
		volume: volume,
		ctrl:   &beep.Ctrl{Streamer: volume},
		done:   make(chan bool),
	}
}

func (s *Stream) finish(completed bool) {
	s.finishOnce.Do((func() {
		_ = s.stream.Close()
		s.done <- completed
	}))
}

func (s *Stream) Stop() {
	s.finish(false)
}

func (s *Stream) Play() {
	_ = speaker.Init(s.format.SampleRate, s.format.SampleRate.N(time.Second/10))
	speaker.Play(beep.Seq(
		s.ctrl,
		beep.Callback(func() {
			s.finish(true)
		}),
	))
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

func (s *Stream) Volume() float64 {
	return s.volume.Volume
}

func (s *Stream) SetVolume(v VolumeInfo) {
	s.volume.Volume = v.Volume()
	s.volume.Silent = v.Silent()
}

func (s *Stream) Progress() float64 {
	return float64(s.stream.Position()) / float64(s.stream.Len())
}
