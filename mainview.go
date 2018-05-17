package main

import "fmt"
import "github.com/nsf/termbox-go"
import "github.com/eugene-eeo/orchid/liborchid"
import "github.com/dhowden/tag"

const ATTR_DIM = termbox.Attribute(0xf0)
const ATTR_DEFAULT = termbox.ColorDefault

// Layout (50x8)
// ┌────────┐
// │        │    Album (Year)
// │  8x16  │ <Play/Pause> Title
// │        │    Artist
// │        │    Track [i/n]
// │        │  Next (1)
// └────────┘  Next (2)
//

type playerView struct {
	current  *liborchid.Song
	image    string
	metadata tag.Metadata
}

func newPlayerView() *playerView {
	return &playerView{
		current: nil,
		image:   DefaultImage.Render(),
	}
}

func (pv *playerView) drawCurrent(y int, paused bool, shuffle bool, repeat bool) {
	name := getPlayingIndicator(paused, shuffle) + " " + getSongTitle(pv.current, pv.metadata)
	attr := termbox.AttrBold
	if repeat {
		attr = termbox.AttrReverse
	}
	drawName(name, 18, y, attr)
}

func (pv *playerView) drawImage() {
	fmt.Print("\033[0;0H" + pv.image + "\u001B[?25l")
}

func (pv *playerView) drawMetaData() {
	meta := pv.metadata
	if meta == nil {
		return
	}
	album := defaultString(meta.Album(), "Unknown album")
	year := defaultString(defaultInt(meta.Year()), "?")
	artist := defaultString(meta.Artist(), "Unknown artist")
	track, total := meta.Track()

	drawName(fmt.Sprintf("%s (%s)", album, year), 20, 1, ATTR_DEFAULT)
	drawName(artist, 20, 3, ATTR_DEFAULT)
	drawName(fmt.Sprintf("Track [%d/%d]", track, total), 20, 4, ATTR_DEFAULT)
}

func (pv *playerView) drawProgress(progress float64) {
	b := int(progress * 32)
	for i := 0; i <= 31; i++ {
		a := ATTR_DIM
		if i <= b {
			a = ATTR_DEFAULT
		}
		termbox.SetCell(18+i, 6, '─', a, ATTR_DEFAULT)
	}
	termbox.SetCell(18+b, 6, '╼', ATTR_DEFAULT, ATTR_DEFAULT)
}

func (pv *playerView) Update(player *liborchid.Queue, progress float64, paused, shuffle, repeat bool) {
	must(termbox.Clear(ATTR_DEFAULT, ATTR_DEFAULT))
	song := player.Song()
	if song != nil && song != pv.current {
		pv.current = song
		pv.metadata = song.Metadata()
		pv.image = getImage(pv.metadata).Render()
	}
	pv.drawMetaData()
	pv.drawCurrent(2, paused, shuffle, repeat)
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
	name := song.Name()
	if metadata != nil {
		name = defaultString(metadata.Title(), name)
	}
	return name
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
