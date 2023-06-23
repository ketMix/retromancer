package game

import "ebijam23/resources"

type Hand struct {
	Shape       CircleShape
	Sprite      *resources.Sprite
	HoverSprite *resources.Sprite
}

func (h *Hand) SetXY(x, y float64) {
	h.Shape.X = x
	h.Shape.Y = y
	h.Sprite.SetXY(x, y)
	h.HoverSprite.SetXY(x, y)
}
