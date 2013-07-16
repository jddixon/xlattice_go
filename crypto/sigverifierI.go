package crypto

// SigVerifier is a translation of a Java abstract class.

type SigVerifierI interface {
	GetAlgorithm() string
	Init(PublicKeyI)
	Update([]byte)
	Verify([]byte) bool
	String() string
}
