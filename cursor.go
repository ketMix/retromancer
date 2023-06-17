package main

import "github.com/hajimehoshi/ebiten/v2"

type Cursor struct {
	image   *ebiten.Image
	enabled bool
}

func (c *Cursor) Enabled() bool {
	return c.enabled
}

func (c *Cursor) Enable() {
	c.enabled = true
}

func (c *Cursor) Disable() {
	c.enabled = false
}
