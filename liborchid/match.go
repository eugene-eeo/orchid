package liborchid

func Match(a, b string) (matched bool, distance int) {
	i := 0
	q := []rune(a)
	r := []rune(b)
	n := len(r)
	// s keeps track of how many gaps are in between consecutive characters
	// m keeps track of when to start counting
	s := 0
	m := false
outer:
	for _, c := range q {
		for i != n {
			if c == r[i] {
				m = true
				i++
				continue outer
			} else if m {
				s++
			}
			i++
		}
		return false, s
	}
	return true, s
}
