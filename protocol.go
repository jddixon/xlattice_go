package xlattice_go

// Abstracts a family of messages.
type Protocol interface {
    Name() string
}
