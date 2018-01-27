package liborchid

import "math/rand"
import "sort"

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

func sortSongs(xs []*Song) {
	sort.Slice(xs, func(i, j int) bool {
		return string(xs[i].Name()) < string(xs[j].Name())
	})
}

func mod(r int, m int) int {
	t := r % m
	if t < 0 {
		t += m
	}
	return t
}

type Player struct {
	Shuffle bool
	Repeat  bool
	Songs   []*Song
	curr    int
}

func NewPlayer(songs []*Song) *Player {
	sortSongs(songs)
	return &Player{
		Shuffle: false,
		Repeat:  false,
		Songs:   songs,
	}
}

func (p *Player) ToggleRepeat() {
	p.Repeat = !p.Repeat
}

func (p *Player) ToggleShuffle() {
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

func (p *Player) Song() *Song {
	return p.Peek(0)
}

func (p *Player) Peek(i int) *Song {
	if len(p.Songs) == 0 {
		return nil
	}
	j := mod(p.curr+i, len(p.Songs))
	return p.Songs[j]
}

func (p *Player) Next(i int, force bool) *Song {
	if len(p.Songs) == 0 {
		return nil
	}
	if !p.Repeat || force {
		p.curr = mod(p.curr+i, len(p.Songs))
	}
	return p.Song()
}

func (p *Player) Remove() {
	p.Songs = remove(p.curr, p.Songs)
}

func (p *Player) SetCurrent(s *Song) {
	for i := 0; i < len(p.Songs); i++ {
		if p.Songs[i] == s {
			p.curr = i
			break
		}
	}
}
