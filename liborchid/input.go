package liborchid

import "github.com/nsf/termbox-go"

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

// Input represents an editable rune buffer with a cursor.
type Input struct {
	buf []rune
	idx int
}

// NewInput creates a new empty Input.
func NewInput() *Input {
	return &Input{[]rune{}, 0}
}

// Insert inserts the given rune at the current cursor position.
func (i *Input) Insert(r rune) {
	i.buf = append(i.buf, r)
	copy(i.buf[i.idx+1:], i.buf[i.idx:])
	i.buf[i.idx] = r
	i.idx++
}

// Delete deletes a rune at the current cursor if possible, else does nothing.
func (i *Input) Delete() {
	if i.idx == 0 || len(i.buf) == 0 {
		return
	}
	i.idx--
	i.buf = append(i.buf[:i.idx], i.buf[i.idx+1:]...)
}

// Move increments the cursor by n (n can be negative). If moving the cursor
// brings it to an invalid position (< 0 or > length of buffer) then it will
// automatically be corrected.
func (i *Input) Move(n int) {
	i.idx = min(max(i.idx+n, 0), len(i.buf))
}

// Feed is a convenience function that interprets the given key, rune, and
// modifier mask and calls the appropriate methods, e.g. if key == left arrow
// then Move(-1) is called.
func (i *Input) Feed(key termbox.Key, ch rune, mod termbox.Modifier) {
	if ch != 0 && mod == 0 {
		i.Insert(ch)
		return
	}
	switch key {
	case termbox.KeyArrowLeft:
		i.Move(-1)
	case termbox.KeyArrowRight:
		i.Move(1)
	case termbox.KeyBackspace2:
		fallthrough
	case termbox.KeyBackspace:
		i.Delete()
	case termbox.KeySpace:
		i.Insert(' ')
	}
}

// String returns the internal buffer as a string.
func (i *Input) String() string {
	return string(i.buf)
}

// Cursor returns the current cursor position. Note that the cursor position
// is not the length of the buffer, and is in the interval [0, len(buffer)].
func (i *Input) Cursor() int {
	return i.idx
}
