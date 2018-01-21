package main

import "bytes"
import "os"
import "github.com/nsf/termbox-go"
import "github.com/eugene-eeo/hubwub/player"
import "github.com/lucasb-eyer/go-colorful"
import "github.com/eliukblau/pixterm/ansimage"

func nextTrack(p *player.Player, i int, force bool, q chan func(*player.Player)) func(*player.Player) {
	var next func(int, bool)
	next = func(i int, f bool) {
		done, err := p.Next(i, f)
		if err != nil {
			p.Remove()
			if _, err := p.Peek(0); err != nil {
				return
			}
			go func() {
				q <- nextTrack(p, 1, true, q)
			}()
		}
		go func() {
			complete := <-done
			if complete {
				q <- nextTrack(p, 1, false, q)
			}
		}()
	}
	return func(*player.Player) {
		next(i, force)
	}
}

func getIndicator(p *player.Player) string {
	if p.Speaker.Paused() {
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

	getImage := func(sng player.Song) image {
		r, ok := sng.Picture()
		if !ok {
			return &defaultImage{}
		}
		bg, _ := colorful.Hex("#000000")
		img, err := ansimage.NewScaledFromReader(bytes.NewReader(r), 16, 16, bg, ansimage.ScaleModeResize, ansimage.NoDithering)
		if err != nil {
			return &defaultImage{}
		}
		return img
	}

	go (func() {
		var currentSong player.Song = player.Song("")
		var img image = &defaultImage{}
		for {
			sng := <-imageQueue
			if sng != currentSong {
				currentSong = sng
				img = getImage(sng)
			}
			termbox.SetCursor(0, 0)
			must(termbox.Sync())
			print(img.Render())
			print("\u001B[?25l")
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

	requests <- nextTrack(app, 0, true, requests)
	go (func() {
		for {
			evt := termbox.PollEvent()
			switch evt.Ch {
			case 'q':
				exit <- struct{}{}
				break
			case 'n':
				requests <- nextTrack(app, 1, true, requests)
			case 'p':
				requests <- nextTrack(app, -1, true, requests)
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
