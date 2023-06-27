package resources

import (
	"ebijam23/states"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
)

type VFXListMode int

const (
	Parallel VFXListMode = iota
	Sequential
)

type VFXList struct {
	items []VFX
	mode  VFXListMode
}

func (v *VFXList) Items() []VFX {
	return v.items
}

func (v *VFXList) SetMode(mode VFXListMode) {
	v.mode = mode
}

func (v *VFXList) Add(vfx VFX) {
	v.items = append(v.items, vfx)
}

func (v *VFXList) RemoveByID(s string) {
	for i, vfx := range v.items {
		if vfx.ID() == s {
			v.items = append(v.items[:i], v.items[i+1:]...)
			return
		}
	}
}

func (v *VFXList) Process(ctx states.DrawContext, opts *ebiten.DrawImageOptions) {
	if v.mode == Sequential {
		if len(v.items) > 0 {
			vfx := v.items[0]
			vfx.Process(ctx, opts)
			if vfx.Done() {
				v.items = v.items[1:]
			}
		}
		return
	}
	for i := 0; i < len(v.items); i++ {
		vfx := v.items[i]
		vfx.Process(ctx, opts)
		if vfx.Done() {
			v.items = append(v.items[:i], v.items[i+1:]...)
			i--
		}
	}
}

func (v *VFXList) Empty() bool {
	return len(v.items) == 0
}

type VFX interface {
	Process(ctx states.DrawContext, opts *ebiten.DrawImageOptions)
	Done() bool
	ID() string
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

func (f *Fade) ID() string {
	return "fade"
}

func (f *Fade) Process(ctx states.DrawContext, opts *ebiten.DrawImageOptions) {
	if opts == nil {
		opts = &ebiten.DrawImageOptions{}
	}
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

func (t *Text) ID() string {
	return "text"
}

func (v *Text) Process(ctx states.DrawContext, opts *ebiten.DrawImageOptions) {
	t := time.Now()
	if v.lastTime.IsZero() {
		v.lastTime = t
	}
	v.elapsed += t.Sub(v.lastTime)
	v.lastTime = t

	if !v.hasDelayed {
		if v.elapsed < v.Delay {
			return
		}
		v.hasDelayed = true
		v.elapsed = 0
	}

	m := float64(v.elapsed)

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

type Hover struct {
	elapsed   time.Duration
	Intensity float64
	Rate      float64
}

func (h *Hover) ID() string {
	return "hover"
}

func (h *Hover) Process(ctx states.DrawContext, opts *ebiten.DrawImageOptions) {
	h.elapsed += time.Duration(float64(time.Second/60) * h.Rate)
	opts.GeoM.Translate(0, h.Intensity*math.Sin(float64(h.elapsed.Seconds())))
}

func (h *Hover) Done() bool {
	return false
}

type Darkness struct {
	elapsed  time.Duration
	Fade     bool
	Duration time.Duration
	img      *ebiten.Image
	Fin      bool
}

func (d *Darkness) ID() string {
	return "darkness"
}

func (d *Darkness) Process(ctx states.DrawContext, opts *ebiten.DrawImageOptions) {
	if opts == nil {
		opts = &ebiten.DrawImageOptions{}
	}
	if d.img == nil {
		d.img = ebiten.NewImage(ctx.Screen.Bounds().Dx(), ctx.Screen.Bounds().Dy())
		d.img.Fill(color.NRGBA{R: 0, G: 0, B: 0, A: 250})
		// Set pixels to make an empty circle.
		for x := 0; x < d.img.Bounds().Dx(); x++ {
			for y := 0; y < d.img.Bounds().Dy(); y++ {
				if math.Sqrt(math.Pow(float64(x-d.img.Bounds().Dx()/2), 2)+math.Pow(float64(y-d.img.Bounds().Dy()/2), 2)) < float64(d.img.Bounds().Dx()/8) {
					d.img.Set(x, y, color.NRGBA{R: 0, G: 0, B: 0, A: 200})
				}
			}
		}
	}
	m := 1.0
	if d.Fade {
		d.elapsed += time.Second / 60
		m = 1.0 - float64(d.elapsed)/float64(d.Duration)
		if m <= 0 {
			m = 0
		}
		if m == 0 {
			d.Fin = true
		}
	}
	c := color.NRGBA{R: 0, G: 0, B: 0, A: uint8(255 * m)}
	opts.ColorScale.ScaleWithColor(c)
	ctx.Screen.DrawImage(d.img, opts)
}

func (d *Darkness) Done() bool {
	return d.Fin
}

type VFXDef struct {
	Type     string
	Duration time.Duration
}
