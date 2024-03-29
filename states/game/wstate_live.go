package game

import (
	"image/color"
	"math"
	"reflect"
	"strings"

	"github.com/ketMix/retromancer/states"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ketMix/retromancer/resources"
)

type WorldStateLive struct {
	signText  *string
	isSignNPC bool
}

func (w *WorldStateLive) Enter(s *World, ctx states.Context) {
}

func (w *WorldStateLive) Leave(s *World, ctx states.Context) {
}

func (w *WorldStateLive) Tick(s *World, ctx states.Context) {
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

	// Process particles
	for _, p := range s.activeMap.particles {
		p.Update()
	}

	// Okay, this is very likely overkill to process actions entirely separately, but whatever.
	for _, actorAction := range actorActions {
		actor := actorAction.Actor
		deflecting := false
		reversing := false
		shielding := false
		for _, action := range actorAction.Actions {
			switch action := action.(type) {
			case ActionMove:
				var checkShape Shape
				if shape, ok := actor.Shape().(*CircleShape); ok {
					checkShape = shape.Clone()
					checkShape.(*CircleShape).X = action.X
					// Stupid -4 to make the visual offset look nicer when bumpin' walls
					checkShape.(*CircleShape).Y = action.Y - 4
				} else if shape, ok := actor.Shape().(*RectangleShape); ok {
					checkShape = shape.Clone()
					checkShape.(*RectangleShape).X = action.X
					// Stupid -4 to make the visual offset look nicer when bumpin' walls
					checkShape.(*RectangleShape).Y = action.Y - 4
				}
				if collision := s.activeMap.Collides(checkShape); collision == nil || !collision.Cell.blockMove {
					actor.SetXY(action.X, action.Y)
				}
				// forgive me.
				if s.tick%4 == 0 {
					if pc, ok := actor.(*PC); ok {
						s.SpawnParticle(ctx, "puff", action.X, action.Y+pc.Sprite.Height()/2-2, 0, 0, 10)
					}
				}
			case ActionReverse:
				reversing = true
				x, y, _, _ := actor.Bounds()
				// Adjust x and y by the direction the actor wishes to shoot.
				a := math.Atan2(action.Y-y, action.X-x)
				x += math.Cos(a) * 6
				y += math.Sin(a) * 6
				if !s.activeMap.DoesLineCollide(x, y, action.X, action.Y, s.activeMap.currentZ) {

					// Reverse bullets
					bullets := s.IntersectingBullets(&CircleShape{
						X:      action.X,
						Y:      action.Y,
						Radius: 20,
					})
					for _, b := range bullets {
						b.Reverse()
					}

					// Reverse actors
					actors := s.IntersectingActors(&CircleShape{
						X:      action.X,
						Y:      action.Y,
						Radius: 20,
					})
					// Reverse interactive actors.
					for _, actor := range actors {
						if a, ok := actor.(*Interactive); ok {
							if a.Reverseable() {
								a.Reverse()
								s.SpawnParticle(ctx, "reverse", action.X, action.Y, rng.Float64()*math.Pi*2, rng.Float64()*2.0, 30)

								// Store saved NPCs for a later tally.
								if a.active && a.npc {
									s.savedNPCs[a.text] = true
								}

								if a.active {
									// FIXME: This isn't the right place for this, as it would be best if the interactive actor created an action containing its VFX remove, but this is the most obvious place.
									for _, v := range a.removeVFX {
										if v == "darkness" {
											for _, vfx := range s.activeMap.vfx.Items() {
												if v2, ok := vfx.(*resources.Darkness); ok {
													v2.Fade = true
												}
											}
										} else {
											s.activeMap.vfx.RemoveByID(v)
										}
									}
								}
							}
						} else if pc, ok := actor.(*PC); ok {
							// More hackiness. Let us resurrect ourself if we haven't been already.
							if !pc.resurrected {
								pc.resurrected = true
							}
						}
					}
				}
			case ActionDeflect:
				deflecting = true
				x, y, _, _ := actor.Bounds()
				// Adjust x and y by the direction the actor wishes to shoot.
				a := math.Atan2(action.Y-y, action.X-x)
				x += math.Cos(a) * 6
				y += math.Sin(a) * 6
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
					bullet.TargetActor = nil
					bullet.aimTime = 0
				}
			case ActionSpawnBullets:
				s.activeMap.bullets = append(s.activeMap.bullets, action.Bullets...)
			case ActionSpawnParticle:
				s.SpawnParticle(ctx, action.Img, action.X, action.Y, action.Angle, action.Speed, action.Life)
			case ActionSpawnEnemy:
				e := CreateEnemy(ctx, action.ID, action.Name)
				e.SetXY(action.X, action.Y)
				s.activeMap.actors = append(s.activeMap.actors, e)
				s.activeMap.enemies = append(s.activeMap.enemies, e)
			case ActionFindNearestActor:
				if e, ok := actor.(*Enemy); ok {
					target := s.FindNearestActor(&e.shape, action.Actor)
					if target != nil {
						tx, ty, _, _ := target.Bounds()
						if !s.activeMap.DoesLineCollide(e.shape.X, e.shape.Y, tx, ty, s.activeMap.currentZ) {
							e.SetTarget(target)
						}
					}
				}
			}
		}
		if a, ok := actor.(*PC); ok {
			a.shielding = shielding
			// FIXME: Probably only SetImage if image is not the expected one.
			if deflecting {
				a.Hand.Sprite.SetImage(ctx.R.GetAs("images", "hand-deflect", (*ebiten.Image)(nil)).(*ebiten.Image))
			} else if reversing {
				a.Hand.Sprite.SetImage(ctx.R.GetAs("images", "hand-reverse", (*ebiten.Image)(nil)).(*ebiten.Image))
			} else if shielding {
				a.Hand.Sprite.SetImage(ctx.R.GetAs("images", "hand-shield", (*ebiten.Image)(nil)).(*ebiten.Image))
			} else {
				a.Hand.Sprite.SetImage(ctx.R.GetAs("images", "hand-normal", (*ebiten.Image)(nil)).(*ebiten.Image))
			}
			// Play the associated audio
			a.PlayAudio(deflecting, reversing, shielding)
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
			case ActionSpawnParticle:
				s.SpawnParticle(ctx, action.Img, action.X, action.Y, action.Angle, action.Speed, action.Life)
			}
		}
	}

	// Okay, this probably isn't great, but let's check bullet collisions here.
	for _, bullet := range s.activeMap.bullets {
		// Check for bullet collisions with actors.
		for _, actor := range s.activeMap.actors {
			// Check player collisions.
			if !bullet.friendly {
				if p, ok := actor.(*PC); ok {
					// Prevent taking damage while shielding and while having invuln ticks.
					if p.InvulnerableTicks > 0 || p.shielding {
						if bullet.Shape.Collides(actor.Shape()) {
							bullet.TargetActor = nil
							bullet.aimTime = 0
						}
						continue
					}
					if bullet.Shape.Collides(actor.Shape()) {
						x, y, _, _ := actor.Bounds()
						for i := 0; i < 6; i++ {
							s.SpawnParticle(ctx, "hurt", x, y, bullet.Angle-math.Pi/4+(math.Pi/2*rng.Float64()), rng.Float64()*2.0, 30)
						}
						bullet.Destroyed = true
						p.Hurtie()
						break
					}
					continue // skip checking other actors
				}
			}

			// Check interactive and enemy collisions.
			// Only applicable to reversed or deflected bullets
			if bullet.reversed || bullet.deflected || bullet.friendly {
				// If the interactive is shootable, hit the interactive.
				if i, ok := actor.(*Interactive); ok {
					if i.Shootable() && bullet.Shape.Collides(i.Shape()) {
						i.Hit()
						bullet.Destroyed = true
						break
					}
					continue // skip checking other actors
				}

				if e, ok := actor.(*Enemy); ok {
					if e.IsAlive() {
						if bullet.Shape.Collides(e.Shape()) {
							if e.Damage(bullet.Damage) {
								bullet.Destroyed = true
							}
							break
						}
					}
				}
			}

			// Check for bullet collisions with walls.
			if collision := s.activeMap.Collides(&bullet.Shape); collision != nil && collision.Cell.blockView {
				bullet.Destroyed = true
				continue
			}
		}
	}

	// Oh boy, yet another loop.
	// Check for collisions between player characters and interactives/snaggables.
	touchingSign := false
	for _, pl := range s.Players {
		if c, ok := pl.Actor().(*Companion); ok {
			// Allow companions to touch interactives.
			for _, actor := range s.activeMap.actors {
				if i, ok := actor.(*Interactive); ok {
					// If touchable, apply reverse to it
					if i.touchable && i.shape.Collides(pl.Actor().Shape()) {
						i.Reverse()
						continue
					}
				}
				if sn, ok := actor.(*Snaggable); ok {
					if sn.shape.Collides(pl.Actor().Shape()) {
						// Allow companion to snarf life.
						switch sn.spriteName {
						case "item-life":
							c.Snarf()
							sn.destroyed = true
						}
					}
				}
			}

		} else if pc, ok := pl.Actor().(*PC); ok {
			for _, actor := range s.activeMap.actors {
				// Check interactive collisions.
				if i, ok := actor.(*Interactive); ok {
					// If the interactive has text and is active, show the text.
					if !touchingSign && i.text != "" && i.Active() && i.shape.Collides(pl.Actor().Shape()) {
						localized := ctx.L.Get(i.text)
						w.signText = &localized
						w.isSignNPC = i.npc
						touchingSign = true
					}

					// If touchable, apply reverse to it
					if i.touchable && i.shape.Collides(pl.Actor().Shape()) {
						i.Reverse()
						continue
					}

					// If nextMap is defined and active, go to next map.
					if i.nextMap != nil && i.Active() && i.shape.Collides(pl.Actor().Shape()) {
						ctx.R.GetAs("sounds", "stairs", (*resources.Sound)(nil)).(*resources.Sound).Play(0.5)
						s.TravelToMap(ctx, *i.nextMap)
					}
				}

				// Check snaggable collisions.
				if sn, ok := actor.(*Snaggable); ok {
					if sn.shape.Collides(pl.Actor().Shape()) {
						// TODO: Bless the player with powers beyond their wildest imaginations (lives, power unlocks, etc.)
						switch sn.spriteName {
						case "item-life":
							if pc.Lives < playerMaxLives {
								sn.destroyed = true
								ctx.R.GetAs("sounds", "item", (*resources.Sound)(nil)).(*resources.Sound).Play(0.5)
								pc.Lives++
							}
						case "item-book":
							pc.HasDeflect = true
							sn.destroyed = true
							ctx.R.GetAs("sounds", "book", (*resources.Sound)(nil)).(*resources.Sound).Play(0.5)
							for _, p := range s.Players {
								if pl, ok := p.(*LocalPlayer); ok {
									if _, ok := pl.actor.(*PC); ok {
										if pl.GamepadID != -1 {
											s.hints.ActivateGroup("p1-controller-deflect")
										} else {
											s.hints.ActivateGroup("p1-keyboard-deflect")
										}
										s.hints.ticker = -30
										s.hints.active = true
									}
								}
							}
						case "item-shield":
							pc.HasShield = true
							sn.destroyed = true
							ctx.R.GetAs("sounds", "book", (*resources.Sound)(nil)).(*resources.Sound).Play(0.5)
							for _, p := range s.Players {
								if pl, ok := p.(*LocalPlayer); ok {
									if _, ok := pl.actor.(*PC); ok {
										if pl.GamepadID != -1 {
											s.hints.ActivateGroup("p1-controller-shield")
										} else {
											s.hints.ActivateGroup("p1-keyboard-shield")
										}
										s.hints.ticker = -20
										s.hints.active = true
									}
								}
							}

						}

						continue
					}
				}

				// Check enemy collisions.
				if pc.InvulnerableTicks <= 0 && !pc.shielding {
					if e, ok := actor.(*Enemy); ok {
						if !e.friendly && e.IsAlive() && e.Shape().Collides(pl.Actor().Shape()) {
							pc.Hurtie()
						}
					}
				}
			}
		}
	}
	if !touchingSign {
		w.signText = nil
	}

	interactives := s.activeMap.interactives
	enemies := s.activeMap.enemies

	// Check the our interactive actor conditions
	for _, actor := range interactives {
		if actor.id == "drawbridge" && actor.active { // Need to not run this every update and not tie to id, but drawbridge is kinda special
			cell := s.activeMap.FindCellById(actor.ID())
			cells := make([]*Cell, 0)
			if cell != nil {
				cells = append(cells, cell)
				cells = append(cells, s.activeMap.GetCell(int(cell.Shape.X/cellW), int(cell.Shape.Y/cellH)+2, s.activeMap.currentZ))
				cells = append(cells, s.activeMap.GetCell(int(cell.Shape.X/cellW)+1, int(cell.Shape.Y/cellH)+2, s.activeMap.currentZ))
			}
			for _, cell := range cells {
				if cell != nil {
					cell.blockMove = false // No
					cell.blockView = false // No
				}
			}
		}
		if !actor.active {
			if CheckConditions(actor.Conditions(), interactives, enemies) {
				actor.IncreaseActivation(nil)
			}
		} else {
			for _, v := range actor.removeVFX {
				if v == "darkness" {
					for _, vfx := range s.activeMap.vfx.Items() {
						if v2, ok := vfx.(*resources.Darkness); ok {
							v2.Fade = true
						}
					}
				} else {
					s.activeMap.vfx.RemoveByID(v)
				}
			}

			cell := s.activeMap.FindCellById(actor.ID())
			if cell != nil {
				cell.blockMove = false // No
				cell.blockView = false // No
			}
		}
	}

	// Bad check for final boss logic.
	if s.activeMap.filename == "3-boss" {
		for _, e := range enemies {
			if e.friendly {
				// Destroy all bullets once the lich b. dead.
				s.activeMap.bullets = make([]*Bullet, 0)
				if len(s.savedNPCs) < 14 { // In theory npcs + stump should = 14
					e.Damage(50)
					if !e.IsAlive() {
						e.friendly = false
					}
				}
			}
		}
	}

	// Check our map conditions if not yet cleared
	if !s.activeMap.cleared {
		if CheckConditions(s.activeMap.conditions, interactives, enemies) {
			s.activeMap.cleared = true
		}
	}

	// Show hints as needed.
	s.hints.Update(ctx)

	if s.activeMap.data.End {
		s.PopState(ctx)
		s.PushState(&WorldStateEnd{}, ctx)
	} else if s.ArePlayersDead() {
		s.PopState(ctx)
		s.PushState(&WorldStateDead{}, ctx)
	}
}

