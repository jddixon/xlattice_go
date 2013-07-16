package crypto

type KeyI interface {
	Algorithm() string
	GetPublicKey() PublicKeyI
	GetSigner() DigSignerI
	String() string
}
