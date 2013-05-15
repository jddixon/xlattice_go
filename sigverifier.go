package xlattice_go

// SigVerifier is a translation of a Java abstract class.

type SigVerifier interface {
	GetAlgorithm() string
	Init(PublicKey)
	Update([]byte)
	Verify([]byte) bool
}
