package menu

import (
	"ebijam23/net"
	"encoding/binary"
)

type HatMessage struct {
	Hat int
}

func (m HatMessage) Ident() uint8 {
	return 10
}

func (m HatMessage) ToBytes() (b []byte) {
	b = append(b, m.Ident())
	b = binary.LittleEndian.AppendUint64(b, uint64(m.Hat))
	return
}

func (m HatMessage) FromBytes(b []byte) (net.Message, int) {
	return HatMessage{Hat: int(binary.LittleEndian.Uint64(b[1:]))}, 9
}

type StartMessage struct {
}

func (m StartMessage) Ident() uint8 {
	return 11
}

func (m StartMessage) ToBytes() (b []byte) {
	b = append(b, m.Ident())
	return
}

func (m StartMessage) FromBytes(b []byte) (net.Message, int) {
	return m, 1
}
