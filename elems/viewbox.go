package elems

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

type Viewbox struct {
	lo     int
	hi     int
	Max    int
	Height int
}

func NewViewBox(max, height int) *Viewbox {
	return &Viewbox{
		lo:     0,
		hi:     min(max, height),
		Max:    max,
		Height: height,
	}
}

func (v *Viewbox) Reset() {
	v.lo = 0
	v.hi = min(v.Max, v.Height)
}

func (v *Viewbox) Update(i int) (int, int) {
	if v.lo == v.hi {
		v.lo = 0
		v.hi = min(v.Max, v.Height)
	} else if i < v.lo {
		v.lo = i
		v.hi = min(v.Max, i+v.Height)
	} else if i >= v.hi {
		v.lo = i - v.Height + 1
		v.hi = i + 1
	}
	return v.lo, v.hi
}

func (v *Viewbox) Lo() int {
	return v.lo
}

func (v *Viewbox) Hi() int {
	return v.hi
}
