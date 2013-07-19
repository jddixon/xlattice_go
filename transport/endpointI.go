package transport

/**
 * An EndPoint is specified by a transport and an Address, including
 * the local part.  If the transport is TCP/IP, for example, the
 * Address includes the IP address and the port number.
 *
 */

type EndPointI interface {
	Clone() (*EndPointI, error)
	GetAddress() AddressI
	GetTransport() string
	String() string
}
