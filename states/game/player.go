package game

// Player represents either a local or remote player.
type Player interface {
	// ???
	Ready(nextTick int) bool
	SetActor(actor Actor)
	Actor() Actor
	Update()
	Tick()
	Impulses() ImpulseSet
	ClearImpulses()
	QueueImpulses(impulses ImpulseSet)
	//
	Thoughts() []Thought
	//
	SetHat(string)
	Hat() string
}
