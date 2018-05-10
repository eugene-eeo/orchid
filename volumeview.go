package main

import "time"
import "fmt"
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
	pctg := (volume.Volume + 2) / 4
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	s := " "

	for i := 0; i < 40; i++ {
		if float64(i)/40 > pctg || pctg == 0 {
			s += " "
			continue
		}
		s += block
	}

	s += " "
	pctgString := fmt.Sprintf("%d%%", int(pctg*100))
	for i := 0; i < 4-len(pctgString); i++ {
		s += " "
	}
	s += pctgString

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
				vol.Volume += 0.125
				vol.Silent = false
				if vol.Volume > 2 {
					vol.Volume = 2
				}
			case termbox.KeyArrowLeft:
				v.timer.Reset(time.Duration(2) * time.Second)
				vol.Volume -= 0.125
				if vol.Volume <= -2 {
					vol.Volume = -2
					vol.Silent = true
				}
			}
		case _ = <-v.timer.C:
			v.done = true
		}
	}
}
