package xlattice_go

type Secret interface {
	Algorithm() string
	Equal(any interface{}) bool
	ToString() string
}
