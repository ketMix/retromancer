package resources

import (
	"ebijam23/states"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type VFX interface {
	Process(ctx states.DrawContext, opts *ebiten.DrawImageOptions)
	Done() bool
}

type Fade struct {
	Out          bool
	Alpha        float64
	Duration     time.Duration
	elapsed      time.Duration
	lastTime     time.Time
	ApplyToImage bool
	fadeInImage  *ebiten.Image
}

func (f *Fade) Process(ctx states.DrawContext, opts *ebiten.DrawImageOptions) {
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
		if f.fadeInImage == nil {
			f.fadeInImage = ebiten.NewImage(ctx.Screen.Bounds().Dx(), ctx.Screen.Bounds().Dy())
		}
		f.fadeInImage.Fill(color.Black)
		opts := &ebiten.DrawImageOptions{}
		opts.ColorScale.ScaleAlpha(1 - float32(f.Alpha*m))
		ctx.Screen.DrawImage(f.fadeInImage, opts)
	} else {
		c := color.NRGBA{R: 255, G: 255, B: 255, A: uint8(f.Alpha * m * 255)}
		opts.ColorScale.ScaleWithColor(c)
	}
}

func (f *Fade) Done() bool {
	return f.elapsed >= f.Duration
}
