package protocol

// Abstracts a family of messages.
type ProtocolI interface {
	Name() string
	String() string
}
