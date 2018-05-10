package main

import "time"
import "github.com/nsf/termbox-go"
import "github.com/eugene-eeo/orchid/liborchid"

// Layout (50x8)
//
//  3x space
//   BLOCK CHARACTERS 46x1
//  4x space
//

type volumeUI struct {
	stream *liborchid.Stream
	bar    *liborchid.ProgressBar
	timer  *time.Timer
	done   bool
}

func newVolumeUI(stream *liborchid.Stream) *volumeUI {
	return &volumeUI{
		bar:    liborchid.NewProgressBar(46, '▊'),
		stream: stream,
		timer:  time.NewTimer(time.Duration(2) * time.Second),
		done:   false,
	}
}

func (v *volumeUI) render() {
	ratio := (v.stream.Volume() - MIN_VOLUME) / (MAX_VOLUME - MIN_VOLUME)
	must(termbox.Clear(termbox.ColorDefault, termbox.ColorDefault))
	unicodeCells(v.bar.Update(ratio), 46, false, func(dx int, r rune) {
		termbox.SetCell(2+dx, 3, r, termbox.ColorDefault, termbox.ColorDefault)
	})
	must(termbox.Sync())
}

func (v *volumeUI) Loop(events <-chan termbox.Event) {
	for !v.done {
		v.render()
		select {
		case evt := <-events:
			switch evt.Key {
			case termbox.KeyEsc:
				if !v.timer.Stop() {
					<-v.timer.C
				}
				v.done = true
				break
			case termbox.KeyArrowLeft:
				v.timer.Reset(time.Duration(2) * time.Second)
				v.stream.SetVolume(
					v.stream.Volume()-0.125,
					MIN_VOLUME,
					MAX_VOLUME,
				)
			case termbox.KeyArrowRight:
				v.timer.Reset(time.Duration(2) * time.Second)
				v.stream.SetVolume(
					v.stream.Volume()+0.125,
					MIN_VOLUME,
					MAX_VOLUME,
				)
			}
		case _ = <-v.timer.C:
			v.done = true
		}
	}
}
