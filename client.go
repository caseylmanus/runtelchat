package runtelchat

import (
	"net"
	"time"
	"fmt"
	"strings"
)
type client struct {
	conn   net.Conn
	name   string
	address string
	inbox  chan Message
	outbox chan Message
	closing chan bool 
	closed bool
}

func (c *client) waitForMessage() {
	buf := make([]byte, 4096)
	n, err := c.conn.Read(buf)
	if err != nil {
		c.closing <- true
	}
	if n > 0 {
		text := string(buf[0:n])
		if strings.HasPrefix(text, ".exit") {
			fmt.Println("Closing Channel?")
			c.closing <- true
		}
		c.outbox <- Message{Text: text, From: c.name, Address: c.address}
		if !c.closed {
			c.waitForMessage()
		}
	}
}

func (c *client) handleConnection() {
	for c.name == "" && !c.closed {
		c.conn.Write([]byte("Enter your chat name:"))
		buf := make([]byte, 4096)
		n, err := c.conn.Read(buf)
		if err != nil {
			c.closing <- true
		}
		if n > 0 {
			c.name = strings.TrimRight(string(buf[0:n]), "\r\n")
		}
	}
	c.conn.Write([]byte("Send .exit to disconnect.\r\n")) 
	go c.waitForMessage()
	for {
		select {
		case msg := <-c.inbox:
			c.conn.SetWriteDeadline(time.Now().Add(time.Second * 1))
			_, err := c.conn.Write([]byte(fmt.Sprintf("%v: %v", msg.From, msg.Text)))
			if err != nil {
				c.closing <- true
			}
		case <- c.closing:
			c.closed = true
			c.outbox <- Message{Text: fmt.Sprintf("%v has left the chat.", c.name), From: "System", Address: c.address, TimeStamp: time.Now()}
			c.conn.Close()
			return
		}
	}
}