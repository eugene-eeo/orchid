package main

import "flag"
import "time"
import "math/rand"
import "github.com/nsf/termbox-go"
import "github.com/eugene-eeo/orchid/liborchid"

const MAX_VOLUME float64 = +0.0
const MIN_VOLUME float64 = -4.0

type request func(*hub)

type hub struct {
	Player   *liborchid.Player
	Stream   *liborchid.Stream
	MWorker  *liborchid.MWorker
	view     *playerView
	requests chan request
	done     chan struct{}
	volume   float64
	focus    bool
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
		focus:    true,
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

func (h *hub) handle(events <-chan termbox.Event, evt termbox.Event, done chan struct{}) {
	if evt.Type != termbox.EventKey {
		return
	}
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
		h.focus = false
		go func() {
			f := newFinderUIFromPlayer(h.Player)
			f.Loop(events)
			h.requests <- func(h *hub) {
				h.focus = true
				song := f.Choice()
				if song != nil {
					h.Player.SetCurrent(song)
					h.Play()
				}
			}
			done <- struct{}{}
		}()
		return
	}
	switch evt.Key {
	case termbox.KeySpace:
		h.Toggle()
	case termbox.KeyArrowLeft:
		fallthrough
	case termbox.KeyArrowRight:
		h.focus = false
		go func() {
			v := newVolumeUI(h.Stream)
			v.Loop(events)
			h.volume = v.volume
			if h.Stream != nil {
				h.Stream.SetVolume(h.volume, MIN_VOLUME, MAX_VOLUME)
			}
			h.focus = true
			done <- struct{}{}
		}()
		return
	}
	done <- struct{}{}
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
			h.Stream = res.Stream
			break
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

func (h *hub) RequestsLoop() chan struct{} {
	wait := make(chan struct{})
	go func() {
	loop:
		for {
			select {
			case r := <-h.requests:
				r(h)
				if h.focus {
					h.Render()
				}
			case <-wait:
				break loop
			}
		}
	}()
	return wait
}

func (h *hub) Loop(events <-chan termbox.Event) {
	h.Play()
	h.Render()
	wait := make(chan struct{})
loop:
	for {
		select {
		case evt := <-events:
			h.requests <- func(h *hub) {
				h.handle(events, evt, wait)
			}
			<-wait
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
	// NOTE: very important that this events stream is passed
	// around and not just used in h.Loop since it will consume
	// other keyboard events as well
	events := make(chan termbox.Event)
	go func() {
		for {
			events <- termbox.PollEvent()
		}
	}()
	w := h.RequestsLoop()
	go h.MWLoop()
	h.Loop(events)
	w <- struct{}{}
}
