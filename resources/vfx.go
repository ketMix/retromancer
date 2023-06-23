package resources

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type VFX interface {
	Process(screen *ebiten.Image, opts *ebiten.DrawImageOptions)
	Done() bool
}

type Fade struct {
	Out          bool
	Alpha        float64
	Duration     time.Duration
	elapsed      time.Duration
	lastTime     time.Time
	ApplyToImage bool
}

func (f *Fade) Process(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	t := time.Now()
	if f.lastTime.IsZero() {
		f.lastTime = t
	}
	f.elapsed += t.Sub(f.lastTime)
	f.lastTime = t

	m := float64(f.elapsed) / float64(f.Duration)

	if f.Out {
		m = 1.0 - m
	}

	if f.ApplyToImage {
		screen.Fill(color.NRGBA{R: 0, G: 0, B: 0, A: uint8(f.Alpha * m * 255)})
	} else {
		c := color.NRGBA{R: 255, G: 255, B: 255, A: uint8(f.Alpha * m * 255)}
		opts.ColorScale.ScaleWithColor(c)
	}
}

func (f *Fade) Done() bool {
	return f.elapsed >= f.Duration
}
