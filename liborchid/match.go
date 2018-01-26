package liborchid

func Match(a, b string) bool {
	i := 0
	n := len(b)
outer:
	for _, c := range a {
		for i != n {
			if c == rune(b[i]) {
				continue outer
			}
			i++
		}
		return false
	}
	return true
}
