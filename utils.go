package main

import "github.com/mattn/go-runewidth"
import "github.com/gobuffalo/packr"
import "github.com/eliukblau/pixterm/ansimage"
import "github.com/lucasb-eyer/go-colorful"
import "bytes"

var DefaultImage image = nil

func init() {
	var err error
	DefaultImage, err = ansimage.NewScaledFromReader(
		bytes.NewReader(packr.NewBox("./assets").Bytes("default.png")),
		16, 16,
		colorful.LinearRgb(0, 0, 0),
		ansimage.ScaleModeResize,
		ansimage.NoDithering,
	)
	must(err)
}

type image interface {
	Render() string
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
