package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Actor interface {
	Player() Player
	SetPlayer(p Player)
	SetImpulses(impulses ImpulseSet)
	Update() []Action
	Draw(screen *ebiten.Image)
	Bounds() (x, y, w, h float64)
	SetXY(x, y float64)
	SetSize(r float64)
}