package liborchid

import "time"
import "sync"

const (
	PlaybackStart = iota
	PlaybackEnd
)

type VolumeInfo struct {
	V   float64
	Min float64
	Max float64
}

type PlaybackResult struct {
	State    int
	Song     *Song
	Stream   *Stream
	Complete bool
	Error    error
}

type MWorker struct {
	mux          sync.Mutex
	stream       *Stream
	VolumeChange chan VolumeInfo
	Results      chan *PlaybackResult
	SongQueue    chan *Song
	Progress     chan float64
	stop         chan struct{}
}

func NewMWorker() *MWorker {
	return &MWorker{
		VolumeChange: make(chan VolumeInfo),
		Results:      make(chan *PlaybackResult),
		SongQueue:    make(chan *Song),
		Progress:     make(chan float64),
		stop:         make(chan struct{}),
	}
}

func (mw *MWorker) report(state int, song *Song, stream *Stream, complete bool, err error) {
	go func() {
		mw.Results <- &PlaybackResult{
			State:    state,
			Song:     song,
			Stream:   stream,
			Complete: complete,
			Error:    err,
		}
	}()
}

func (mw *MWorker) setStream(stream *Stream) {
	mw.mux.Lock()
	defer mw.mux.Unlock()
	mw.stream = stream
}

func (mw *MWorker) Stream() *Stream {
	mw.mux.Lock()
	defer mw.mux.Unlock()
	return mw.stream
}

func (mw *MWorker) Stop() {
	mw.stop <- struct{}{}
}

func (mw *MWorker) Play() {
	vol := VolumeInfo{0, -1, 0}
	interval := time.NewTicker(time.Duration(1) * time.Second)
	for {
		select {
		case song := <-mw.SongQueue:
			stream, err := song.Stream()
			if err != nil {
				mw.report(PlaybackEnd, song, nil, false, err)
				break
			}
			stream.Play()
			stream.SetVolume(vol.V, vol.Min, vol.Max)
			mw.setStream(stream)
			go func() {
				mw.Progress <- 0.0
				mw.report(PlaybackEnd, song, stream, <-stream.Complete(), nil)
				mw.setStream(nil)
			}()
		case vol = <-mw.VolumeChange:
			if s := mw.Stream(); s != nil {
				s.SetVolume(vol.V, vol.Min, vol.Max)
			}
		case <-interval.C:
			if s := mw.Stream(); s != nil {
				mw.Progress <- s.Progress()
			}
		case <-mw.stop:
			return
		}
	}
}
