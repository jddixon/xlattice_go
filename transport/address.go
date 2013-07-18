package transport

// An Address provides enough information to identify an endpoint.
// The information needed depends upon the communications protocol
// used.
type Address struct {
	repr string
}

func (a *Address) String() string {
	return a.repr
}
