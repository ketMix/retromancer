package game

import (
	"ebijam23/resources"

	"github.com/hajimehoshi/ebiten/v2"
)

type PC struct {
	player Player
	Sprite *resources.Sprite
	Shape  Shape
}

func (p *PC) SetPlayer(player Player) {
	p.player = player
}

func (p *PC) Player() Player {
	return p.player
}

func (p *PC) Update() error {
	return nil
}

func (p *PC) Draw(screen *ebiten.Image) {
	p.Sprite.Draw(screen)
}

func (p *PC) Bounds() (x, y, w, h float64) {
	// Return the radius of the shape as width and height in diameter.
	return p.Shape.X, p.Shape.Y, p.Shape.Radius * 2, p.Shape.Radius * 2
}

func (p *PC) SetXY(x, y float64) {
	p.Shape.X = x
	p.Shape.Y = y
	p.Sprite.X = x
	p.Sprite.Y = y
}

func (p *PC) SetSize(r float64) {
	p.Shape.Radius = r
}
