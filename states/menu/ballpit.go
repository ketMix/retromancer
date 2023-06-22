package menu

import (
	"ebijam23/states"
	"ebijam23/states/game"
)

type Ballpit struct {
}

func (s *Ballpit) Init(ctx states.Context) error {
	return nil
}

func (s *Ballpit) Finalize(ctx states.Context) error {
	return nil
}

func (s *Ballpit) Update(ctx states.Context) error {
	ctx.StateMachine.PopState()
	world := game.World{
		StartingMap: "ballpit",
		Players: []game.Player{
			game.NewLocalPlayer(),
		},
	}
	ctx.StateMachine.PushState(&world)
	return nil
}

func (s *Ballpit) Draw(ctx states.DrawContext) {
}
