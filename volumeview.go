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
	timer  *time.Timer
	done   bool
}

func newVolumeUI(stream *liborchid.Stream) *volumeUI {
	return &volumeUI{
		stream: stream,
		timer:  time.NewTimer(time.Duration(2) * time.Second),
		done:   false,
	}
}

func (v *volumeUI) render() {
	block := "▊"
	volume := v.stream.Volume()
	s := ""
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for i := 0; i < 46; i++ {
		if float64(i)/46 <= (volume.Volume+5)/10 {
			s += block
		} else {
			s += " "
		}
	}
	unicodeCells(s, 46, false, func(dx int, r rune) {
		termbox.SetCell(2+dx, 3, r, termbox.ColorDefault, termbox.ColorDefault)
	})
	termbox.Sync()
}

func (v *volumeUI) Loop(events <-chan termbox.Event) {
	for !v.done {
		v.render()
		vol := v.stream.Volume()
		select {
		case evt := <-events:
			switch evt.Key {
			case termbox.KeyEsc:
				if !v.timer.Stop() {
					<-v.timer.C
				}
				v.done = true
			case termbox.KeyArrowRight:
				v.timer.Reset(time.Duration(2) * time.Second)
				vol.Volume += 0.25
				vol.Silent = false
				if vol.Volume > 5 {
					vol.Volume = 5
				}
			case termbox.KeyArrowLeft:
				v.timer.Reset(time.Duration(2) * time.Second)
				vol.Volume -= 0.25
				if vol.Volume <= -5 {
					vol.Volume = -5
					vol.Silent = true
				}
			}
		case _ = <-v.timer.C:
			v.done = true
		}
	}
}
