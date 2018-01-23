package player

import "os"
import "strings"
import "errors"

var NoMoreSongs error = errors.New("No more songs")

func FindSongs(dir string) (songs []Song, err error) {
	f, err := os.Open(dir)
	if err != nil {
		return
	}
	files, err := f.Readdirnames(-1)
	if err != nil {
		return
	}
	songs = []Song{}
	for _, name := range files {
		if strings.HasSuffix(name, ".mp3") {
			songs = append(songs, Song(name))
		}
	}
	return
}

type Player struct {
	Shuffle bool
	Repeat  bool
	Speaker *Speaker
	index   Indexer
	order   *Seq
	songs   []Song
}

func NewPlayer(songs []Song) *Player {
	seq := NewSeq(len(songs))
	return &Player{
		Shuffle: false,
		Repeat:  false,
		Speaker: NewSpeaker(),
		index:   seq,
		order:   seq,
		songs:   songs,
	}
}

func (p *Player) ToggleRepeat() {
	p.Repeat = !p.Repeat
	if !p.Repeat {
		p.index = p.order
	} else {
		p.index = &Repeat{p.order}
	}
}

func (p *Player) ToggleShuffle() {
	p.Shuffle = !p.Shuffle
	if p.Shuffle {
		p.order.Shuffle()
	} else {
		p.order.Sort()
	}
}

func (p *Player) Song() (Song, error) {
	return p.Peek(0)
}

func (p *Player) Peek(i int) (Song, error) {
	j := p.index.Peek(i)
	if j == -1 {
		return Song(""), NoMoreSongs
	}
	return p.songs[j], nil
}

func (p *Player) Next(i int, force bool) (Song, error) {
	p.index.Next(i, force)
	return p.Song()
}

func (p *Player) Remove() {
	p.order.Pop()
}

func (p *Player) Toggle() {
	p.Speaker.Toggle()
}

func (p *Player) SetCurrent(s Song) {
	i := -1
	p.order.Each(func(j int) bool {
		if p.songs[j] == s {
			i = j
			return false
		}
		return true
	})
	if i != -1 {
		p.order.cursor = i
	}
}

func (p *Player) Songs() []Song {
	a := []Song{}
	p.order.Each(func(i int) bool {
		a = append(a, p.songs[i])
		return true
	})
	return a
}
