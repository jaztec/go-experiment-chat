package network

import (
	"bufio"
	"fmt"
	"net"
)

// Connection to a client with a reference to a write channel
type Connection struct {
	ID              string
	conn            *net.Conn
	Writes          chan Message
	Reads           chan Message
	CloseConnection chan byte
	reader          *bufio.Reader
	writer          *bufio.Writer
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
		buf, err := c.reader.ReadBytes('\n')
		if err != nil {
			fmt.Printf("Error occured: %v", err)
			SleepMS(1)
			break
		}

		c.Reads <- Message{ID: c.ID, Raw: buf, Type: IncomingMessage}
	}
}

// watchWrites coming into the write channel and send them to to client
func (c Connection) watchWrites() {
	for {
		select {
		case msg := <-c.Writes:
			if msg.Raw[len(msg.Raw)-1] != '\n' {
				msg.Raw = append(msg.Raw, '\n')
			}
			_, err := c.writer.Write(msg.Raw)
			if err != nil {
				continue
			}
			err = c.writer.Flush()
			if err != nil {
				continue
			}
		default:
			SleepMS(1)
		}
	}
}

// NewConnection for clients and/or servers
func NewConnection(conn *net.Conn) (*Connection, error) {
	fmt.Print("\nNew connection\n")
	id := RandString(32)
	c := &Connection{
		ID:              string(id),
		conn:            conn,
		CloseConnection: make(chan byte),
		Writes:          make(chan Message, 5),
		Reads:           make(chan Message, 5),
		reader:          bufio.NewReader(*conn),
		writer:          bufio.NewWriter(*conn)}

	go c.watchClose()
	go c.watchReads()
	go c.watchWrites()
	return c, nil
}
