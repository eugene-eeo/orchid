package player

import "math/rand"
import "sort"

func mod(r int, m int) int {
	t := r % m
	if t < 0 {
		t += m
	}
	return t
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

func remove(i int, xs []int) []int {
	return append(xs[:i], xs[i+1:]...)
}

type Indexer interface {
	Peek(i int) int
	Next(i int, force bool) int
}

type Repeat struct {
	Indexer Indexer
}

func (r *Repeat) Peek(i int) int {
	return r.Indexer.Peek(i)
}

func (r *Repeat) Next(i int, force bool) int {
	if !force {
		i = 0
	}
	return r.Indexer.Next(i, force)
}

type Seq struct {
	seq    []int
	cursor int
}

func NewSeq(n int) *Seq {
	b := make([]int, n)
	for i := 0; i < n; i++ {
		b[i] = i
	}
	return &Seq{seq: b, cursor: 0}
}

func (s *Seq) Peek(i int) int {
	if len(s.seq) == 0 {
		return -1
	}
	return s.seq[mod(s.cursor+i, len(s.seq))]
}

func (s *Seq) Next(i int, force bool) int {
	if len(s.seq) == 0 {
		return -1
	}
	s.cursor = mod(s.cursor+i, len(s.seq))
	return s.seq[s.cursor]
}

func (s *Seq) Sort() {
	n := s.Peek(0)
	sort.Ints(s.seq)
	s.cursor = sort.SearchInts(s.seq, n)
}

func (s *Seq) Shuffle() {
	s.cursor = shuffle(s.seq, s.cursor)
}

func (s *Seq) Pop() {
	if len(s.seq) > 0 {
		s.seq = remove(s.cursor, s.seq)
	}
}

func (s *Seq) Each(f func(int) bool) {
	for i := 0; i < len(s.seq); i++ {
		if !f(s.Peek(i)) {
			break
		}
	}
}
