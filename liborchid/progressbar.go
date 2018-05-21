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
	percentage := fmt.Sprintf(" %3d%%", int(f*100))
	available := pg.maxWidth - len(percentage)
	total := int(f * float64(available))
	blocks := ""
	for i := 1; i <= available; i++ {
		r := ' '
		if i <= total {
			r = pg.symbol
		}
		blocks += string(r)
	}
	return blocks + percentage
}
