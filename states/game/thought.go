package game

import "encoding/gob"

// Thoughts are a poorly named notion that sits above Impulses. These can effectively be viewed as player desires/commands and are sent across the wire between clients.
type Thought interface{}

type Thoughts struct {
	Thoughts []Thought
}

type ImpulsesThought struct {
	Impulses ImpulseSet
}

type ResetThought struct{}

type QuitThought struct{}

func init() {
	gob.Register(Thoughts{})
	gob.Register(ImpulsesThought{})
	gob.Register(ResetThought{})
	gob.Register(QuitThought{})
	gob.Register(TickState{})
}

type TickState struct {
	Thoughts Thoughts
	Impulses ImpulseSet
}
