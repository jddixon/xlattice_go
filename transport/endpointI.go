package transport

/**
 * An EndPoint is specified by a transport and an Address, including
 * the local part.  If the transport is TCP/IP, for example, the
 * Address includes the IP address and the port number.
 *
 */

type EndPointI interface {
	Address() AddressI
	Clone() (EndPointI, error)
	Equal(any interface{}) bool
	String() string
	Transport() string
}
