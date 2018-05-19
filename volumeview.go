package main

import (
	"math"
	"time"

	"github.com/eugene-eeo/orchid/liborchid"
	"github.com/nsf/termbox-go"
)

const MAX_VOLUME float64 = +0.0
const MIN_VOLUME float64 = -12.0

type volumeUI struct {
	volume float64
	mw     *liborchid.MWorker
	bar    *liborchid.ProgressBar
	timer  *time.Timer
}

func newVolumeUI(mw *liborchid.MWorker) *volumeUI {
	timer := time.NewTimer(time.Duration(2) * time.Second)
	go func() {
		<-timer.C
		REACTOR.Focus(nil)
	}()
	return &volumeUI{
		bar:    liborchid.NewProgressBar(46, '▊'),
		volume: mw.VolumeInfo().Volume(),
		timer:  timer,
		mw:     mw,
	}
}

// Layout (50x8)
//
//  3x space
//   BLOCK CHARACTERS 46x1
//  4x space
//
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
	must(termbox.Flush())
}

func (v *volumeUI) resetTimer() {
	v.timer.Reset(time.Duration(2) * time.Second)
}

func (v *volumeUI) changeVolume(diff float64) {
	v.volume = math.Max(math.Min(v.volume+diff, MAX_VOLUME), MIN_VOLUME)
	v.mw.VolumeChange <- liborchid.VolumeInfo{
		V:   v.volume,
		Min: MIN_VOLUME,
		Max: MAX_VOLUME,
	}
}

func (v *volumeUI) OnFocus() {
	v.render()
}

func (v *volumeUI) Handle(evt termbox.Event) {
	switch evt.Key {
	case termbox.KeyEsc:
		v.timer.Stop()
		REACTOR.Focus(nil)
		return
	case termbox.KeyArrowLeft:
		v.changeVolume(-0.125)
		v.resetTimer()
	case termbox.KeyArrowRight:
		v.changeVolume(+0.125)
		v.resetTimer()
	}
	v.render()
}
