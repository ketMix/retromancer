package game

import "ebijam23/net"

// RemotePlayer is a networked player.
type RemotePlayer struct {
	connection *net.Conn // Remote connection, wowee.
	lastTick   int
	actor      Actor
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
