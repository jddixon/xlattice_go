package crypto

type DigSignerI interface {
	// Returns the name of the algorithm used to sign.
	Algorithm(any interface{}) string
	// Return the length in bytes of the digital signature generated.
	Length() int
	Update([]byte)
	// Generates a digital signature and implicitly resets.
	Sign() []byte
	String() string
}
