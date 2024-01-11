package game

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"

	"github.com/ketMix/retromancer/net"
	"github.com/ketMix/retromancer/resources"

	"github.com/ketMix/retromancer/states"

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
	Net         net.ServerClient
	Seed        int64
	savedNPCs   map[string]bool
	Difficulty  *states.Difficulty
}

var (
	rng *rand.Rand
)

func (s *World) PushState(state WorldState, ctx states.Context) {
	// Mmmmm
	if s.Difficulty != nil {
		ctx.Difficulty = *s.Difficulty
	}
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

	// Initialize the the game package specific RNG with the passed in sneed.
	rng = rand.New(rand.NewSource(s.Seed))

	s.savedNPCs = make(map[string]bool)

	// Create actors for our players.
	for i, p := range s.Players {
		if (!s.Net.Running && i == 0) || (s.Net.Hosting && i == 0) || (!s.Net.Hosting && s.Net.Running && i == 1) {
			pc := s.NewPC(ctx)

			pc.Hat = resources.NewSprite(ctx.R.GetAs("images", p.Hat(), (*ebiten.Image)(nil)).(*ebiten.Image))

			// If the starting map is not start, then set the player as resurrected.
			if s.StartingMap != "start" {
				pc.resurrected = true
			}

			p.SetActor(pc)
		} else {
			c := s.NewCompanion(ctx)

			c.Hat = resources.NewSprite(ctx.R.GetAs("images", p.Hat(), (*ebiten.Image)(nil)).(*ebiten.Image))

			p.SetActor(c)
		}
	}

	// Init the hints.
	s.hints.active = s.ShowHints
	s.hints.ticker = -60
	// TODO: Read these from a hints file.
	s.hints.hintGroup = make(map[string]HintGroup)
	prefix := ""
	offset := 0.0
	if len(s.Players) > 1 {
		prefix = ctx.L.Get("Player 1:")
		offset = 20
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
		Prefix:  ctx.L.Get("Player 2:"),
		OffsetY: offset,
		Items:   []string{"p2-controller-hint-1", "p2-controller-hint-2", "p2-controller-hint-3", "p2-controller-hint-4"},
	})
	s.hints.AddHintGroup("p2-keyboard-start", HintGroup{
		Prefix:  ctx.L.Get("Player 2:"),
		OffsetY: offset,
		Items:   []string{"p2-keyboard-hint-1", "p2-keyboard-hint-2"},
	})
	s.hints.AddHintGroup("p1-controller-deflect", HintGroup{
		Prefix: prefix,
		Items:  []string{"p1-controller-hint-deflect"},
	})
	s.hints.AddHintGroup("p1-keyboard-deflect", HintGroup{
		Prefix: prefix,
		Items:  []string{"p1-keyboard-hint-deflect"},
	})
	s.hints.AddHintGroup("p1-controller-shield", HintGroup{
		Prefix: prefix,
		Items:  []string{"p1-controller-hint-shield"},
	})
	s.hints.AddHintGroup("p1-keyboard-shield", HintGroup{
		Prefix: prefix,
		Items:  []string{"p1-keyboard-hint-shield"},
	})

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

func (s *World) Enter(ctx states.Context, v interface{}) error {
	return nil
}

func (s *World) Update(ctx states.Context) error {
	if s.Difficulty != nil {
		ctx.Difficulty = *s.Difficulty
	}
	s.overlay.Update(ctx)

	if s.Net.Running {
		select {
		case ev := <-s.Net.EventChan:
			switch e := ev.(type) {
			case net.EventMessage:
				player := s.PlayerFromPeer(e.Peer)
				if player != nil {
					switch msg := e.Message.(type) {
					case Thoughts:
						player.thoughts = msg
					case TickState:
						if s.tick == player.lastTick+1 {
							player.QueueImpulses(msg.Impulses)
							player.ClearImpulses()
							player.lastTick++
						}
					}
				}
			default:
				fmt.Println("uhoh", e)
			}
		default:
		}
		// Send our local thoughts?
		if local, ok := s.Players[0].(*LocalPlayer); ok {
			if local.hasNewThoughts {
				local.hasNewThoughts = false
				for _, p := range s.Players {
					if remote, ok := p.(*RemotePlayer); ok {
						remote.peer.Send(local.Thoughts())
					}
				}
			}
		}
	}

	s.ebitenTicks++
	readyCount := 0
	for _, player := range s.Players {
		player.Update()
		if player.Ready(s.tick + 1) {
			readyCount++
		}
	}
	if s.ebitenTicks >= 2 { // Basically tick every 3 ebiten ticks.
		if readyCount == len(s.Players) {
			//fmt.Println("now ticking", s.tick)
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

			// Queue up the local player's impulses for the next tick!
			for _, player := range s.Players {
				if _, ok := player.(*LocalPlayer); ok {
					player.QueueImpulses(player.Impulses())
					for _, p := range s.Players {
						if remote, ok := p.(*RemotePlayer); ok {
							remote.peer.Send(TickState{
								Impulses: player.Impulses(),
							})
						}
					}
					player.ClearImpulses()
					// TODO: Send network message to peers with our impulses!
				}
			}

			s.HandleTrash()
			s.tick++
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
		for _, t := range p.Thoughts().Thoughts {
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
