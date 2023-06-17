package game

import "ebijam23/states"

// Player represents either a local or remote player.
type Player interface {
	// ???
	Ready(nextTick int) bool
	SetActor(actor Actor)
	Actor() Actor
	Update(ctx states.Context)
	Tick()
	Impulses() ImpulseSet
	ClearImpulses()
	QueueImpulses(impulses ImpulseSet)
}
