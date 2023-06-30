package game

import (
	"ebijam23/net"
	"encoding/binary"
	"math"
)

func init() {
	net.RegisterMessage(ImpulseMove{})
	net.RegisterMessage(ImpulseSet{})
	net.RegisterMessage(ImpulseReverse{})
	net.RegisterMessage(ImpulseDeflect{})
	net.RegisterMessage(ImpulseShield{})
	net.RegisterMessage(ImpulseShoot{})
}

type ImpulseSet struct {
	Move        *ImpulseMove
	Interaction Impulse
}

type Impulse interface {
	Cost() int
	Ident() uint8
	ToBytes() []byte
	FromBytes([]byte) (net.Message, int)
}

type ImpulseMove struct {
	Direction float64
}

type ImpulseReverse struct {
	X, Y float64 // X and Y are the current cursor coordinates.
}

func (i ImpulseReverse) Cost() int {
	return 1
}

type ImpulseDeflect struct {
	X, Y float64 // X and Y are the current cursor coordinates.
}

func (i ImpulseDeflect) Cost() int {
	return 2
}

type ImpulseShield struct {
}

func (i ImpulseShield) Cost() int {
	return 1
}

type ImpulseShoot struct {
	X, Y float64 // X and Y are the current cursor coordinates.
}

func (i ImpulseShoot) Cost() int {
	return 6
}

// Networking crap

func (i ImpulseMove) Ident() uint8 {
	return 30
}

func (i ImpulseMove) ToBytes() (b []byte) {
	b = append(b, i.Ident())
	b = binary.LittleEndian.AppendUint64(b, uint64(math.Float64bits(i.Direction)))
	return
}

func (i ImpulseMove) FromBytes(b []byte) (net.Message, int) {
	i.Direction = math.Float64frombits(binary.LittleEndian.Uint64(b[1:]))
	return i, 9
}

func (i ImpulseReverse) Ident() uint8 {
	return 31
}

func (i ImpulseReverse) ToBytes() (b []byte) {
	b = append(b, i.Ident())
	b = binary.LittleEndian.AppendUint64(b, uint64(math.Float64bits(i.X)))
	b = binary.LittleEndian.AppendUint64(b, uint64(math.Float64bits(i.Y)))
	return
}

func (i ImpulseReverse) FromBytes(b []byte) (net.Message, int) {
	i.X = math.Float64frombits(binary.LittleEndian.Uint64(b[1:]))
	i.Y = math.Float64frombits(binary.LittleEndian.Uint64(b[9:]))
	return i, 17
}

func (i ImpulseDeflect) Ident() uint8 {
	return 32
}

func (i ImpulseDeflect) ToBytes() (b []byte) {
	b = append(b, i.Ident())
	b = binary.LittleEndian.AppendUint64(b, uint64(math.Float64bits(i.X)))
	b = binary.LittleEndian.AppendUint64(b, uint64(math.Float64bits(i.Y)))
	return
}

func (i ImpulseDeflect) FromBytes(b []byte) (net.Message, int) {
	i.X = math.Float64frombits(binary.LittleEndian.Uint64(b[1:]))
	i.Y = math.Float64frombits(binary.LittleEndian.Uint64(b[9:]))
	return i, 17
}

func (i ImpulseShield) Ident() uint8 {
	return 33
}

func (i ImpulseShield) ToBytes() (b []byte) {
	b = append(b, i.Ident())
	return
}

func (i ImpulseShield) FromBytes(b []byte) (net.Message, int) {
	return i, 1
}

func (i ImpulseShoot) Ident() uint8 {
	return 34
}

func (i ImpulseShoot) ToBytes() (b []byte) {
	b = append(b, i.Ident())
	b = binary.LittleEndian.AppendUint64(b, uint64(math.Float64bits(i.X)))
	b = binary.LittleEndian.AppendUint64(b, uint64(math.Float64bits(i.Y)))
	return
}

func (i ImpulseShoot) FromBytes(b []byte) (net.Message, int) {
	i.X = math.Float64frombits(binary.LittleEndian.Uint64(b[1:]))
	i.Y = math.Float64frombits(binary.LittleEndian.Uint64(b[9:]))
	return i, 17
}

func (i ImpulseSet) Ident() uint8 {
	return 35
}

func (i ImpulseSet) ToBytes() (b []byte) {
	b = append(b, i.Ident())
	if i.Move != nil {
		b = append(b, byte(1))
		b = append(b, i.Move.ToBytes()...)
	} else {
		b = append(b, byte(0))
	}
	if i.Interaction != nil {
		b = append(b, byte(1))
		b = append(b, i.Interaction.ToBytes()...)
	}
	return
}

func (i ImpulseSet) FromBytes(b []byte) (net.Message, int) {
	i.Move = nil
	i.Interaction = nil
	offset := 1
	if b[offset] == 1 {
		offset++
		m, n := (ImpulseMove{}).FromBytes(b[offset:])
		argh := m.(ImpulseMove)
		i.Move = &argh
		offset += n
	} else {
		offset++
	}
	if len(b) == offset {
		return i, offset
	}
	if b[offset] == 1 {
		offset++
		switch b[offset] {
		case (ImpulseReverse{}).Ident():
			m, n := (ImpulseReverse{}).FromBytes(b[offset:])
			i.Interaction = m.(Impulse)
			offset += n
		case (ImpulseDeflect{}).Ident():
			m, n := (ImpulseDeflect{}).FromBytes(b[offset:])
			i.Interaction = m.(Impulse)
			offset += n
		case (ImpulseShield{}).Ident():
			m, n := (ImpulseShield{}).FromBytes(b[offset:])
			i.Interaction = m.(Impulse)
			offset += n
		case (ImpulseShoot{}).Ident():
			m, n := (ImpulseShoot{}).FromBytes(b[offset:])
			i.Interaction = m.(Impulse)
			offset += n
		}
	}
	return i, offset
}
