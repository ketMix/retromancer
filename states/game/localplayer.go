package game

import (
	"ebijam23/net"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// LocalPlayer is a player on the local computer.
type LocalPlayer struct {
	connection     net.ServerClient // Only used if the player is a server.
	actor          Actor
	thoughts       []Thought
	impulses       ImpulseSet
	queuedImpulses ImpulseSet
	GamepadID      int // Target gamepad for this player to use.
	hat            string
}

// TEMP: These need to be tied to a gamepad type or something.
var (
	cx              = 0.0
	cy              = 0.0
	ca              = math.Pi / 2
	cd              = 40.0
	handRotateSpeed = 0.075
)

func (p *LocalPlayer) Update() {
	// FIXME: All of this is pretty rough, but I wanted to test controller usage.
	if p.GamepadID > 0 && len(ebiten.GamepadIDs()) >= p.GamepadID {
		lr := ebiten.GamepadAxisValue(ebiten.GamepadID(p.GamepadID), 0)
		ud := ebiten.GamepadAxisValue(ebiten.GamepadID(p.GamepadID), 1)

		if math.Abs(lr) > 0.01 || math.Abs(ud) > 0.01 {
			a := math.Atan2(ud, lr)
			p.impulses.Move = &ImpulseMove{Direction: a}
		}

		r1 := ebiten.GamepadAxisValue(ebiten.GamepadID(p.GamepadID), 3)
		r2 := ebiten.GamepadAxisValue(ebiten.GamepadID(p.GamepadID), 4)
		if math.Abs(r1) > 0.01 || math.Abs(r2) > 0.01 {
			a := math.Atan2(r2, r1)

			// Increment ca towards a in increments of 0.1 in the range [-pi, pi]:
			d := a - ca
			if d > math.Pi {
				d -= 2 * math.Pi
			}
			if d < -math.Pi {
				d += 2 * math.Pi
			}
			accel := 1.0
			if math.Abs(d) > 1 {
				accel = math.Abs(d)
			}
			speed := handRotateSpeed

			speed /= cd / 40.0

			if d > speed {
				ca += speed * accel
			}
			if d < -speed {
				ca -= speed * accel
			}
			if d > -speed && d < speed {
				ca = a
			}
		}

		if ebiten.IsGamepadButtonPressed(ebiten.GamepadID(p.GamepadID), ebiten.GamepadButton5) {
			cd += 3
			if cd > 300 {
				cd = 300
			}
		} else if ebiten.IsGamepadButtonPressed(ebiten.GamepadID(p.GamepadID), ebiten.GamepadButton4) {
			if cd-3 > 0 {
				cd -= 3
			}
		}

		cx = math.Cos(ca) * cd
		cy = math.Sin(ca) * cd

		if a, ok := p.actor.(*PC); ok {
			cx += float64(a.shape.X)
			cy += float64(a.shape.Y)
			a.Hand.SetXY(cx, cy)
		}

		if ebiten.IsGamepadButtonPressed(ebiten.GamepadID(p.GamepadID), ebiten.GamepadButton1) {
			p.impulses.Interaction = ImpulseShield{}
		} else if ebiten.IsGamepadButtonPressed(ebiten.GamepadID(p.GamepadID), ebiten.GamepadButton7) {
			p.impulses.Interaction = ImpulseReflect{
				X: float64(cx),
				Y: float64(cy),
			}
		} else if ebiten.IsGamepadButtonPressed(ebiten.GamepadID(p.GamepadID), ebiten.GamepadButton6) {
			p.impulses.Interaction = ImpulseDeflect{
				X: float64(cx),
				Y: float64(cy),
			}
		} else {
			p.impulses.Interaction = nil
		}

	} else {
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

		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			p.impulses.Interaction = ImpulseShield{}
		} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			p.impulses.Interaction = ImpulseReflect{
				X: float64(x),
				Y: float64(y),
			}
		} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
			p.impulses.Interaction = ImpulseDeflect{
				X: float64(x),
				Y: float64(y),
			}
		} else {
			p.impulses.Interaction = nil
		}
	}

	// Thoughts
	p.thoughts = []Thought{}
	if p.actor != nil && p.actor.Dead() {
		if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsGamepadButtonPressed(ebiten.GamepadID(p.GamepadID), ebiten.GamepadButton0) {
			p.thoughts = append(p.thoughts, ResetThought{})
		} else if ebiten.IsKeyPressed(ebiten.KeyEscape) || ebiten.IsGamepadButtonPressed(ebiten.GamepadID(p.GamepadID), ebiten.GamepadButton1) {
			p.thoughts = append(p.thoughts, QuitThought{})
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

func (p *LocalPlayer) Thoughts() []Thought {
	return p.thoughts
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

func (p *LocalPlayer) Hat() string {
	return p.hat
}

func (p *LocalPlayer) SetHat(hat string) {
	p.hat = hat
}
