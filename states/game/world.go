package game

import (
	"ebijam23/resources"
	"ebijam23/states"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
)

type World struct {
	overlay     Overlay
	Players     []Player // Exposed so singleplayer/multiplayer can set it.
	tick        int      // tick represents the current processed game tick. This is used to lockstep the players.
	ebitenTicks int      // Elapsed ebiten ticks.
	StartingMap string
	activeMap   *Map
	states      []WorldState
}

func (s *World) PushState(state WorldState) {
	s.states = append(s.states, state)
}

func (s *World) PopState() {
	s.states = s.states[:len(s.states)-1]
}

func (s *World) CurrentState() WorldState {
	return s.states[len(s.states)-1]
}

func (s *World) Init(ctx states.Context) error {
	// Disable the global cursor.
	ctx.Cursor.Disable()

	// Init the overlay.
	if err := s.overlay.Init(ctx); err != nil {
		return err
	}

	// Create actors for our players.
	for _, p := range s.Players {
		pc := s.NewPC(ctx)

		pc.Hat = resources.NewSprite(ctx.Manager.GetAs("images", p.Hat(), (*ebiten.Image)(nil)).(*ebiten.Image))

		p.SetActor(pc)
	}

	// Set our starting state.
	s.PushState(&WorldStateLive{})

	// Travel to the starting map.
	s.TravelToMap(ctx, s.StartingMap)

	return nil
}

func (s *World) Finalize(ctx states.Context) error {
	// Renable the global cursor.
	ctx.Cursor.Enable()
	return nil
}

func (s *World) Update(ctx states.Context) error {
	s.overlay.Update(ctx)

	s.ebitenTicks++
	readyCount := 0
	for _, player := range s.Players {
		// Passing context to players seems a bit of a violation.
		player.Update()
		if player.Ready(s.tick + 1) {
			readyCount++
		}
	}
	if s.ebitenTicks >= 2 { // Basically tick every 3 ebiten ticks.
		if readyCount == len(s.Players) {
			s.tick++
			// Process the players' current tick think -- this also sends impulses to their respective actors.
			for _, player := range s.Players {
				player.Tick()

				if pc, ok := player.Actor().(*PC); ok {
					hoveringInteractable := false
					for _, a := range s.activeMap.GetInteractiveActors() {
						if !a.Reverseable() {
							continue
						}
						if a.shape.Collides(&CircleShape{
							X:      pc.Hand.Shape.X,
							Y:      pc.Hand.Shape.Y,
							Radius: 20,
						}) {
							hoveringInteractable = true
							break
						}
					}
					pc.Hand.HoverSprite.Hidden = !hoveringInteractable
				}
			}

			// Process the world!!!
			s.CurrentState().Tick(s, ctx)
		}
		s.ebitenTicks = 0
	}

	return nil
}

func (s *World) Draw(ctx states.DrawContext) {
	s.CurrentState().Draw(s, ctx)
	s.overlay.Draw(ctx)
}

func (s *World) ArePlayersDead() bool {
	for _, p := range s.Players {
		if a, ok := p.Actor().(*PC); ok {
			if !a.Dead() {
				return false
			}
		}
	}
	return true
}

func (s *World) DoPlayersShareThought(thought Thought) bool {
	for _, p := range s.Players {
		match := false
		for _, t := range p.Thoughts() {
			// More reflection, woo.
			if reflect.TypeOf(t) == reflect.TypeOf(thought) {
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}
	return true
}

func (s *World) HandleTrash() {
	newBullets := make([]*Bullet, 0)
	for _, b := range s.activeMap.bullets {
		if !b.OutOfBounds() && !b.Destroyed {
			newBullets = append(newBullets, b)
		}
	}
	s.activeMap.bullets = newBullets
}

func (s *World) IntersectingBullets(sh Shape) []*Bullet {
	var bullets []*Bullet
	for _, b := range s.activeMap.bullets {
		if b.Shape.Collides(sh) {
			bullets = append(bullets, b)
		}
	}
	return bullets
}

func (s *World) IntersectingActors(sh Shape) []Actor {
	var actors []Actor
	for _, a := range s.activeMap.actors {
		if a.Shape().Collides(sh) {
			actors = append(actors, a)
		}
	}
	return actors
}
