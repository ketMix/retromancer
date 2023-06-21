package menu

import (
	"ebijam23/states"
	"ebijam23/states/game"
)

type SinglePlayer struct {
}

func (s *SinglePlayer) Init(ctx states.Context) error {
	return nil
}

func (s *SinglePlayer) Finalize(ctx states.Context) error {
	return nil
}

func (s *SinglePlayer) Update(ctx states.Context) error {
	ctx.StateMachine.PopState()
	ctx.StateMachine.PushState(&game.World{
		StartingMap: "start",
		Players: []game.Player{
			&game.LocalPlayer{},
		},
	})
	return nil
}

func (s *SinglePlayer) Draw(ctx states.DrawContext) {
}
