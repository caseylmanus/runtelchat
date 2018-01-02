package runtelchat

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/pkg/errors"
)

//ServeTCP opens tcp listeners for the given the specified configuration
func ServeTCP(config Config) error {
	clients := []*client{}
	outbox := make(chan Message)
	for _, port := range config.Ports {
		url := fmt.Sprint(config.Host, ":", port)
		listener, err := net.Listen("tcp", url)
		if err != nil {
			return errors.Wrap(err, fmt.Sprint("Failed to open listener on ", url))
		}
		defer listener.Close()

		go func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					log.Fatal(errors.Wrap(err, "Failed to handle new connection."))
				}
				c := &client{
					conn:   conn,
					name:   "",
					inbox:  make(chan Message),
					outbox: outbox,
				}
				clients = append(clients, c)
				go handleConnection(c)

			}
		}()
	}
	for {
		select {
		case msg := <-outbox:
			log.Println(msg)
			for i := range clients {
				if clients[i].name != msg.From {
					clients[i].inbox <- msg
				}
			}
		}
	}
}

func handleConnection(c *client) {
	for c.name == "" {
		c.conn.Write([]byte("Enter your chat name:"))
		buf := make([]byte, 4096)
		n, err := c.conn.Read(buf)
		if err != nil {
			log.Println(errors.Wrap(err, "Failed to read from channel"))
		}
		if n > 0 {
			c.name = strings.TrimRight(string(buf[0:n]), "\r\n")
		}
	}
	go waitForMessage(c)
	for {

		select {
		case msg := <-c.inbox:
			c.conn.Write([]byte(fmt.Sprintf("%v: %v", msg.From, msg.Text)))
		}
	}
}
func waitForMessage(c *client) {
	buf := make([]byte, 4096)
	n, err := c.conn.Read(buf)
	if err != nil {
		log.Println(errors.Wrap(err, "Failed to read from channel"))
	}
	if n > 0 {
		text := string(buf[0:n])
		c.outbox <- Message{Text: text, From: c.name, Address: c.conn.RemoteAddr().String()}
		waitForMessage(c)
	}
}

type client struct {
	conn   net.Conn
	name   string
	inbox  chan Message
	outbox chan Message
}

//Message represents a chat message
type Message struct {
	Text      string
	From      string
	Address   string
	TimeStamp time.Time
}

func (m Message) String() string {
	data, _ := json.Marshal(m)
	return string(data)
}
