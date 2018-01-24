package main

import "time"
import "fmt"
import "bytes"
import "os"
import "github.com/nsf/termbox-go"
import "github.com/eugene-eeo/orchid/player"
import "github.com/lucasb-eyer/go-colorful"
import "github.com/eliukblau/pixterm/ansimage"

type request func(*player.Player)

func nextTrack(i int, force bool, q chan request) request {
	return request(func(p *player.Player) {
		_, err := p.Next(i, force)
		// check if there is a next song so that we don't
		// loop infinitely with nothing, since after each
		// request we make a render call
		if err != nil {
			return
		}
		done, err := play(p)
		if err != nil {
			p.Remove()
			go func() { q <- nextTrack(0, true, q) }()
			return
		}
		go (func() {
			if <-done {
				q <- nextTrack(1, false, q)
			}
		})()
	})
}

func getIndicator(p *player.Player) rune {
	if p.Speaker.Paused() {
		return 'Ⅱ'
	}
	if p.Shuffle {
		return '⥮'
	}
	return '⏵'
}

func getImage(sng player.Song) image {
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

func drawName(name string, y int, color termbox.Attribute) {
	unicodeCells(name, 30, true, func(dx int, r rune) {
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
	requests := make(chan request)
	imageQueue := make(chan player.Song, 1)

	go (func() {
		var currentSong player.Song = player.Song("")
		var img image = &defaultImage{}
		for {
			sng := <-imageQueue
			if sng != currentSong {
				currentSong = sng
				img = getImage(sng)
			}
			fmt.Print("\033[0;0H")
			fmt.Print(img.Render())
			fmt.Print("\u001B[?25l")
		}
	})()

	render := func(app *player.Player) {
		must(termbox.Clear(termbox.ColorDefault, termbox.ColorDefault))
		color := termbox.ColorDefault
		if app.Repeat {
			color = termbox.AttrReverse
		}
		s, err := app.Song()
		name := "<No songs>"
		if err == nil {
			name = s.Name()
		}
		drawName(string(getIndicator(app))+" "+name, 2, color)
		updateQueue(app)
		must(termbox.Sync())
		if err == nil {
			imageQueue <- s
		}
	}

	go (func() {
		t := 5 * time.Millisecond
		after := time.AfterFunc(t, func() {
			render(app)
		})
		for {
			select {
			case req := <-requests:
				req(app)
				after.Reset(t)
			}
		}
	})()

	requests <- nextTrack(0, true, requests)
	go (func() {
		for {
			evt := termbox.PollEvent()
			if evt.Type != termbox.EventKey {
				continue
			}
			switch evt.Ch {
			case 'q':
				exit <- struct{}{}
				break
			case 'n':
				requests <- nextTrack(1, true, requests)
			case 'p':
				requests <- nextTrack(-1, true, requests)
			case 's':
				requests <- func(h *player.Player) { h.ToggleShuffle() }
			case 'r':
				requests <- func(h *player.Player) { h.ToggleRepeat() }
			case 'f':
				hang := make(chan struct{})
				requests <- func(h *player.Player) {
					f := newFinderUIFromPlayer(h)
					go f.Loop()
					song := <-f.choice
					if song != nil {
						h.SetCurrent(*song)
						go func() { requests <- nextTrack(0, true, requests) }()
					}
					hang <- struct{}{}
				}
				<-hang
			}
			if evt.Key == termbox.KeySpace {
				requests <- func(h *player.Player) { h.Toggle() }
			}
		}
	})()

	<-exit
}
