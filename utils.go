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

func unicodeCells(s string, width int, f func(int, rune)) {
	x := 0
	R := []rune(s)
	n := len(R)
	for i := 0; x <= width; i++ {
		r := ' '
		if x == width && n > i {
			// if we are at the final width and string is
			// too long then end with ellipsis
			r = '…'
		} else if i < n {
			r = R[i]
		}
		f(x, r)
		x += runewidth.RuneWidth(r)
	}
}
