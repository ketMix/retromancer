package net

import (
	"encoding/binary"
	"reflect"
)

var messageRegistry = map[uint8]Message{}

func RegisterMessage(m Message) {
	messageRegistry[m.Ident()] = m
}

func MessageFromBytes(b []byte) (Message, int) {
	if messageRegistry[b[0]] != nil {
		msg := reflect.New(reflect.ValueOf(messageRegistry[b[0]]).Type()).Interface().(Message)
		msg, n := msg.FromBytes(b)
		return msg, n
	}
	return nil, 0
}

type Message interface {
	Ident() uint8
	ToBytes() []byte
	FromBytes(b []byte) (Message, int)
}

type MessageID struct {
	ID uint32
}

func (m MessageID) Type() string {
	return "ID"
}

func (m MessageID) Ident() uint8 {
	return 1
}

func (m MessageID) ToBytes() (b []byte) {
	b = append(b, m.Ident())
	b = binary.LittleEndian.AppendUint32(b, m.ID)
	return b
}

func (m MessageID) FromBytes(b []byte) (Message, int) {
	return MessageID{ID: binary.LittleEndian.Uint32(b[1:])}, 5
}

type MessageClose struct {
}

func (m MessageClose) Ident() uint8 {
	return 2
}

func (m MessageClose) FromBytes(b []byte) (Message, int) {
	return m, 1
}

func (m MessageClose) ToBytes() []byte {
	return []byte{2}
}

type MessageRaw struct {
	Data []byte
}

func (m MessageRaw) Ident() uint8 {
	return 3
}

func (m MessageRaw) ToBytes() (b []byte) {
	b = append(b, m.Ident())
	b = append(b, m.Data...)
	return
}

func (m MessageRaw) FromBytes(b []byte) (Message, int) {
	return MessageRaw{Data: b[1:]}, 1 + len(b)
}

func init() {
	RegisterMessage(MessageID{})
	RegisterMessage(MessageClose{})
	RegisterMessage(MessageRaw{})
}
