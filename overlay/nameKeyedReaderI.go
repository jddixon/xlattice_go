package overlay

// xlattice_go/overlay/nameKeyedReaderI.go

type NameKeyedReaderI interface {
    Get (key string, listener GetCallBackI)
}
