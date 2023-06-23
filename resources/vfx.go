package resources

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type VFX interface {
	PreProcess(screen *ebiten.Image, opts *ebiten.DrawImageOptions)
	PostProcess(screen *ebiten.Image, opts *ebiten.DrawImageOptions)
	Done() bool
}

type Fade struct {
	Out      bool
	Alpha    float64
	Duration time.Duration
	elapsed  time.Duration
	lastTime time.Time
}

func (f *Fade) PreProcess(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
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

	c := color.NRGBA{R: 255.0, G: 255.0, B: 255.0, A: uint8(f.Alpha * m * 255)}
	opts.ColorScale.ScaleWithColor(c)
}

func (f *Fade) PostProcess(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
}

func (f *Fade) Done() bool {
	return f.elapsed >= f.Duration
}
