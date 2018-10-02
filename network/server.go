package network

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	appErr "github.com/jaztec/go-experiment-chat/error"
)

// ServerConfig containing settings for the server
type ServerConfig struct {
	MessageBufferSize int16
	Port              int16
}

// ServerClass server object
type ServerClass struct {
	messageBufferSize int16
	channel           chan Message
	conns             map[string]*net.Conn
	messages          []byte
	port              int16
	server            net.Listener
}

// ServerInterface defined in the server object
type ServerInterface interface {
	Listen() (chan Message, error)
	Close() bool
}

// Listen to TCP connections
func (c *ServerClass) Listen() (chan Message, error) {
	var err error
	// Listen on all interfaces to port
	c.server, err = net.Listen("tcp", ":"+strconv.FormatInt(int64(c.port), 10))
	if appErr.HasError(err) {
		return nil, err
	}
	c.channel = make(chan Message, c.messageBufferSize)
	go c.run(c.channel)
	fmt.Printf("Listening to port %d\n", c.port)
	return c.channel, nil
}

// Close the server
func (c ServerClass) Close() bool {
	fmt.Println("Close the server")
	defer c.server.Close()
	return true
}

func (c *ServerClass) run(messages chan Message) {
	for {
		conn, err := c.server.Accept()
		if appErr.HasError(err) {
			continue
		}
		go c.handleConnection(conn, messages)
	}
}

func (c *ServerClass) handleConnection(conn net.Conn, messages chan Message) {
	id := randString(32)
	c.conns[id] = &conn
	defer func() {
		conn.Close()
		delete(c.conns, id)
	}()
	reader := bufio.NewReader(conn)

	for {
		// conn.SetReadDeadline(time.Now().Add(time.Second * 5))

		buff, err := reader.ReadBytes('\n')
		if appErr.HasError(err) {
			break
		}

		messages <- Message{id: id, raw: buff}

		time.Sleep(time.Millisecond)
	}
}

// NewServer network.ServerInterface instance pointer
func NewServer(config ServerConfig) ServerInterface {
	c := new(ServerClass)
	c.messageBufferSize = config.MessageBufferSize
	c.port = config.Port
	c.conns = make(map[string]*net.Conn)
	return c
}

// COPIED CODE

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
