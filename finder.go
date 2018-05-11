package main

import "strings"
import "sort"
import "github.com/eugene-eeo/orchid/liborchid"
import "github.com/nsf/termbox-go"

type item struct {
	str   string
	idx   int
	score int
}

func matchAll(query string, haystack []*item) []*item {
	matching := []*item{}
	for _, x := range haystack {
		if ok, score := liborchid.Match(query, x.str); ok {
			matching = append(matching, x)
			x.score = score
		}
	}
	sort.Slice(matching, func(i, j int) bool {
		return matching[i].score < matching[j].score
	})
	return matching
}

type FinderUI struct {
	songs   []*liborchid.Song
	items   []*item
	results []*item
	input   *liborchid.Input
	viewbox *liborchid.Viewbox
	cursor  int
}

func newFinderUIFromPlayer(p *liborchid.Player) *FinderUI {
	items := make([]*item, len(p.Songs))
	for i, song := range p.Songs {
		items[i] = &item{
			str: strings.ToLower(song.Name()),
			idx: i,
		}
	}
	return &FinderUI{
		songs:   p.Songs,
		items:   items,
		results: items,
		input:   liborchid.NewInput(),
		viewbox: liborchid.NewViewbox(len(items), 7),
		cursor:  0,
	}
}

func (f *FinderUI) Get(i *item) *liborchid.Song {
	return f.songs[i.idx]
}

// > ________________ 48x1 => 46x1 query
// 1                  48x7
// 2
// ...

func (f *FinderUI) RenderQuery() {
	query := f.input.String()
	termbox.SetCell(1, 0, '⏵', termbox.ColorBlue, ATTR_DEFAULT)
	m := f.input.Cursor()
	u := ' '
	unicodeCells(query, 46, false, func(x int, r rune) {
		termbox.SetCell(3+x, 0, r, ATTR_DEFAULT, ATTR_DEFAULT)
		if x == f.input.Cursor() {
			u = r
		}
	})
	termbox.SetCell(3+m, 0, u, termbox.AttrReverse, ATTR_DEFAULT)
}

func (f *FinderUI) RenderResults() {
	j := 0
	for i := f.viewbox.Lo(); i < f.viewbox.Hi(); i++ {
		song := f.Get(f.results[i])
		color := ATTR_DEFAULT
		if i == f.cursor {
			color = termbox.AttrReverse
		}
		drawName(song.Name(), 1, j+1, color)
		j++
	}
}

func (f *FinderUI) Render() {
	must(termbox.Clear(ATTR_DEFAULT, ATTR_DEFAULT))
	f.RenderResults()
	f.RenderQuery()
	must(termbox.Sync())
}

func (f *FinderUI) Choice() *liborchid.Song {
	if len(f.results) == 0 || f.cursor < 0 {
		return nil
	}
	return f.Get(f.results[f.cursor])
}

func (f *FinderUI) Loop(events <-chan termbox.Event) {
	exit := false
	for !exit {
		f.viewbox.Update(f.cursor)
		f.Render()
		ev := <-events
		switch ev.Key {
		case termbox.KeyArrowUp:
			f.cursor--
			if f.cursor < 0 {
				f.cursor = len(f.results) - 1
			}
		case termbox.KeyTab:
			fallthrough
		case termbox.KeyArrowDown:
			f.cursor++
			if f.cursor > len(f.results)-1 {
				f.cursor = 0
			}
		case termbox.KeyEsc:
			f.cursor = -1
			exit = true
		case termbox.KeyEnter:
			exit = true
		default:
			f.input.Feed(ev.Key, ev.Ch, ev.Mod)
			f.results = matchAll(strings.ToLower(f.input.String()), f.items)
			f.viewbox = liborchid.NewViewbox(len(f.results), 7)
			f.cursor = 0
		}
	}
}
