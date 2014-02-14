package node

// xlattice_go/node/baseNode.go

import (
	// "code.google.com/p/go.crypto/ssh"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xo "github.com/jddixon/xlattice_go/overlay"
	"strings"
)

/**
 * Basic abstraction underlying Peer and Node
 */

type BaseNode struct {
	name        string // convenience for testing
	nodeID      *xi.NodeID
	commsPubKey *rsa.PublicKey
	sigPubKey   *rsa.PublicKey
	overlays    []xo.OverlayI
}

func NewNewBaseNode(name string, id *xi.NodeID) (*BaseNode, error) {
	return NewBaseNode(name, id, nil, nil, nil)
}

func NewBaseNode(name string, id *xi.NodeID,
	ck *rsa.PublicKey, sk *rsa.PublicKey,
	o []xo.OverlayI) (p *BaseNode, err error) {

	// IDENTITY /////////////////////////////////////////////////////
	if id == nil {
		return nil, NilNodeID
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
func (p *BaseNode) GetNodeID() *xi.NodeID {
	return p.nodeID
}
func (p *BaseNode) GetCommsPublicKey() *rsa.PublicKey {
	return p.commsPubKey
}
// XXX Commented out 2014-02-14 because ssh API has changed
//func (p *BaseNode) GetSSHCommsPublicKey() string {
//	out := ssh.MarshalAuthorizedKey(p.commsPubKey)
//	return string(out)
//}

func (p *BaseNode) GetSigPublicKey() *rsa.PublicKey {
	return p.sigPubKey
}

// OVERLAYS /////////////////////////////////////////////////////////
//func (p *BaseNode) addOverlayI(o xo.OverlayI) error {
//	if o == nil {
//		return NilOverlay
//	}
//	p.overlays = append(p.overlays, o)
//	return nil
//} // FOO
func (n *BaseNode) AddOverlay(o xo.OverlayI) (ndx int, err error) {
	ndx = -1
	if o == nil {
		err = NilOverlay
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
	case *BaseNode:
		_ = v
	default:
		return false
	}
	other := any.(*BaseNode) // type assertion

	// THINK ABOUT publicKey.equals(any.publicKey)

	return xi.SameNodeID(p.nodeID, other.nodeID)
}

// SERIALIZATION ////////////////////////////////////////////////////

func addStringlet(slice *[]string, s string) {
	*slice = append(*slice, s)
}
func (p *BaseNode) Strings() []string {
	// DEBUG
	if p == nil {
		panic("BaseNode.Strings: nil p !")
	}
	if p.commsPubKey == nil {
		panic("BaseNode.Strings: nil p.commsPubKey !")
	}
	// END
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
func (p *BaseNode) String() string {
	return strings.Join(p.Strings(), "\n")
}

// DESERIALIZATION //////////////////////////////////////////////////

// Return the next non-blank line in the slice of strings, trimmed.
// This line and any preceding blank lines are removed from the slice.
func NextNBLine(lines *[]string) string {
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

// Parse a serialized BaseNode, ignoring blank lines and leading and
// trailing whitespace.  Expect the first line to be like "TYPE {"

func ParseBaseNode(data, whichType string) (bn *BaseNode, rest []string, err error) {
	ss := strings.Split(data, "\n") // yields a slice of strings
	return ParseBNFromStrings(ss, whichType)
}

// Version of the above which consumes a slice of strings.  XXX Copies the
// slice unnecessarily.
func ParseBNFromStrings(ss []string, whichType string) (bn *BaseNode, rest []string, err error) {
	var (
		name        string
		nodeID      *xi.NodeID
		commsPubKey *rsa.PublicKey
		sigPubKey   *rsa.PublicKey
		overlays    []xo.OverlayI
	)
	s := NextNBLine(&ss)
	opener := fmt.Sprintf("%s {", whichType) // "peer {" or "node {"
	if s != opener {
		err = NotExpectedOpener
	}
	if err == nil {
		s := NextNBLine(&ss)
		if strings.HasPrefix(s, "name: ") {
			name = s[6:]
		} else {
			err = NotABaseNode
		}
	}
	if err == nil {
		s = NextNBLine(&ss)
		if strings.HasPrefix(s, "nodeID: ") {
			var val []byte
			val, err = hex.DecodeString(s[8:])
			if err == nil {
				nodeID, err = xi.NewNodeID(val)
			}
		} else {
			err = NotABaseNode
		}
	}
	if err == nil {
		s = NextNBLine(&ss)
		if strings.HasPrefix(s, "commsPubKey: ") {
			// XXX we do not verify that the next line is empty
			ckSSH := []byte(s[13:] + "\n")
			commsPubKey, err = xc.RSAPubKeyFromDisk(ckSSH)
		} else {
			err = NotABaseNode
		}
	}
	if err == nil {
		s = NextNBLine(&ss)
		if strings.HasPrefix(s, "sigPubKey: ") {
			skSSH := []byte(s[11:] + "\n")
			sigPubKey, err = xc.RSAPubKeyFromDisk(skSSH)
		} else {
			err = NotABaseNode
		}
	}
	if err == nil {
		s = NextNBLine(&ss)
		if s == "overlays {" {
			for {
				s = NextNBLine(&ss)
				if s == "" { // end of strings
					err = NotABaseNode
					break
				} else if s == "}" {
					prepend := []string{s}
					ss = append(prepend, ss...)
					break
				}
				var o xo.OverlayI
				o, err = xo.Parse(s)
				if err == nil {
					overlays = append(overlays, o)
				}
			}
		} else {
			err = NotABaseNode
		}
	}
	if err == nil {
		s = NextNBLine(&ss)
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
