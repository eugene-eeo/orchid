package main

import "fmt"
import "github.com/nsf/termbox-go"
import "github.com/eugene-eeo/orchid/liborchid"
import "github.com/eliukblau/pixterm/ansimage"
import "github.com/dhowden/tag"

// Layout (50x8)
// ┌────────┐
// │        │ Album (Year)
// │ 16x16  │ <Play/Pause> Song Title / Name
// │        │ Artist
// │        │ Track [i/n]
// │        │
// └────────┘
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
	fmt.Print("\u001B[0;0H" + pv.image + "\u001B[?25l")
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
	drawName(fmt.Sprintf("%s (%s)", album, year), 18, 1, 0xf0)
	drawName(artist, 18, 3, 0xf0)
	drawName(fmt.Sprintf("Track [%d/%d]", track, total), 18, 4, 0xf0)
}

func (pv *playerView) Update(player *liborchid.Player, paused, shuffle, repeat bool) {
	must(termbox.Clear(termbox.ColorDefault, termbox.ColorDefault))
	song := player.Song()
	name := "<No name>"
	if song != nil {
		if song != pv.current {
			pv.current = song
			pv.metadata = song.Metadata()
			pv.image = getImage(pv.metadata).Render()
		}
	}
	if pv.metadata != nil {
		name = pv.metadata.Title()
	}
	pv.drawMetaData()
	pv.drawCurrent(defaultString(name, song.Name()), 2, paused, shuffle, repeat)
	must(termbox.Sync())
	pv.drawImage()
}

func drawName(name string, x int, y int, color termbox.Attribute) {
	unicodeCells(name, 30, true, func(dx int, r rune) {
		termbox.SetCell(x+dx, y, r, color, termbox.ColorDefault)
	})
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
