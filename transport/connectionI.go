package transport

import xc "github.com/jddixon/xlattice_go/crypto"

/**
 * A Connection is a relationship between two EndPoints.  In XLattice,
 * one *EndPoint will have an Address on this Node.  The other EndPoint
 * will have an Address associated with a second Node.  There is always
 * a transport protocol associated with the connection; both EndPoints
 * use this same protocol.
 *
 * Connections are established to allow one or more Messages to be passed
 * between the Nodes at the two EndPoints.
 *
 * XXX Connections could be homogeneous or heterogeneous.  In the first
 * XXX case, each EndPoint would use the same transport.  In the second case,
 * XXX the heterogeneous Connection, the two EndPoints would use different
 * XXX transports.
 *
 * A connection passes through a set of states.  This progress is
 * irreversible.  Each state has an associated numeric value.  The
 * order of these values is guaranteed.  That is, UNBOUND is less
 * than BOUND, which is less than PENDING, and so forth.  Therefore
 * it is reasonable to test on the relative value of states.
 *
 * XXX If the Transport is Udp, then it is likely that we will want
 * XXX to be able to bind and unbind the connection, allowing us to use
 * XXX it with more than one remote *EndPoint.  In this case, the
 * XXX progression through numbered states would not be irreversible.
 *
 * If new states are defined, they should adhere to this contract.
 * That is, the passage of a connection through a sequence of states
 * must be irreversible, and the numeric value of any later state
 * must be greater than that associated with any earlier state.
 *
 * XXX Any application can encrypt data passing over a connection.
 * XXX Is it reasonable to mandate what follows as part of the
 * XXX interface?
 *
 * There may be a secret key associated with the Connection.  This
 * will be used with a symmetric cipher to encrypt and decrypt
 * traffic over the Connection.  Such secret keys are negotiated
 * between the *EndPoint Nodes and possibly periodically renegotiated.
 *
 * Connections exist at various levels of abstraction.  TCP, for
 * example, is layered on top of IP, and BGP4 on top of TCP.  It is
 * possible for a connection to be in clear, but used for carrying
 * encrypted messages at a higher protocol level.  It is equally
 * possible that data passing over a connection will be encrypted
 * at more than one level.
 *
 * @author Jim Dixon
 */

// STATE ////////////////////////////////////////////////////////
/** neither end point is set */
const UNBOUND = 100

/** near end point is set */
const BOUND = 200

/** connection to far end point has been requested */
const PENDING = 300

/** both end points have been set, connection is established */
const CONNECTED = 400

/** connection has been closed */
const DISCONNECTED = 500

type ConnectionI interface {

	/**
	 * Return the current state index.  In the current implementation,
	 * this is not necessarily reliable, but the actual state index
	 * is guaranteed to be no less than the value reported.
	 *
	 * @return one of the values above
	 */
	GetState() int

	/**
	 * Set the near end point of a connection.  If either the
	 * near or far end point has already been set, this will
	 * cause an exception.  If successful, the connection's
	 * state becomes BOUND.
	 */
	BindNearEnd(e *EndPoint) (err error) // throws IOException

	/**
	 * Set the far end point of a connection.  If the near end
	 * point has NOT been set or if the far end point has already
	 * been set -- in other words, if the connection is already
	 * beyond state BOUND -- this will cause an exception.
	 * If the operation is successful, the connection's state
	 * becomes either PENDING or CONNECTED.
	 *
	 * XXX The state should become CONNECTED if the far end is on
	 * XXX the same host and PENDING if it is on a remoted host.
	 */
	BindFarEnd(e *EndPoint) (err error) // throws IOException

	/**
	 * Bring the connection to the DISCONNECTED state.
	 */
	Close() (err error) // throws IOException

	/**
	 * This should be considered deprecated.  Test on whether the
	 * state is DISCONNECTED instead.
	 *
	 * @return whether the connection state is DISCONNECTED.
	 */
	IsClosed() bool

	// END POINTS ///////////////////////////////////////////////////
	GetNearEnd() *EndPoint

	GetFarEnd() *EndPoint

	// I/O //////////////////////////////////////////////////////////
	IsBlocking() bool

	// ///////////////////////////////////////////////////////////////////
	// XXX CONFUSION BETWEEN PACKET vs STREAM AND BLOCKING vs NON-BLOCKING
	// ///////////////////////////////////////////////////////////////////
	// non-blocking

	// blocking
	//  GetInputStream(i *InputStream, e error)     // throws IOException
	//  GetOutputStream(o *OutputStream, e error)   // throws IOException

	// ENCRYPTION ///////////////////////////////////////////////////
	/** @return whether the connection is encrypted */
	IsEncrypted() bool

	/**
	 * (Re)negotiate the Secret used to encrypt traffic over the
	 * connection.
	 *
	 * @param myKey  this Node's asymmetric key
	 * @param hisKey Peer's public key
	 */
	Negotiate(myKey xc.KeyI, hisKey xc.PublicKeyI) (s xc.SecretI, e error)
	// throws CryptoException

	Equal(any interface{}) bool
	String() string
}
