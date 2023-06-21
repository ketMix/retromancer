package states

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type State interface {
	Init(ctx Context) error
	Finalize(ctx Context) error
	Update(ctx Context) error
	Draw(screen *ebiten.Image)
}

type StateMachine interface {
	PushState(state State)
	PopState()
}
