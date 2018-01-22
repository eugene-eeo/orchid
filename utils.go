package main

import "github.com/mattn/go-runewidth"
import "github.com/eugene-eeo/orchid/player"

func play(p *player.Player) (<-chan bool, error) {
	song, err := p.Song()
	if err != nil {
		return nil, err
	}
	done, err := p.Speaker.Play(song)
	if err != nil {
		return nil, err
	}
	return done, nil
}

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
