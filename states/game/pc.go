package game

import (
	"ebijam23/resources"
	"ebijam23/states"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type PC struct {
	player Player
	//
	Arrow                     *resources.Sprite
	Sprite                    *resources.Sprite
	Phylactery                *resources.Sprite
	Hat                       *resources.Sprite
	shape                     CircleShape
	Hand                      Hand
	InvulnerableTicks         int // Ticks the player should be invulnerable for
	TicksSinceLastInteraction int
	Energy                    int
	MaxEnergy                 int
	EnergyRestoreRate         int
	//
	previousInteraction Action
	//
	impulses ImpulseSet
}

func (s *World) NewPC(ctx states.Context) *PC {
	pc := &PC{
		Sprite:            resources.NewSprite(ctx.Manager.GetAs("images", "player", (*ebiten.Image)(nil)).(*ebiten.Image)),
		Phylactery:        resources.NewSprite(ctx.Manager.GetAs("images", "phylactery", (*ebiten.Image)(nil)).(*ebiten.Image)),
		Arrow:             resources.NewSprite(ctx.Manager.GetAs("images", "direction-arrow", (*ebiten.Image)(nil)).(*ebiten.Image)),
		Energy:            100,
		MaxEnergy:         100,
		EnergyRestoreRate: 2,
	}
	pc.shape.Radius = 2
	//pc.Sprite.Interpolate = true
	pc.Sprite.Centered = true
	pc.Hand.Sprite = resources.NewSprite(ctx.Manager.GetAs("images", "hand-normal", (*ebiten.Image)(nil)).(*ebiten.Image))
	pc.Hand.Sprite.Centered = true

	return pc
}

func (p *PC) SetPlayer(player Player) {
	p.player = player
}

func (p *PC) Player() Player {
	return p.player
}

func (p *PC) Update() (actions []Action) {
	p.TicksSinceLastInteraction++
	if p.TicksSinceLastInteraction > 20 {
		if p.Energy+p.EnergyRestoreRate <= p.MaxEnergy {
			p.Energy += p.EnergyRestoreRate
		}
	}
	if p.impulses.Move != nil {
		x := 5 * math.Cos((*p.impulses.Move).Direction)
		y := 5 * math.Sin((*p.impulses.Move).Direction)
		if x < 0 {
			p.Sprite.Flipped = true
		} else if x > 0 {
			p.Sprite.Flipped = false
		}
		actions = append(actions, ActionMove{
			X: p.shape.X + x,
			Y: p.shape.Y + y,
		})
	}

	p.previousInteraction = nil
	if p.impulses.Interaction != nil {
		switch imp := p.impulses.Interaction.(type) {
		case ImpulseReflect:
			if p.HasEnergyFor(imp) {
				p.Energy -= imp.Cost()
				p.TicksSinceLastInteraction = 0
				actions = append(actions, ActionReflect{
					X: imp.X,
					Y: imp.Y,
				})
				p.previousInteraction = ActionReflect{}
			}
		case ImpulseDeflect:
			if p.HasEnergyFor(imp) {
				p.Energy -= imp.Cost()
				p.TicksSinceLastInteraction = 0
				actions = append(actions, ActionDeflect{
					Direction: math.Atan2(imp.Y-p.shape.Y, imp.X-p.shape.X),
					X:         imp.X,
					Y:         imp.Y,
				})
				p.previousInteraction = ActionDeflect{}
			}
		case ImpulseShield:
			if p.HasEnergyFor(imp) {
				p.Energy -= imp.Cost()
				p.TicksSinceLastInteraction = 0
				actions = append(actions, ActionShield{})
				p.previousInteraction = ActionShield{}
			}
		default:
			// Do nothing.
		}
	}

	return actions
}

func (p *PC) HasEnergyFor(imp Impulse) bool {
	return p.Energy-imp.Cost() >= 0
}

func (p *PC) SetImpulses(impulses ImpulseSet) {
	p.impulses = impulses
}

func (p *PC) Draw(screen *ebiten.Image) {
	if _, ok := p.previousInteraction.(ActionDeflect); ok {
		vector.DrawFilledCircle(screen, float32(p.Hand.Shape.X), float32(p.Hand.Shape.Y), 20, color.NRGBA{0xff, 0x66, 0x99, 0x33}, false)
	} else if _, ok := p.previousInteraction.(ActionReflect); ok {
		vector.DrawFilledCircle(screen, float32(p.Hand.Shape.X), float32(p.Hand.Shape.Y), 20, color.NRGBA{0x66, 0x99, 0xff, 0x33}, false)
	} else if _, ok := p.previousInteraction.(ActionShield); ok {
		vector.DrawFilledCircle(screen, float32(p.shape.X), float32(p.shape.Y), 20, color.NRGBA{0x66, 0xff, 0x99, 0x33}, false)
	}

	opts := &ebiten.DrawImageOptions{}

	p.InvulnerableTicks--

	p.Hand.Sprite.Draw(screen)

	if p.InvulnerableTicks <= 0 || p.InvulnerableTicks%20 >= 10 {
		p.Sprite.Draw(screen)

		// Draw the player's phylactery (hit box representation).
		opts.GeoM.Translate(p.shape.X-float64(int(p.Phylactery.Width())/2), p.shape.Y-float64(int(p.Phylactery.Height())/2))
		opts.ColorScale.Scale(0.5, 0.5, 1.0, 1.0)
		screen.DrawImage(p.Phylactery.Image(), opts)

		// Draw the player's dumb hat.
		opts = &ebiten.DrawImageOptions{}
		if p.Sprite.Flipped {
			opts.GeoM.Scale(-1, 1)
			opts.GeoM.Translate(p.Hat.Width(), 0)
		}
		opts.GeoM.Translate(p.shape.X-float64(int(p.Hat.Width())/2), p.Sprite.Y-p.Sprite.Height()/2-p.Hat.Height()+3)
		screen.DrawImage(p.Hat.Image(), opts)

	}

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
	screen.DrawImage(p.Arrow.Image(), opts)
}

func (p *PC) Bounds() (x, y, w, h float64) {
	// Return the radius of the shape as width and height in diameter.
	return p.shape.X, p.shape.Y, p.shape.Radius * 2, p.shape.Radius * 2
}

func (p *PC) Shape() Shape {
	return &p.shape
}

func (p *PC) SetXY(x, y float64) {
	p.shape.X = x
	p.shape.Y = y
	p.Sprite.X = x
	p.Sprite.Y = y + 2 // We lightly offset the sprite so the phylactery is in a nicer visual position.
}

func (p *PC) SetSize(r float64) {
	p.shape.Radius = r
}
