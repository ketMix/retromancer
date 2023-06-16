package game

import (
	"ebijam23/states"

	"github.com/hajimehoshi/ebiten/v2"
)

type World struct {
}

func (s *World) Init(ctx states.Context) error {
	return nil
}

func (s *World) Update(ctx states.Context) error {
	return nil
}

func (s *World) Draw(screen *ebiten.Image) {
}
