package game

// Player represents either a local or remote player.
type Player interface {
	// ???
	Ready(nextTick int) bool
}
