package transport

/**
 * Used to establish a Connection with another entity (Node).
 *
 * The notion is that a node has a collection of Connectors used
 * for establishing Connections with Peers, neighboring nodes.
 *
 * @author Jim Dixon
 */
type ConnectorI interface {

	/**
	 * Establish a Connection with another entity using the transport
	 * and address in the EndPoint.
	 *
	 * @param nearEnd  local end point to use for connection
	 * @param blocking whether the new Connection is to be blocking
	 */
	Connect(near *EndPointI, blocking bool) (
		c ConnectionI, e error) // throws IOException

	/**
	 * @return the Acceptor EndPoint that this Connector is used to
	 *          establish connections to
	 */
	GetFarEnd() EndPointI
	String() string
}
