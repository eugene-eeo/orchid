package main

import (
	"fmt"

	"github.com/dhowden/tag"
	"github.com/eugene-eeo/orchid/liborchid"
	"github.com/nsf/termbox-go"
)

const ATTR_DIM = termbox.Attribute(0xf0)
const ATTR_DEFAULT = termbox.ColorDefault

// Layout (50x8)
// ┌────────┐
// │        │    Album (Year)
// │  8x16  │ <Play/Pause> Title
// │        │    Artist
// │        │    Track [i/n]
// │        │
// └────────┘ Progress
//

type playerView struct {
	current  *liborchid.Song
	image    string
	metadata tag.Metadata
}

func newPlayerView() *playerView {
	return &playerView{
		current: nil,
		image:   DefaultImage,
	}
}

func (pv *playerView) drawCurrent(paused, shuffle, repeat bool) {
	name := getPlayingIndicator(paused, shuffle) + " " + getSongTitle(pv.current, pv.metadata)
	attr := termbox.AttrBold
	if repeat {
		attr = termbox.AttrReverse
	}
	drawName(name, 18, 2, attr)
}

func (pv *playerView) drawImage() {
	fmt.Print("\033[0;0H" + pv.image + "\u001B[?25l")
}

func (pv *playerView) drawMetaData() {
	m := pv.metadata
	if m == nil {
		return
	}
	album := defaultString(m.Album(), "Unknown album")
	year := defaultString(defaultInt(m.Year()), "?")
	artist := defaultString(m.Artist(), "Unknown artist")
	track, total := m.Track()

	drawName(fmt.Sprintf("%s (%s)", album, year), 20, 1, ATTR_DEFAULT)
	drawName(artist, 20, 3, ATTR_DEFAULT)
	drawName(fmt.Sprintf("Track [%d/%d]", track, total), 20, 4, ATTR_DEFAULT)
}

func (pv *playerView) drawProgress(progress float64) {
	b := int(progress*31) + 18
	for x := 18; x < 50; x++ {
		a := ATTR_DIM
		if x <= b {
			a = ATTR_DEFAULT
		}
		termbox.SetCell(x, 6, '─', a, ATTR_DEFAULT)
	}
	termbox.SetCell(b, 6, '╼', ATTR_DEFAULT, ATTR_DEFAULT)
}

func (pv *playerView) Update(player *liborchid.Queue, progress float64, paused, shuffle, repeat bool) {
	must(termbox.Clear(ATTR_DEFAULT, ATTR_DEFAULT))
	song := player.Song()
	if song != nil && song != pv.current {
		pv.current = song
		pv.metadata = song.Metadata()
		pv.image = getImage(pv.metadata)
	}
	pv.drawMetaData()
	pv.drawCurrent(paused, shuffle, repeat)
	pv.drawProgress(progress)
	must(termbox.Flush())
	pv.drawImage()
}

func drawName(name string, x int, y int, color termbox.Attribute) {
	unicodeCells(name, 49-x, true, func(dx int, r rune) {
		termbox.SetCell(x+dx, y, r, color, ATTR_DEFAULT)
	})
}

func getSongTitle(song *liborchid.Song, metadata tag.Metadata) string {
	if song == nil {
		return "<No Songs>"
	}
	if metadata != nil {
		return defaultString(metadata.Title(), song.Name())
	}
	return song.Name()
}

func getPlayingIndicator(paused, shuffle bool) string {
	if paused {
		return "Ⅱ"
	}
	if shuffle {
		return "⥮"
	}
	return "▶"
}
