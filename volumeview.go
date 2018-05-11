package main

import "time"
import "math"
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
	events chan termbox.Event
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
		events: make(chan termbox.Event),
	}
}

func (v *volumeUI) render() {
	ratio := (v.volume - MIN_VOLUME) / (MAX_VOLUME - MIN_VOLUME)
	must(termbox.Clear(ATTR_DEFAULT, ATTR_DEFAULT))
	unicodeCells(
		v.bar.Update(ratio), 46, false,
		func(dx int, r rune) {
			termbox.SetCell(
				2+dx, 3, r,
				ATTR_DEFAULT, ATTR_DEFAULT)
		})
	must(termbox.Sync())
}

func (v *volumeUI) resetTimer() {
	v.timer.Reset(time.Duration(2) * time.Second)
}

func (v *volumeUI) changeVolume(diff float64) {
	v.volume = math.Max(math.Min(v.volume+diff, MAX_VOLUME), MIN_VOLUME)
	if v.stream != nil {
		v.stream.SetVolume(
			v.volume,
			MIN_VOLUME,
			MAX_VOLUME,
		)
	}
}

func (v *volumeUI) Sink() chan termbox.Event {
	return v.events
}

func (v *volumeUI) Loop() {
loop:
	for {
		v.render()
		select {
		case evt := <-v.events:
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
	REACTOR.Focus(nil)
}
