package game

import (
	"ebijam23/resources"
	"ebijam23/states"
	"fmt"
	"image/color"
	"math"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type World struct {
	Players     []Player // Exposed so singleplayer/multiplayer can set it.
	tick        int      // tick represents the current processed game tick. This is used to lockstep the players.
	ebitenTicks int      // Elapsed ebiten ticks.
	actors      []Actor
	bullets     []*Bullet
}

func (s *World) Init(ctx states.Context) error {
	// Disable the global cursor.
	ctx.Cursor.Disable()

	// Create actors for our players.
	for _, p := range s.Players {
		// TODO: Move to a NewPC ctor
		pc := &PC{
			Sprite:            resources.NewSprite(ctx.Manager.GetAs("images", "player", (*ebiten.Image)(nil)).(*ebiten.Image)),
			Arrow:             resources.NewSprite(ctx.Manager.GetAs("images", "direction-arrow", (*ebiten.Image)(nil)).(*ebiten.Image)),
			Energy:            0,
			MaxEnergy:         100,
			EnergyRestoreRate: 1,
		}
		pc.Sprite.Interpolate = true
		pc.Sprite.Centered = true
		pc.Hand.Sprite = resources.NewSprite(ctx.Manager.GetAs("images", "hand-normal", (*ebiten.Image)(nil)).(*ebiten.Image))
		pc.Hand.Sprite.Centered = true
		p.SetActor(pc)
		// Add to the world. FIXME: This should be done in some sort of sub-game state.
		s.actors = append(s.actors, p.Actor())
	}

	return nil
}

func (s *World) Update(ctx states.Context) error {
	s.ebitenTicks++
	readyCount := 0
	for _, player := range s.Players {
		// Passing context to players seems a bit of a violation.
		player.Update(ctx)
		if player.Ready(s.tick + 1) {
			readyCount++
		}
	}
	if s.ebitenTicks >= 3 { // Basically tick every 3 ebiten ticks.
		if readyCount == len(s.Players) {
			s.tick++
			// Process the players' current tick think -- this also sends impulses to their respective actors.
			for _, player := range s.Players {
				player.Tick()
			}

			// Process the world!!!
			var actorActions []ActorActions
			for _, actor := range s.actors {
				actorActions = append(actorActions, ActorActions{
					Actor:   actor,
					Actions: actor.Update(),
				})
			}

			// Process bulleets
			var bulletActions []BulletActions
			for _, b := range s.bullets {
				bulletActions = append(bulletActions, BulletActions{
					Bullet:  b,
					Actions: b.Update(),
				})
			}

			// Okay, this is very likely overkill to process actions entirely separately, but whatever.
			for _, actorAction := range actorActions {
				actor := actorAction.Actor
				for _, action := range actorAction.Actions {
					switch action := action.(type) {
					case ActionMove:
						actor.SetXY(action.X, action.Y)
					case ActionReflect:
						bullets := s.IntersectingBullets(&CircleShape{
							X:      action.X,
							Y:      action.Y,
							Radius: 20,
						})
						for _, bullet := range bullets {
							bullet.Reflect()
						}
					case ActionDeflect:
						x, y, _, _ := actor.Bounds()
						bullets := s.IntersectingBullets(&CircleShape{
							X:      x,
							Y:      y,
							Radius: 20,
						})
						for _, bullet := range bullets {
							bullet.Deflect(action.Direction)
						}
					case ActionSpawnBullets:
						s.bullets = append(s.bullets, action.Bullets...)
					}
				}
			}
			// Even more overkill for the bullets.
			for _, bulletAction := range bulletActions {
				bullet := bulletAction.Bullet
				for _, action := range bulletAction.Actions {
					switch action := action.(type) {
					case ActionFindNearestActor:
						var closestActor Actor
						var closestDistance float64
						for _, actor := range s.actors {
							// Reflect isn't great to use here, but it beats nested type switches.
							if reflect.TypeOf(actor) == reflect.TypeOf(action.Actor) {
								x, y, _, _ := actor.Bounds()
								distance := math.Sqrt(math.Pow(bullet.Shape.X-x, 2) + math.Pow(bullet.Shape.Y-y, 2))
								if closestActor == nil || distance < closestDistance {
									closestActor = actor
									closestDistance = distance
								}
							}
						}
						bullet.TargetActor = closestActor
					}
				}
			}

			// Okay, this probably isn't great, but let's check bullet collisions here.
			for _, bullet := range s.bullets {
				for _, actor := range s.actors {
					if _, ok := actor.(*PC); ok {
						if bullet.Shape.Collides(actor.Shape()) {
							fmt.Println("TODO: bullet hit player!")
							break
						}
					}
				}
			}

			// Queue up the local player's impulses for the next tick!
			for _, player := range s.Players {
				if _, ok := player.(*LocalPlayer); ok {
					player.QueueImpulses(player.Impulses())
					player.ClearImpulses()
					// TODO: Send network message to peers with our impulses!
				}
			}

			s.HandleTrash()
		}
		s.ebitenTicks = 0
	}

	return nil
}

func (s *World) Draw(screen *ebiten.Image) {
	// Draw bullets first.
	for _, b := range s.bullets {
		b.Draw(screen)
	}
	for _, a := range s.actors {
		a.Draw(screen)
	}
	for _, p := range s.Players {
		y := screen.Bounds().Max.Y - 100
		if _, ok := p.(*LocalPlayer); !ok {
			continue
		}
		if a, ok := p.Actor().(*PC); ok {
			w := 100
			h := 5
			x := screen.Bounds().Max.X/2 - w/2
			vector.StrokeRect(screen, float32(x), float32(y), float32(w), float32(h), 1, color.White, false)
			w2 := int(float32(w-3) * (float32(a.Energy) / float32(a.MaxEnergy)))
			vector.DrawFilledRect(screen, float32(x+1), float32(y+1), float32(w2), float32(h-3), color.White, false)

			y += h + 5
		}
	}
}

// TODO: This should be handled by map loading.
func (s *World) AddActor(a Actor) {
	s.actors = append(s.actors, a)
}

func (s *World) HandleTrash() {
	newBullets := make([]*Bullet, 0)
	for _, b := range s.bullets {
		if !b.OutOfBounds() {
			newBullets = append(newBullets, b)
		}
	}
	s.bullets = newBullets
}

func (s *World) IntersectingBullets(sh Shape) []*Bullet {
	var bullets []*Bullet
	for _, b := range s.bullets {
		if b.Shape.Collides(sh) {
			bullets = append(bullets, b)
		}
	}
	return bullets
}
