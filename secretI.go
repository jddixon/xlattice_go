package xlattice_go

type SecretI interface {
	Algorithm() string
	Equal(any interface{}) bool
	String() string
}
