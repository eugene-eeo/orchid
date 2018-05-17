package main

import "sync"
import "github.com/nsf/termbox-go"

type Component interface {
	Handle(termbox.Event)
	OnFocus()
}

type Reactor struct {
	mutex   sync.Mutex
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
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if c == nil {
		c = r.root
	}
	r.focused = c
	c.OnFocus()
}

func (r *Reactor) Focused() Component {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.focused
}

func (r *Reactor) InFocus(c Component) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
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
