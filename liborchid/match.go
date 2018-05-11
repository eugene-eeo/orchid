package liborchid

func Match(a, b string) (matched bool, distance int) {
	i := 0
	r := []rune(b)
	n := len(r)
	// s keeps track of how many gaps are in between consecutive characters
	// m keeps trach of when to start reducing s
	s := 0
	m := false
outer:
	for _, c := range a {
		for i != n {
			if c == r[i] {
				m = true
				continue outer
			}
			i++
			if m {
				s++
			}
		}
		return false, s
	}
	return true, s
}
