package game

import (
	"ebijam23/resources"
	"ebijam23/states"
	"math"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
)

type World struct {
	overlay     Overlay
	Players     []Player // Exposed so singleplayer/multiplayer can set it.
	tick        int      // tick represents the current processed game tick. This is used to lockstep the players.
	ebitenTicks int      // Elapsed ebiten ticks.
	StartingMap string
	ShowHints   bool
	currentHint string
	lastHint    string
	hintTicks   int
	hints       Hints
	activeMap   *Map
	states      []WorldState
}

func (s *World) PushState(state WorldState, ctx states.Context) {
	s.states = append(s.states, state)
	state.Enter(s, ctx)
}

func (s *World) PopState(ctx states.Context) {
	if len(s.states) == 0 {
		return
	}
	s.states[len(s.states)-1].Leave(s, ctx)
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
	for i, p := range s.Players {
		if i == 0 {
			pc := s.NewPC(ctx)

			pc.Hat = resources.NewSprite(ctx.Manager.GetAs("images", p.Hat(), (*ebiten.Image)(nil)).(*ebiten.Image))

			// If the starting map is not start, then set the player as resurrected.
			if s.StartingMap != "start" {
				pc.resurrected = true
			}

			p.SetActor(pc)
		} else {
			c := s.NewCompanion(ctx)

			c.Hat = resources.NewSprite(ctx.Manager.GetAs("images", p.Hat(), (*ebiten.Image)(nil)).(*ebiten.Image))

			p.SetActor(c)
		}
	}

	// Init the hints.
	s.hints.active = s.ShowHints
	s.hints.ticker = -60
	// TODO: Read these from a hints file.
	s.hints.hintGroup = make(map[string]HintGroup)
	prefix := ""
	if len(s.Players) > 1 {
		prefix = ctx.L("Player 1:")
	}
	s.hints.AddHintGroup("p1-controller-start", HintGroup{
		Prefix: prefix,
		Items:  []string{"p1-controller-hint-1", "p1-controller-hint-2", "p1-controller-hint-3", "p1-controller-hint-4"},
	})
	s.hints.AddHintGroup("p1-keyboard-start", HintGroup{
		Prefix: prefix,
		Items:  []string{"p1-keyboard-hint-1", "p1-keyboard-hint-2", "p1-keyboard-hint-3"},
	})
	s.hints.AddHintGroup("p2-controller-start", HintGroup{
		Prefix:  ctx.L("Player 2:"),
		OffsetY: 32,
		Items:   []string{"p2-controller-hint-1", "p2-controller-hint-2", "p2-controller-hint-3", "p2-controller-hint-4"},
	})
	s.hints.AddHintGroup("p2-keyboard-start", HintGroup{
		Prefix:  ctx.L("Player 2:"),
		OffsetY: 32,
		Items:   []string{"p2-keyboard-hint-1", "p2-keyboard-hint-2"},
	})
	if s.StartingMap == "start" {
		for _, p := range s.Players {
			if pl, ok := p.(*LocalPlayer); ok {
				if _, ok := pl.actor.(*PC); ok {
					if pl.GamepadID != -1 {
						s.hints.ActivateGroup("p1-controller-start")
					} else {
						s.hints.ActivateGroup("p1-keyboard-start")
					}
				} else {
					if pl.GamepadID != -1 {
						s.hints.ActivateGroup("p2-controller-start")
					} else {
						s.hints.ActivateGroup("p2-keyboard-start")
					}
				}
			}
		}
	}

	// Set our starting state.
	if len(s.states) == 0 {
		s.PushState(&WorldStateBegin{}, ctx)
	}

	// Travel to the starting map.
	s.TravelToMap(ctx, s.StartingMap)

	return nil
}

func (s *World) Finalize(ctx states.Context) error {
	// Renable the global cursor.
	ctx.Cursor.Enable()
	return nil
}

func (s *World) Enter(ctx states.Context) error {
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
					for _, a := range s.activeMap.interactives {
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
					// This feels bad, but show the hover sprite over the player if they haven't resurrected yet.
					if !pc.resurrected {
						if pc.shape.Collides(&CircleShape{
							X:      pc.Hand.Shape.X,
							Y:      pc.Hand.Shape.Y,
							Radius: 20,
						}) {
							hoveringInteractable = true
						}
					}
					pc.Hand.HoverSprite.Hidden = !hoveringInteractable
				}
			}

			// Process the world!!!
			s.CurrentState().Tick(s, ctx)
			s.HandleTrash()
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
	newActors := s.activeMap.actors[:0]
	for _, a := range s.activeMap.actors {
		if !a.Destroyed() {
			newActors = append(newActors, a)
		}
	}
	s.activeMap.actors = newActors

	newBullets := make([]*Bullet, 0)
	for _, b := range s.activeMap.bullets {
		if !s.activeMap.OutOfBounds(&b.Shape) && !b.Destroyed {
			newBullets = append(newBullets, b)
		}
	}
	s.activeMap.bullets = newBullets

	newParticles := s.activeMap.particles[:0]
	for _, p := range s.activeMap.particles {
		if !p.Dead() {
			newParticles = append(newParticles, p)
		}
	}
	s.activeMap.particles = newParticles
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

func (s *World) FindNearestActor(from Shape, a Actor) Actor {
	var closestActor Actor
	var closestDistance float64
	for _, actor := range s.activeMap.actors {
		// Skip dead actors.
		if actor.Dead() {
			continue
		}
		// Reflect isn't great to use here, but it beats nested type switches.
		if reflect.TypeOf(actor) == reflect.TypeOf(a) {
			x, y, _, _ := actor.Bounds()
			distance := 0.0
			if circle, ok := from.(*CircleShape); ok {
				distance = math.Sqrt(math.Pow(circle.X-x, 2) + math.Pow(circle.Y-y, 2))
			} else if rect, ok := from.(*RectangleShape); ok {
				distance = math.Sqrt(math.Pow(rect.X-x, 2) + math.Pow(rect.Y-y, 2))
			}
			if closestActor == nil || distance < closestDistance {
				closestActor = actor
				closestDistance = distance
			}
		}
	}
	return closestActor
}
