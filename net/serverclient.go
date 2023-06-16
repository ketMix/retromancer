package net

import (
	"encoding/binary"
	"fmt"
	"math/rand"

	"github.com/xtaci/kcp-go"
)

type ServerClient struct {
	id          int32 // Our ID used to advertise to others.
	serverID    int32 // The ID of the server we're connected to.
	listener    *kcp.Listener
	session     *kcp.UDPSession
	connections []*Conn
}

func (s *ServerClient) Send(b []byte) error {
	if s.listener != nil {
		for _, c := range s.connections {
			n, err := c.session.Write(b)
			if err != nil {
				return err
			}
			if n != len(b) {
				return fmt.Errorf("failed to send all bytes")
			}
		}
	} else if s.session != nil {
		n, err := s.session.Write(b)
		if err != nil {
			return err
		}
		if n != len(b) {
			return fmt.Errorf("failed to send all bytes")
		}
	}
	return fmt.Errorf("not connected")
}

func (s *ServerClient) Listen(address string) error {
	if s.listener != nil {
		return fmt.Errorf("server already listening")
	}
	listener, err := kcp.ListenWithOptions(address, nil, 10, 3)
	if err != nil {
		return err
	}
	s.listener = listener

	go s.AcceptLoop()

	return nil
}

func (s *ServerClient) Connect(address string) (err error) {
	if s.session != nil {
		return fmt.Errorf("already connected")
	}
	s.session, err = kcp.DialWithOptions(address, nil, 10, 3)
	if err != nil {
		return err
	}

	// Await the response ID from the server.
	s.serverID, err = s.readID(s.session)
	if err != nil {
		return err
	}

	// Send a random number as our first message to the server to identify ourselves as that particular number. WHY NOT
	s.id = rand.Int31()

	// Write the ID as an integer in little endian to the session using the binary package.
	binary.Write(s.session, binary.LittleEndian, s.id)

	go s.ReadLoop()

	return nil
}

func (s *ServerClient) readID(session *kcp.UDPSession) (int32, error) {
	data := make([]byte, 4)
	n, err := session.Read(data)
	if err != nil {
		return 0, err
	}
	if n != 4 {
		return 0, fmt.Errorf("failed to read 4 bytes for the id")
	}
	return int32(binary.LittleEndian.Uint32(data)), nil
}

func (s *ServerClient) AcceptLoop() {
	// Create our server id.
	s.id = rand.Int31()

	defer s.listener.Close()
	for s.listener != nil {
		session, err := s.listener.AcceptKCP()
		if err != nil {
			fmt.Println(err)
			return
		}
		conn := &Conn{
			session: session,
		}

		// Send our ID to the client.
		binary.Write(session, binary.LittleEndian, s.id)

		// Read in the ID of the client.
		id, err := s.readID(session)
		if err != nil {
			session.Close()
			continue
		}

		conn.ID = id

		s.connections = append(s.connections, conn)

		go s.ReadConnLoop(conn)
	}
}

func (s *ServerClient) ReadConnLoop(c *Conn) {
	for c.connected {
		var buf [4096]byte
		n, err := c.session.Read(buf[:])
		if err != nil {
			c.session.Close()
		}
		fmt.Println("read", n, buf[:n])
	}
}

func (s *ServerClient) ReadLoop() {
	for s.session != nil {
		var buf [4096]byte
		n, err := s.session.Read(buf[:])
		if err != nil {
			s.session.Close()
			return
		}
		fmt.Println("read", n, buf[:n])
	}
}
