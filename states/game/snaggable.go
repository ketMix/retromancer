package game

import (
	"ebijam23/resources"
	"ebijam23/states"
	"math"
	"math/rand"
)

type Snaggable struct {
	shape        CircleShape
	sprite       *resources.Sprite
	destroyed    bool
	nextParticle int
}

func CreateSnaggable(x, y float64, sprite *resources.Sprite) *Snaggable {
	sprite.VFX.Add(&resources.Hover{
		Intensity: 2.0,
		Rate:      1.25,
	})
	sprite.Centered = true
	return &Snaggable{
		shape:  CircleShape{X: x, Y: y, Radius: 3}, // FIXME: don't hardcode radius
		sprite: sprite,
	}
}

func (s *Snaggable) Update() (actions []Action) {
	s.nextParticle++
	if s.nextParticle >= 0 {
		actions = append(actions, ActionSpawnParticle{
			Img:   "life",
			X:     s.shape.X,
			Y:     s.shape.Y,
			Angle: math.Pi + rand.Float64()*math.Pi,
			Speed: rand.Float64() * 0.5,
			Life:  40,
		})
		s.nextParticle = -10
	}
	return
}

func (s *Snaggable) Draw(ctx states.DrawContext) {
	s.sprite.X = s.shape.X
	s.sprite.Y = s.shape.Y
	s.sprite.Draw(ctx)
}

func (s *Snaggable) Destroyed() bool {
	return s.destroyed
}

func (s *Snaggable) Shape() Shape                    { return &s.shape }
func (s *Snaggable) Save()                           {}
func (s *Snaggable) Restore()                        {}
func (s *Snaggable) Player() Player                  { return nil }
func (s *Snaggable) SetPlayer(p Player)              {}
func (s *Snaggable) SetImpulses(impulses ImpulseSet) {}
func (s *Snaggable) Bounds() (x, y, w, h float64)    { return 0, 0, 0, 0 }
func (s *Snaggable) SetXY(x, y float64)              {}
func (s *Snaggable) SetSize(r float64)               {}
func (s *Snaggable) Dead() bool                      { return false }
