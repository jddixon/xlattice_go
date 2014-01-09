package overlay

// xlattice_go/overlay/nameKeyedWriterI.go

type NameKeyedWriterI interface {
	Delete(key string, listener DelCallBackI)
	Put (key string, buffer []byte, listener PutCallBackI)
}
