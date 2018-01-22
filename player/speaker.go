package player

type Speaker struct {
	Stream *Stream
}

func NewSpeaker() *Speaker {
	return &Speaker{nil}
}

func (s *Speaker) Toggle() {
	if s.Stream != nil {
		s.Stream.Toggle()
	}
}

func (s *Speaker) Paused() bool {
	if s.Stream != nil {
		return s.Stream.Paused()
	}
	return true
}

func (s *Speaker) Stop() {
	if s.Stream != nil {
		s.Stream.Teardown(false)
	}
}

func (s *Speaker) getStream(song Song, res chan bool) (*Stream, error) {
	stream, err := song.Stream(func(graceful bool) {
		go func() {
			res <- graceful
			close(res)
		}()
	})
	if err != nil {
		return stream, err
	}
	return stream, stream.Play()
}

func (s *Speaker) Play(song Song) (<-chan bool, error) {
	s.Stop()
	q := make(chan bool)
	stream, err := s.getStream(song, q)
	if err != nil {
		if stream != nil {
			stream.Teardown(false)
		}
		s.Stream = nil
		return nil, err
	}
	s.Stream = stream
	return q, nil
}
