package net

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"github.com/xtaci/kcp-go"
)

var NetDataShards = 10
var NetParityShards = 3
var NetBufferSize = 2048
var NetChannelSize = 10

type ServerClient struct {
	id uint32 // Our ID used to advertise to others.
	//
	Matchmaker     string
	matchmakerAddr *net.UDPAddr
	UseMatchmaker  bool
	Hosting        bool
	Running        bool
	localAddr      *net.UDPAddr
	localConn      *net.UDPConn
	closeChan      chan struct{}
	rawChan        chan Packet
	peerChan       chan PeerPacket
	EventChan      chan Event
	peers          []*Peer
}

func (s *ServerClient) Init() {
	s.id = uint32(rand.Int31())
	s.closeChan = make(chan struct{})
	s.rawChan = make(chan Packet, NetChannelSize*2)
	s.peerChan = make(chan PeerPacket, NetChannelSize)
	s.EventChan = make(chan Event, NetChannelSize)
	s.Matchmaker = "gamu.group:20220"
}

func (s *ServerClient) ID() uint32 {
	return s.id
}

func (s *ServerClient) Open(address string) error {
	// Set up our matchmaker address.
	matchMakerAddr, err := net.ResolveUDPAddr("udp", s.Matchmaker)
	if err == nil {
		s.matchmakerAddr = matchMakerAddr
	} else {
		fmt.Println("matchmaker's address is bad", err)
	}

	// Get our UDP address.
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}

	// Open our UDP connection.
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	s.localAddr = addr
	s.localConn = conn
	s.Running = true

	fmt.Println("...now listening on", conn.LocalAddr().String())

	s.EventChan <- EventHosting{}

	go s.LogicLoop()
	go s.ReadLoop()

	return nil
}

func (s *ServerClient) ConnectTo(address string) error {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}
	peer := NewPeer(addr, s.localConn)
	s.peers = append(s.peers, peer)

	session, err := kcp.NewConn3(0, addr, nil, NetDataShards, NetParityShards, peer)
	if err != nil {
		panic(err)
	}
	peer.session = session
	go peer.loop(s.peerChan)

	s.EventChan <- EventJoining{}

	peer.Send(MessageID{ID: s.id})

	return nil
}

type Packet struct {
	buffer    []byte
	addr      *net.UDPAddr
	readBytes int
}

// LogicLoop is the main logic loop that handles raw packets and otherwise.
func (s *ServerClient) LogicLoop() {
	for s.Running {
		select {
		case <-s.closeChan:
			for _, p := range s.peers {
				p.session.Close()
			}
			s.localConn.Close()
			s.Running = false
			s.EventChan <- EventClosed{}
			return
		case msg := <-s.peerChan:
			switch msg.msg.(type) {
			case MessageID:
				if msg.peer.id == 0 {
					if !s.Hosting {
						s.EventChan <- EventJoined{}
					}
				}
				msg.peer.id = msg.msg.(MessageID).ID
				s.EventChan <- EventConnect{
					Peer: msg.peer,
					ID:   msg.peer.id,
				}
			case MessageClose:
				s.EventChan <- EventDisconnect{
					Peer: msg.peer,
					ID:   msg.peer.id,
				}
			default:
				s.EventChan <- EventMessage{
					Peer:    msg.peer,
					Message: msg.msg,
				}
			}
		case packet := <-s.rawChan:
			// If the packet is from the matchmaker, handle it.
			if s.matchmakerAddr != nil && packet.addr.String() == s.matchmakerAddr.String() {
				fmt.Println("handling as matchmaker")
				// TODO: Handle matchmaker packets.
				return
			}
			// If we've never received from this address before, add it to our list of peers.
			var peer *Peer
			for _, p := range s.peers {
				if p.addr.String() == packet.addr.String() {
					peer = p
					break
				}
			}
			if peer == nil {
				peer = NewPeer(packet.addr, s.localConn)
				s.peers = append(s.peers, peer)
			}

			if packet.readBytes > 0 {
				peer.writeToPacketBuffer(packet.buffer[:packet.readBytes])
			}

			// Session is nil, try to set up a kcp session.
			if peer.session == nil {
				session, err := kcp.NewConn3(0, packet.addr, nil, NetDataShards, NetParityShards, peer)
				if err != nil {
					panic(err)
				}
				peer.session = session
				go peer.loop(s.peerChan)

				peer.Send(MessageID{ID: s.id})
			} else {
				// Send packet data to peer's virtual buffer.
				peer.writeToPacketBuffer(packet.buffer[:packet.readBytes])
			}
		}
	}
}

// ReadLoop runs the network read loop, handling new connections as necessary.
func (s *ServerClient) ReadLoop() {
	for s.Running {
		buffer := make([]byte, NetBufferSize)
		n, addr, err := s.localConn.ReadFromUDP(buffer)
		if err != nil {
			if !os.IsTimeout(err) {
				if !(strings.HasSuffix(err.Error(), "closed network connection")) {
					fmt.Println(err)
				}
			}
			break
		}
		packet := Packet{
			buffer:    buffer,
			addr:      addr,
			readBytes: n,
		}
		s.rawChan <- packet
	}
}

func (s *ServerClient) Close() {
	if !s.Running {
		return
	}
	for _, peer := range s.peers {
		peer.Send(MessageClose{})
		s.EventChan <- EventDisconnect{
			Peer: peer,
			ID:   peer.id,
		}
	}
	<-time.After(1 * time.Second)
	s.closeChan <- struct{}{}
}

func (s *ServerClient) Peers() []*Peer {
	return s.peers
}
