package elems

import "github.com/nsf/termbox-go"

type Input struct {
	buf []rune
	idx int
}

func NewInput() *Input {
	return &Input{[]rune{}, 0}
}

func (i *Input) Insert(r rune) {
	i.buf = append(i.buf, r)
	copy(i.buf[i.idx+1:], i.buf[i.idx:])
	i.buf[i.idx] = r
	i.idx++
}

func (i *Input) Delete() {
	if i.idx == 0 || len(i.buf) == 0 {
		return
	}
	i.idx--
	i.buf = append(i.buf[:i.idx], i.buf[i.idx+1:]...)
}

func (i *Input) Feed(key termbox.Key, ch rune, mod termbox.Modifier) {
	if ch != 0 && mod == 0 {
		i.Insert(ch)
		return
	}
	switch key {
	case termbox.KeyArrowLeft:
		if i.idx > 0 {
			i.idx--
		}
	case termbox.KeyArrowRight:
		if i.idx < len(i.buf) {
			i.idx++
		}
	case termbox.KeyBackspace2:
		fallthrough
	case termbox.KeyBackspace:
		i.Delete()
	case termbox.KeySpace:
		i.Insert(' ')
	}
}

func (i *Input) String() string {
	return string(i.buf)
}

func (i *Input) Cursor() int {
	return i.idx
}
