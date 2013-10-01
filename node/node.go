package node

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xo "github.com/jddixon/xlattice_go/overlay"
	xt "github.com/jddixon/xlattice_go/transport"
	xu "github.com/jddixon/xlattice_go/util"
	"hash"
	"strings"
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
	lfs         string
	commsKey    *rsa.PrivateKey // private
	sigKey      *rsa.PrivateKey // private
	endPoints   []xt.EndPointI
	acceptors   []xt.AcceptorI // volatile, do not serialize
	peers       []Peer
	connections []xt.ConnectionI // volatile
	gateways    []Gateway
	peerMap     *BNIMap
	BaseNode    // listed last, but serialize first
}

func NewNew(name string, id *xi.NodeID, lfs string) (*Node, error) {
	// XXX create default 2K bit RSA key
	return New(name, id, lfs, nil, nil, nil, nil, nil)
}

// XXX Creating a Node with a list of live connections seems nonsensical.
func New(name string, id *xi.NodeID, lfs string,
	commsKey, sigKey *rsa.PrivateKey,
	o []xo.OverlayI, e []xt.EndPointI, p []Peer) (n *Node, err error) {

	// lfs should be a well-formed POSIX path; if the directory does
	// not exist we should create it.
	err = xu.CheckLFS(lfs)

	// The commsKey is an RSA key used to encrypt short messages.
	if err == nil {
		if commsKey == nil {
			commsKey, err = rsa.GenerateKey(rand.Reader, 2048)
		}
		if err == nil {
			// The sigKey is an RSA key used to create digital signatures.
			if sigKey == nil {
				sigKey, err = rsa.GenerateKey(rand.Reader, 2048)
			}
		}
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
	var acceptors []xt.AcceptorI // each must share index with endPoint
	var overlays []xo.OverlayI

	if err == nil {
		if o != nil {
			count := len(o)
			for i := 0; i < count; i++ {
				overlays = append(overlays, o[i])
			}
		}
		if e != nil {
			count := len(e)
			for i := 0; i < count; i++ {
				_, err = addEndPoint(e[i], &endPoints, &acceptors, &overlays)
			}
		}
	}
	var pm BNIMap // empty BNIMap
	pmPtr := &pm
	var peers []Peer // an empty slice
	if err == nil {
		if p != nil {
			count := len(p)
			for i := 0; i < count; i++ {
				err = pmPtr.AddToBNIMap(&p[i])
				if err != nil {
					break
				}
				peers = append(peers, p[i])
			}
		}
	}
	if err == nil {
		commsPubKey := &(*commsKey).PublicKey
		sigPubKey := &(*sigKey).PublicKey

		var baseNode *BaseNode
		baseNode, err = NewBaseNode(name, id, commsPubKey, sigPubKey, overlays)
		if err == nil {
			n = &Node{commsKey: commsKey,
				sigKey:    sigKey,
				acceptors: acceptors,
				endPoints: endPoints,
				peers:     peers,
				gateways:  nil,
				lfs:       lfs,
				peerMap:   pmPtr,
				BaseNode:  *baseNode}
		} 
	}
	return
}

// ENDPOINTS ////////////////////////////////////////////////////////

// Add an endPoint to a node and open an acceptor.  If a compatible
// overlay does not exist, add the default for the endPoint.
func addEndPoint(e xt.EndPointI, endPoints *[]xt.EndPointI,
	acceptors *[]xt.AcceptorI, overlays *[]xo.OverlayI) (ndx int, err error) {
	ndx = -1
	foundOverlay := false
	if len(*overlays) > 0 {
		for j := 0; j < len(*overlays); j++ {
			overlay := (*overlays)[j]
			if overlay.IsElement(e) {
				foundOverlay = true
				break
			}
		}
	}
	if !foundOverlay {
		// create a suitable overlay
		var newO xo.OverlayI
		newO, err = xo.DefaultOverlay(e)
		if err != nil {
			return
		}
		// add it to our collection
		*overlays = append(*overlays, newO)
	}
	var acc *xt.TcpAcceptor
	if e.Transport() == "tcp" {
		// XXX HACK ON ADDRESS
		strAddr := e.String()[13:]
		acc, err = xt.NewTcpAcceptor(strAddr)
		if err != nil {
			return
		}
		e = acc.GetEndPoint()

		if *endPoints != nil {
			for i := 0; i < len(*endPoints); i++ {
				if e.Equal((*endPoints)[i]) {
					ndx = i
					acc.Close()
					break
				}
			}
		}
	}
	if ndx == -1 {
		*acceptors = append(*acceptors, acc)
		*endPoints = append(*endPoints, e)
		ndx = len(*endPoints) - 1
	}
	return
}

// Returns an instance of a DigSigner which can be run in a separate
// goroutine.  This allows the Node to calculate more than one
// digital signature at the same time.
//
// XXX would prefer that *DigSigner be returned
func (n *Node) getSigner() *signer {
	return newSigner(n.sigKey)
}

func (n *Node) AddEndPoint(e xt.EndPointI) (ndx int, err error) {
	if e == nil {
		return -1, NilEndPoint
	}
	return addEndPoint(e, &n.endPoints, &n.acceptors, &n.overlays)
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

// ACCEPTORS ////////////////////////////////////////////////////////
// no accAcceptor() function; add the endPoint instead

// return a count of the number of acceptors the node listens on
func (n *Node) SizeAcceptors() int {
	return len(n.acceptors)
}

// Return the Nth acceptor, should it exist, or nil.

func (n *Node) GetAcceptor(x int) (acc xt.AcceptorI) {
	if x >= 0 && x < len(n.acceptors) {
		acc = n.acceptors[x]
	}
	return
}

// OVERLAYS /////////////////////////////////////////////////////////
//func (n *Node) AddOverlay(o xo.OverlayI) (ndx int, err error) {
//	ndx = -1
//	if o == nil {
//		err = NilOverlay
//	} else {
//		for i := 0; i < len(n.overlays); i++ {
//			if n.overlays[i].Equal(o) {
//				ndx = i
//				break
//			}
//		}
//		if ndx == -1 {
//			n.overlays = append(n.overlays, o)
//			ndx = len(n.overlays) - 1
//		}
//	}
//	return
//}
//
//func (n *Node) SizeOverlays() int {
//	return len(n.overlays)
//}
//FOO
/////** @return how to access the peer (transport, protocol, address) */
////func (n *Node) GetOverlay(x int) xo.OverlayI {
////	return n.overlays[x]
////} // GEEP

// PEERS ////////////////////////////////////////////////////////////
func (n *Node) AddPeer(peer *Peer) (ndx int, err error) {
	ndx = -1
	if peer == nil {
		err = NilPeer
	} else {
		if n.peers != nil {
			for i := 0; i < len(n.peers); i++ {
				if n.peers[i].Equal(peer) {
					ndx = i
					break
				}
			}
		}
		if ndx == -1 {
			err = n.peerMap.AddToBNIMap(peer)
			if err == nil {
				n.peers = append(n.peers, *peer)
				ndx = len(n.peers) - 1
			}
		}
	}
	return
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
func (n *Node) FindPeer(id []byte) *Peer {
	// XXX should return copy
	return n.peerMap.FindBNI(id).(*Peer)
}

// CONNECTIONS //////////////////////////////////////////////////////
func (n *Node) addConnection(c xt.ConnectionI) (ndx int, err error) {
	if c == nil {
		return -1, NilConnection
	}
	n.connections = append(n.connections, c)
	ndx = len(n.connections) - 1
	return
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
 * @return the Nth Connection
 */
func (n *Node) GetConnection(x int) xt.ConnectionI {
	return n.connections[x]
}

// LOCAL FILE SYSTEM ////////////////////////////////////////////////

// Return the path to the Node's local file system, its private
// persistent storage.  Conventionally there is a .xlattice subdirectory
// for storage of the Node's configuration information.
func (n *Node) GetLFS() string {
	return n.lfs
}

// Sets the path to the node's local storage.  If the directory does
// not exist, it creates it.  

// XXX Note possible race condition!  What is the justification for
// this function??

func (n *Node) setLFS(val string) (err error) {

	if val == "" {
		err = NilLFS
	} else {
		err = xu.CheckLFS(val) 
	}
	if err == nil {
		n.lfs = val
	}
	return
}

// CLOSE ////////////////////////////////////////////////////////////
func (n *Node) Close() {
	// XXX should run down list of connections and close each,

	// XXX STUB

	// then run down list of endpoints and close any active acceptors.
	if n.acceptors != nil {
		for i := 0; i < len(n.acceptors); i++ {
			if n.acceptors[i] != nil {
				n.acceptors[i].Close()
			}
		}
	}
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

// SERIALIZATION ////////////////////////////////////////////////////
func (n *Node) Strings() []string {
	ss := []string{"node {"}
	bns := n.BaseNode.Strings()
	for i := 0; i < len(bns); i++ {
		ss = append(ss, "    "+bns[i])
	}
	addStringlet(&ss, fmt.Sprintf("    lfs: %s", n.lfs))

	cPriv, _ := xc.RSAPrivateKeyToDisk(n.commsKey)
	addStringlet(&ss, "    commsKey: "+string(cPriv))

	sPriv, _ := xc.RSAPrivateKeyToDisk(n.sigKey)
	addStringlet(&ss, "    sigKey: "+string(sPriv))

	addStringlet(&ss, "    endPoints {")
	for i := 0; i < len(n.endPoints); i++ {
		addStringlet(&ss, "        "+n.GetEndPoint(i).String())
	}
	addStringlet(&ss, "    }")

	// peers
	addStringlet(&ss, "    peers {")
	for i := 0; i < len(n.peers); i++ {
		p := n.GetPeer(i).Strings()
		for j := 0; j < len(p); j++ {
			addStringlet(&ss, "        "+p[j])
		}
	}
	addStringlet(&ss, "    }")

	// gateways ?

	addStringlet(&ss, "}")
	return ss
}
func (n *Node) String() string {
	return strings.Join(n.Strings(), "\n")
}

// Collect an RSA private key in string form.  Only call this if
// '-----BEGIN -----' has already been seen

func collectKey(rest *[]string) (key *rsa.PrivateKey, err error) {
	ss := []string{"-----BEGIN -----"}
	for {
		// NOT ROBUST; should detect end of rest, blank line, any other errors
		line := (*rest)[0]
		*rest = (*rest)[1:]
		ss = append(ss, line)
		if line == "-----END -----" {
			break
		}
	}
	if err == nil {
		text := strings.Join(ss, "\n")
		key, err = xc.RSAPrivateKeyFromDisk([]byte(text))
	}
	return
}
func Parse(s string) (node *Node, rest []string, err error) {
	bn, rest, err := ParseBaseNode(s, "node")
	if err == nil {
		node = &Node{BaseNode: *bn}
		var pm BNIMap
		node.peerMap = &pm

		line := NextNBLine(&rest)
		parts := strings.Split(line, ": ")
		if parts[0] == "lfs" {
			node.lfs = strings.TrimSpace(parts[1])
		} else {
			fmt.Println("MISSING LFS")
			err = NotASerializedNode
		}

		var commsKey, sigKey *rsa.PrivateKey
		if err == nil {
			// move some of this into collectKey() !
			line = NextNBLine(&rest)
			parts = strings.Split(line, ": ")
			if parts[0] == "commsKey" && parts[1] == "-----BEGIN -----" {
				commsKey, err = collectKey(&rest)
				node.commsKey = commsKey
			} else {
				fmt.Println("MISSING OR ILL-FORMED COMMS_KEY")
				err = NotASerializedNode
			}
		} // FOO

		if err == nil {
			// move some of this into collectKey() !
			line = NextNBLine(&rest)
			parts = strings.Split(line, ": ")
			if parts[0] == "sigKey" && parts[1] == "-----BEGIN -----" {
				sigKey, err = collectKey(&rest)
				node.sigKey = sigKey
			} else {
				fmt.Println("MISSING OR ILL-FORMED SIG_KEY")
				err = NotASerializedNode
			}
		} // FOO

		// endPoints
		if err == nil {
			line = NextNBLine(&rest)
			if line == "endPoints {" {
				for {
					line = NextNBLine(&rest)
					if line == "}" {
						// prepend := []string{line}
						// rest = append(prepend, rest...)
						break
					}
					var ep xt.EndPointI
					ep, err = xt.ParseEndPoint(line)
					if err != nil {
						break
					}
					_, err = node.AddEndPoint(ep)
					if err != nil {
						break
					}
				}
			} else {
				fmt.Println("MISSING END_POINTS BLOCK")
				fmt.Printf("    EXPECTED 'endPoints {', GOT: '%s'\n", line)
				err = NotASerializedNode
			}
		}

		// peers
		if err == nil {
			line = NextNBLine(&rest)
			if line == "peers {" {
				for {
					line = strings.TrimSpace(rest[0])
					if line == "}" { // ZZZ
						break
					}
					var peer *Peer
					peer, rest, err = parsePeerFromStrings(rest)
					if err != nil {
						break
					}
					_, err = node.AddPeer(peer)
					if err != nil {
						break
					}
				}
			} else {
				fmt.Println("MISSING PEERS BLOCK")
				fmt.Printf("    EXPECTED 'peers {', GOT: '%s'\n", line)
				err = NotASerializedNode
			}
			line = NextNBLine(&rest) // discard the ZZZ }

		}
		// gateways, but not yet

		// expect closing brace for node {
		// XXX we need an expect(&rest)

		line = NextNBLine(&rest)
		if line != "}" {
			fmt.Printf("extra text at end of node declaration: '%s'\n", line)
		}
	}
	if err != nil {
		node = nil
	}
	return
}

// DIG SIGNER ///////////////////////////////////////////////////////

func (n *Node) Sign(chunks [][]byte) (sig []byte, err error) {
	if chunks == nil {
		err = NothingToSign
	} else {
		s := newSigner(n.sigKey)
		for i := 0; i < len(chunks); i++ {
			s.digest.Write(chunks[i])
		}
		h := s.digest.Sum(nil)
		sig, err = rsa.SignPKCS1v15(rand.Reader, s.key, crypto.SHA1, h)
	}
	return
}

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

///////////////////////////////////////////////////
// XXX This stuff needs to be cleaned up or dropped
///////////////////////////////////////////////////

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
