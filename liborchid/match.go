package liborchid

func Match(a, b string) bool {
	i := 0
	r := []rune(b)
	n := len(r)
outer:
	for _, c := range a {
		for i != n {
			if c == r[i] {
				continue outer
			}
			i++
		}
		return false
	}
	return true
}
