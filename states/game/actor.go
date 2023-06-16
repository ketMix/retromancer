package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Actor interface {
	Player() Player
	SetPlayer(p Player)
	Update() error
	Draw(screen *ebiten.Image)
	Bounds() (x, y, w, h float64)
	SetXY(x, y float64)
	SetSize(r float64)
}
