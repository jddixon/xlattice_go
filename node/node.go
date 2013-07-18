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
	endPoints   []*xt.EndPoint
	peers       []*Peer
	connections []*xt.ConnectionI
	gateways    []*Gateway
	BaseNode
}

func NewNew(id *NodeID) (*Node, error) {
	// XXX create default 2K bit RSA key
	return New(id, nil, nil, nil, nil, nil)
}

func New(id *NodeID, commsKey, sigKey *rsa.PrivateKey,
	e *[]*xt.EndPoint, p *[]*Peer, c *[]*xt.ConnectionI) (*Node, error) {

	///////////////////////////////
	// XXX STUB: switch on key type
	// * extract the public key
	// * build the DigSigner
	///////////////////////////////
	if commsKey == nil {
		k, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, err
		}
		commsKey = k
	}
	if sigKey == nil {
		k, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, err
		}
		sigKey = k
	}

	var endPoints []*xt.EndPoint // an empty slice
	var overlays []*xo.OverlayI
	if e != nil {
		count := len(*e)
		for i := 0; i < count; i++ {
			endPoints = append(endPoints, (*e)[i])
			// XXX get the overlay from the endPoint
			// overlays = append(overlays, (*o)[i])
		}
	} // FOO
	var peers []*Peer // an empty slice
	if p != nil {
		count := len(*p)
		for i := 0; i < count; i++ {
			peers = append(peers, (*p)[i])
		}
	}
	var cnxs []*xt.ConnectionI // another empty slice
	if c != nil {
		count := len(*c)
		for i := 0; i < count; i++ {
			cnxs = append(cnxs, (*c)[i])
		}
	}

	commsPubKey := &(*commsKey).PublicKey
	sigPubKey := &(*sigKey).PublicKey

	baseNode, err := NewBaseNode(id, commsPubKey, sigPubKey, &overlays)
	if err == nil {
		p := Node{commsKey, sigKey, endPoints, peers, cnxs, nil, *baseNode}
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

// OVERLAYS /////////////////////////////////////////////////////////
func (n *Node) addOverlay(o *xo.OverlayI) error {
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
func (n *Node) GetOverlay(x int) *xo.OverlayI {
	return n.overlays[x]
}

// PEERS ////////////////////////////////////////////////////////////
func (n *Node) addPeer(o *Peer) error {
	if o == nil {
		return errors.New("IllegalArgument: nil Peer")
	}
	n.peers = append(n.peers, o)
	return nil
}

/**
 * @return a count,  the number of peers
 */
func (n *Node) SizePeers() int {
	return len(n.peers)
}
func (n *Node) GetPeer(x int) *Peer {
	return n.peers[x]
}

// CONNECTORS ///////////////////////////////////////////////////////
func (n *Node) addConnectionI(c *xt.ConnectionI) error {
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
 * preferred.
 *
 * XXX Could as easily return an EndPoint.
 *
 * @return the Nth Connection
 */
func (n *Node) GetConnection(x int) *xt.ConnectionI {
	return n.connections[x]
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
