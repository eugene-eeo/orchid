package main

import "os"
import "github.com/nsf/termbox-go"
import "github.com/lucasb-eyer/go-colorful"
import "github.com/eliukblau/pixterm/ansimage"
import "github.com/faiface/beep/speaker"
import "github.com/mattn/go-runewidth"

const (
	RepeatPlaylist = 0
	RepeatSong     = 1
)

func repeatSymbol(r int) rune {
	switch r {
	case RepeatPlaylist:
		return '⚬'
	case RepeatSong:
		return '∞'
	}
	return ' '
}

func songNames(app *state, list []*song) []string {
	s := make([]string, len(list))
	for i, song := range list {
		s[i] = app.NameOf(song)
	}
	return s
}

func formatNowPlaying(a string) string {
	if len(a) > 30 {
		return a[:29] + "…"
	}
	for len(a) < 30 {
		a = a + " "
	}
	return a
}

func main() {
	speaker.Init(48000, 4800)
	app, err := newState(".")
	if err != nil {
		os.Exit(1)
	}
	/*
		+---+ <Now Playing>
		|   | <Up Next>
		+---+
	*/
	termbox.Init()
	defer termbox.Close()

	repeatMode := 0
	exit := make(chan struct{})

	updateRepeat := func() {
		termbox.SetCell(18+31, 1, rune(repeatSymbol(repeatMode)), termbox.ColorGreen, termbox.ColorDefault)
	}

	updateQueue := func() {
		for i := 1; i <= 4; i++ {
			s := formatNowPlaying(app.NameOf(app.Peek(i)))
			x := 0
			for _, c := range formatNowPlaying(s) {
				termbox.SetCell(18+x, 1+i, rune(c), termbox.ColorDefault, termbox.ColorDefault)
				x += runewidth.RuneWidth(rune(c))
			}
		}
	}

	go (func() {
		for {
			sng := <-app.nowPlaying
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			x := 0
			for _, c := range formatNowPlaying(app.NameOf(sng)) {
				termbox.SetCell(18+x, 1, rune(c), termbox.AttrReverse, termbox.ColorDefault)
				x += runewidth.RuneWidth(rune(c))
			}
			updateRepeat()
			updateQueue()
			termbox.Sync()
			r, ok := sng.Picture()
			if !ok {
				continue
			}
			bg, _ := colorful.Hex("#000000")
			img, err := ansimage.NewScaledFromReader(r, 16, 16, bg, ansimage.ScaleModeResize, ansimage.NoDithering)
			if err != nil {
				continue
			}
			termbox.SetCursor(0, 0)
			termbox.Sync()
			print(img.Render())
			print("\u001B[?25l")
			termbox.HideCursor()
		}
	})()

	go app.Loop()
	if len(app.queue) > 0 {
		app.Next(0)
	}

	go (func() {
		for {
			evt := termbox.PollEvent()
			switch evt.Ch {
			case 'q':
				exit <- struct{}{}
				break
			case 'n':
				app.Next(1)
			case 'p':
				app.Next(-1)
			case 's':
				app.Shuffle()
				app.Next(0)
			case 'r':
				repeatMode = mod(repeatMode+1, 2)
				switch repeatMode {
				case RepeatPlaylist:
					app.queue = app.songs
					app.Next(0)
					break
				case RepeatSong:
					s := app.currentSong()
					app.queue = []*song{s}
					app.Next(0)
				}
			}
			if evt.Key == termbox.KeySpace {
				app.TogglePlay()
			}
		}
	})()

	<-exit
}
