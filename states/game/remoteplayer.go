package game

import (
	"ebijam23/net"
)

// RemotePlayer is a networked player.
type RemotePlayer struct {
	connection     *net.Conn // Remote connection, wowee.
	lastTick       int
	actor          Actor
	impulses       ImpulseSet
	queuedImpulses ImpulseSet
	thoughts       []Thought
	hat            string
}

func (p *RemotePlayer) Update() {
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

func (p *RemotePlayer) Thoughts() []Thought {
	return p.thoughts
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

func (p *RemotePlayer) Hat() string {
	return p.hat
}

func (p *RemotePlayer) SetHat(hat string) {
	p.hat = hat
}
