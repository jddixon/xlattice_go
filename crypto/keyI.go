package crypto

// This interface is used in several places; DO NOT DELETE.

type KeyI interface {
	Algorithm() string
	GetPublicKey() PublicKeyI
	GetSigner() DigSignerI
	String() string
}
