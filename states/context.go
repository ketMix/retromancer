package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
)

type DrawContext struct {
	Screen *ebiten.Image
	Text   *etxt.Renderer
	R      Resource
	L      Localizer
}

type Context struct {
	StateMachine StateMachine
	Cursor       Cursor
	MusicPlayer  MusicPlayer
	Difficulty   Difficulty
	R            Resource
	L            Localizer
}
