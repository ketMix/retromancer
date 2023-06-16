package menu

import (
	"ebijam23/states"
	"ebijam23/states/game"

	"github.com/hajimehoshi/ebiten/v2"
)

type SinglePlayer struct {
}

func (s *SinglePlayer) Init(ctx states.Context) error {
	return nil
}

func (s *SinglePlayer) Update(ctx states.Context) error {
	ctx.StateMachine.PopState()
	ctx.StateMachine.PushState(&game.World{
		Players: []game.Player{
			&game.LocalPlayer{},
		},
	})
	return nil
}

func (s *SinglePlayer) Draw(screen *ebiten.Image) {
}
