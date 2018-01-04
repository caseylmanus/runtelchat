package runtelchat

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/pkg/errors"
)

//ChatServer is a TCP chat server capable of running on multiple ports
type ChatServer struct {
	clients []*client
	outbox  chan Message
	Errors  chan error
	closing chan *client
	sync.RWMutex
}

//NewChatServer instatiates a new chat server
func NewChatServer() *ChatServer {
	return &ChatServer{
		clients: []*client{},
		outbox:  make(chan Message),
		Errors:  make(chan error),
		closing: make(chan *client),
	}
}

//Serve blocks and coordinates between listeners the server
func (server *ChatServer) Serve() error {
	for {
		select {
		case msg := <-server.outbox:
			log.Println(msg)
			server.RLock()
			for i := range server.clients {
				if server.clients[i] != msg.client {
					server.clients[i].inbox <- msg
				}
			}
			server.RUnlock()
		case err := <-server.Errors:
			return err
		case c := <-server.closing:
			server.closeClient(c)
		}
	}
}

//Listen will accept connections on a the listener
func (server *ChatServer) Listen(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			server.Errors <- err
		}
		server.registerClient(conn, conn.RemoteAddr().String())
	}
}

//ListenAndServe opens tcp listeners for the given the specified configuration
func ListenAndServe(config Config) error {
	server := NewChatServer()
	url := fmt.Sprint(config.Host, ":", config.Port)
	listener, err := net.Listen("tcp", url)
	if err != nil {
		return errors.Wrap(err, fmt.Sprint("Failed to open listener on ", url))
	}
	go server.Listen(listener)

	return server.Serve()
}

func (server *ChatServer) registerClient(conn net.Conn, remoteAddress string) {
	server.Lock()

	c := &client{
		conn:    conn,
		name:    "",
		inbox:   make(chan Message),
		outbox:  server.outbox,
		closing: server.closing,
	}

	server.clients = append(server.clients, c) 
	server.Unlock()
	go c.handleConnection()

}

func (server *ChatServer) closeClient(c *client) {
	server.Lock()
	copy := []*client{} 
	for i := range server.clients {
		if server.clients[i] != c  {
			copy = append(copy, server.clients[i])
		}
	}
	server.clients = copy
	server.Unlock()
	go func() {
		server.outbox <- Message {
			Text: fmt.Sprintf("%v has left the building.\r\n", c.name),
			TimeStamp: time.Now(),
			Address: c.address,
			client: c,
			From: "system",
		}
	}()
	c.conn.Close()
}
