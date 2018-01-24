package main

import "strings"
import "sort"
import "github.com/eugene-eeo/orchid/player"
import "github.com/eugene-eeo/orchid/elems"
import "github.com/nsf/termbox-go"

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type item struct {
	str string
	idx int
}

func match(a, b string) bool {
	B := []rune(b)
	n := len(B)
	j := 0
	for _, c := range a {
		found := false
		for j < n && !found {
			found = c == B[j]
			j++
		}
		if !found {
			return false
		}
	}
	return true
}

func rankMatch(a, b string) float64 {
	return float64(len(a)) / float64(len(b))
}

func matchAll(query string, haystack []*item) []*item {
	matching := []*item{}
	for _, x := range haystack {
		if match(query, x.str) {
			matching = append(matching, x)
		}
	}
	sort.Slice(matching, func(i, j int) bool {
		return rankMatch(query, matching[i].str) > rankMatch(query, matching[j].str)
	})
	return matching
}

type Finder struct {
	items []*item
	songs []player.Song
}

func finderFromPlayer(p *player.Player) *Finder {
	all := p.Songs()
	songs := make([]player.Song, len(all))
	items := make([]*item, len(all))
	for i, song := range all {
		songs[i] = song
		items[i] = &item{
			str: strings.ToLower(song.Name()),
			idx: i,
		}
	}
	return &Finder{
		items: items,
		songs: songs,
	}
}

func (f *Finder) Find(q string) []*item {
	return matchAll(strings.ToLower(q), f.items)
}

func (f *Finder) Get(i *item) player.Song {
	return f.songs[i.idx]
}

type FinderUI struct {
	finder   *Finder
	requests chan func(*FinderUI)
	choice   chan *player.Song
	results  []*item
	input    *elems.Input
	viewbox  *elems.Viewbox
	cursor   int
}

func newFinderUIFromPlayer(p *player.Player) *FinderUI {
	finder := finderFromPlayer(p)
	return &FinderUI{
		finder:   finder,
		requests: make(chan func(*FinderUI)),
		choice:   make(chan *player.Song),
		results:  finder.items,
		input:    elems.NewInput(),
		viewbox:  elems.NewViewBox(len(finder.items), 7),
		cursor:   0,
	}
}

// > ________________ 48x1 => 46x1 query
// 1                  48x7
// 2
// ...

func (f *FinderUI) RenderQuery() {
	query := f.input.String()
	termbox.SetCell(1, 0, '>', termbox.ColorRed, termbox.ColorDefault)
	m := -1
	unicodeCells(query, 46, false, func(x int, r rune) {
		termbox.SetCell(3+x, 0, r, termbox.AttrBold, termbox.ColorDefault)
		m = x
	})
	termbox.SetCursor(3+m+1, 0)
}

func (f *FinderUI) RenderResults() {
	j := 0
	for i := f.viewbox.Lo(); i < f.viewbox.Hi(); i++ {
		song := f.finder.Get(f.results[i])
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

func (f *FinderUI) Choice() *player.Song {
	if len(f.results) == 0 || f.cursor < 0 {
		return nil
	}
	song := f.finder.Get(f.results[f.cursor])
	return &song
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
			if f.cursor > 0 {
				f.cursor--
			}
		case termbox.KeyArrowDown:
			if f.cursor < len(f.results)-1 {
				f.cursor++
			}
		case termbox.KeyEsc:
			f.cursor = -1
			exit = true
		case termbox.KeyEnter:
			exit = true
		default:
			f.input.Feed(ev)
			f.results = f.finder.Find(f.input.String())
			f.viewbox = elems.NewViewBox(len(f.results), 7)
			f.cursor = 0
		}
	}
	f.choice <- f.Choice()
}
