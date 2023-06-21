package menu

import (
	"ebijam23/states"
	"ebijam23/states/game"

	"github.com/hajimehoshi/ebiten/v2"
)

type MultiPlayer struct {
}

func (s *MultiPlayer) Init(ctx states.Context) error {
	return nil
}

func (s *MultiPlayer) Finalize(ctx states.Context) error {
	return nil
}

func (s *MultiPlayer) Update(ctx states.Context) error {
	ctx.StateMachine.PopState()
	ctx.StateMachine.PushState(&game.World{
		StartingMap: "start",
	})
	return nil
}

func (s *MultiPlayer) Draw(screen *ebiten.Image) {
}
