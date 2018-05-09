package main

import "time"
import "math/rand"
import "fmt"
import "github.com/nsf/termbox-go"
import "github.com/eugene-eeo/orchid/liborchid"
import "github.com/lucasb-eyer/go-colorful"
import "github.com/eliukblau/pixterm/ansimage"
import "bytes"

type request func(*hub)

type hub struct {
	Player   *liborchid.Player
	Stream   *liborchid.Stream
	Song     *liborchid.Song
	Requests chan request
	rendered *liborchid.Song
	image    image
}

func (h *hub) Paused() bool {
	if h.Stream == nil {
		return true
	}
	return h.Stream.Paused()
}

func (h *hub) Toggle() {
	if h.Stream != nil {
		h.Stream.Toggle()
	}
}

func newHub(p *liborchid.Player) *hub {
	h := &hub{
		Player:   p,
		Stream:   nil,
		Requests: make(chan request),
		image:    &defaultImage{},
	}
	return h
}

func (h *hub) Render() {
	must(termbox.Clear(termbox.ColorDefault, termbox.ColorDefault))
	s := h.Player.Peek(-1)
	if s == nil {
		return
	}
	currentSong := h.Player.Song()
	name := "<No songs>"
	if currentSong != nil {
		name = currentSong.Name()
	}
	color := termbox.ColorDefault
	if h.Player.Repeat {
		color = termbox.AttrReverse
	}
	drawName(s.Name(), 1, 0xf0)
	drawName(string(getIndicator(h))+" "+name, 2, color)
	for i := 1; i <= 3; i++ {
		s := h.Player.Peek(i)
		if s == nil {
			break
		}
		drawName(s.Name(), 2+i, 0xf0)
	}
	must(termbox.Sync())
	if h.rendered != currentSong {
		h.image = getImage(currentSong)
	}
	drawImage(h.image)
}

func (h *hub) Play() {
	h.Song = h.Player.Peek(0)
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

func getIndicator(h *hub) rune {
	if h.Paused() {
		return 'Ⅱ'
	}
	if h.Player.Shuffle {
		return '⥮'
	}
	return '⏵'
}

func drawName(name string, y int, color termbox.Attribute) {
	unicodeCells(name, 30, true, func(dx int, r rune) {
		termbox.SetCell(18+dx, y, r, color, termbox.ColorDefault)
	})
}

func getImage(song *liborchid.Song) (img image) {
	img = &defaultImage{}
	defer func() {
		// sometimes getting tags raises a panic;
		// no idea why but this is an okay fix since images
		// should not crash the application
		if r := recover(); r != nil {
		}
	}()
	if song == nil {
		return
	}
	p := song.Image()
	if p == nil {
		return
	}
	rv, err := ansimage.NewScaledFromReader(
		bytes.NewReader(p.Data),
		16, 16,
		colorful.LinearRgb(0, 0, 0),
		ansimage.ScaleModeResize,
		ansimage.NoDithering,
	)
	if err == nil {
		return rv
	}
	return
}

func drawImage(img image) {
	if img == nil {
		img = &defaultImage{}
	}
	termbox.SetCursor(0, 0)
	must(termbox.Sync())
	fmt.Print(img.Render())
	fmt.Print("\u001B[?25l")
}

/*
	/-------\
	| 16x16 | <Prev>
	|       | <Now Playing>
	|       | <Up Next>
	\-------/
*/

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
