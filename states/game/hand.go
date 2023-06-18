package game

import "ebijam23/resources"

type Hand struct {
	Shape  Shape
	Sprite *resources.Sprite
}

func (h *Hand) SetXY(x, y float64) {
	h.Shape.X = x
	h.Shape.Y = y
	h.Sprite.SetXY(x, y)
}