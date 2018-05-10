package main

import "fmt"
import "github.com/nsf/termbox-go"
import "github.com/eugene-eeo/orchid/liborchid"

type updatable interface {
	Update(player *liborchid.Player, paused bool, shuffle bool, repeat bool)
}

type infoView struct {
	playerView
}

func newInfoView() *infoView {
	return &infoView{}
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

func (i *infoView) Update(player *liborchid.Player, paused bool, shuffle bool, repeat bool) {
	must(termbox.Clear(termbox.ColorDefault, termbox.ColorDefault))
	song := player.Song()
	name := "<No name>"
	if song != nil {
		if meta := song.Metadata(); meta != nil {
			name = defaultString(meta.Title(), song.Name())
			album := defaultString(meta.Album(), "Unknown album")
			year := defaultString(defaultInt(meta.Year()), "?")
			artist := defaultString(meta.Artist(), "Unknown artist")
			track, total := meta.Track()
			drawName(fmt.Sprintf("%s (%s)", album, year), 18, 1, 0xf0)
			drawName(artist, 18, 3, 0xf0)
			drawName(fmt.Sprintf("Track [%d/%d]", track, total), 18, 4, 0xf0)
		}
	}
	i.drawCurrent(name, 2, paused, shuffle, repeat)
	must(termbox.Sync())
	i.drawImage(song)
}
