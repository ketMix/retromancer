package resources

import (
	"ebijam23/states"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
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

type Text struct {
	Text         string
	X, Y         float64
	Scale        float64
	Outline      bool
	OutlineColor color.NRGBA
	Delay        time.Duration
	hasDelayed   bool
	InDuration   time.Duration
	HoldDuration time.Duration
	OutDuration  time.Duration
	elapsed      time.Duration
	lastTime     time.Time
}

func (v *Text) Process(ctx states.DrawContext, opts *ebiten.DrawImageOptions) {
	t := time.Now()
	if v.lastTime.IsZero() {
		v.lastTime = t
	}
	v.elapsed += t.Sub(v.lastTime)
	v.lastTime = t

	m := float64(v.elapsed)

	if !v.hasDelayed {
		if v.elapsed < v.Delay {
			return
		}
		v.hasDelayed = true
		v.elapsed = 0
	}

	if v.elapsed < v.InDuration {
		m /= float64(v.InDuration)
	} else if v.elapsed < v.InDuration+v.HoldDuration {
		m = 1.0
	} else if v.elapsed < v.OutDuration+v.InDuration+v.HoldDuration {
		m -= float64(v.InDuration + v.HoldDuration)
		m /= float64(v.OutDuration)
		m = 1.0 - m
	}

	if v.Scale != 0 {
		ctx.Text.SetScale(v.Scale)
	}

	if v.Outline {
		ctx.Text.SetColor(color.NRGBA{R: v.OutlineColor.R, G: v.OutlineColor.G, B: v.OutlineColor.B, A: uint8(float64(v.OutlineColor.A) * m)})
		DrawTextOutline(ctx.Text, ctx.Screen, v.Text, int(v.X), int(v.Y), int(v.Scale))
	}

	ctx.Text.SetAlign(etxt.XCenter | etxt.YCenter)
	ctx.Text.SetColor(color.NRGBA{R: 255, G: 255, B: 255, A: uint8(255 * m)})
	ctx.Text.Draw(ctx.Screen, v.Text, int(v.X), int(v.Y))
}

func (v *Text) Done() bool {
	return v.elapsed >= v.InDuration+v.HoldDuration+v.OutDuration
}
