package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Actor interface {
	Dead() bool
	Player() Player
	SetPlayer(p Player)
	SetImpulses(impulses ImpulseSet)
	Update() []Action
	Draw(screen *ebiten.Image)
	Shape() Shape
	Bounds() (x, y, w, h float64)
	SetXY(x, y float64)
	SetSize(r float64)
}
