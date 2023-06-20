package game

import (
	"ebijam23/resources"
	"ebijam23/states"
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
)

type World struct {
	Players     []Player // Exposed so singleplayer/multiplayer can set it.
	tick        int      // tick represents the current processed game tick. This is used to lockstep the players.
	ebitenTicks int      // Elapsed ebiten ticks.
	StartingMap string
	activeMap   *Map
}

func (s *World) Init(ctx states.Context) error {
	// Disable the global cursor.
	ctx.Cursor.Disable()

	// Create actors for our players.
	for _, p := range s.Players {
		pc := s.NewPC(ctx)

		// TODO: Read this in from the player's desired hat. Also, the hats should be dynamically built from the any hat- prefixed file in the manager's images.
		hats := []string{"hat-ebiten", "hat-wizard", "hat-gopher", "hat-tux", "hat-max", "hat-pep"}
		pc.Hat = resources.NewSprite(ctx.Manager.GetAs("images", hats[rand.Int31n(int32(len(hats)))], (*ebiten.Image)(nil)).(*ebiten.Image))

		p.SetActor(pc)
	}

	// Travel to the starting map.
	s.TravelToMap(ctx, s.StartingMap)

	return nil
}

func (s *World) Update(ctx states.Context) error {
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
			}

			// Process the world!!!
			var actorActions []ActorActions
			for _, actor := range s.activeMap.actors {
				actorActions = append(actorActions, ActorActions{
					Actor:   actor,
					Actions: actor.Update(),
				})
			}

			// Process bulleets
			var bulletActions []BulletActions
			for _, b := range s.activeMap.bullets {
				bulletActions = append(bulletActions, BulletActions{
					Bullet:  b,
					Actions: b.Update(),
				})
			}

			// Okay, this is very likely overkill to process actions entirely separately, but whatever.
			for _, actorAction := range actorActions {
				actor := actorAction.Actor
				deflecting := false
				reflecting := false
				for _, action := range actorAction.Actions {
					switch action := action.(type) {
					case ActionMove:
						checkShape := actor.Shape().Clone().(*CircleShape)
						checkShape.X = action.X
						checkShape.Y = action.Y + 4 // Stupid -4 to make the visual offset look nicer when bumpin' walls
						if !s.activeMap.Collides(checkShape) {
							actor.SetXY(action.X, action.Y)
						}
					case ActionReflect:
						reflecting = true
						x, y, _, _ := actor.Bounds()
						if !s.activeMap.DoesLineCollide(x, y, action.X, action.Y, s.activeMap.currentZ) {
							bullets := s.IntersectingBullets(&CircleShape{
								X:      action.X,
								Y:      action.Y,
								Radius: 20,
							})
							for _, bullet := range bullets {
								bullet.Reflect()
							}
						}
					case ActionDeflect:
						deflecting = true
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
						s.activeMap.bullets = append(s.activeMap.bullets, action.Bullets...)
					}
				}
				if a, ok := actor.(*PC); ok {
					// FIXME: Probably only SetImage if image is not the expected one.
					if deflecting {
						a.Hand.Sprite.SetImage(ctx.Manager.GetAs("images", "hand-deflect", (*ebiten.Image)(nil)).(*ebiten.Image))
					} else if reflecting {
						a.Hand.Sprite.SetImage(ctx.Manager.GetAs("images", "hand-reflect", (*ebiten.Image)(nil)).(*ebiten.Image))
					} else {
						a.Hand.Sprite.SetImage(ctx.Manager.GetAs("images", "hand-normal", (*ebiten.Image)(nil)).(*ebiten.Image))
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
						for _, actor := range s.activeMap.actors {
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
			for _, bullet := range s.activeMap.bullets {
				// Check for bullet collisions with walls.
				if s.activeMap.Collides(&bullet.Shape) {
					bullet.Destroyed = true
					continue
				}
				for _, actor := range s.activeMap.actors {
					if _, ok := actor.(*PC); ok {
						if bullet.Shape.Collides(actor.Shape()) {
							bullet.Destroyed = true
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
	s.activeMap.Draw(screen)

	for _, p := range s.Players {
		//y := screen.Bounds().Max.Y - 100
		if _, ok := p.(*LocalPlayer); !ok {
			continue
		}
		if a, ok := p.Actor().(*PC); ok {
			// Draw the hand's current energy.
			resources.DrawArc(screen, a.Hand.Shape.X, a.Hand.Shape.Y, 12, 0, 2*math.Pi*float64(a.Energy)/float64(a.MaxEnergy), color.RGBA{0xa0, 0x20, 0xf0, 0xaa})

			// Old energy bar
			/*w := 100
			h := 5
			x := screen.Bounds().Max.X/2 - w/2
			vector.StrokeRect(screen, float32(x), float32(y), float32(w), float32(h), 1, color.White, false)
			w2 := int(float32(w-3) * (float32(a.Energy) / float32(a.MaxEnergy)))
			vector.DrawFilledRect(screen, float32(x+1), float32(y+1), float32(w2), float32(h-3), color.White, false)

			y += h + 5*/
		}
	}
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
