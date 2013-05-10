package protocol

// Abstracts a family of messages.
type Protocol interface {
    Name() string
}
