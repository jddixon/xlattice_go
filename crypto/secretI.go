package crypto

type SecretI interface {
	Algorithm() string
	Equal(any interface{}) bool
	String() string
}
