package game

import (
	"ebijam23/resources"
	"ebijam23/states"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Snaggable struct {
	id           string
	spriteName   string
	shape        CircleShape
	sprite       *resources.Sprite
	destroyed    bool
	nextParticle int
}

func CreateSnaggable(ctx states.Context, id, spriteName string) *Snaggable {
	imageNames := ctx.R.GetNamesWithPrefix("images", spriteName)
	images := make([]*ebiten.Image, 0)
	for _, s := range imageNames {
		images = append(images, ctx.R.GetAs("images", s, (*ebiten.Image)(nil)).(*ebiten.Image))
	}
	sprite := resources.NewAnimatedSprite(images)

	sprite.VFX.Add(&resources.Hover{
		Intensity: 2.0,
		Rate:      1.25,
	})
	sprite.Centered = true
	return &Snaggable{
		id:         id,
		spriteName: spriteName,
		shape:      CircleShape{Radius: 3}, // FIXME: don't hardcode radius
		sprite:     sprite,
	}
}
func (s *Snaggable) SetXY(x, y float64) {
	s.shape.X = x
	s.shape.Y = y
	s.sprite.X = x
	s.sprite.Y = y
}

func (s *Snaggable) Update() (actions []Action) {
	s.nextParticle++
	if s.nextParticle >= 0 {
		// Set Img property here
		img := ""
		switch s.spriteName {
		case "item-life":
			img = "life"
		case "item-book":
			img = "book"
		case "item-shield":
			img = "shield"
		}
		actions = append(actions, ActionSpawnParticle{
			Img:   img,
			X:     s.shape.X,
			Y:     s.shape.Y,
			Angle: math.Pi + rng.Float64()*math.Pi,
			Speed: rng.Float64() * 0.5,
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
func (s *Snaggable) SetSize(r float64)               {}
func (s *Snaggable) Dead() bool                      { return false }
