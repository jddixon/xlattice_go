package xlattice_go

type PublicKeyI interface {
	Equal(any interface{}) bool
	String() string
}
