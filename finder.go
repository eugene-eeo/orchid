package main

import "strings"
import "sort"
import "github.com/eugene-eeo/orchid/liborchid"
import "github.com/nsf/termbox-go"

type item struct {
	str string
	idx int
}

func matchAll(query string, haystack []*item) []*item {
	matching := []*item{}
	for _, x := range haystack {
		if liborchid.Match(query, x.str) {
			matching = append(matching, x)
		}
	}
	sort.Slice(matching, func(i, j int) bool {
		return len(matching[i].str) < len(matching[j].str)
	})
	return matching
}

type FinderUI struct {
	songs   []*liborchid.Song
	items   []*item
	results []*item
	choice  chan *liborchid.Song
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
		choice:  make(chan *liborchid.Song),
		input:   liborchid.NewInput(),
		viewbox: liborchid.NewViewbox(len(items), 7),
		cursor:  0,
	}
}

func (f *FinderUI) Find(q string) []*item {
	return matchAll(strings.ToLower(q), f.items)
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
	termbox.SetCell(1, 0, '>', termbox.ColorRed, termbox.ColorDefault)
	m := f.input.Cursor()
	u := ' '
	unicodeCells(query, 46, false, func(x int, r rune) {
		termbox.SetCell(3+x, 0, r, termbox.ColorDefault, termbox.ColorDefault)
		if x == f.input.Cursor() {
			u = r
		}
	})
	termbox.SetCell(3+m, 0, u, termbox.AttrReverse, termbox.ColorDefault)
}

func (f *FinderUI) RenderResults() {
	j := 0
	for i := f.viewbox.Lo(); i < f.viewbox.Hi(); i++ {
		song := f.Get(f.results[i])
		color := termbox.ColorDefault
		if i == f.cursor {
			color = termbox.AttrReverse
		}
		unicodeCells(song.Name(), 48, true, func(x int, r rune) {
			termbox.SetCell(1+x, 1+j, r, color, termbox.ColorDefault)
		})
		j++
	}
}

func (f *FinderUI) Render() {
	must(termbox.Clear(termbox.ColorDefault, termbox.ColorDefault))
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

func (f *FinderUI) Loop() {
	exit := false
	for !exit {
		f.viewbox.Update(f.cursor)
		f.Render()
		ev := termbox.PollEvent()
		if ev.Type != termbox.EventKey {
			continue
		}
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
			f.results = f.Find(f.input.String())
			f.viewbox = liborchid.NewViewbox(len(f.results), 7)
			f.cursor = 0
		}
	}
	f.choice <- f.Choice()
}
