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
	 */
	Connect(near EndPointI) (c ConnectionI, e error)

	/**
	 * @return the Acceptor EndPoint that this Connector is used to
	 *          establish connections to
	 */
	GetFarEnd() EndPointI
	String() string
}
