package main

import "fmt"
import "github.com/nsf/termbox-go"
import "github.com/eugene-eeo/orchid/liborchid"

// Layout (50x8)
// ┌────────┐
// │        │ Previous Song
// │ 16x16  │ <Play/Pause> Current Song
// │        │ Next 3 songs
// │        │ ...
// │        │ ...
// └────────┘
//

type playerView struct {
	rendered *liborchid.Song
	image    image
}

func newPlayerView() *playerView {
	return &playerView{
		rendered: nil,
		image:    DefaultImage,
	}
}

func (pv *playerView) drawOld(song *liborchid.Song, y int) {
	if song != nil {
		drawName(song.Name(), 18, y, 0xf0)
	}
}

func (pv *playerView) drawCurrent(song *liborchid.Song, y int, paused bool, shuffle bool, repeat bool) {
	name := "<No Songs>"
	if song != nil {
		name = string(getPlayingIndicator(paused, shuffle)) + " " + song.Name()
	}
	attr := termbox.AttrBold
	if repeat {
		attr = termbox.AttrReverse
	}
	drawName(name, 18, y, attr)
}

func (pv *playerView) drawImage(song *liborchid.Song) {
	if song != pv.rendered {
		pv.rendered = song
		pv.image = getImage(song)
	}
	termbox.SetCursor(0, 0)
	must(termbox.Flush())
	fmt.Print(pv.image.Render() + "\u001B[?25l")
}

func (pv *playerView) Update(player *liborchid.Player, paused bool, shuffle bool, repeat bool) {
	must(termbox.Clear(termbox.ColorDefault, termbox.ColorDefault))
	pv.drawOld(player.Peek(-1), 1)
	pv.drawCurrent(player.Song(), 2, paused, shuffle, repeat)
	// can be encapsulated into a loop, but meh.
	pv.drawOld(player.Peek(1), 3)
	pv.drawOld(player.Peek(2), 4)
	pv.drawOld(player.Peek(3), 5)
	must(termbox.Sync())
	pv.drawImage(player.Song())
}

func drawName(name string, x int, y int, color termbox.Attribute) {
	unicodeCells(name, 30, true, func(dx int, r rune) {
		termbox.SetCell(x+dx, y, r, color, termbox.ColorDefault)
	})
}

func getPlayingIndicator(paused bool, shuffle bool) rune {
	if paused {
		return 'Ⅱ'
	}
	if shuffle {
		return '⥮'
	}
	return '⏵'
}

func getImage(song *liborchid.Song) (img image) {
	img = DefaultImage
	// sometimes getting tags raises a panic;
	// no idea why but this is an okay fix since images
	// should not crash the application
	defer recover()
	if song == nil {
		return
	}
	p := song.Image()
	if p == nil {
		return
	}
	if rv, err := bytesToImage(p.Data); err == nil {
		return rv
	}
	return
}
