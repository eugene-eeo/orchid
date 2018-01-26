package main

import "github.com/mattn/go-runewidth"

const DefaultImage string = "\u001B[38;5;147m" + `        _
    _ (` + " - " + `) _
  /` + "` '.\\ /.' `" + `\
  ` + "``" + `'-.,=,.-'` + "``" + `
    .'//v\\'.
   (_/\ " /\_)
       '-'`

type image interface {
	Render() string
}

type defaultImage struct{}

func (d *defaultImage) Render() string {
	return DefaultImage
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func unicodeCells(s string, width int, fill bool, f func(int, rune)) {
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
		} else if !fill {
			break
		}
		f(x, r)
		x += runewidth.RuneWidth(r)
	}
}
