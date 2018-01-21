package main

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
