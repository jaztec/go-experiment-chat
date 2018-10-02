package network

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"net"

	"github.com/jaztec/go-experiment-chat/errors"
)

// Connection to a client with a reference to a write channel
type Connection struct {
	ID              string
	conn            *net.Conn
	reader          *bufio.Reader
	writer          *bufio.Writer
	Writes          chan Message
	Reads           chan Message
	CloseConnection chan byte
}

// watchClose this connection upon message in close channel
func (c Connection) watchClose() {
	defer func() {
		(*c.conn).Close()
	}()
	<-c.CloseConnection
}

// watchReads coming into the read channel
func (c Connection) watchReads() {
	for {
		buff, err := c.reader.ReadBytes('\n')
		if errors.HasError(err) {
			break
		}

		c.Reads <- Message{ID: c.ID, Raw: buff, Type: IncomingMessage}
	}
}

// watchWrites coming into the write channel and send them to to client
func (c Connection) watchWrites() {
	for {
		select {
		case msg := <-c.Writes:
			c.writer.Write(msg.Raw)
		}
	}
}

// NewConnection for clients and/or servers
func NewConnection(conn *net.Conn) (*Connection, error) {
	fmt.Print("\nNew connection\n")
	id, err := generateRandomBytes(32)
	if errors.HasError(err) {
		return nil, err
	}
	c := &Connection{
		ID:              string(id),
		conn:            conn,
		writer:          bufio.NewWriter(*conn),
		reader:          bufio.NewReader(*conn),
		CloseConnection: make(chan byte),
		Writes:          make(chan Message, 16),
		Reads:           make(chan Message, 16)}

	go c.watchClose()
	go c.watchReads()
	go c.watchWrites()
	return c, nil
}

// generateRandomBytes to the amount of parameter n
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
