package net

type Event interface {
}

type EventClosed struct {
}

type EventHosting struct {
}

type EventJoining struct {
}

type EventJoined struct {
}

type EventConnect struct {
	ID   uint32
	Peer *Peer
}

type EventDisconnect struct {
	ID   uint32
	Peer *Peer
}

type EventMessage struct {
	ID      uint32
	Peer    *Peer
	Message Message
}
