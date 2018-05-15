package liborchid

// Match takes two strings a, b and returns whether they
// match (fuzzy matching) and the distance between
// a and b. If match is false then the distance returned
// should not be trusted.
func Match(a, b string) (matched bool, distance int) {
	i := 0
	q := []rune(a)
	r := []rune(b)
	n := len(r)
	// s = # of gaps in between consecutive characters
	// m = when to start counting
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
