package main

type Input struct {
	idx int
	buf []rune
}

func newInput() *Input {
	return &Input{idx: 0, buf: []rune{}}
}
