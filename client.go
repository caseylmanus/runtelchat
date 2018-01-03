package runtelchat

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

const welcomePrompt = "Enter your chat name:"

type client struct {
	conn    net.Conn
	name    string
	address string
	inbox   chan Message
	outbox  chan Message
	closing chan bool
	closed  bool
}

func (c *client) waitForMessage() {

	reader := bufio.NewReader(c.conn)
	text, err := reader.ReadString('\n')
	if err != nil && !c.closed {
		c.closing <- true
	}

	if strings.HasPrefix(text, ".exit") && !c.closed {
		c.closing <- true
	}
	if !c.closed {
		c.outbox <- Message{Text: text, From: c.name, Address: c.address}
		c.waitForMessage()
	}
}

func (c *client) handleConnection() {
	for c.name == "" && !c.closed {
		c.conn.Write([]byte(welcomePrompt))
		buf := make([]byte, 4096)
		n, err := c.conn.Read(buf)
		if err != nil && !c.closed {
			c.closing <- true
		}
		if n > 0 {
			c.name = strings.TrimRight(string(buf[0:n]), "\r\n")
		}
	}
	c.outbox <- Message{
		Text:      "has joined.\r\n",
		TimeStamp: time.Now(),
		Address:   c.address,
		From:      c.name,
	}
	c.conn.Write([]byte("Send .exit to disconnect.\r\n"))
	go c.waitForMessage()
	for {
		select {
		case msg := <-c.inbox:
			c.sendMessage(msg)
		case <-c.closing:
			c.closed = true
			c.conn.Close()
			c.outbox <- Message{
				Text:      "has left the building.",
				From:      c.name,
				Address:   c.address,
				TimeStamp: time.Now(),
			}
		}
	}
}
func (c *client) sendMessage(msg Message) {
	c.conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
	_, err := c.conn.Write([]byte(fmt.Sprintf("%v: %v", msg.From, msg.Text)))
	if err != nil && !c.closed {
		c.closing <- true
	}
}
