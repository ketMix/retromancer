package net

import "encoding/gob"

type Message interface {
}

type MessageID struct {
	ID uint32
}

func (m MessageID) Type() string {
	return "ID"
}

type MessageClose struct {
}

func init() {
	gob.Register(MessageID{})
	gob.Register(MessageClose{})
}
