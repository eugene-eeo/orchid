package main

import "time"
import "math/rand"
import "github.com/nsf/termbox-go"
import "github.com/eugene-eeo/orchid/liborchid"

type request func(*hub)

type hub struct {
	Player     *liborchid.Player
	Stream     *liborchid.Stream
	Song       *liborchid.Song
	Requests   chan request
	playerView *playerView
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
	h := &hub{
		Player:     p,
		Stream:     nil,
		Requests:   make(chan request),
		playerView: newPlayerView(),
	}
	return h
}

func (h *hub) Render() {
	h.playerView.Update(
		h.Player,
		h.Paused(),
		h.Player.Shuffle,
		h.Player.Repeat,
	)
}

func (h *hub) Play() {
	h.Song = h.Player.Song()
	if h.Song == nil {
		return
	}
	stream, err := h.Song.Stream()
	if err != nil {
		h.Stream = nil
		h.Player.Remove()
		go func() {
			h.Requests <- func(c *hub) { c.Play() }
		}()
		return
	}
	h.Stream = stream
	stream.Play()
	go func() {
		complete := <-stream.Complete()
		if complete {
			h.Requests <- func(c *hub) {
				c.Player.Next(1, false)
				c.Play()
			}
		}
	}()
}

func (h *hub) Loop() {
	for {
		select {
		case req := <-h.Requests:
			req(h)
			h.Render()
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	songs := liborchid.FindSongs(".")

	must(termbox.Init())
	termbox.SetOutputMode(termbox.Output256)
	defer termbox.Close()

	exit := make(chan struct{})

	h := newHub(liborchid.NewPlayer(songs))
	go h.Loop()
	h.Requests <- func(h *hub) {
		h.Play()
	}
	go (func() {
		for {
			evt := termbox.PollEvent()
			if evt.Type != termbox.EventKey {
				continue
			}
			switch evt.Ch {
			case 'q':
				exit <- struct{}{}
			case 'n':
				h.Requests <- func(h *hub) {
					h.Player.Next(1, true)
					h.Play()
				}
			case 'p':
				h.Requests <- func(h *hub) {
					h.Player.Next(-1, true)
					h.Play()
				}
			case 's':
				h.Requests <- func(h *hub) { h.Player.ToggleShuffle() }
			case 'r':
				h.Requests <- func(h *hub) { h.Player.ToggleRepeat() }
			case 'f':
				hang := make(chan struct{})
				h.Requests <- func(h *hub) {
					f := newFinderUIFromPlayer(h.Player)
					go f.Loop()
					song := <-f.choice
					if song != nil {
						h.Player.SetCurrent(song)
						go func() { h.Requests <- func(h *hub) { h.Play() } }()
					}
					hang <- struct{}{}
				}
				<-hang
			}
			if evt.Key == termbox.KeySpace {
				h.Requests <- func(h *hub) {
					h.Toggle()
				}
			}
		}
	})()

	<-exit
}
