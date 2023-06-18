package menu

import (
	"ebijam23/states"
	"ebijam23/states/game"

	"github.com/hajimehoshi/ebiten/v2"
)

type BulletBox struct {
}

func (s *BulletBox) Init(ctx states.Context) error {
	return nil
}

func (s *BulletBox) Update(ctx states.Context) error {
	ctx.StateMachine.PopState()
	world := game.World{
		Players: []game.Player{
			&game.LocalPlayer{},
		},
	}
	w, h := ebiten.WindowSize()
	spawner := game.CreateSpawner(float64(w/4), float64(h/4))
	world.AddActor(spawner)
	ctx.StateMachine.PushState(&world)
	return nil
}

func (s *BulletBox) Draw(screen *ebiten.Image) {
}
