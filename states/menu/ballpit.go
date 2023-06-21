package menu

import (
	"ebijam23/states"
	"ebijam23/states/game"

	"github.com/hajimehoshi/ebiten/v2"
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
			&game.LocalPlayer{},
		},
	}
	ctx.StateMachine.PushState(&world)
	return nil
}

func (s *Ballpit) Draw(screen *ebiten.Image) {
}
