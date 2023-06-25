package game

import (
	"ebijam23/resources"
	"ebijam23/states"
	"image/color"
	"math"
	"math/rand"
	"reflect"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type WorldStateLive struct {
	signText *string
}

func (w *WorldStateLive) Enter(s *World, ctx states.Context) {
}

func (w *WorldStateLive) Leave(s *World, ctx states.Context) {
}

func (w *WorldStateLive) Tick(s *World, ctx states.Context) {
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

	// Process particles
	for _, p := range s.activeMap.particles {
		p.Update()
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
				if collision := s.activeMap.Collides(checkShape); collision == nil || !collision.Cell.BlockMove {
					actor.SetXY(action.X, action.Y)
				}
				// forgive me.
				if s.tick%4 == 0 {
					if pc, ok := actor.(*PC); ok {
						s.SpawnParticle(ctx, "puff", action.X, action.Y+pc.Sprite.Height()/2-2, 0, 0, 10)
					}
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
							if a.Reverseable() {
								a.Reverse()
								s.SpawnParticle(ctx, "reverse", action.X, action.Y, rand.Float64()*math.Pi*2, rand.Float64()*2.0, 30)
							}
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
			case ActionSpawnParticle:
				s.SpawnParticle(ctx, action.Img, action.X, action.Y, action.Angle, action.Speed, action.Life)
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
			case ActionSpawnParticle:
				s.SpawnParticle(ctx, action.Img, action.X, action.Y, action.Angle, action.Speed, action.Life)
			}
		}
	}

	// Okay, this probably isn't great, but let's check bullet collisions here.
	for _, bullet := range s.activeMap.bullets {
		// Check for bullet collisions with walls.
		if collision := s.activeMap.Collides(&bullet.Shape); collision != nil && collision.Cell.BlockView {
			bullet.Destroyed = true
			continue
		}

		// Check for bullet collisions with actors.
		for _, actor := range s.activeMap.actors {
			// Check player collisions.
			if p, ok := actor.(*PC); ok {
				if p.InvulnerableTicks > 0 {
					continue
				}
				if bullet.Shape.Collides(actor.Shape()) {
					x, y, _, _ := actor.Bounds()
					for i := 0; i < 6; i++ {
						s.SpawnParticle(ctx, "hurt", x, y, bullet.Angle-math.Pi/4+(math.Pi/2*rand.Float64()), rand.Float64()*2.0, 30)
					}
					bullet.Destroyed = true
					p.Hurtie()
					break
				}
				continue // skip checking other actors
			}

			// Check interactive and enemy collisions.
			// Only applicable to reflected or deflected bullets
			if bullet.reflected || bullet.deflected {
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
							bullet.Destroyed = true
							e.Damage(1)
							break
						}
					}
				}
			}
		}
	}

	// Oh boy, yet another loop.
	// Check for collisions between player characters and interactives/snaggables.
	touchingSign := false
	for _, pl := range s.Players {
		if pc, ok := pl.Actor().(*PC); ok {
			for _, actor := range s.activeMap.actors {
				// Check interactive collisions.
				if i, ok := actor.(*Interactive); ok {
					// If the interactive has text and is active, show the text.
					if !touchingSign && i.text != "" && i.Active() && i.shape.Collides(pl.Actor().Shape()) {
						localized := ctx.L(i.text)
						w.signText = &localized
						touchingSign = true
					}

					// If touchable, apply reverse to it
					if i.touchable && i.shape.Collides(pl.Actor().Shape()) {
						i.Reverse()
						continue
					}

					// If nextMap is defined and active, go to next map.
					if i.nextMap != nil && i.Active() && i.shape.Collides(pl.Actor().Shape()) {
						ctx.Manager.GetAs("sounds", "stairs", (*resources.Sound)(nil)).(*resources.Sound).Play(0.5)
						s.TravelToMap(ctx, *i.nextMap)
					}
				}

				// Check snaggable collisions.
				if s, ok := actor.(*Snaggable); ok {
					if s.shape.Collides(pl.Actor().Shape()) && pc.Lives < PLAYER_MAX_LIVES {
						s.destroyed = true
						// TODO: Bless the player with powers beyond their wildest imaginations (lives, power unlocks, etc.)
						switch s.spriteName {
						case "item-life":
							pc.Lives++
						}
						continue
					}
				}

				// Check enemy collisions.
				if e, ok := actor.(*Enemy); ok {
					if e.IsAlive() && e.Shape().Collides(pl.Actor().Shape()) {
						pc.Hurtie()
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
			cells := make([]*resources.Cell, 0)
			if cell != nil {
				cells = append(cells, cell)
				cells = append(cells, s.activeMap.GetCell(int(cell.Sprite.X/cellW), int(cell.Sprite.Y/cellH)+1, s.activeMap.currentZ))
				cells = append(cells, s.activeMap.GetCell(int(cell.Sprite.X/cellW)+1, int(cell.Sprite.Y/cellH)+1, s.activeMap.currentZ))
			}
			for _, cell := range cells {
				if cell != nil {
					cell.BlockMove = false // No
					cell.BlockView = false // No
				}
			}
		}
		if !actor.active {
			if CheckConditions(actor.Conditions(), interactives, enemies) {
				actor.IncreaseActivation(nil)
				cell := s.activeMap.FindCellById(actor.ID())
				if cell != nil {
					cell.BlockMove = false // No
					cell.BlockView = false // No
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

	// Queue up the local player's impulses for the next tick!
	for _, player := range s.Players {
		if _, ok := player.(*LocalPlayer); ok {
			player.QueueImpulses(player.Impulses())
			player.ClearImpulses()
			// TODO: Send network message to peers with our impulses!
		}
	}

	if s.ArePlayersDead() {
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

	// Draw the sign text if it exists.
	if w.signText != nil {
		centerX := float32(ctx.Screen.Bounds().Max.X / 2)
		centerY := float32(ctx.Screen.Bounds().Max.Y / 2)
		boardSizeX := float32(ctx.Screen.Bounds().Max.X) * 0.85
		boardSizeY := float32(ctx.Screen.Bounds().Max.Y) * 0.75
		text := strings.Split(*w.signText, "\n")

		boardX := centerX - boardSizeX/2
		boardY := centerY - boardSizeY/2

		// Draw the stake
		vector.DrawFilledRect(
			ctx.Screen,
			centerX-boardSizeX*0.05,
			centerY-boardSizeY*0.05,
			boardSizeX*0.05,
			float32(ctx.Screen.Bounds().Max.Y),
			color.RGBA{0x8b, 0x45, 0x13, 0xff},
			false,
		)

		// Draw the sign board
		vector.DrawFilledRect(
			ctx.Screen,
			boardX,
			boardY,
			boardSizeX,
			boardSizeY,
			color.RGBA{0x8b, 0x45, 0x13, 0xff},
			false,
		)

		// Draw the paper
		vector.DrawFilledRect(
			ctx.Screen,
			boardX+boardSizeX*0.05,
			boardY+boardSizeY*0.1,
			boardSizeX*0.9,
			boardSizeY*0.8,
			color.White,
			false,
		)

		// Draw the text
		x := int(centerX)
		y := int(centerY)
		splitText := make([]string, 0)
		maxLen := 35
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
				ctx.Text.SetScale(1.5)
				ctx.Text.SetColor(color.Black)
				ctx.Text.Draw(ctx.Screen, line, x, y)
			}
			y += int(ctx.Text.Utils().GetLineHeight())
		}
	}
}
