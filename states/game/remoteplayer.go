package game

// RemotePlayer is a networked player.
type RemotePlayer struct {
	lastTick int
}

func (p *RemotePlayer) Ready(nextTick int) bool {
	return p.lastTick == nextTick-1
}
