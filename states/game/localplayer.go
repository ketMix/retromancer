package game

// LocalPlayer is a player on the local computer.
type LocalPlayer struct {
}

func (p *LocalPlayer) Ready(nextTick int) bool {
	return true
}
