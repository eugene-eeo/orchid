package elems

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type Viewbox struct {
	lo     int
	hi     int
	max    int
	height int
}

func NewViewBox(max, height int) *Viewbox {
	return &Viewbox{
		lo:     0,
		hi:     min(max, height),
		max:    max,
		height: height,
	}
}

func (v *Viewbox) Update(i int) (int, int) {
	if v.lo < v.hi {
		if i < v.lo {
			v.lo = i
			v.hi = min(v.max, i+v.height)
		} else if i >= v.hi {
			v.lo = i - v.height + 1
			v.hi = i + 1
		}
	}
	return v.lo, v.hi
}

func (v *Viewbox) Lo() int {
	return v.lo
}

func (v *Viewbox) Hi() int {
	return v.hi
}
