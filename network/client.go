package network

import (
	"errors"
)

// ClientConfig containing settings for the client
type ClientConfig struct {
}

// ClientClass client object
type ClientClass struct {
}

// ClientInterface defining client functions
type ClientInterface interface {
	Dial() (chan Message, error)
	Close() bool
}

// Dial in with the server
func (c ClientClass) Dial() (chan Message, error) {
	return nil, errors.New("Not implemented yet")
}

// Close the connection to the server
func (c ClientClass) Close() bool {
	return true
}

// NewClient returns a pointer to a new client
func NewClient() (ClientInterface, error) {
	c := new(ClientClass)
	return c, nil
}
