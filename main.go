package main

import "os"
import "github.com/nsf/termbox-go"
import "github.com/lucasb-eyer/go-colorful"
import "github.com/eliukblau/pixterm/ansimage"
import "github.com/mattn/go-runewidth"

func must(err error) {
	if err != nil {
		panic(err)
	}
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
	must(termbox.Init())
	termbox.SetOutputMode(termbox.Output256)
	defer termbox.Close()

	exit := make(chan struct{})

	drawName := func(name string, y int, color termbox.Attribute) {
		unicodeCells(fit(name, 30), func(dx int, r rune) {
			termbox.SetCell(18+dx, y, r, color, termbox.ColorDefault)
		})
	}

	updateQueue := func(state *playerState) {
		for i := 1; i <= 3; i++ {
			drawName(state.Peek(i).Name(), 2+i, 0xf0)
		}
		drawName(state.Peek(-1).Name(), 1, 0xf0)
	}

	imageQueue := make(chan *song, 1)

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
				must(termbox.Sync())
				print(image.Render())
				print("\u001B[?25l")
			}
		}
	})()

	go (func() {
		for {
			state := <-app.State
			must(termbox.Clear(termbox.ColorDefault, termbox.ColorDefault))
			color := termbox.Attribute(0x1ff)
			if state.Repeat {
				color = termbox.AttrReverse
			}
			symbol := "⏵ "
			if state.Paused() {
				symbol = "Ⅱ "
			}
			drawName(symbol+state.Song().Name(), 2, color)
			updateQueue(state)
			must(termbox.Sync())
			imageQueue <- state.Song()
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
			case 'r':
				app.ToggleRepeat()
			}
			if evt.Key == termbox.KeySpace {
				app.TogglePlay()
			}
		}
	})()

	<-exit
}
