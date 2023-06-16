package game

import (
	"ebijam23/resources"
	"ebijam23/states"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type World struct {
	Players     []Player // Exposed so singleplayer/multiplayer can set it.
	tick        int      // tick represents the current processed game tick. This is used to lockstep the players.
	ebitenTicks int      // Elapsed ebiten ticks.
	actors      []Actor
}

func (s *World) Init(ctx states.Context) error {
	// Create actors for our players.
	for _, p := range s.Players {
		p.SetActor(&PC{
			Sprite: resources.NewSprite(ctx.Manager.GetAs("images", "player", (*ebiten.Image)(nil)).(*ebiten.Image)),
		})
		// Add to the world. FIXME: This should be done in some sort of sub-game state.
		s.actors = append(s.actors, p.Actor())
	}

	return nil
}

func (s *World) Update(ctx states.Context) error {
	s.ebitenTicks++
	readyCount := 0
	for _, player := range s.Players {
		if player.Ready(s.tick + 1) {
			readyCount++
		}
	}
	if s.ebitenTicks >= 3 { // Basically tick every 3 ebiten ticks.
		if readyCount == len(s.Players) {
			s.tick++
			// Process the world!!!
			for _, actor := range s.actors {
				actor.Update()
			}
			fmt.Println(s.tick)

			// Send our tick to the other player!
		}
		s.ebitenTicks = 0
	}

	return nil
}

func (s *World) Draw(screen *ebiten.Image) {
	for _, a := range s.actors {
		a.Draw(screen)
	}
}
