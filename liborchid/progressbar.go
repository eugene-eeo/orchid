package liborchid

import "fmt"

type ProgressBar struct {
	maxWidth int
	symbol   rune
}

func NewProgressBar(maxWidth int, symbol rune) *ProgressBar {
	return &ProgressBar{
		maxWidth: maxWidth,
		symbol:   symbol,
	}
}

func (pg *ProgressBar) Update(f float64) string {
	percentage := fmt.Sprintf("%d%%", int(f*100))
	available := pg.maxWidth - 5

	diff := 1 / float64(available)
	blocks := ""
	total := diff

	for i := 0; i < available; i++ {
		r := ' '
		if total <= f {
			r = pg.symbol
		}
		blocks += string(r)
		total += diff
	}

	n := pg.maxWidth - available - len(percentage)
	spaces := ""
	for i := 0; i < n; i++ {
		spaces += " "
	}

	return blocks + spaces + percentage
}
