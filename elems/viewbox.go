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

// Viewbox represents the bounds [a,b) for a scrollable list with height h and
// maximum index m. It can be used to find the bounds for a scrollable list as
// the user is scrolling through it.
type Viewbox struct {
	lo     int
	hi     int
	max    int
	height int
}

// NewViewbox returns a new Viewbox.
func NewViewBox(max, height int) *Viewbox {
	return &Viewbox{
		lo:     0,
		hi:     min(max, height),
		max:    max,
		height: height,
	}
}

// Update updates the bounds so that i fits in [a',b'), taking into account the
// maximum value and the height. The new bounds are returned.
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

// Lo returns the lower bound.
func (v *Viewbox) Lo() int {
	return v.lo
}

// Hi returns the upper bound.
func (v *Viewbox) Hi() int {
	return v.hi
}
