package network

import (
	"errors"
	"fmt"
	"net"
	"strconv"
)

var debug bool

func init() {
	// Make sure to include some better env variable integration
	debug = true
}

// ServerConfig containing settings for the server
type ServerConfig struct {
	MessageBufferSize int16
	Port              int16
}

// ServerClass server object
type ServerClass struct {
	messageBufferSize int16
	conns             map[string]*Connection
	messages          []byte
	port              int16
	server            net.Listener
	close             chan bool
	reads             chan Message
}

// ServerInterface defined in the server object
type ServerInterface interface {
	Listen() (chan Message, error)
	Close() error
	Send(string, Message) error
	Broadcast(Message)
}

// Listen to TCP connections
func (s *ServerClass) Listen() (chan Message, error) {
	var err error
	// Listen on all interfaces to port
	s.server, err = net.Listen("tcp", ":"+strconv.FormatInt(int64(s.port), 10))
	if err != nil {
		return nil, err
	}
	s.reads = make(chan Message, s.messageBufferSize)
	go s.run()
	fmt.Printf("Listening to port %d\n", s.port)
	return s.reads, nil
}

// Close the server
func (s ServerClass) Close() error {
	fmt.Println("Close the server")
	s.close <- true
	err := s.server.Close()
	for id := range s.conns {
		s.conns[id].CloseConnection <- byte(1)
	}
	return err
}

// Send message to a client
func (s ServerClass) Send(id string, message Message) error {
	conn := s.conns[id]
	if conn == nil {
		return errors.New(id + ": ID not found")
	}
	conn.Writes <- message
	return nil
}

// Broadcast message to all clients
func (s ServerClass) Broadcast(message Message) {
	for id := range s.conns {
		s.conns[id].Writes <- message
	}
}

func (s *ServerClass) run() {
	for {
		conn, err := s.server.Accept()
		if err != nil {
			continue
		}

		print("Create new connection from server")
		connection, err := NewConnection(&conn)
		if err != nil {
			continue
		}

		print(fmt.Sprintf("\nNew connection: %s\n", connection.ID))
		s.conns[connection.ID] = connection

		go func() {
			for {
				select {
				case msg := <-connection.Reads:
					s.reads <- msg
				}
			}
		}()
	}
}

// NewServer network.ServerInterface instance pointer
func NewServer(config ServerConfig) ServerInterface {
	s := &ServerClass{
		messageBufferSize: config.MessageBufferSize,
		port:              config.Port,
		conns:             make(map[string]*Connection)}
	return s
}

func print(message string) {
	if debug {
		fmt.Print(message)
	}
}
