package liborchid

import (
	"math/rand"
	"sort"
)

func shuffle(xs []*Song, i int) int {
	x := xs[i]
	for j := 0; j < len(xs); j++ {
		r := rand.Intn(len(xs))
		xs[j], xs[r] = xs[r], xs[j]
		if xs[j] == x {
			i = j
		}
		if xs[r] == x {
			i = r
		}
	}
	return i
}

func remove(i int, xs []*Song) []*Song {
	return append(xs[:i], xs[i+1:]...)
}

func sortSongs(xs []*Song) []*Song {
	sort.Slice(xs, func(i, j int) bool {
		return string(xs[i].Name()) < string(xs[j].Name())
	})
	return xs
}

func mod(r int, m int) int {
	t := r % m
	if t < 0 {
		t += m
	}
	return t
}

type Queue struct {
	Shuffle bool
	Repeat  bool
	Songs   []*Song
	curr    int
}

func NewQueue(songs []*Song) *Queue {
	return &Queue{
		Songs:   sortSongs(songs),
		Shuffle: false,
		Repeat:  false,
	}
}

func (p *Queue) ToggleRepeat() {
	p.Repeat = !p.Repeat
}

func (p *Queue) ToggleShuffle() {
	p.Shuffle = !p.Shuffle
	if p.Shuffle {
		p.curr = shuffle(p.Songs, p.curr)
	} else {
		song := p.Song()
		sortSongs(p.Songs)
		if song != nil {
			p.SetCurrent(song)
		}
	}
}

func (p *Queue) Song() *Song {
	return p.Peek(0)
}

func (p *Queue) Peek(i int) *Song {
	if len(p.Songs) == 0 {
		return nil
	}
	j := mod(p.curr+i, len(p.Songs))
	return p.Songs[j]
}

func (p *Queue) Next(i int, force bool) *Song {
	if len(p.Songs) == 0 {
		return nil
	}
	if !p.Repeat || force {
		p.curr = mod(p.curr+i, len(p.Songs))
	}
	return p.Song()
}

func (p *Queue) Remove(s *Song) {
	j := -1
	for i, song := range p.Songs {
		if song == s {
			j = i
			break
		}
	}
	if j != -1 {
		p.Songs = remove(p.curr, p.Songs)
	}
}

func (p *Queue) SetCurrent(s *Song) {
	for i, song := range p.Songs {
		if song == s {
			p.curr = i
			break
		}
	}
}
