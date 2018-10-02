package network

// Message to be used by servers and clients to map things around
type Message struct {
	raw []byte
	id  string
}
