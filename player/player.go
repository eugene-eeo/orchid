package player

import "os"
import "strings"
import "errors"
import "math/rand"
import "sort"

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

func remove(i int, xs []int) []int {
	return append(xs[:i], xs[i+1:]...)
}

func mod(r int, m int) int {
	t := r % m
	if t < 0 {
		t += m
	}
	return t
}

func seq(n int) []int {
	b := make([]int, n)
	for i := 0; i < n; i++ {
		b[i] = i
	}
	return b
}

func shuffle(xs []int, c int) int {
	n := len(xs)
	m := c
	for i := 0; i < n; i++ {
		j := rand.Intn(n)
		xs[i], xs[j] = xs[j], xs[i]
		if c == xs[i] {
			m = i
		}
		if c == xs[j] {
			m = j
		}
	}
	return m
}

type Player struct {
	Shuffle bool
	Repeat  bool
	Stream  *Stream
	cursor  int
	order   []int
	Queue   []Song
}

func NewPlayer(songs []Song) *Player {
	return &Player{
		Shuffle: false,
		cursor:  0,
		order:   seq(len(songs)),
		Queue:   songs,
	}
}

func (p *Player) ToggleRepeat() {
	p.Repeat = !p.Repeat
}

func (p *Player) ToggleShuffle() {
	p.Shuffle = !p.Shuffle
	if p.Shuffle {
		p.cursor = shuffle(p.order, p.cursor)
	} else {
		if len(p.order) == 0 {
			return
		}
		p.cursor = p.order[p.cursor]
		sort.Ints(p.order)
	}
}

func (p *Player) Song() (Song, error) {
	return p.Peek(0)
}

func (p *Player) Next(i int, force bool) (chan bool, error) {
	if len(p.order) == 0 {
		return nil, NoMoreSongs
	}
	if !p.Repeat || force {
		p.cursor += i
	}
	p.cursor = mod(p.cursor, len(p.order))
	return p.Play()
}

func (p *Player) Peek(i int) (Song, error) {
	if len(p.order) == 0 {
		return Song(""), NoMoreSongs
	}
	return p.Queue[p.order[mod(p.cursor+i, len(p.order))]], nil
}

func (p *Player) play(done chan bool) (*Stream, error) {
	s, err := p.Song()
	if err != nil {
		return nil, err
	}
	stream, err := s.Stream(func(graceful bool) {
		go func() {
			done <- graceful
			close(done)
		}()
	})
	if err != nil {
		return stream, err
	}
	return stream, stream.Play()
}

func (p *Player) Play() (done chan bool, err error) {
	if p.Stream != nil {
		p.Stream.Teardown(false)
	}
	done = make(chan bool)
	stream, err := p.play(done)
	if err != nil {
		if stream != nil {
			stream.Teardown(false)
			<-done
		}
		p.Stream = nil
		return
	}
	p.Stream = stream
	return
}

func (p *Player) Remove() {
	if len(p.order) > 0 {
		p.order = remove(p.cursor, p.order)
	}
}

func (p *Player) Toggle() {
	if p.Stream != nil {
		p.Stream.Toggle()
	}
}