func (w *WorldStateLive) Draw(s *World, ctx states.DrawContext) {
	s.activeMap.Draw(ctx)

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
			} else if a, ok := p.Actor().(*Companion); ok {
				resources.DrawArc(ctx.Screen, a.Hand.Shape.X, a.Hand.Shape.Y, 12, 0, 2*math.Pi*float64(a.Energy)/float64(a.MaxEnergy), color.RGBA{0xa0, 0x20, 0xf0, 0xaa})
			}
		} else if a, ok := p.Actor().(*Companion); ok {
			resources.DrawArc(ctx.Screen, a.Hand.Shape.X, a.Hand.Shape.Y, 8, 0, 2*math.Pi*float64(a.Energy)/float64(a.MaxEnergy), color.RGBA{0xa0, 0x20, 0xf0, 0xaa})
		}
	}

	// Draw the sign text if it exists.
	if w.signText != nil {
		centerX := float32(ctx.Screen.Bounds().Max.X / 2)
		centerY := float32(ctx.Screen.Bounds().Max.Y - ctx.Screen.Bounds().Max.Y/4)
		boardSizeX := float32(ctx.Screen.Bounds().Max.X) * 0.6
		boardSizeY := float32(ctx.Screen.Bounds().Max.Y) * 0.3
		text := strings.Split(*w.signText, "\n")

		boardX := centerX - boardSizeX/2
		boardY := centerY - boardSizeY/2

		boardColor := color.NRGBA{0x8b, 0x45, 0x13, 0xcc}
		paperColor := color.NRGBA{0xff, 0xf4, 0xd4, 0xcc}
		textColor := color.NRGBA{0x00, 0x00, 0x00, 0xff}
		outlineColor := color.NRGBA{0x00, 0x00, 0x00, 0x00}

		if w.isSignNPC {
			boardColor = color.NRGBA{0x00, 0x00, 0x00, 0x44}
			paperColor = boardColor
			textColor = color.NRGBA{0xff, 0xff, 0xff, 0xff}
			outlineColor = color.NRGBA{0xa0, 0x20, 0xf0, 0xff}
		}

		// Draw the sign board
		vector.DrawFilledRect(
			ctx.Screen,
			boardX,
			boardY,
			boardSizeX,
			boardSizeY,
			boardColor,
			false,
		)

		margin := boardSizeX * 0.03

		// Draw the paper
		vector.DrawFilledRect(
			ctx.Screen,
			boardX+margin,
			boardY+margin,
			boardSizeX-margin*2,
			boardSizeY-margin*2,
			paperColor,
			false,
		)

		// Draw the text
		ctx.Text.SetScale(1)
		x := int(centerX)
		y := int(centerY)
		splitText := make([]string, 0)
		maxLen := 45
		for _, line := range text {
			if len(line) <= maxLen {
				splitText = append(splitText, line)
				continue
			}
			// split text into lines that are less than maxLen
			// split on spaces
			words := strings.Split(line, " ")
			currentLine := ""
			for _, word := range words {
				if len(currentLine)+len(word) > maxLen {
					splitText = append(splitText, currentLine)
					currentLine = ""
				}
				currentLine += word + " "
			}
			splitText = append(splitText, currentLine)
		}
		y = int(centerY) - (len(splitText)/2)*int(ctx.Text.Utils().GetLineHeight())
		for _, line := range splitText {
			{
				ctx.Text.SetColor(outlineColor)
				resources.DrawTextOutline(ctx.Text, ctx.Screen, line, x, y, 1.0)
				ctx.Text.SetColor(textColor)
				ctx.Text.Draw(ctx.Screen, line, x, y)
			}
			y += int(ctx.Text.Utils().GetLineHeight())
		}
	}
	s.hints.Draw(ctx)
}
