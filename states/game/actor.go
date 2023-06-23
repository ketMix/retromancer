package game

import (
	"ebijam23/states"
)

type Actor interface {
	Save()
	Restore()
	Dead() bool
	Destroyed() bool
	Player() Player
	SetPlayer(p Player)
	SetImpulses(impulses ImpulseSet)
	Update() []Action
	Draw(ctx states.DrawContext)
	Shape() Shape
	Bounds() (x, y, w, h float64)
	SetXY(x, y float64)
	SetSize(r float64)
}
