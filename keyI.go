package xlattice_go

type KeyI interface {
	Algorithm() string
	GetPublicKey() PublicKeyI
	GetSigner() DigSignerI
	String() string
}
