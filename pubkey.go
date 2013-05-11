package xlattice_go

type PublicKey interface {
    Equal(any interface{}) bool
    ToString() string
}
