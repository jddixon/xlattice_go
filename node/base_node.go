package node

// xlattice_go/node/base_node.go

import (
	"code.google.com/p/go.crypto/ssh"
	"crypto/rsa"
	"errors"
	"fmt" // DEBUG
	xo "github.com/jddixon/xlattice_go/overlay"
)

var _ = fmt.Print

/**
 * Basic abstraction underlying Peer and Node
 */

type BaseNode struct {
	name        string // convenience for testing
	nodeID      *NodeID
	commsPubkey *rsa.PublicKey
	sigPubkey   *rsa.PublicKey
	overlays    []xo.OverlayI
}

func NewNewBaseNode(id *NodeID) (*BaseNode, error) {
	return NewBaseNode(id, nil, nil, nil)
}

func NewBaseNode(id *NodeID,
	ck *rsa.PublicKey, sk *rsa.PublicKey,
	o []xo.OverlayI) (*BaseNode, error) {

	// IDENTITY /////////////////////////////////////////////////////
	if id == nil {
		err := errors.New("IllegalArgument: nil NodeID")
		return nil, err
	}
	nodeID := (*id).Clone()
	commsPubkey := ck
	sigPubkey := sk
	var overlays []xo.OverlayI // an empty slice
	if o != nil {
		count := len(o)
		for i := 0; i < count; i++ {
			overlays = append(overlays, o[i])
		}
	} // FOO
	p := new(BaseNode)
	p.nodeID = nodeID // the clone
	p.commsPubkey = commsPubkey
	p.sigPubkey = sigPubkey
	p.overlays = overlays
	return p, nil
}
func (p *BaseNode) GetNodeID() *NodeID {
	return p.nodeID
}
func (p *BaseNode) GetCommsPublicKey() *rsa.PublicKey {
	return p.commsPubkey
}
func (p *BaseNode) GetSSHCommsPublicKey() string {
	out := ssh.MarshalAuthorizedKey(p.commsPubkey)
	// PLAYING AROUND
	outAgain, comment, options, rest, ok := ssh.ParseAuthorizedKey(out)
	pubKey := outAgain.(*rsa.PublicKey)
	fmt.Printf("outAgain: %v\n", pubKey)
	fmt.Printf("comment: '%s'\n", comment)
	fmt.Printf("len(options) = %d\n", len(options))
	fmt.Printf("len(rest) = %d\n", len(rest))
	fmt.Printf("OK = %v\n", ok)
	// END PLAYING
	return string(out)
}

func (p *BaseNode) GetSigPublicKey() *rsa.PublicKey {
	return p.sigPubkey
}

// OVERLAYS /////////////////////////////////////////////////////////
func (p *BaseNode) addOverlayI(o xo.OverlayI) error {
	if o == nil {
		return errors.New("IllegalArgument: nil OverlayI")
	}
	p.overlays = append(p.overlays, o)
	return nil
}

/**
 * @return a count of the number of overlays the peer can be
 *         accessed through
 */
func (p *BaseNode) sizeOverlays() int {
	return len(p.overlays)
}

/** @return how to access the peer (transport, protocol, address) */
func (p *BaseNode) GetOverlay(n int) xo.OverlayI {
	return p.overlays[n]
}

// EQUAL ////////////////////////////////////////////////////////////
func (p *BaseNode) Equal(any interface{}) bool {
	if any == p {
		return true
	}
	if any == nil {
		return false
	}
	switch v := any.(type) {
	case BaseNode:
		_ = v
	default:
		return false
	}
	other := any.(Peer) // type assertion
	// THINK ABOUT publicKey.equals(any.publicKey)
	if p.nodeID == other.nodeID {
		return true
	}
	if p.nodeID.Length() != other.nodeID.Length() {
		return false
	}
	myVal := p.nodeID.Value()
	otherVal := other.nodeID.Value()
	for i := 0; i < p.nodeID.Length(); i++ {
		if myVal[i] != otherVal[i] {
			return false
		}
	}
	return false
}
func (p *BaseNode) String() string {
	return "NOT IMPLEMENTED"
}
