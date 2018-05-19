package main

import (
	"sync"

	"github.com/nsf/termbox-go"
)

type Component interface {
	Handle(termbox.Event)
	OnFocus()
}

type Reactor struct {
	sync.Mutex
	focused Component
	root    Component
}

func NewReactor(root Component) *Reactor {
	return &Reactor{
		focused: root,
		root:    root,
	}
}

func (r *Reactor) Focus(c Component) {
	r.Lock()
	defer r.Unlock()
	if c == nil {
		c = r.root
	}
	r.focused = c
	c.OnFocus()
}

func (r *Reactor) Focused() Component {
	r.Lock()
	defer r.Unlock()
	return r.focused
}

func (r *Reactor) InFocus(c Component) bool {
	r.Lock()
	defer r.Unlock()
	return r.focused == c
}

func (r *Reactor) Loop() {
	for {
		evt := termbox.PollEvent()
		if evt.Type != termbox.EventKey {
			continue
		}
		r.Focused().Handle(evt)
	}
}
