package game

import (
	"ebijam23/net"
	"encoding/binary"
)

// Thoughts are a poorly named notion that sits above Impulses. These can effectively be viewed as player desires/commands and are sent across the wire between clients.
type Thought interface {
	Ident() uint8
	ToBytes() []byte
	FromBytes(b []byte) (net.Message, int)
}

type Thoughts struct {
	Thoughts []Thought
}

func (t Thoughts) Ident() uint8 {
	return 20
}

func (t Thoughts) ToBytes() (b []byte) {
	b = append(b, t.Ident())
	b = binary.LittleEndian.AppendUint64(b, uint64(len(t.Thoughts)))
	for _, thought := range t.Thoughts {
		b = append(b, thought.ToBytes()...)
	}
	return
}

func (t Thoughts) FromBytes(b []byte) (net.Message, int) {
	offset := 1
	t.Thoughts = make([]Thought, binary.LittleEndian.Uint64(b[offset:]))
	offset += 8
	for i := 0; i < len(t.Thoughts); i++ {
		thought, n := net.MessageFromBytes(b[offset:])
		t.Thoughts[i] = thought.(Thought)
		offset += n
	}
	return t, offset
}

type ResetThought struct{}

func (t ResetThought) Ident() uint8 {
	return 21
}

func (t ResetThought) ToBytes() (b []byte) {
	b = append(b, t.Ident())
	return
}

func (t ResetThought) FromBytes(b []byte) (net.Message, int) {
	return t, 1
}

type QuitThought struct{}

func (t QuitThought) Ident() uint8 {
	return 22
}

func (t QuitThought) ToBytes() (b []byte) {
	b = append(b, t.Ident())
	return
}

func (t QuitThought) FromBytes(b []byte) (net.Message, int) {
	return t, 1
}

func init() {
	net.RegisterMessage(Thoughts{})
	net.RegisterMessage(ResetThought{})
	net.RegisterMessage(QuitThought{})
	net.RegisterMessage(TickState{})
}

type TickState struct {
	Thoughts Thoughts
	Impulses ImpulseSet
}

func (t TickState) Ident() uint8 {
	return 23
}

func (t TickState) ToBytes() (b []byte) {
	b = append(b, t.Ident())
	b = append(b, t.Thoughts.ToBytes()...)
	b = append(b, t.Impulses.ToBytes()...)
	return
}

func (t TickState) FromBytes(b []byte) (net.Message, int) {
	offset := 1
	thoughts, n := net.MessageFromBytes(b[offset:])
	t.Thoughts = thoughts.(Thoughts)
	offset += n
	impulses, n := net.MessageFromBytes(b[offset:])
	t.Impulses = impulses.(ImpulseSet)
	offset += n
	return t, offset
}
