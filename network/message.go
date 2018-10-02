package network

// MessageType declares if a message is incoming or outgoing
type MessageType uint8

const (
	// IncomingMessage type
	IncomingMessage MessageType = 1
	// OutgoingMessage type
	OutgoingMessage MessageType = 2
)

func (mt MessageType) String() string {
	switch mt {
	case IncomingMessage:
		return "Incoming"
	case OutgoingMessage:
		return "Outgoing"
	}
	return ""
}

// Message to be used by servers and clients to map things around
type Message struct {
	ID           string
	Raw          []byte
	Type         MessageType
	ConnectionID string
}
