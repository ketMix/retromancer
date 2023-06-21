package game

// Thoughts are a poorly named notion that sits above Impulses. These can effectively be viewed as player desires/commands and are sent across the wire between clients.
type Thought interface{}

type ImpulsesThought struct {
	Impulses ImpulseSet
}

type ResetThought struct{}

type QuitThought struct{}
