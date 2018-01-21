package main

import "fmt"
import "bytes"
import "os"
import "github.com/nsf/termbox-go"
import "github.com/eugene-eeo/hubwub/player"
import "github.com/lucasb-eyer/go-colorful"
import "github.com/eliukblau/pixterm/ansimage"

func nextTrack(p *player.Player, i int, force bool, q chan func(*player.Player)) {
	d, err := p.Next(i, force)
	if err != nil {
		fmt.Println(err)
		p.Remove()
		go func() {
			if _, err := p.Peek(1); err != nil {
				return
			}
			q <- func(p *player.Player) {
				nextTrack(p, 1, true, q)
			}
		}()
		return
	}
	go (func() {
		graceful := <-d
		if !graceful {
			return
		}
		q <- func(p *player.Player) {
			nextTrack(p, 1, false, q)
		}
	})()
}

func getIndicator(p *player.Player) string {
	if p.Stream != nil && p.Stream.Paused() {
		return "Ⅱ"
	}
	if p.Shuffle {
		return "⥮"
	}
	return "⏵"
}

func drawName(name string, y int, color termbox.Attribute) {
	unicodeCells(fit(name, 30), func(dx int, r rune) {
		termbox.SetCell(18+dx, y, r, color, termbox.ColorDefault)
	})
}

func updateQueue(app *player.Player) {
	for i := 1; i <= 3; i++ {
		s, err := app.Peek(i)
		if err != nil {
			break
		}
		drawName(s.Name(), 2+i, 0xf0)
	}
	s, err := app.Peek(-1)
	if err != nil {
		return
	}
	drawName(s.Name(), 1, 0xf0)
}

/*
	+-------+
	|       | <Prev>
	| 16x16 | <Now Playing>
	|       | <Up Next>
	+-------+
*/

func main() {
	songs, err := player.FindSongs(".")
	if err != nil {
		os.Exit(1)
	}
	app := player.NewPlayer(songs)

	must(termbox.Init())
	termbox.SetOutputMode(termbox.Output256)
	defer termbox.Close()

	exit := make(chan struct{})
	requests := make(chan func(*player.Player))

	imageQueue := make(chan player.Song, 1)

	go (func() {
		var currentSong player.Song = player.Song("")
		var image *ansimage.ANSImage = nil
		for {
			sng := <-imageQueue
			if sng != currentSong {
				currentSong = sng
				r, ok := sng.Picture()
				if !ok {
					image = nil
					continue
				}
				bg, _ := colorful.Hex("#000000")
				image, err = ansimage.NewScaledFromReader(bytes.NewReader(r), 16, 16, bg, ansimage.ScaleModeResize, ansimage.NoDithering)
				if err != nil {
					image = nil
					continue
				}
			}
			if image != nil {
				termbox.SetCursor(0, 0)
				must(termbox.Sync())
				print(image.Render())
				print("\u001B[?25l")
			}
		}
	})()

	render := func(app *player.Player) {
		must(termbox.Clear(termbox.ColorDefault, termbox.ColorDefault))
		color := termbox.Attribute(0x1ff)
		if app.Repeat {
			color = termbox.AttrReverse
		}
		s, err := app.Song()
		name := "<No songs>"
		if err == nil {
			name = s.Name()
		}
		drawName(getIndicator(app)+" "+name, 2, color)
		updateQueue(app)
		must(termbox.Sync())
		if err == nil {
			imageQueue <- s
		}
	}

	go (func() {
		for {
			r := <-requests
			r(app)
			render(app)
		}
	})()

	if len(app.Queue) > 0 {
		requests <- func(app *player.Player) { nextTrack(app, 0, true, requests) }
	}
	go (func() {
		for {
			evt := termbox.PollEvent()
			switch evt.Ch {
			case 'q':
				exit <- struct{}{}
				break
			case 'n':
				requests <- func(h *player.Player) { nextTrack(app, 1, true, requests) }
			case 'p':
				requests <- func(h *player.Player) { nextTrack(app, -1, true, requests) }
			case 's':
				requests <- func(h *player.Player) { h.ToggleShuffle() }
			case 'r':
				requests <- func(h *player.Player) { h.ToggleRepeat() }
			}
			if evt.Key == termbox.KeySpace {
				requests <- func(h *player.Player) { h.Toggle() }
			}
		}
	})()

	<-exit
}
