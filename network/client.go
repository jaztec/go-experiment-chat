package network

import (
	"net"

	"github.com/jaztec/go-experiment-chat/errors"
)

// ClientConfig containing settings for the client
type ClientConfig struct {
}

// ClientClass client object
type ClientClass struct {
	network string
	address string
	conn    Connection
}

// ClientInterface defining client functions
type ClientInterface interface {
	Dial() (chan Message, error)
	Close()
	Send(Message)
	CreateMessage([]byte) Message
}

// Dial in with the server
func (c *ClientClass) Dial() (chan Message, error) {
	conn, err := net.Dial(c.network, c.address)
	if errors.HasError(err) {
		return nil, err
	}

	connection, err := NewConnection(&conn)
	if errors.HasError(err) {
		return nil, err
	}

	return connection.Reads, nil
}

// Close the connection to the server
func (c ClientClass) Close() {
	c.conn.CloseConnection <- byte(1)
}

// Send a message to the server
func (c ClientClass) Send(message Message) {
	c.conn.Writes <- message
}

// CreateMessage to be send to the server
func (c ClientClass) CreateMessage(raw []byte) Message {
	return Message{
		Raw:  raw,
		Type: OutgoingMessage}
}

// NewClient returns a pointer to a new client
func NewClient(network string, address string) (ClientInterface, error) {
	c := &ClientClass{
		network: network,
		address: address}
	_, err := c.Dial()
	if errors.HasError(err) {
		return nil, err
	}
	return c, nil
}
