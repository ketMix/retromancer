package game

import (
	"ebijam23/states"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type World struct {
	players     []Player
	tick        int // tick represents the current processed game tick. This is used to lockstep the players.
	ebitenTicks int // Elapsed ebiten ticks.
}

func (s *World) Init(ctx states.Context) error {
	return nil
}

func (s *World) Update(ctx states.Context) error {
	s.ebitenTicks++
	readyCount := 0
	for _, player := range s.players {
		if player.Ready(s.tick + 1) {
			readyCount++
		}
	}
	if s.ebitenTicks >= 3 { // Basically tick every 3 ebiten ticks.
		if readyCount == len(s.players) {
			s.tick++
			// Process the world!!!
			fmt.Println(s.tick)

			// Send our tick to the other player!
		}
		s.ebitenTicks = 0
	}

	return nil
}

func (s *World) Draw(screen *ebiten.Image) {
}
