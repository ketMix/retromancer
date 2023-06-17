package game

import (
	"ebijam23/net"
	"ebijam23/states"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// LocalPlayer is a player on the local computer.
type LocalPlayer struct {
	connection     net.ServerClient // Only used if the player is a server.
	actor          Actor
	impulses       ImpulseSet
	queuedImpulses ImpulseSet
}

func (p *LocalPlayer) Update(ctx states.Context) {
	// FIXME: Use a bind system.
	ydir := 0
	xdir := 0
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		ydir = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		ydir = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		xdir = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		xdir = 1
	}
	if xdir != 0 || ydir != 0 {
		a := math.Atan2(float64(ydir), float64(xdir))
		p.impulses.Move = &ImpulseMove{Direction: a}
	} else {
		p.impulses.Move = nil
	}

	// TODO: Constrain/convert x, y to world coordinates.
	x, y := ebiten.CursorPosition()
	// This feels a bit wrong to set the player actor's hand position directly, but the hand position is just for visual indication as to where interactions go.
	if a, ok := p.actor.(*PC); ok {
		a.Hand.SetXY(float64(x), float64(y))
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		p.impulses.Interaction = ImpulseReflect{
			X: float64(x),
			Y: float64(y),
		}
		// FIXME: Change the sprite _during_ world update.
		if a, ok := p.actor.(*PC); ok {
			a.Hand.Sprite.SetImage(ctx.Manager.GetAs("images", "hand-reflect", (*ebiten.Image)(nil)).(*ebiten.Image))
		}

	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		p.impulses.Interaction = ImpulseDeflect{
			X: float64(x),
			Y: float64(y),
		}
		if a, ok := p.actor.(*PC); ok {
			a.Hand.Sprite.SetImage(ctx.Manager.GetAs("images", "hand-deflect", (*ebiten.Image)(nil)).(*ebiten.Image))
		}
	} else {
		p.impulses.Interaction = nil
		if a, ok := p.actor.(*PC); ok {
			a.Hand.Sprite.SetImage(ctx.Manager.GetAs("images", "hand-normal", (*ebiten.Image)(nil)).(*ebiten.Image))
		}
	}
}

// Tick is called on actual world tick.
func (p *LocalPlayer) Tick() {
	if p.actor != nil {
		p.actor.SetImpulses(p.queuedImpulses)
		p.queuedImpulses = ImpulseSet{}
	}
}

// Impulses is a list of impulses that the player currently desires their actor to process.
func (p *LocalPlayer) Impulses() ImpulseSet {
	return p.impulses
}

func (p *LocalPlayer) QueueImpulses(impulses ImpulseSet) {
	p.queuedImpulses = impulses
}

func (p *LocalPlayer) ClearImpulses() {
	p.impulses = ImpulseSet{}
}

func (p *LocalPlayer) Ready(nextTick int) bool {
	return true
}

func (p *LocalPlayer) Actor() Actor {
	return p.actor
}

func (p *LocalPlayer) SetActor(actor Actor) {
	p.actor = actor
	actor.SetPlayer(p)
}
