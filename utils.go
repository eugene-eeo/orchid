package main

import "fmt"
import "github.com/mattn/go-runewidth"
import "github.com/gobuffalo/packr"
import "github.com/eugene-eeo/orchid/ansimage"
import "github.com/dhowden/tag"
import "image/color"
import "bytes"

var DefaultImage string = ""

func init() {
	img, err := bytesToImage(packr.NewBox("./assets").Bytes("default.png"))
	must(err)
	DefaultImage = img.Render()
}

func bytesToImage(data []byte) (*ansimage.ANSImage, error) {
	return ansimage.NewScaledFromReader(
		bytes.NewReader(data),
		16, 16,
		color.RGBA{0, 0, 0, 0xff},
		ansimage.ScaleModeResize,
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

func getImage(metadata tag.Metadata) (img string) {
	img = DefaultImage
	if metadata == nil {
		return
	}
	p := metadata.Picture()
	if p == nil {
		return
	}
	if rv, err := bytesToImage(p.Data); err == nil {
		return rv.Render()
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
