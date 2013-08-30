package node

import (
	"crypto/rsa"
	"fmt"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xo "github.com/jddixon/xlattice_go/overlay"
	xt "github.com/jddixon/xlattice_go/transport"
	"strings"
	"sync"
	"time"
)

/**
 * A Peer is another Node, a neighbor.
 */

type Peer struct {
	connectors []xt.ConnectorI // to reach the peer
	timeout    int64           // ns from epoch
	contacted  int64           // last contact from this peer, ns from epoch
	ndx        int             // order in which added
	up         bool            // set to false if considered unreachable
	mu         sync.Mutex
	BaseNode
}

func NewNewPeer(name string, id *xi.NodeID) (*Peer, error) {
	return NewPeer(name, id, nil, nil, nil, nil)
}

func NewPeer(name string, id *xi.NodeID,
	ck *rsa.PublicKey, sk *rsa.PublicKey,
	o []xo.OverlayI, c []xt.ConnectorI) (*Peer, error) {

	baseNode, err := NewBaseNode(name, id, ck, sk, o)

	if err == nil {
		var ctors []xt.ConnectorI // another empty slice
		if c != nil {
			count := len(c)
			for i := 0; i < count; i++ {
				ctors = append(ctors, c[i])
			}
		}
		p := Peer{connectors: ctors, BaseNode: *baseNode}
		return &p, nil // FOO
	} else {
		return nil, err
	}
}

// CONNECTORS ///////////////////////////////////////////////////////
func (p *Peer) addConnector(c xt.ConnectorI) error {
	if c == nil {
		return NilConnector
	}
	p.connectors = append(p.connectors, c)
	return nil
}

/** @return a count of known Connectors for this Peer */
func (p *Peer) SizeConnectors() int {
	return len(p.connectors)
}

/**
 * Return a Connector, an Address-Protocol pair identifying
 * an Acceptor for the Peer.  Connectors are arranged in order
 * of preference, with the zero-th Connector being the most
 * preferred.
 *
 * XXX Could as easily return an EndPoint.
 *
 * @return the Nth Connector
 */
func (p *Peer) GetConnector(n int) xt.ConnectorI {
	return p.connectors[n]
}

// EQUAL ////////////////////////////////////////////////////////////
// func (p *Peer) Equal(any interface{}) bool {
//     XXX Uses BaseNode.Equal()
// }

func (p *Peer) Strings() []string {
	ss := []string{"peer {"}
	bns := p.BaseNode.Strings()
	for i := 0; i < len(bns); i++ {
		ss = append(ss, "    "+bns[i])
	}
	ss = append(ss, "    connectors {")
	for i := 0; i < len(p.connectors); i++ {
		ss = append(ss, fmt.Sprintf("        %s", p.connectors[i].String()))
	}
	ss = append(ss, "    }")
	ss = append(ss, "}")
	return ss
}

func (p *Peer) String() string {
	return strings.Join(p.Strings(), "\n")
}

func collectConnectors(peer *Peer, ss []string) (rest []string, err error) {
	rest = ss
	line := nextLine(&rest)
	if line == "connectors {" {
		for {
			line = nextLine(&rest)
			if line == "}" {
				break
			}
			var ctor xt.ConnectorI
			ctor, err = xt.ParseConnector(line)
			if err != nil {
				return
			}
			err = peer.addConnector(ctor)
			if err != nil {
				return
			}
		}
		line = nextLine(&rest)
		if line != "}" {
			err = NotASerializedPeer
		}
	} else {
		// no connectors, not a very useful peer
		err = NotASerializedPeer
		peer = nil
	}
	return
}
func ParsePeer(s string) (peer *Peer, rest []string, err error) {
	bn, rest, err := ParseBaseNode(s, "peer")
	if err == nil {
		peer = &Peer{BaseNode: *bn}
		rest, err = collectConnectors(peer, rest)
	}
	return
}
func parsePeerFromStrings(ss []string) (peer *Peer, rest []string, err error) {
	bn, rest, err := parseBNFromStrings(ss, "peer")
	if err == nil {
		peer = &Peer{BaseNode: *bn}
		rest, err = collectConnectors(peer, rest)
	}
	return
}

// LIVENESS /////////////////////////////////////////////////////////

// Return the time (ns from the Epoch) of the last communication with
// this peer.
func (p *Peer) LastContact() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.contacted
}

// A communication with the peer has occurred: mark the time.
func (p *Peer) StillAlive() {
	t := time.Now().UnixNano()
	p.mu.Lock()
	p.contacted = t
	p.mu.Unlock()
}

// Return whether the peer is considered reachable.
func (p *Peer) IsUp() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.up
}

// Clear the peer's up flag.  It is no longer considered reachable.
// Return the flag's previous state.
func (p *Peer) MarkDown() (prevState bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	prevState = p.up
	p.up = false
	return
}

// Set the peer's up flag.  It is now considered reachable.  Return
// the flag's previous state.
func (p *Peer) MarkUp() (prevState bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	prevState = p.up
	p.up = true
	return
}
