package game

import (
	"ebijam23/resources"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type PC struct {
	player Player
	//
	Sprite                    *resources.Sprite
	Shape                     Shape
	Hand                      Hand
	TicksSinceLastInteraction int
	Energy                    int
	MaxEnergy                 int
	EnergyRestoreRate         int
	//
	impulses ImpulseSet
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
		actions = append(actions, ActionMove{
			X: p.Shape.X + x,
			Y: p.Shape.Y + y,
		})
	}

	if p.impulses.Interaction != nil {
		switch imp := p.impulses.Interaction.(type) {
		case ImpulseReflect:
			if p.HasEnergyFor(imp) {
				p.Energy -= imp.Cost()
			}
			p.TicksSinceLastInteraction = 0
			actions = append(actions, ActionReflect{
				X: imp.X,
				Y: imp.Y,
			})
		case ImpulseDeflect:
			if p.HasEnergyFor(imp) {
				p.Energy -= imp.Cost()
			}
			p.TicksSinceLastInteraction = 0
			actions = append(actions, ActionDeflect{
				Direction: math.Atan2(imp.Y-p.Shape.Y, imp.X-p.Shape.X),
			})
		default:
			// Do nothing.
		}
		// TODO: Handle interactions.
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
	p.Sprite.Draw(screen)
	p.Hand.Sprite.Draw(screen)
}

func (p *PC) Bounds() (x, y, w, h float64) {
	// Return the radius of the shape as width and height in diameter.
	return p.Shape.X, p.Shape.Y, p.Shape.Radius * 2, p.Shape.Radius * 2
}

func (p *PC) SetXY(x, y float64) {
	p.Shape.X = x
	p.Shape.Y = y
	p.Sprite.X = x
	p.Sprite.Y = y
}

func (p *PC) SetSize(r float64) {
	p.Shape.Radius = r
}
