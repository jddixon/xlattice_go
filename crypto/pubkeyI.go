package crypto

type PublicKeyI interface {
	Equal(any interface{}) bool
	String() string
}
