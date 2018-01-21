package main

import "os"
import "github.com/nsf/termbox-go"
import "github.com/lucasb-eyer/go-colorful"
import "github.com/eliukblau/pixterm/ansimage"
import "github.com/mattn/go-runewidth"

type updateRequest struct {
	song   *song
	paused bool
}

func fit(a string, width int) string {
	if runewidth.StringWidth(a) > width {
		return a[:29] + "…"
	}
	for runewidth.StringWidth(a) < width {
		a = a + " "
	}
	return a
}

func unicodeCells(s string, f func(int, rune)) {
	x := 0
	for _, c := range s {
		r := rune(c)
		f(x, r)
		x += runewidth.RuneWidth(r)
	}
}

func main() {
	app, err := newState(".")
	if err != nil {
		os.Exit(1)
	}
	/*
		+-------+
		|       | <Prev>
		| 16x16 | <Now Playing>
		|       | <Up Next>
		+-------+
	*/
	termbox.Init()
	termbox.SetOutputMode(termbox.Output256)
	defer termbox.Close()

	exit := make(chan struct{})

	drawName := func(name string, y int, color termbox.Attribute) {
		unicodeCells(fit(name, 30), func(dx int, r rune) {
			termbox.SetCell(18+dx, y, r, color, termbox.ColorDefault)
		})
	}

	updateQueue := func() {
		for i := 1; i <= 3; i++ {
			drawName(app.Peek(i).Name(), 2+i, 0xf0)
		}
		drawName(app.Peek(-1).Name(), 1, 0xf0)
	}

	imageQueue := make(chan *song, 1)
	uiUpdateQueue := make(chan updateRequest)

	go (func() {
		var currentSong *song = nil
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
				image, err = ansimage.NewScaledFromReader(r, 16, 16, bg, ansimage.ScaleModeResize, ansimage.NoDithering)
				if err != nil {
					image = nil
					continue
				}
			}
			if image != nil {
				termbox.SetCursor(0, 0)
				termbox.Sync()
				print(image.Render())
				print("\u001B[?25l")
			}
		}
	})()

	go (func() {
		for {
			req := <-uiUpdateQueue
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			color := termbox.Attribute(0x1ff)
			if app.repeat {
				color = termbox.AttrReverse
			}
			symbol := "⏵ "
			if req.paused {
				symbol = "Ⅱ "
			}
			drawName(symbol+req.song.Name(), 2, color)
			updateQueue()
			termbox.Sync()
			imageQueue <- req.song
		}
	})()

	go (func() {
		var currentSong *song = nil
		for {
			select {
			case p := <-app.Paused:
				uiUpdateQueue <- updateRequest{
					paused: p,
					song:   currentSong,
				}
			case currentSong = <-app.NowPlaying:
				uiUpdateQueue <- updateRequest{
					paused: false,
					song:   currentSong,
				}
			}
		}
	})()

	go app.Loop()
	if len(app.songs) > 0 {
		app.Next(0, true)
	}

	go (func() {
		for {
			evt := termbox.PollEvent()
			switch evt.Ch {
			case 'q':
				exit <- struct{}{}
				break
			case 'n':
				app.Next(1, true)
			case 'p':
				app.Next(-1, true)
			case 's':
				app.Shuffle()
				app.NowPlaying <- app.currentSong()
			case 'r':
				app.ToggleRepeat()
				app.NowPlaying <- app.currentSong()
			}
			if evt.Key == termbox.KeySpace {
				app.TogglePlay()
			}
		}
	})()

	<-exit
}
