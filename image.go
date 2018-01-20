package main

import "fmt"
import "github.com/gizak/termui"

type Image struct {
	rawData string
	width   int
	height  int
	X       int
	Y       int
}

func (i *Image) Buffer() termui.Buffer {
	b := termui.NewFilledBuffer(i.X, i.Y, i.X+i.width, i.Y+i.height, ' ', termui.ColorDefault, termui.ColorDefault)
	fmt.Println(b.Bounds())
	for y := 0; y < i.height; y++ {
		for x := 0; x < i.width; x++ {
			b.Set(i.X+x, i.Y+y, termui.NewCell(
				rune(i.rawData[y*i.height+x]),
				termui.ColorDefault,
				termui.ColorDefault,
			))
		}
	}
	return b
}
