package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
)

type DrawContext struct {
	Screen *ebiten.Image
	Text   *etxt.Renderer
}

type Context struct {
	Manager      ResourceManager
	StateMachine StateMachine
	Cursor       Cursor
}
