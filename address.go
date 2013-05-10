package xlattice_go

// An Address provides enough information to identify an endpoint.
// The information needed depends upon the communications protocol
// used.
type Address interface {
    Equal(any interface{}) bool
    ToString() string
}
