package net

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/xtaci/kcp-go"
)

type Peer struct {
	id      uint32
	addr    *net.UDPAddr // Address of the peer
	conn    *net.UDPConn // Pointer to the serverclient's conn.
	session *kcp.UDPSession
	// Packet reading.
	packetBuffer  []byte
	readLock      sync.Mutex
	readReadyChan chan bool
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

func (p *Peer) loop(ch chan PeerPacket) {
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
}

// Send sends a Payload to the given peer.
func (p *Peer) Send(msg Message) error {
	if p.session == nil {
		return fmt.Errorf("no sesssion")
	}
	//p.conn.WriteTo(msg.ToBytes(), p.addr)
	p.session.Write(msg.ToBytes())
	//p.encoder.Encode(&msg)
	return nil
}

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
