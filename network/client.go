package network

import (
	"net"
)

// ClientConfig containing settings for the client
type ClientConfig struct {
}

// ClientClass client object
type ClientClass struct {
	network string
	address string
	conn    *Connection
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
	if err != nil {
		return nil, err
	}

	print("Create new connection from client")
	connection, err := NewConnection(&conn)
	if err != nil {
		return nil, err
	}
	c.conn = connection

	return connection.Reads, nil
}

// Close the connection to the server
func (c ClientClass) Close() {
	if c.conn != nil {
		c.conn.CloseConnection <- byte(1)
	}
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
	return c, nil
}
