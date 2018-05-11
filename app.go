package main

import "flag"
import "time"
import "math/rand"
import "github.com/nsf/termbox-go"
import "github.com/eugene-eeo/orchid/liborchid"
import "github.com/eugene-eeo/orchid/reactor"

const MAX_VOLUME float64 = +0.0
const MIN_VOLUME float64 = -4.0

var REACTOR *reactor.Reactor = nil

type request func(*hub)

type hub struct {
	Player   *liborchid.Player
	Stream   *liborchid.Stream
	MWorker  *liborchid.MWorker
	view     *playerView
	requests chan request
	done     chan struct{}
	volume   float64
	stream   chan termbox.Event
}

func (h *hub) Paused() bool {
	return h.Stream != nil && h.Stream.Paused()
}

func (h *hub) Toggle() {
	if h.Stream != nil {
		h.Stream.Toggle()
	}
}

func newHub(p *liborchid.Player) *hub {
	mw := liborchid.NewMWorker()
	return &hub{
		Player:   p,
		MWorker:  mw,
		view:     newPlayerView(),
		requests: make(chan request),
		done:     make(chan struct{}),
		stream:   make(chan termbox.Event),
	}
}

func (h *hub) Render() {
	h.view.Update(
		h.Player,
		h.Paused(),
		h.Player.Shuffle,
		h.Player.Repeat,
	)
}

func (h *hub) Play() {
	song := h.Player.Song()
	if song == nil {
		return
	}
	h.MWorker.SongQueue <- song
}

func (h *hub) handle(evt termbox.Event) {
	switch evt.Ch {
	case 'q':
		go func() { h.done <- struct{}{} }()
	case 'n':
		h.Player.Next(1, true)
		h.Play()
	case 'p':
		h.Player.Next(-1, true)
		h.Play()
	case 's':
		h.Player.ToggleShuffle()
	case 'r':
		h.Player.ToggleRepeat()
	case 'f':
		f := newFinderUIFromPlayer(h.Player)
		REACTOR.Focus(f)
		go func() {
			f.Loop()
			h.requests <- func(h *hub) {
				song := f.Choice()
				if song != nil {
					h.Player.SetCurrent(song)
					h.Play()
				}
			}
		}()
	}
	switch evt.Key {
	case termbox.KeySpace:
		h.Toggle()
	case termbox.KeyArrowLeft:
		fallthrough
	case termbox.KeyArrowRight:
		v := newVolumeUI(h.Stream)
		REACTOR.Focus(v)
		go func() {
			v.Loop()
			h.requests <- func(h *hub) {
				h.volume = v.volume
				if h.Stream != nil {
					h.Stream.SetVolume(h.volume, MIN_VOLUME, MAX_VOLUME)
				}
			}
		}()
	}
}

func (h *hub) MWLoop() {
	go h.MWorker.Play()
	for {
		res := <-h.MWorker.Results
		// exit signal
		if res == nil {
			break
		}
		switch res.State {
		case liborchid.PlaybackStart:
			res.Stream.SetVolume(h.volume, MIN_VOLUME, MAX_VOLUME)
			h.requests <- func(h *hub) {
				h.Stream = res.Stream
			}
		case liborchid.PlaybackEnd:
			h.requests <- func(h *hub) {
				h.Stream = nil
				if res.Error != nil {
					// playback error
					h.Player.Remove()
				} else {
					// current song interrupted, jump to next song
					h.Player.Next(1, false)
				}
				h.Play()
			}
		}
	}
}

func (h *hub) Sink() chan termbox.Event {
	return h.stream
}

func (h *hub) Loop() {
	h.Play()
loop:
	for {
		if REACTOR.InFocus(h) {
			h.Render()
		}
		select {
		case evt := <-h.stream:
			h.handle(evt)
		case r := <-h.requests:
			r(h)
		case <-h.done:
			h.MWorker.Stop <- liborchid.SIGNAL
			break loop
		}
	}
}

func main() {
	recursive := flag.Bool("r", true, "Whether orchid looks recursively for .mp3 files")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	songs := liborchid.FindSongs(".", *recursive)

	must(termbox.Init())
	termbox.SetOutputMode(termbox.Output256)
	defer termbox.Close()

	h := newHub(liborchid.NewPlayer(songs))
	REACTOR = reactor.NewReactor(h)

	go REACTOR.Loop()
	go h.MWLoop()
	h.Loop()
}
