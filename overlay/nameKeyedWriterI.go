package overlay

// xlattice_go/overlay/nameKeyedWriterI.go

// This call is synchronous: it blocks.
type NameKeyedWriterI interface {
	Delete(key string) error
	Put(key string, buffer []byte) error
}
