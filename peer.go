package xlattice_go

import "errors"

/**
 * A Peer is another Node, a neighbor.
 */

type Peer struct {
	nodeID     *NodeID
	pubkey     *PublicKeyI
	overlays   []*OverlayI
	connectors []*ConnectorI
}

func NewNewPeer(id *NodeID) (*Peer, error) {
	return NewPeer(id, nil, nil, nil)
}

func NewPeer(id *NodeID, k *PublicKeyI, o *[]*OverlayI,
	c *[]*ConnectorI) (*Peer, error) {

	// IDENTITY /////////////////////////////////////////////////////
	if id == nil {
		err := errors.New("IllegalArgument: nil NodeID")
		return nil, err
	}
	nodeID := (*id).Clone()
	pubkey := k
	var overlays []*OverlayI // an empty slice
	if o != nil {
		count := len(*o)
		for i := 0; i < count; i++ {
			overlays = append(overlays, (*o)[i])
		}
	} // FOO
	var ctors []*ConnectorI // another empty slice
	if c != nil {
		count := len(*c)
		for i := 0; i < count; i++ {
			ctors = append(ctors, (*c)[i])
		}
	} // FOO
	p := new(Peer)
	p.nodeID = nodeID // the clone
	p.pubkey = pubkey
	p.overlays = overlays
	p.connectors = ctors
	return p, nil
}
func (p *Peer) GetNodeID() *NodeID {
	return p.nodeID
}
func (p *Peer) GetPublicKeyI() *PublicKeyI {
	return p.pubkey
}

// XXX Is this a good idea ??
func (p *Peer) SetPublicKeyI(k *PublicKeyI) error {
	err := error(nil)
	if k == nil {
		err = errors.New("IllegalArgument: nil PublicKeyI")
	} else {
		p.pubkey = k
	}
	return err
}

// OVERLAYS /////////////////////////////////////////////////////////
func (p *Peer) addOverlayI(o *OverlayI) error {
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
func (p *Peer) sizeOverlays() int {
	return len(p.overlays)
}

/** @return how to access the peer (transport, protocol, address) */
func (p *Peer) GetOverlay(n int) *OverlayI {
	return p.overlays[n]
}

// CONNECTORS ///////////////////////////////////////////////////////
func (p *Peer) addConnector(c *ConnectorI) error {
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
func (p *Peer) GetConnector(n int) *ConnectorI {
	return p.connectors[n]
}

// EQUAL ////////////////////////////////////////////////////////////
func (p *Peer) Equal(any interface{}) bool {
	if any == p {
		return true
	}
	if any == nil {
		return false
	}
	switch v := any.(type) {
	case Peer:
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
