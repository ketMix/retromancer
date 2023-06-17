package game

import (
	"ebijam23/net"
	"ebijam23/states"
)

// RemotePlayer is a networked player.
type RemotePlayer struct {
	connection     *net.Conn // Remote connection, wowee.
	lastTick       int
	actor          Actor
	impulses       ImpulseSet
	queuedImpulses ImpulseSet
}

func (p *RemotePlayer) Update(ctx states.Context) {
}

func (p *RemotePlayer) Tick() {
}

func (p *RemotePlayer) Impulses() ImpulseSet {
	return p.impulses
}

func (p *RemotePlayer) QueueImpulses(impulses ImpulseSet) {
	p.queuedImpulses = impulses
}

func (p *RemotePlayer) ClearImpulses() {
	p.impulses = ImpulseSet{}
}

func (p *RemotePlayer) Ready(nextTick int) bool {
	return p.lastTick == nextTick-1
}

func (p *RemotePlayer) Actor() Actor {
	return p.actor
}

func (p *RemotePlayer) SetActor(actor Actor) {
	p.actor = actor
	actor.SetPlayer(p)
}
