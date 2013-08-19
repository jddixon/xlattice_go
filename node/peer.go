package node

import (
	"crypto/rsa"
	"errors"
	"fmt"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xo "github.com/jddixon/xlattice_go/overlay"
	xt "github.com/jddixon/xlattice_go/transport"
	"strings"
	"sync"
)

var (
	NotASerializedPeer = errors.New("not a serialized peer")
)

/**
 * A Peer is another Node, a neighbor.
 */

type Peer struct {
	connectors []xt.ConnectorI // to reach the peer
	timeout    int64           // ns from epoch
	prev       int64           // last contact from this peer, ns from epoch
	ndx        int             // order in which added
	down       bool            // set to true if considered unreachable
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
		return errors.New("IllegalArgument: nil Connector")
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
//func (p *Peer) Equal(any interface{}) bool {
//	if any == p {
//		return true
//	}
//	if any == nil {
//		return false
//	}
//	switch v := any.(type) {
//	case Peer:
//		_ = v
//	default:
//		return false
//	}
//	other := any.(Peer) // type assertion
//	// THINK ABOUT publicKey.equals(any.publicKey)
//	if p.nodeID == other.nodeID {
//		return true
//	}
//	if p.nodeID.Length() != other.nodeID.Length() {
//		return false
//	}
//	myVal := p.nodeID.Value()
//	otherVal := other.nodeID.Value()
//	for i := 0; i < p.nodeID.Length(); i++ {
//		if myVal[i] != otherVal[i] {
//			return false
//		}
//	}
//	return false
//} // GEEP

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
