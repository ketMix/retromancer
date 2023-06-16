package net

import (
	"fmt"

	"github.com/xtaci/kcp-go"
)

type Conn struct {
	ID        int32
	session   *kcp.UDPSession
	connected bool
}

func (c *Conn) Loop() {
	for c.connected {
		var buf [4096]byte
		n, err := c.session.Read(buf[:])
		if err != nil {
			c.session.Close()
			// TODO: Send closed on channel.
			return
		}
		// TODO: Convert buf to an actual message.
		fmt.Println("got message: ", string(buf[:n]))
	}
}
