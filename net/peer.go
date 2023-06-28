package net

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"
)

type Peer struct {
	id   uint32
	addr *net.UDPAddr // Address of the peer
	conn *net.UDPConn // Pointer to the serverclient's conn.
	//session *kcp.UDPSession
	shook bool
	// Packet reading.
	packetBuffer     []byte
	messageIndex     int
	outboundMessages []*Envelope
	reconfirms       []int
	inboundConfirms  []int
	inboundMessages  []*Envelope
	messages         []Message
	readLock         sync.Mutex
	readReadyChan    chan bool
}

type PeerPacket struct {
	peer *Peer
	msg  Message
}

func NewPeer(addr *net.UDPAddr, conn *net.UDPConn) *Peer {
	return &Peer{
		addr:          addr,
		conn:          conn,
		readReadyChan: make(chan bool, NetChannelSize),
	}
}

func (p *Peer) ID() uint32 {
	return p.id
}

// writeToPacketBuffer is used internally to write from the single UDP connection to a virtual packet buffer for use by the Peer.
func (p *Peer) writeToPacketBuffer(b []byte) {
	p.readLock.Lock()
	p.packetBuffer = append(p.packetBuffer, b...)
	p.readLock.Unlock()
	p.readReadyChan <- true
}

/*func (p *Peer) loop(ch chan PeerPacket) {
	for {
		b := make([]byte, NetBufferSize)
		n, err := p.session.Read(b)

		msg, _ := MessageFromBytes(b[:n])
		if msg == nil {
			fmt.Println("unknown message ID:", b[0], "passing as raw message.")
			msg, _ = MessageRaw{}.FromBytes(b[:n])
		}

		if err != nil {
			fmt.Println(err)
			return
		}
		ch <- PeerPacket{
			peer: p,
			msg:  msg,
		}
	}
}*/

type Envelope struct {
	ID        int
	lastTime  time.Time
	bytes     []byte
	confirmed bool
}

func (e Envelope) ToBytes() (b []byte) {
	b = append(b, 0)
	b = binary.LittleEndian.AppendUint64(b, uint64(e.ID))
	b = append(b, e.bytes...)
	return b
}

type Confirm struct {
	ID int
}

func (c Confirm) ToBytes() (b []byte) {
	b = append(b, 1)
	b = binary.LittleEndian.AppendUint64(b, uint64(c.ID))
	return b
}

type Reconfirm struct {
	ID int
}

func (c Reconfirm) ToBytes() (b []byte) {
	b = append(b, 2)
	b = binary.LittleEndian.AppendUint64(b, uint64(c.ID))
	return b
}

// Send sends a Payload to the given peer.
func (p *Peer) Send(msg Message) error {
	p.readLock.Lock()
	defer p.readLock.Unlock()
	envelope := Envelope{
		ID:    p.messageIndex,
		bytes: msg.ToBytes(),
	}
	p.messageIndex++
	p.outboundMessages = append(p.outboundMessages, &envelope)
	return nil
}

func (p *Peer) Receive(b []byte) (err error) {
	//fmt.Println("Receive", b)
	p.readLock.Lock()
	defer p.readLock.Unlock()
	if len(b) == 0 {
		return
	}
	if b[0] == 0 {
		var envelope Envelope
		envelope.ID = int(binary.LittleEndian.Uint64(b[1:9]))
		envelope.bytes = b[9:]

		// If we already have this envelope, mark it is confirmed.
		/*for _, env := range p.inboundMessages {
			if env.ID == envelope.ID {
				env.confirmed = true
				return
			}
		}*/

		p.inboundMessages = append(p.inboundMessages, &envelope)
	} else if b[0] == 1 {
		var confirm Confirm
		confirm.ID = int(binary.LittleEndian.Uint64(b[1:9]))
		p.inboundConfirms = append(p.inboundConfirms, confirm.ID)
		//fmt.Println("got inbound confirm", confirm.ID)
	} else if b[0] == 2 {
		var confirm Reconfirm
		confirm.ID = int(binary.LittleEndian.Uint64(b[1:9]))
		p.reconfirms = append(p.reconfirms, confirm.ID)
		//fmt.Println("got inbound reconfirm")
	}
	return nil
}

/*func (p *Peer) Receive(b []byte) (msg Message, err error) {
	msg, _ = MessageFromBytes(b)
	if msg == nil {
		fmt.Println("unknown message ID:", b[0], "passing as raw message.")
		msg, _ = MessageRaw{}.FromBytes(b)
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return msg, nil
}*/

// ReadFrom is used to read from the peer's virtual packet buffer.
func (p *Peer) ReadFrom(b []byte) (n int, addr net.Addr, err error) {
	<-p.readReadyChan
	addr = p.addr
	p.readLock.Lock()
	n = copy(b, p.packetBuffer)
	p.packetBuffer = p.packetBuffer[n:]
	p.readLock.Unlock()
	return
}

// WriteTo writes the bytes to the given address.
func (p *Peer) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	n, err = p.conn.WriteTo(b, addr)
	return
}

// Close closes the peer connection. (it actually does nothing!)
func (p *Peer) Close() error {
	return nil
}

// LocalAddr returns the connection's local address.
func (p *Peer) LocalAddr() net.Addr {
	return p.addr
}

// SetDeadline does nothing.
func (p *Peer) SetDeadline(t time.Time) error {
	fmt.Println("TODO: Peer.SetDeadline")
	return nil
}

// SetReadDeadline does nothing.
func (p *Peer) SetReadDeadline(t time.Time) error {
	fmt.Println("TODO: Peer.SetReadDeadline")
	return nil
}

// SetWriteDeadline does nothing.
func (p *Peer) SetWriteDeadline(t time.Time) error {
	fmt.Println("TODO: Peer.SetWriteDeadline")
	return nil
}
