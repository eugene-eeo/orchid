package main

import "time"
import "math/rand"
import "github.com/nsf/termbox-go"
import "github.com/eugene-eeo/orchid/liborchid"

type updatable interface {
	Update(player *liborchid.Player, paused, shuffle, repeat bool)
}

type request func(*hub)

type hub struct {
	Player   *liborchid.Player
	Stream   *liborchid.Stream
	requests chan request
	view     updatable
	done     bool
	isInfo   bool
	volume   float64
}

func (h *hub) Paused() bool {
	return h.Stream == nil || h.Stream.Paused()
}

func (h *hub) Toggle() {
	if h.Stream != nil {
		h.Stream.Toggle()
	}
}

func newHub(p *liborchid.Player) *hub {
	return &hub{
		Player:   p,
		requests: make(chan request),
		view:     newPlayerView(),
		done:     false,
		isInfo:   false,
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
	stream, err := song.Stream()
	if err != nil {
		h.Stream = nil
		h.Player.Remove()
		go func() {
			h.requests <- func(c *hub) { c.Play() }
		}()
		return
	}
	h.Stream = stream
	stream.Volume().Volume = h.volume
	if h.volume == -4 {
		stream.Volume().Silent = true
	}
	stream.Play()
	go func() {
		if <-stream.Complete() {
			h.requests <- func(c *hub) {
				c.Player.Next(1, false)
				c.Play()
			}
		}
	}()
}

func (h *hub) handle(events <-chan termbox.Event, evt termbox.Event) {
	if evt.Type != termbox.EventKey {
		return
	}
	switch evt.Ch {
	case 'q':
		if h.Stream != nil {
			h.Stream.Stop()
		}
		h.done = true
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
		f.Loop(events)
		song := f.Choice()
		if song != nil {
			h.Player.SetCurrent(song)
			h.Play()
		}
	case 'i':
		h.isInfo = !h.isInfo
		if h.isInfo {
			h.view = newInfoView()
		} else {
			h.view = newPlayerView()
		}
	}
	switch evt.Key {
	case termbox.KeySpace:
		h.Toggle()
	case termbox.KeyArrowLeft:
		fallthrough
	case termbox.KeyArrowRight:
		v := newVolumeUI(h.Stream)
		v.Loop(events)
		h.volume = h.Stream.Volume().Volume
	default:
	}
}

func (h *hub) Loop(events <-chan termbox.Event) {
	h.Play()
	for !h.done {
		h.Render()
		select {
		case evt := <-events:
			h.handle(events, evt)
		case req := <-h.requests:
			req(h)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	songs := liborchid.FindSongs(".")

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
	h.Loop(events)
}
