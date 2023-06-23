package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
)

type DrawContext struct {
	Screen *ebiten.Image
	Text   *etxt.Renderer
	L      func(key string) string
}

type Context struct {
	Manager      ResourceManager
	StateMachine StateMachine
	L            func(key string) string
	Cursor       Cursor
}
