package node

import (
	"crypto/rsa"
	"errors"
	"fmt"
	xo "github.com/jddixon/xlattice_go/overlay"
	xt "github.com/jddixon/xlattice_go/transport"
	"strings"
)

/**
 * A Peer is another Node, a neighbor.
 */

type Peer struct {
	connectors []xt.ConnectorI // to reach the peer
	BaseNode
}

func NewNewPeer(name string, id *NodeID) (*Peer, error) {
	return NewPeer(name, id, nil, nil, nil, nil)
}

func NewPeer(name string, id *NodeID,
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
		p := Peer{ctors, *baseNode}
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
	bns := p.BaseNode.Strings()
	bns = append(bns, "connectors {")
	for i := 0; i < len(p.connectors); i++ {
		bns = append(bns, fmt.Sprintf("    %s", p.connectors[i].String()))
	}
	bns = append(bns, "}")
	return bns
}
func (p *Peer) String() string {
	return strings.Join(p.Strings(), "\n")
}
