package main

import "fmt"
import "github.com/nsf/termbox-go"
import "github.com/eugene-eeo/orchid/liborchid"
import "github.com/eliukblau/pixterm/ansimage"
import "github.com/dhowden/tag"

const ATTR_DIM = termbox.Attribute(0xf0)
const ATTR_DEFAULT = termbox.ColorDefault

// Layout (50x8)
// ┌────────┐ Prev (-1)
// │        │    Album (Year)
// │  14x8  │ <Play/Pause> Song Title / Name
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

func (pv *playerView) drawCurrent(title string, y int, paused bool, shuffle bool, repeat bool) {
	name := getPlayingIndicator(paused, shuffle) + " " + title
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

	drawName(fmt.Sprintf("%s (%s)", album, year), 18, 1, ATTR_DEFAULT)
	drawName(artist, 18, 3, ATTR_DEFAULT)
	drawName(fmt.Sprintf("Track [%d/%d]", track, total), 18, 4, ATTR_DEFAULT)
}

func (pv *playerView) drawOld(song *liborchid.Song, y int) {
	if song != nil {
		drawName(song.Name(), 16, y, 0xf0)
	}
}

func (pv *playerView) Update(player *liborchid.Player, paused, shuffle, repeat bool) {
	must(termbox.Clear(ATTR_DEFAULT, ATTR_DEFAULT))
	song := player.Song()
	if song != nil && song != pv.current {
		pv.current = song
		pv.metadata = song.Metadata()
		pv.image = getImage(pv.metadata).Render()
	}
	pv.drawOld(player.Peek(-1), 1)
	pv.drawMetaData()
	pv.drawCurrent(getSongTitle(song, pv.metadata), 2, paused, shuffle, repeat)
	pv.drawOld(player.Peek(1), 5)
	pv.drawOld(player.Peek(2), 6)
	must(termbox.Sync())
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
	return "⏵"
}

func getImage(metadata tag.Metadata) (img *ansimage.ANSImage) {
	img = DefaultImage
	// sometimes getting tags raises a panic;
	// no idea why but this is an okay fix since images
	// should not crash the application
	defer recover()
	if metadata == nil {
		return
	}
	p := metadata.Picture()
	if p == nil {
		return
	}
	if rv, err := bytesToImage(p.Data); err == nil {
		return rv
	}
	return
}

func defaultInt(a int) string {
	if a <= 0 {
		return ""
	}
	return fmt.Sprintf("%d", a)
}

func defaultString(a string, b string) string {
	if a == "" {
		return b
	}
	return a
}
