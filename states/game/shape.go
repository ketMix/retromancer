package game

import (
	"math"
)

type Shape interface {
	Clone() Shape
	Bounds() (x, y, width, height float64)
	Collides(s Shape) bool
}

type CircleShape struct {
	X, Y   float64
	Radius float64
}

func (c *CircleShape) Clone() Shape {
	return &CircleShape{
		X:      c.X,
		Y:      c.Y,
		Radius: c.Radius,
	}
}

func (c *CircleShape) Bounds() (x, y, width, height float64) {
	return c.X, c.Y, c.Radius * 2, c.Radius * 2
}

func (c *CircleShape) Collides(s Shape) bool {
	if c2, ok := s.(*CircleShape); ok {
		dx := c.X - c2.X
		dy := c.Y - c2.Y
		d := dx*dx + dy*dy
		r := c.Radius + c2.Radius
		return d <= r*r
	} else if r2, ok := s.(*RectangleShape); ok {
		edgeX := c.X
		edgeY := c.Y

		if c.X < r2.X {
			edgeX = r2.X
		} else if c.X > r2.X+r2.Width {
			edgeX = r2.X + r2.Width
		}
		if c.Y < r2.Y {
			edgeY = r2.Y
		} else if c.Y > r2.Y+r2.Height {
			edgeY = r2.Y + r2.Height
		}

		dx := c.X - edgeX
		dy := c.Y - edgeY
		d := math.Sqrt(dx*dx + dy*dy)
		return d <= c.Radius
	}
	return false
}

type RectangleShape struct {
	X, Y, Width, Height float64
}

func (r *RectangleShape) Clone() Shape {
	return &RectangleShape{
		X:      r.X,
		Y:      r.Y,
		Width:  r.Width,
		Height: r.Height,
	}
}

func (r *RectangleShape) Bounds() (x, y, width, height float64) {
	return r.X, r.Y, r.Width, r.Height
}

func (r *RectangleShape) Collides(s Shape) bool {
	if c2, ok := s.(*CircleShape); ok {
		return c2.Collides(r)
	} else if r2, ok := s.(*RectangleShape); ok {
		return r.X < r2.X+r2.Width &&
			r.X+r.Width > r2.X &&
			r.Y < r2.Y+r2.Height &&
			r.Y+r.Height > r2.Y
	}
	return false
}
