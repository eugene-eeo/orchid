package main

import "github.com/nsf/termbox-go"

type Input struct {
	buf string
}

func newInput() *Input {
	return &Input{buf: ""}
}

func (i *Input) Feed(ev termbox.Event) {
	if ev.Type != termbox.EventKey {
		return
	}
	switch ev.Key {
	case termbox.KeyBackspace2:
		fallthrough
	case termbox.KeyBackspace:
		if len(i.buf) == 0 {
			return
		}
		i.buf = i.buf[:len(i.buf)-1]
	case termbox.KeySpace:
		i.buf += " "
	default:
		i.buf += string(ev.Ch)
	}
}

func (i *Input) String() string {
	return i.buf
}
