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
	//
	impulses ImpulseSet
}

func (p *PC) SetPlayer(player Player) {
	p.player = player
}

func (p *PC) Player() Player {
	return p.player
}

func (p *PC) Update() error {
	p.TicksSinceLastInteraction++
	if p.TicksSinceLastInteraction > 20 {
		p.Energy++
	}
	if p.impulses.Move != nil {
		x := 5 * math.Cos((*p.impulses.Move).Direction)
		y := 5 * math.Sin((*p.impulses.Move).Direction)
		p.SetXY(p.Shape.X+x, p.Shape.Y+y)
	}

	if p.impulses.Interaction != nil {
		switch p.impulses.Interaction.(type) {
		case ImpulseReflect:
			// TODO: Process reflection within the world.
			if p.Energy-1 >= 0 {
				p.Energy--
			}
			p.TicksSinceLastInteraction = 0
		case ImpulseDeflect:
			if p.Energy-4 >= 0 {
				p.Energy -= 4
			}
			p.TicksSinceLastInteraction = 0
			//r := math.Atan2(imp.Y-p.Shape.Y, imp.X-p.Shape.X)
			// TODO: Process deflection within the world.
		default:
			// Do nothing.
		}
		// TODO: Handle interactions.
	}

	return nil
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
