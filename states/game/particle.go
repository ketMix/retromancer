package game

import (
	"ebijam23/resources"
	"ebijam23/states"
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Particle struct {
	Sprite *resources.Sprite
	X, Y   float64
	VX, VY float64
	Age    int
	Life   int
}

func (w *World) SpawnParticle(ctx states.Context, img string, x, y float64, angle float64, speed float64, life int) {
	vx := speed * math.Cos(angle)
	vy := speed * math.Sin(angle)

	p := Particle{
		Sprite: resources.NewSprite(ctx.Manager.GetAs("images", fmt.Sprintf("particle-%s", img), (*ebiten.Image)(nil)).(*ebiten.Image)),
		X:      x,
		Y:      y,
		VX:     vx,
		VY:     vy,
		Age:    0,
		Life:   life,
	}
	w.activeMap.particles = append(w.activeMap.particles, &p)
}

func (p *Particle) Update() {
	p.X += p.VX
	p.Y += p.VY
	p.Sprite.X = p.X
	p.Sprite.Y = p.Y
	p.Age++
}

func (p *Particle) Draw(ctx states.DrawContext) {
	opts := &ebiten.DrawImageOptions{}
	opts.ColorScale.Reset()
	m := 1.0 - float32(p.Age)/float32(p.Life)
	opts.ColorScale.Scale(m, m, m, m)
	p.Sprite.DrawWithOptions(ctx, opts)
}

func (p *Particle) Dead() bool {
	return p.Age >= p.Life
}
