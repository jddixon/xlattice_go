package xlattice_go

type DigSigner interface {
    Algorithm() string
    Length() int
    Update([]byte)
    Sign() []byte
}
