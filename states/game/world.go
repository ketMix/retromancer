package game

import (
	"ebijam23/resources"
	"ebijam23/states"
	"image/color"
	"math"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
)

type World struct {
	overlay     Overlay
	Players     []Player // Exposed so singleplayer/multiplayer can set it.
	tick        int      // tick represents the current processed game tick. This is used to lockstep the players.
	ebitenTicks int      // Elapsed ebiten ticks.
	StartingMap string
	activeMap   *Map
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
			var actorActions []ActorActions
			for _, actor := range s.activeMap.actors {
				// Don't process spawners if the map is cleared.
				if _, ok := actor.(*Spawner); ok {
					if s.activeMap.cleared {
						continue
					}
				}

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
				shielding := false
				for _, action := range actorAction.Actions {
					switch action := action.(type) {
					case ActionMove:
						checkShape := actor.Shape().Clone().(*CircleShape)
						checkShape.X = action.X
						checkShape.Y = action.Y + 4 // Stupid -4 to make the visual offset look nicer when bumpin' walls
						if collision := s.activeMap.Collides(checkShape); collision == nil || !collision.Cell.Blocks {
							actor.SetXY(action.X, action.Y)
						}
					case ActionReflect:
						reflecting = true
						x, y, _, _ := actor.Bounds()
						if !s.activeMap.DoesLineCollide(x, y, action.X, action.Y, s.activeMap.currentZ) {

							// Reflect bullets
							bullets := s.IntersectingBullets(&CircleShape{
								X:      action.X,
								Y:      action.Y,
								Radius: 20,
							})
							for _, b := range bullets {
								b.Reflect()
							}

							// Reverse actors
							actors := s.IntersectingActors(&CircleShape{
								X:      action.X,
								Y:      action.Y,
								Radius: 20,
							})
							// Reverse interactive actors.
							for _, a := range actors {
								if a, ok := a.(*Interactive); ok {
									a.Reverse()
								}
							}
						}
					case ActionDeflect:
						deflecting = true
						x, y, _, _ := actor.Bounds()
						if !s.activeMap.DoesLineCollide(x, y, action.X, action.Y, s.activeMap.currentZ) {
							bullets := s.IntersectingBullets(&CircleShape{
								X:      action.X,
								Y:      action.Y,
								Radius: 20,
							})
							for _, bullet := range bullets {
								bullet.Deflect(action.Direction)
							}
						}
					case ActionShield:
						shielding = true
						x, y, _, _ := actor.Bounds()
						bullets := s.IntersectingBullets(&CircleShape{
							X:      x,
							Y:      y,
							Radius: 20,
						})
						for _, bullet := range bullets {
							bullet.holdFor = 30
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
					} else if shielding {
						a.Hand.Sprite.SetImage(ctx.Manager.GetAs("images", "hand-shield", (*ebiten.Image)(nil)).(*ebiten.Image))
					} else {
						a.Hand.Sprite.SetImage(ctx.Manager.GetAs("images", "hand-normal", (*ebiten.Image)(nil)).(*ebiten.Image))
					}
				}
			}
			// Even more overkill for the bullets.
			for _, bulletAction := range bulletActions {
				bullet := bulletAction.Bullet
				// FIXME: This is dumb.
				if bullet.Destroyed {
					continue
				}
				for _, action := range bulletAction.Actions {
					switch action := action.(type) {
					case ActionFindNearestActor:
						var closestActor Actor
						var closestDistance float64
						for _, actor := range s.activeMap.actors {
							// Skip dead actors.
							if actor.Dead() {
								continue
							}
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
				if collision := s.activeMap.Collides(&bullet.Shape); collision != nil && collision.Cell.Blocks {
					bullet.Destroyed = true
					continue
				}
				for _, actor := range s.activeMap.actors {
					if p, ok := actor.(*PC); ok {
						if p.InvulnerableTicks > 0 {
							continue
						}
						if bullet.Shape.Collides(actor.Shape()) {
							bullet.Destroyed = true
							p.InvulnerableTicks = 40
							p.Lives--
							break
						}
					}
				}
			}

			// Check the our interactive actor conditions
			interactiveActors := s.activeMap.GetInteractiveActors()
			for _, actor := range interactiveActors {
				conditions := actor.Conditions()
				for _, condition := range conditions {
					args := condition.Args
					switch condition.Type {
					case resources.Active:
						if CheckActiveCondition(args, interactiveActors) {
							actor.IncreaseActivation(nil)
							cell, _ := s.activeMap.FindCellById(actor.ID())
							cell.Blocks = false // No
						}
					}
				}
			}

			// Check our map conditions if not yet cleared
			if !s.activeMap.cleared {
				conditions := s.activeMap.conditions
				for _, condition := range conditions {
					args := condition.Args
					switch condition.Type {
					case resources.Active:
						if CheckActiveCondition(args, interactiveActors) {
							s.activeMap.cleared = true
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

			if s.ArePlayersDead() {
				if s.DoPlayersShareThought(ResetThought{}) {
					s.ResetActiveMap(ctx)
				} else if s.DoPlayersShareThought(QuitThought{}) {
					ctx.StateMachine.PopState()
				}
			}
		}
		s.ebitenTicks = 0
	}

	return nil
}

func (s *World) Draw(ctx states.DrawContext) {
	s.activeMap.Draw(ctx.Screen)

	for _, p := range s.Players {
		//y := screen.Bounds().Max.Y - 100
		if _, ok := p.(*LocalPlayer); !ok {
			continue
		}
		if a, ok := p.Actor().(*PC); ok {
			// Draw the hand's current energy.
			resources.DrawArc(ctx.Screen, a.Hand.Shape.X, a.Hand.Shape.Y, 12, 0, 2*math.Pi*float64(a.Energy)/float64(a.MaxEnergy), color.RGBA{0xa0, 0x20, 0xf0, 0xaa})
			// Also draw the energy around the player if they shielded.
			if _, ok := a.previousInteraction.(ActionShield); ok {
				resources.DrawArc(ctx.Screen, a.shape.X, a.shape.Y, 12, 0, 2*math.Pi*float64(a.Energy)/float64(a.MaxEnergy), color.RGBA{0xa0, 0x20, 0xf0, 0xaa})
			}

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

	if s.ArePlayersDead() {
		// :)
		poorMansOutline := func(text string, x, y int, scale int) {
			ctx.Text.Draw(ctx.Screen, text, x-scale, y)
			ctx.Text.Draw(ctx.Screen, text, x+scale, y)
			ctx.Text.Draw(ctx.Screen, text, x, y-scale)
			ctx.Text.Draw(ctx.Screen, text, x, y+scale)
		}

		ctx.Text.SetAlign(etxt.YCenter | etxt.XCenter)
		x := ctx.Screen.Bounds().Max.X / 2
		y := float64(ctx.Screen.Bounds().Max.Y / 2)
		y -= ctx.Text.Utils().GetLineHeight() / 2
		// Death title
		{
			ctx.Text.SetScale(2.0)
			ctx.Text.SetColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
			poorMansOutline("Morte", x, int(y), 2)
			ctx.Text.SetColor(color.RGBA{0xff, 0x00, 0x00, 0xff})
			ctx.Text.Draw(ctx.Screen, "Morte", x, int(y))
		}
		y += ctx.Text.Utils().GetLineHeight()
		// Subtitles
		{
			ctx.Text.SetScale(1.0)
			ctx.Text.SetColor(color.Black)
			poorMansOutline("<Enter> to restart room", x, int(y), 1)
			ctx.Text.SetColor(color.White)
			ctx.Text.Draw(ctx.Screen, "<Enter> to restart room", x, int(y))
			y += ctx.Text.Utils().GetLineHeight()
			ctx.Text.SetColor(color.Black)
			poorMansOutline("<Escape> to quit", x, int(y), 1)
			ctx.Text.SetColor(color.White)
			ctx.Text.Draw(ctx.Screen, "<Escape> to quit", ctx.Screen.Bounds().Max.X/2, int(y))
		}
	}

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
