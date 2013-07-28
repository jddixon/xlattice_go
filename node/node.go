package node

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"errors"
	"fmt"
	xo "github.com/jddixon/xlattice_go/overlay"
	xt "github.com/jddixon/xlattice_go/transport"
	"hash"
)

var _ = fmt.Print

/**
 * A Node is uniquely identified by a NodeID and can satisfy an
 * identity test constructed using its public key.  That is, it
 * can prove that it holds the private key materials corresponding
 * to the public key.
 *
 * @author Jim Dixon
 */
type Node struct {
	commsKey    *rsa.PrivateKey // private
	sigKey      *rsa.PrivateKey // private
	overlays    []xo.OverlayI
	endPoints   []xt.EndPointI
	peers       []Peer
	connections []xt.ConnectionI
	gateways    []Gateway
	BaseNode
}

func NewNew(id *NodeID) (*Node, error) {
	// XXX create default 2K bit RSA key
	return New(id, nil, nil, nil, nil, nil, nil)
}

// XXX Creating a Node with a list of live connections seems nonsensical.
func New(id *NodeID, commsKey, sigKey *rsa.PrivateKey, o []xo.OverlayI,
	e []xt.EndPointI, p []Peer, c []xt.ConnectionI) (*Node, error) {

	// The commsKey is an RSA key used to encrypt short messages.
	if commsKey == nil {
		k, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, err
		}
		commsKey = k
	}
	// The sigKey is an RSA key used to create digital signatures.
	if sigKey == nil {
		k, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, err
		}
		sigKey = k
	}
	// The node communicates through its endpoints.  These are
	// contained in overlays.  If an endpoint in 127.0.0.0/8
	// is in the list of endpoints, that overlay is automatically
	// added to the list of overlays with the name "localhost".
	// Other IPv4 endpoints are assumed to be in 0.0.0.0/0
	// ("globalV4") unless there is another containing overlay
	// except that endpoints in private address space are treated
	// differently.  Unless there is an overlay with a containing
	// address space, addresses in 10/8 are assigned to "privateA",
	// addresses in 172.16/12 are assigned to "privateB", and
	// any in 192.168/16 are assigned to "privateC".  All of these
	// overlays are automatically created unless there is a
	// pre-existing overlay whose address range is the same as one
	// of these are contained within one of them.

	var endPoints []xt.EndPointI 
	var overlays []xo.OverlayI

	if o != nil {
		count := len(o)
		for i := 0; i < count; i++ {
			overlays = append(overlays, o[i])
		}
	}
	if e != nil {
		count := len(e)
		for i := 0; i < count; i++ {
			endPoints = append(endPoints, e[i])
			foundIt := false
			if len(overlays) > 0 {
				for j := 0; j < len(overlays); j++ {
					if overlays[j].IsElement(e[i]) {
						foundIt = true
						break
					}
				}
			}
			if ! foundIt {
				// create a suitable overlay
				// WORKING HERE 
				newO,err := xo.DefaultOverlay(e[i])
				if err != nil {
					return nil, err
				}
				// add it to our collection
				overlays = append(overlays, newO)

			}
		}
	}
	var peers []Peer // an empty slice
	if p != nil {
		count := len(p)
		for i := 0; i < count; i++ {
			peers = append(peers, p[i])
		}
	}
	var cnxs []xt.ConnectionI // another empty slice
	if c != nil {
		count := len(c)
		for i := 0; i < count; i++ {
			cnxs = append(cnxs, c[i])
		}
	}

	commsPubKey := &(*commsKey).PublicKey
	sigPubKey := &(*sigKey).PublicKey

	baseNode, err := NewBaseNode(id, commsPubKey, sigPubKey, overlays)
	if err == nil {
		p := Node{commsKey, sigKey, overlays, endPoints, peers, cnxs, nil, *baseNode}
		return &p, nil
	} else {
		return nil, err
	}
}

// Returns an instance of a DigSigner which can be run in a separate
// goroutine.  This allows the Node to calculate more than one
// digital signature at the same time.
//
// XXX would prefer that *DigSigner be returned
func (n *Node) getSigner() *signer {
	return newSigner(n.sigKey)
}

// ENDPOINTS ////////////////////////////////////////////////////////
func (n *Node) addEndPoint(e xt.EndPointI) error {
	if e == nil {
		return errors.New("IllegalArgument: nil EndPoint")
	}
	// XXX ATTEMPT TO LISTEN ON THE ENDPOINT XXX
	n.endPoints = append(n.endPoints, e)
	return nil
}

