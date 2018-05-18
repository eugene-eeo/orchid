package main

import "flag"
import "time"
import "math/rand"
import "github.com/nsf/termbox-go"
import "github.com/eugene-eeo/orchid/liborchid"

const MAX_VOLUME float64 = +0.0
const MIN_VOLUME float64 = -8.0

var REACTOR *Reactor = nil

type request func(*hub)

type hub struct {
	Player   *liborchid.Queue
	MWorker  *liborchid.MWorker
	view     *playerView
	requests chan request
	done     chan struct{}
	progress float64
}

func (h *hub) Paused() bool {
	if stream := h.Stream(); stream != nil {
		return stream.Paused()
	}
	return false
}

func (h *hub) Toggle() {
	if stream := h.Stream(); stream != nil {
		stream.Toggle()
	}
}

func (h *hub) Stream() *liborchid.Stream {
	return h.MWorker.Stream()
}

func newHub(p *liborchid.Queue) *hub {
	return &hub{
		Player:   p,
		MWorker:  liborchid.NewMWorker(),
		view:     newPlayerView(),
		requests: make(chan request),
		done:     make(chan struct{}),
	}
}

func (h *hub) Render() {
	h.view.Update(
		h.Player,
		h.progress,
		h.Paused(),
		h.Player.Shuffle,
		h.Player.Repeat,
	)
}

func (h *hub) Play() {
	song := h.Player.Song()
	if song != nil {
		go func() {
			h.MWorker.SongQueue <- song
		}()
	}
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
		must(termbox.Sync())
		f := newFinderUIFromPlayer(h.Player)
		REACTOR.Focus(f)
		go func() {
			song := <-f.Choice
			h.requests <- func(h *hub) {
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
	case termbox.KeyArrowLeft, termbox.KeyArrowRight:
		must(termbox.Sync())
		v := newVolumeUI(h.MWorker)
		REACTOR.Focus(v)
	}
}

func (h *hub) Handle(e termbox.Event) {
	h.requests <- func(h *hub) { h.handle(e) }
}

func (h *hub) OnFocus() {
	h.requests <- func(h *hub) {}
}

func (h *hub) Loop() {
	for {
		if REACTOR.InFocus(h) {
			h.Render()
		}
		select {
		case r := <-h.requests:
			r(h)
		case res := <-h.MWorker.Results:
			if res.State == liborchid.PlaybackEnd {
				if res.Error != nil {
					h.Player.Remove(res.Song)
				} else {
					h.Player.Next(1, false)
				}
				h.Play()
			}
		case f := <-h.MWorker.Progress:
			h.progress = f
		case <-h.done:
			h.MWorker.Stop()
			return
		}
	}
}

func main() {
	recursive := flag.Bool("r", true, "look recursively for .mp3 files")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	songs := liborchid.FindSongs(".", *recursive)

	must(termbox.Init())
	termbox.SetOutputMode(termbox.Output256)
	defer termbox.Close()

	h := newHub(liborchid.NewQueue(songs))
	REACTOR = NewReactor(h)

	go REACTOR.Loop()
	go h.MWorker.Play()
	h.Play()
	h.Loop()
}
