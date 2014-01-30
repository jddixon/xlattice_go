package overlay

// xlattice_go/overlay/nameKeyedReaderI.go

// This call is synchronous: it blocks
type NameKeyedReaderI interface {
	Get(key string) error
}
