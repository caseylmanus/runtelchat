package runtelchat

import (
	"fmt"
	"log"
	"net"

	"github.com/pkg/errors"
)


type ChatServer struct {
	clients []*client
	outbox chan Message
	Errors chan error
}

func NewChatServer() *ChatServer { 
	return &ChatServer{[]*client{}, make(chan Message), make(chan error)}
}
func(server *ChatServer) registerClient(conn net.Conn, remoteAddress string) {
	c := &client{
		conn : conn,
		name : "",
		inbox: make(chan Message),
		outbox : server.outbox,
	}
	server.clients = append(server.clients, c)
	go c.handleConnection()
}
func(server *ChatServer) Serve() error {
	for {
		select {
		case msg := <-server.outbox:
			log.Println(msg)
			for i := range server.clients {
				if server.clients[i].name != msg.From && !server.clients[i].closed {
					server.clients[i].inbox <- msg
				}
			}
		case err := <- server.Errors :
			return err
		}
	}
}
func(server *ChatServer) Listen(listener net.Listener) {
	for {
		conn , err := listener.Accept() 
		if err != nil {
			server.Errors <- err
		}
		server.registerClient(conn, conn.RemoteAddr().String())
	}
}



//ServeTCP opens tcp listeners for the given the specified configuration
func ServeTCP(config Config) error {
	server := NewChatServer() 
	for _, port := range config.Ports {
		url := fmt.Sprint(config.Host, ":", port)
		listener, err := net.Listen("tcp", url)
		if err != nil {
			return errors.Wrap(err, fmt.Sprint("Failed to open listener on ", url))
		}
		go server.Listen(listener)
	}
	return server.Serve()
}


