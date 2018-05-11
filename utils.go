package main

import "github.com/mattn/go-runewidth"
import "github.com/gobuffalo/packr"
import "github.com/eliukblau/pixterm/ansimage"
import "github.com/lucasb-eyer/go-colorful"
import "bytes"

var DefaultImage *ansimage.ANSImage = nil

func init() {
	img, err := bytesToImage(packr.NewBox("./assets").Bytes("default.png"))
	must(err)
	DefaultImage = img
}

func bytesToImage(data []byte) (*ansimage.ANSImage, error) {
	return ansimage.NewScaledFromReader(
		bytes.NewReader(data),
		16, 14,
		colorful.LinearRgb(0, 0, 0),
		ansimage.ScaleModeResize,
		ansimage.NoDithering,
	)
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
