package liborchid

var SIGNAL struct{} = struct{}{}

const (
	PlaybackStart = iota
	PlaybackEnd
)

type PlaybackResult struct {
	State    int
	Song     *Song
	Stream   *Stream
	Complete bool
	Error    error
}

type MWorker struct {
	Results   chan *PlaybackResult
	SongQueue chan *Song
	Stop      chan struct{}
}

func NewMWorker() *MWorker {
	return &MWorker{
		Results:   make(chan *PlaybackResult),
		SongQueue: make(chan *Song),
		Stop:      make(chan struct{}),
	}
}

func (mw *MWorker) report(state int, song *Song, stream *Stream, complete bool, err error) {
	mw.Results <- &PlaybackResult{
		State:    state,
		Song:     song,
		Stream:   stream,
		Complete: complete,
		Error:    err,
	}
}

func (mw *MWorker) Play() {
loop:
	for {
		select {
		case song := <-mw.SongQueue:
			stream, err := song.Stream()
			if err != nil {
				mw.report(PlaybackEnd, song, nil, false, err)
				break
			}
			stream.Play()
			mw.report(PlaybackStart, song, stream, false, nil)
			go func() {
				mw.report(
					PlaybackEnd,
					song,
					stream,
					<-stream.Complete(),
					nil,
				)
			}()
		case <-mw.Stop:
			mw.Results <- nil
			close(mw.Results)
			break loop
		}
	}
}
