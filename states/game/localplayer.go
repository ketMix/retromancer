package game

import "ebijam23/net"

// LocalPlayer is a player on the local computer.
type LocalPlayer struct {
	connection net.ServerClient // Only used if the player is a server.
	actor      Actor
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
