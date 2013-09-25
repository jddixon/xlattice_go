package reg

// xlattice_go/reg/regData.go

// This file contains functions and structures used to describe
// and manage the cluster data managed by the registry.

import (
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	"strings"
)

var _ = fmt.Print

// cluster member bit flags
const (
	MEMBER_DELETED = 1 << iota
	FOO
	BAR
)

type ClusterMember struct {
	attrs       uint64
	myEnds      []string // serialized EndPointI
	xn.BaseNode          // name and ID must be unique
}

func NewClusterMember(name string, id *xi.NodeID,
	commsPubKey, sigPubKey *rsa.PublicKey, attrs uint64, myEnds []string) (
	member *ClusterMember, err error) {

	// all attrs bits are zero by default

	base, err := xn.NewBaseNode(name, id, commsPubKey, sigPubKey, nil)
	if err == nil {
		member = &ClusterMember{
			attrs:    attrs,
			myEnds:   myEnds,
			BaseNode: *base,
		}
	}
	return
}

// Return the XLRegMsg_Token corresponding to this cluster member.
func (cm *ClusterMember) Token() (token *XLRegMsg_Token, err error) {

	var ckBytes, skBytes []byte

	ckBytes, err = xc.RSAPubKeyToWire(cm.GetCommsPublicKey())
	if err == nil {
		skBytes, err = xc.RSAPubKeyToWire(cm.GetSigPublicKey())
		if err == nil {
			token = &XLRegMsg_Token{
				Attrs:    &cm.attrs,
				ID:       cm.GetNodeID().Value(),
				CommsKey: ckBytes,
				SigKey:   skBytes,
				MyEnds:   cm.myEnds,
			}
		}
	}
	return
}

// EQUAL ////////////////////////////////////////////////////////////

func (cm *ClusterMember) Equal(any interface{}) bool {

	if any == cm {
		return true
	}
	if any == nil {
		return false
	}
	switch v := any.(type) {
	case *ClusterMember:
		_ = v
	default:
		return false
	}
	other := any.(*ClusterMember) // type assertion
	if cm.attrs != other.attrs {
		return false
	}
	if cm.myEnds == nil {
		if other.myEnds != nil {
			return false
		}
	} else {
		if other.myEnds == nil {
			return false
		}
		if len(cm.myEnds) != len(other.myEnds) {
			return false
		}
		for i := 0; i < len(cm.myEnds); i++ {
			if cm.myEnds[i] != other.myEnds[i] {
				return false
			}
		}
	}
	// WARNING: panics without the ampersand !
	return cm.BaseNode.Equal(&other.BaseNode)
}

// SERIALIZATION ////////////////////////////////////////////////////

func (cm *ClusterMember) Strings() (ss []string) {
	ss = []string{"clusterMember {"}
	bns := cm.BaseNode.Strings()
	for i := 0; i < len(bns); i++ {
		ss = append(ss, "    "+bns[i])
	}
	ss = append(ss, fmt.Sprintf("    attrs: 0x%016x", cm.attrs))
	ss = append(ss, "    endPoints {")
	for i := 0; i < len(cm.myEnds); i++ {
		ss = append(ss, "        "+cm.myEnds[i])
	}
	ss = append(ss, "    }")
	ss = append(ss, "}")
	return
}

func (cm *ClusterMember) String() string {
	return strings.Join(cm.Strings(), "\n")
}
func collectAttrs(cm *ClusterMember, ss []string) (rest []string, err error) {
	rest = ss
	line := xn.NextNBLine(&rest) // trims
	// attrs line looks like "attrs: 0xHHHH..." where H is a hex digit
	if strings.HasPrefix(line, "attrs: 0x") {
		var val []byte
		var attrs uint64
		line := line[9:]
		val, err = hex.DecodeString(line)
		if err == nil {
			if len(val) != 8 {
				err = WrongNumberOfBytesInAttrs
			} else {
				for i := 0; i < 8; i++ {
					// assume little-endian ; but printf has put
					// high order bytes first - ie, it's big-endian
					attrs |= uint64(val[i]) << uint(8*(7-i))
				}
				cm.attrs = attrs
			}
		}
	} else {
		err = BadAttrsLine
	}
	return
}
func collectMyEnds(cm *ClusterMember, ss []string) (rest []string, err error) {
	rest = ss
	line := xn.NextNBLine(&rest)
	if line == "endPoints {" {
		for {
			line = strings.TrimSpace(rest[0]) // peek
			if line == "}" {
				break
			}
			line = xn.NextNBLine(&rest)
			// XXX NO CHECK THAT THIS IS A VALID ENDPOINT
			cm.myEnds = append(cm.myEnds, line)
		}
		line = xn.NextNBLine(&rest)
		if line != "}" {
			err = MissingClosingBrace
		}
	} else {
		err = MissingEndPointsSection
	}
	return
}
func ParseClusterMember(s string) (
	cm *ClusterMember, rest []string, err error) {

	ss := strings.Split(s, "\n")
	return ParseClusterMemberFromStrings(ss)
}

func ParseClusterMemberFromStrings(ss []string) (
	cm *ClusterMember, rest []string, err error) {

	bn, rest, err := xn.ParseBNFromStrings(ss, "clusterMember")
	if err == nil {
		cm = &ClusterMember{BaseNode: *bn}
		rest, err = collectAttrs(cm, rest)
		if err == nil {
			rest, err = collectMyEnds(cm, rest)
		}
		if err == nil {
			// expect and consume a closing brace
			line := xn.NextNBLine(&rest)
			if line != "}" {
				err = MissingClosingBrace
			}
		}
	}
	return
}
