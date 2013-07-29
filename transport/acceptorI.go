package transport

/**
 * An Acceptor is used by a Node or Peer to accept connection requests.
 * It is an advertisement for a service within a Overlay, that is,
 * within a given address space and using a particular transport
 * protocol.
 *
 * An Acceptor is an abstraction of a TCP/IP ServerSocket.  It is a
 * single EndPoint whose Address may be well known.  Other entities on
 * the network send messages to the Acceptor in order to establish
 * Connections.  The Acceptor may in some cases NOT be one of the
 * EndPoints involved in the new Connection; the Connection might
 * be between the requesting remote EndPoint and a new, ephemeral
 * local EndPoint.
 *
 * The transport protocol understood by the Acceptor need not be
 * the same as the transport protocol of Connections created.  That is,
 * the new Connection need not be in the same Overlay as the Acceptor.
 *
 * @author Jim Dixon
 */
type AcceptorI interface {
	Accept() (ConnectionI, error)
	Close() error
	IsClosed() bool
	GetEndPoint() EndPointI
	String() string
}
