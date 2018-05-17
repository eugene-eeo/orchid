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

	total := int(f * float64(available))
	blocks := ""

	for i := 1; i <= available; i++ {
		r := ' '
		if i <= total {
			r = pg.symbol
		}
		blocks += string(r)
	}

	n := pg.maxWidth - available - len(percentage)
	spaces := ""
	for i := 0; i < n; i++ {
		spaces += " "
	}

	return blocks + spaces + percentage
}
