package runtelchat

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/pkg/errors"
)

//ChatServer is a TCP chat server capable of running on multiple ports
type ChatServer struct {
	clients []*client
	outbox  chan Message
	Errors  chan error
	sync.Mutex
}

//NewChatServer instatiates a new chat server
func NewChatServer() *ChatServer {
	return &ChatServer{
		clients: []*client{},
		outbox:  make(chan Message),
		Errors:  make(chan error),
	}
}

//Serve blocks and coordinates between listeners the server
func (server *ChatServer) Serve() error {
	for {
		select {
		case msg := <-server.outbox:
			log.Println(msg)
			for i := range server.clients {
				if server.clients[i].name != msg.From && !server.clients[i].closed {
					server.clients[i].inbox <- msg
				}
			}
		case err := <-server.Errors:
			return err
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
	defer server.Unlock()
	c := &client{
		conn:    conn,
		name:    "",
		inbox:   make(chan Message),
		outbox:  server.outbox,
		closing: make(chan bool),
	}

	server.clients = append(server.clients, c)
	go c.handleConnection()

}