/**
 * @return a count of the number of endPoints the peer can be
 *         accessed through
 */
func (n *Node) SizeEndPoints() int {
	return len(n.endPoints)
}

func (n *Node) GetEndPoint(x int) xt.EndPointI {
	return n.endPoints[x]
}

// OVERLAYS /////////////////////////////////////////////////////////
func (n *Node) addOverlay(o xo.OverlayI) error {
	if o == nil {
		return errors.New("IllegalArgument: nil Overlay")
	}
	n.overlays = append(n.overlays, o)
	return nil
}

/**
 * @return a count of the number of overlays the peer can be
 *         accessed through
 */
func (n *Node) SizeOverlays() int {
	return len(n.overlays)
}

/** @return how to access the peer (transport, protocol, address) */
func (n *Node) GetOverlay(x int) xo.OverlayI {
	return n.overlays[x]
} // GEEP

// PEERS ////////////////////////////////////////////////////////////
func (n *Node) addPeer(o *Peer) error {
	if o == nil {
		return errors.New("IllegalArgument: nil Peer")
	}
	n.peers = append(n.peers, *o)
	return nil
}

/**
 * @return a count,  the number of peers
 */
func (n *Node) SizePeers() int {
	return len(n.peers)
}
func (n *Node) GetPeer(x int) *Peer {
	// XXX should return copy
	return &n.peers[x]
}

// CONNECTORS ///////////////////////////////////////////////////////
func (n *Node) addConnectionI(c xt.ConnectionI) error {
	if c == nil {
		return errors.New("IllegalArgument: nil ConnectionI")
	}
	n.connections = append(n.connections, c)
	return nil
}

/** @return a count of known Connections for this Peer */
func (n *Node) SizeConnections() int {
	return len(n.connections)
}

/**
 * Return a ConnectionI, an Address-Protocol pair identifying
 * an Acceptor for the Peer.  Connections are arranged in order
 * of preference, with the zero-th ConnectionI being the most
 * preferred.  THESE ARE OPEN, LIVE CONNECTIONS.
 *
 * XXX Could as easily return an EndPoint.
 *
 * @return the Nth Connection
 */
func (n *Node) GetConnection(x int) xt.ConnectionI {
	return n.connections[x]
}

// CLOSE ////////////////////////////////////////////////////////////
func (n *Node) Close() {
	// XXX should run down list of connections and close each,
	// then run down list of endpoints and close any active acceptors.
}

// EQUAL ////////////////////////////////////////////////////////////
func (n *Node) Equal(any interface{}) bool {
	if any == n {
		return true
	}
	if any == nil {
		return false
	}
	switch v := any.(type) {
	case *Node:
		_ = v
	default:
		return false
	}
	other := any.(*Node) // type assertion
	// THINK ABOUT publicKey.equals(any.publicKey)
	if n.nodeID == other.nodeID {
		return true
	}
	if n.nodeID.Length() != other.nodeID.Length() {
		return false
	}
	myVal := n.nodeID.Value()
	otherVal := other.nodeID.Value()
	for i := 0; i < n.nodeID.Length(); i++ {
		if myVal[i] != otherVal[i] {
			return false
		}
	}
	return false
}
func (n *Node) String() string {
	return "NOT IMPLEMENTED"
}

// DIG SIGNER ///////////////////////////////////////////////////////

type signer struct {
	key    *rsa.PrivateKey
	digest hash.Hash
}

func newSigner(key *rsa.PrivateKey) *signer {
	// XXX some validation, please
	h := sha1.New()
	ds := signer{key: key, digest: h}
	return &ds
}
func (s *signer) Algorithm() string {
	return "SHA1+RSA" // XXX NOT THE PROPER NAME
}
func (s *signer) Length() int {
	return 42 // XXX NOT THE PROPER VALUE
}
func (s *signer) Update(chunk []byte) {
	s.digest.Write(chunk)
}

// XXX 2013-07-15 Golang crypto package currently does NOT support SHA3 (Keccak)
func (s *signer) Sign() ([]byte, error) {
	h := s.digest.Sum(nil)
	sig, err := rsa.SignPKCS1v15(rand.Reader, s.key, crypto.SHA1, h)
	return sig, err
}

func (s *signer) String() string {
	return "NOT IMPLEMENTED"
}
