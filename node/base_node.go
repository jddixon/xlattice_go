package node

// xlattice_go/node/base_node.go

import (
	"code.google.com/p/go.crypto/ssh"
	"crypto/rsa"
	"encoding/hex"
	"errors"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	xo "github.com/jddixon/xlattice_go/overlay"
	"strings"
)

var _ = fmt.Print

/**
 * Basic abstraction underlying Peer and Node
 */

type BaseNode struct {
	name        string // convenience for testing
	nodeID      *NodeID
	commsPubKey *rsa.PublicKey
	sigPubKey   *rsa.PublicKey
	overlays    []xo.OverlayI
}

func NewNewBaseNode(name string, id *NodeID) (*BaseNode, error) {
	return NewBaseNode(name, id, nil, nil, nil)
}

func NewBaseNode(name string, id *NodeID,
	ck *rsa.PublicKey, sk *rsa.PublicKey,
	o []xo.OverlayI) (p *BaseNode, err error) {

	// IDENTITY /////////////////////////////////////////////////////
	if id == nil {
		err := errors.New("IllegalArgument: nil NodeID")
		return nil, err
	}
	nodeID, err := (*id).Clone()
	if err != nil {
		return
	}
	commsPubKey := ck
	sigPubKey := sk
	var overlays []xo.OverlayI // an empty slice
	if o != nil {
		count := len(o)
		for i := 0; i < count; i++ {
			overlays = append(overlays, o[i])
		}
	} // FOO
	p = new(BaseNode)
	p.name = name
	p.nodeID = nodeID // the clone
	p.commsPubKey = commsPubKey
	p.sigPubKey = sigPubKey
	p.overlays = overlays
	return p, nil
}
func (p *BaseNode) GetName() string {
	return p.name
}
func (p *BaseNode) GetNodeID() *NodeID {
	return p.nodeID
}
func (p *BaseNode) GetCommsPublicKey() *rsa.PublicKey {
	return p.commsPubKey
}
func (p *BaseNode) GetSSHCommsPublicKey() string {
	out := ssh.MarshalAuthorizedKey(p.commsPubKey)
	return string(out)
}

func (p *BaseNode) GetSigPublicKey() *rsa.PublicKey {
	return p.sigPubKey
}

// OVERLAYS /////////////////////////////////////////////////////////
//func (p *BaseNode) addOverlayI(o xo.OverlayI) error {
//	if o == nil {
//		return errors.New("IllegalArgument: nil OverlayI")
//	}
//	p.overlays = append(p.overlays, o)
//	return nil
//} // FOO
func (n *BaseNode) AddOverlay(o xo.OverlayI) (ndx int, err error) {
	ndx = -1
	if o == nil {
		err = errors.New("IllegalArgument: nil Overlay")
	} else {
		for i := 0; i < len(n.overlays); i++ {
			if n.overlays[i].Equal(o) {
				ndx = i
				break
			}
		}
		if ndx == -1 {
			n.overlays = append(n.overlays, o)
			ndx = len(n.overlays) - 1
		}
	}
	return
} // FOO

// Return a count of the number of overlays.
func (p *BaseNode) SizeOverlays() int {
	return len(p.overlays)
}

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

// SERIALIZATION ////////////////////////////////////////////////////

func addStringlet(slice *[]string, s string) {
	*slice = append(*slice, s)
}
func (p *BaseNode) Strings() []string {
	ckSSH, err := xc.RSAPubKeyToDisk(p.commsPubKey)
	if err != nil {
		panic(err)
	}
	skSSH, err := xc.RSAPubKeyToDisk(p.sigPubKey)
	if err != nil {
		panic(err)
	}

	var s []string
	addStringlet(&s, fmt.Sprintf("name: %s", p.name))
	addStringlet(&s, fmt.Sprintf("nodeID: %s", p.nodeID.String()))
	addStringlet(&s, fmt.Sprintf("commsPubKey: %s", ckSSH))
	addStringlet(&s, fmt.Sprintf("sigPubKey: %s", skSSH))
	addStringlet(&s, fmt.Sprintf("overlays {"))
	for i := 0; i < len(p.overlays); i++ {
		addStringlet(&s, fmt.Sprintf("    %s", p.overlays[i].String()))
	}
	addStringlet(&s, fmt.Sprintf("}"))
	return s
}

// DESERIALIZATION //////////////////////////////////////////////////

var (
	NotABaseNode = errors.New("not a serialized BaseNode - missing bits")
)

// Return the next non-blank line in the slice of strings
func nextLine(lines *[]string) string {
	if lines != nil {
		for len(*lines) > 0 {
			s := strings.TrimSpace((*lines)[0])
			*lines = (*lines)[1:]
			if s != "" {
				return s
			}
		}
	}
	return ""
}

// Parse a serialized BaseNode, ignoring blank lines and leading
// and trailing whitespace
func ParseBaseNode(data string) (bn *BaseNode, rest []string, err error) {
	ss := strings.Split(data, "\n") // yields a slice of strings
	var (
		name        string
		nodeID      *NodeID
		commsPubKey *rsa.PublicKey
		sigPubKey   *rsa.PublicKey
		overlays    []xo.OverlayI
	)
	s := nextLine(&ss)
	if strings.HasPrefix(s, "name: ") {
		name = s[6:]
	} else {
		err = NotABaseNode
	}
	if err == nil {
		s = nextLine(&ss)
		if strings.HasPrefix(s, "nodeID: ") {
			var val []byte
			val, err = hex.DecodeString(s[7:])
			if err == nil {
				nodeID, err = NewNodeID(val)
			}
		} else {
			err = NotABaseNode
		}
	}
	if err == nil {
		s = nextLine(&ss)
		if strings.HasPrefix(s, "commsPubKey: ") {
			// XXX we do not verify that the next line is empty
			ckSSH := []byte(s[13:] + "\n")
			commsPubKey, err = xc.RSAPubKeyFromDisk(ckSSH)
		} else {
			err = NotABaseNode
		}
	}
	if err == nil {
		s = nextLine(&ss)
		if strings.HasPrefix(s, "SigPubKey: ") {
			skSSH := []byte(s[12:] + "\n")
			commsPubKey, err = xc.RSAPubKeyFromDisk(skSSH)
		} else {
			err = NotABaseNode
		}
	}
	if err == nil {
		s = nextLine(&ss)
		if s == "overlays {" {
			for {
				s = nextLine(&ss)
				if s == "" { // end of strings
					err = NotABaseNode
					break
				} else if s == "}" {
					prepend := []string{s}
					ss = append(prepend, ss...)
					break
				} else if strings.HasPrefix("IPOverlay: ", s) {
					var ar *xo.AddrRange
					ar, err = xo.NewCIDRAddrRange(s[11:])
					if err != nil {
						break
					}
					var o xo.OverlayI
					// XXX WE DON'T HAVE ITS NAME OR COST!
					o, err = xo.NewIPOverlay("", ar, "ip", 1.0)
					if err != nil {
						break
					}
					overlays = append(overlays, o)
				} else {
					// we can only handle IP overlays
					fmt.Printf("not a recognized IP overlay: %s\n", s)
					err = NotABaseNode
					break
				}
			}
		} else {
			err = NotABaseNode
		}
	}
	if err == nil {
		s = nextLine(&ss)
		if s != "}" {
			err = NotABaseNode
		}
	}
	if err == nil {
		var bn = BaseNode{name, nodeID, commsPubKey, sigPubKey, overlays}
		return &bn, ss, nil
	} else {
		return nil, nil, err
	}
}
