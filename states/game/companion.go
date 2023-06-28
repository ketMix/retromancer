package game

import (
	"ebijam23/resources"
	"ebijam23/states"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Companion struct {
	player Player
	saved  *Companion
	//
	Arrow             *resources.Sprite
	Sprite            *resources.Sprite
	Hat               *resources.Sprite
	shape             CircleShape
	Hand              Hand
	Energy            int
	MaxEnergy         int
	EnergyRestoreRate int
	snarfTicks        int
	//
	fireAllow                 int
	TicksSinceLastInteraction int
	previousInteraction       Action
	//
	impulses ImpulseSet
	//
	momentumX float64
	momentumY float64
}

func (s *World) NewCompanion(ctx states.Context) *Companion {
	pc := &Companion{
		Sprite:            resources.NewSprite(ctx.Manager.GetAs("images", "companion", (*ebiten.Image)(nil)).(*ebiten.Image)),
		Arrow:             resources.NewSprite(ctx.Manager.GetAs("images", "direction-arrow", (*ebiten.Image)(nil)).(*ebiten.Image)),
		Energy:            50,
		MaxEnergy:         50,
		EnergyRestoreRate: 3,
	}

	pc.shape.Radius = 1
	//pc.Sprite.Interpolate = true
	pc.Sprite.Centered = true
	pc.Hand.Sprite = resources.NewSprite(ctx.Manager.GetAs("images", "companion-hand", (*ebiten.Image)(nil)).(*ebiten.Image))
	pc.Hand.Sprite.Centered = true
	pc.Hand.HoverSprite = resources.NewSprite(ctx.Manager.GetAs("images", "companion-hand-glow", (*ebiten.Image)(nil)).(*ebiten.Image))
	pc.Hand.HoverSprite.Hidden = true
	pc.Hand.HoverSprite.Centered = true

	return pc
}

func (p *Companion) Save() {
	saved := *p
	p.saved = &saved
}

func (p *Companion) Restore() {
	*p = *p.saved
}

func (p *Companion) SetPlayer(player Player) {
	p.player = player
}

func (p *Companion) Player() Player {
	return p.player
}

func (p *Companion) Update() (actions []Action) {
	p.snarfTicks--
	multiplier := 1.0
	if p.snarfTicks > 0 {
		multiplier = 2.0
		actions = append(actions, ActionSpawnParticle{
			Img:   "life",
			X:     p.shape.X,
			Y:     p.shape.Y,
			Angle: math.Pi + rng.Float64()*math.Pi,
			Speed: rng.Float64() * 0.5,
			Life:  40,
		})
	}

	p.fireAllow++

	p.TicksSinceLastInteraction++
	if p.TicksSinceLastInteraction > 20/int(multiplier) {
		if p.Energy+p.EnergyRestoreRate <= p.MaxEnergy {
			p.Energy += int(float64(p.EnergyRestoreRate) * multiplier)
		}
	}
	// Do not handle movement until we are resurrected.
	if p.impulses.Move != nil {
		p.momentumX = 0.3*p.momentumX + 4.7*math.Cos((*p.impulses.Move).Direction)*multiplier
		p.momentumY = 0.3*p.momentumY + 4.7*math.Sin((*p.impulses.Move).Direction)*multiplier
		/*x := 5 * math.Cos((*p.impulses.Move).Direction)
		y := 5 * math.Sin((*p.impulses.Move).Direction)
		if x < 0 {
			p.Sprite.Flipped = true
		} else if x > 0 {
			p.Sprite.Flipped = false
		}
		actions = append(actions, ActionMove{
			X: p.shape.X + x,
			Y: p.shape.Y + y,
		})*/
	}
	if p.momentumX != 0 || p.momentumY != 0 {
		actions = append(actions, ActionMove{
			X: p.shape.X + p.momentumX,
			Y: p.shape.Y + p.momentumY,
		})
	}
	p.momentumX *= 0.3
	p.momentumY *= 0.3
	if math.Abs(p.momentumX) < 0.1 {
		p.momentumX = 0
	}
	if math.Abs(p.momentumY) < 0.1 {
		p.momentumY = 0
	}

	if p.Hand.Shape.X < p.shape.X {
		p.Sprite.Flipped = true
	} else {
		p.Sprite.Flipped = false
	}

	p.previousInteraction = nil
	if p.impulses.Interaction != nil {
		switch imp := p.impulses.Interaction.(type) {
		case ImpulseShoot:
			if p.HasEnergyFor(imp) && p.fireAllow > 0 {
				p.Energy -= imp.Cost()
				p.TicksSinceLastInteraction = 0
				angle := math.Atan2(imp.Y-p.shape.Y, imp.X-p.shape.X)
				bullet := CreateBullet(Circular, color.RGBA{0xa0, 0x20, 0xf0, 0xaa}, 1, 4, angle, 1, 2, 0, 2, 1, 1, 0)
				bullet.friendly = true
				bullet.Deathtime = 200
				bullet.Damage = 1
				bullet.SetXY(p.shape.X+math.Cos(angle)*9, p.shape.Y+2)
				actions = append(actions, ActionSpawnBullets{
					Bullets: []*Bullet{
						bullet,
					},
				})
				p.previousInteraction = ActionSpawnBullets{}
				p.fireAllow = -10
				if p.snarfTicks > 0 {
					p.fireAllow = -1
				}
			}
		default:
			// Do nothing.
		}
	}

	return actions
}

func (p *Companion) Dead() bool {
	return false
}

func (p *Companion) HasEnergyFor(imp Impulse) bool {
	return p.Energy-imp.Cost() >= 0
}

func (p *Companion) SetImpulses(impulses ImpulseSet) {
	p.impulses = impulses
}

// DrawHat draws the player's dumb hat.
func (p *Companion) DrawHat(screen *ebiten.Image, x, y float64) {
	opts := &ebiten.DrawImageOptions{}
	if p.Sprite.Flipped {
		opts.GeoM.Scale(-1, 1)
		opts.GeoM.Translate(p.Hat.Width(), 0)
	}
	opts.GeoM.Translate(p.shape.X-float64(int(p.Hat.Width())/2)+x, p.Sprite.Y-p.Sprite.Height()/2-p.Hat.Height()+y)
	screen.DrawImage(p.Hat.Image(), opts)
}

func (p *Companion) DrawHand(ctx states.DrawContext) {
	p.Hand.Sprite.Draw(ctx)
	p.Hand.HoverSprite.Draw(ctx)
}

func (p *Companion) Draw(ctx states.DrawContext) {
	if _, ok := p.previousInteraction.(ActionDeflect); ok {
		vector.DrawFilledCircle(ctx.Screen, float32(p.Hand.Shape.X), float32(p.Hand.Shape.Y), 20, color.NRGBA{0xff, 0x66, 0x99, 0x33}, false)
	} else if _, ok := p.previousInteraction.(ActionReflect); ok {
		vector.DrawFilledCircle(ctx.Screen, float32(p.Hand.Shape.X), float32(p.Hand.Shape.Y), 20, color.NRGBA{0x66, 0x99, 0xff, 0x33}, false)
	} else if _, ok := p.previousInteraction.(ActionShield); ok {
		vector.DrawFilledCircle(ctx.Screen, float32(p.shape.X), float32(p.shape.Y), 20, color.NRGBA{0x66, 0xff, 0x99, 0x33}, false)
	}

	opts := &ebiten.DrawImageOptions{}

	p.Sprite.Draw(ctx)

	p.DrawHat(ctx.Screen, 0, 3)

	r := math.Atan2(p.shape.Y-p.Hand.Shape.Y, p.shape.X-p.Hand.Shape.X)

	// TODO: The arrow image should change based on if we're reflecting or deflecting.
	// Draw direction arrow
	opts = &ebiten.DrawImageOptions{}
	// Rotate about its center.
	opts.GeoM.Translate(-p.Arrow.Width()/2, -p.Arrow.Height()/2)
	opts.GeoM.Rotate(r)
	// Position.
	//opts.GeoM.Translate(p.Arrow.Width()/2, p.Arrow.Height()/2)
	opts.GeoM.Translate(p.shape.X, p.shape.Y)
	// Draw from center.
	ctx.Screen.DrawImage(p.Arrow.Image(), opts)
}

func (p *Companion) Bounds() (x, y, w, h float64) {
	// Return the radius of the shape as width and height in diameter.
	return p.shape.X, p.shape.Y, p.shape.Radius * 2, p.shape.Radius * 2
}

func (p *Companion) Shape() Shape {
	return &p.shape
}

func (p *Companion) SetXY(x, y float64) {
	p.shape.X = x
	p.shape.Y = y
	p.Sprite.X = x
	p.Sprite.Y = y + 2 // We lightly offset the sprite so the phylactery is in a nicer visual position.
}

func (p *Companion) SetSize(r float64) {
	p.shape.Radius = r
}

func (p *Companion) Destroyed() bool {
	return false
}

func (p *Companion) Snarf() {
	if p.snarfTicks <= 0 {
		p.snarfTicks = 300
	} else {
		p.snarfTicks += 300
	}
}
