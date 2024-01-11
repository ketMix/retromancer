package menu

import (
	"github.com/ketMix/retromancer/states"
	"github.com/ketMix/retromancer/states/game"
)

type Ballpit struct {
}

func (s *Ballpit) Init(ctx states.Context) error {
	return nil
}

func (s *Ballpit) Finalize(ctx states.Context) error {
	return nil
}

func (s *Ballpit) Enter(ctx states.Context, v interface{}) error {
	return nil
}

func (s *Ballpit) Update(ctx states.Context) error {
	ctx.StateMachine.PopState(nil)
	world := game.World{
		StartingMap: "ballpit",
		Players: []game.Player{
			game.NewLocalPlayer(),
		},
	}
	world.PushState(&game.WorldStateLive{}, ctx) // Skip to actual gameplay state.
	ctx.StateMachine.PushState(&world)
	return nil
}

func (s *Ballpit) Draw(ctx states.DrawContext) {
}
