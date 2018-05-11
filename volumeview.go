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
	volume float64
	stream *liborchid.Stream
	bar    *liborchid.ProgressBar
	timer  *time.Timer
}

func newVolumeUI(stream *liborchid.Stream) *volumeUI {
	vol := MAX_VOLUME
	if stream != nil {
		vol = stream.Volume()
	}
	return &volumeUI{
		bar:    liborchid.NewProgressBar(46, '▊'),
		timer:  time.NewTimer(time.Duration(2) * time.Second),
		stream: stream,
		volume: vol,
	}
}

func (v *volumeUI) render() {
	ratio := (v.volume - MIN_VOLUME) / (MAX_VOLUME - MIN_VOLUME)
	must(termbox.Clear(termbox.ColorDefault, termbox.ColorDefault))
	unicodeCells(v.bar.Update(ratio), 46, false, func(dx int, r rune) {
		termbox.SetCell(2+dx, 3, r, termbox.ColorDefault, termbox.ColorDefault)
	})
	must(termbox.Sync())
}

func (v *volumeUI) resetTimer() {
	v.timer.Reset(time.Duration(2) * time.Second)
}

func (v *volumeUI) changeVolume(diff float64) {
	vol := v.volume + diff
	if vol > MAX_VOLUME {
		vol = MAX_VOLUME
	}
	if vol < MIN_VOLUME {
		vol = MIN_VOLUME
	}
	v.volume = vol
	if v.stream != nil {
		v.stream.SetVolume(
			vol,
			MIN_VOLUME,
			MAX_VOLUME,
		)
	}
}

func (v *volumeUI) Loop(events <-chan termbox.Event) {
loop:
	for {
		v.render()
		select {
		case evt := <-events:
			switch evt.Key {
			case termbox.KeyEsc:
				if !v.timer.Stop() {
					<-v.timer.C
				}
				break loop
			case termbox.KeyArrowLeft:
				v.changeVolume(-0.125)
				v.resetTimer()
			case termbox.KeyArrowRight:
				v.changeVolume(+0.125)
				v.resetTimer()
			}
		case _ = <-v.timer.C:
			break loop
		}
	}
}
