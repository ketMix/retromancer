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
	// controller vars
	cx, cy, ca, cd  float64
	handRotateSpeed float64
}

func NewLocalPlayer() *LocalPlayer {
	return &LocalPlayer{
		impulses:        ImpulseSet{},
		queuedImpulses:  ImpulseSet{},
		ca:              math.Pi / 2,
		cd:              40.0,
		handRotateSpeed: 0.075,
		hat:             "hat-",
	}
}

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
			d := a - p.ca
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
			speed := p.handRotateSpeed

			speed /= p.cd / 40.0

			if d > speed {
				p.ca += speed * accel
			}
			if d < -speed {
				p.ca -= speed * accel
			}
			if d > -speed && d < speed {
				p.ca = a
			}
		}

		if ebiten.IsGamepadButtonPressed(ebiten.GamepadID(p.GamepadID), ebiten.GamepadButton5) {
			p.cd += 3
			if p.cd > 300 {
				p.cd = 300
			}
		} else if ebiten.IsGamepadButtonPressed(ebiten.GamepadID(p.GamepadID), ebiten.GamepadButton4) {
			if p.cd-3 > 0 {
				p.cd -= 3
			}
		}

		p.cx = math.Cos(p.ca) * p.cd
		p.cy = math.Sin(p.ca) * p.cd

		if a, ok := p.actor.(*PC); ok {
			p.cx += float64(a.shape.X)
			p.cy += float64(a.shape.Y)
			a.Hand.SetXY(p.cx, p.cy)
		}

		if ebiten.IsGamepadButtonPressed(ebiten.GamepadID(p.GamepadID), ebiten.GamepadButton1) {
			p.impulses.Interaction = ImpulseShield{}
		} else if ebiten.IsGamepadButtonPressed(ebiten.GamepadID(p.GamepadID), ebiten.GamepadButton7) {
			p.impulses.Interaction = ImpulseReflect{
				X: float64(p.cx),
				Y: float64(p.cy),
			}
		} else if ebiten.IsGamepadButtonPressed(ebiten.GamepadID(p.GamepadID), ebiten.GamepadButton6) {
			p.impulses.Interaction = ImpulseDeflect{
				X: float64(p.cx),
				Y: float64(p.cy),
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
