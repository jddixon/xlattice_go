package reg

// xlattice_go/reg/member_info.go

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

type MemberInfo struct {
	Attrs       uint64   //  bit flags are defined in const.go
	MyEnds      []string // serialized EndPointI
	xn.BaseNode          // name and ID must be unique
}

func NewMemberInfo(name string, id *xi.NodeID,
	commsPubKey, sigPubKey *rsa.PublicKey, attrs uint64, myEnds []string) (
	member *MemberInfo, err error) {

	// all attrs bits are zero by default

	// DEBUG
	// fmt.Printf("NewMemberInfo for server %s\n", name)
	// END
	base, err := xn.NewBaseNode(name, id, commsPubKey, sigPubKey, nil)
	if err == nil {
		member = &MemberInfo{
			Attrs:    attrs,
			MyEnds:   myEnds,
			BaseNode: *base,
		}
	}
	return
}

// Create the MemberInfo corresponding to the token passed.

func NewMemberInfoFromToken(token *XLRegMsg_Token) (
	m *MemberInfo, err error) {

	var nodeID *xi.NodeID
	if token == nil {
		err = NilToken
	} else {
		nodeID, err = xi.New(token.GetID())
		if err == nil {
			ck, err := xc.RSAPubKeyFromWire(token.GetCommsKey())
			if err == nil {
				sk, err := xc.RSAPubKeyFromWire(token.GetSigKey())
				if err == nil {
					m, err = NewMemberInfo(token.GetName(), nodeID,
						ck, sk, token.GetAttrs(), token.GetMyEnds())
				}
			}
		}
	}
	return
}

// Return the XLRegMsg_Token corresponding to this cluster member.
func (mi *MemberInfo) Token() (token *XLRegMsg_Token, err error) {

	var ckBytes, skBytes []byte

	ck := mi.GetCommsPublicKey()
	// DEBUG
	if ck == nil {
		fmt.Printf("MemberInfo.Token: %s commsPubKey is nil\n", mi.GetName())
	}
	// END
	ckBytes, err = xc.RSAPubKeyToWire(ck)
	if err == nil {
		skBytes, err = xc.RSAPubKeyToWire(mi.GetSigPublicKey())
		if err == nil {
			name := mi.GetName()
			token = &XLRegMsg_Token{
				Name:     &name,
				Attrs:    &mi.Attrs,
				ID:       mi.GetNodeID().Value(),
				CommsKey: ckBytes,
				SigKey:   skBytes,
				MyEnds:   mi.MyEnds,
			}
		}
	}
	return
}

// EQUAL ////////////////////////////////////////////////////////////

func (mi *MemberInfo) Equal(any interface{}) bool {

	if any == mi {
		return true
	}
	if any == nil {
		return false
	}
	switch v := any.(type) {
	case *MemberInfo:
		_ = v
	default:
		return false
	}
	other := any.(*MemberInfo) // type assertion
	if mi.Attrs != other.Attrs {
		return false
	}
	if mi.MyEnds == nil {
		if other.MyEnds != nil {
			return false
		}
	} else {
		if other.MyEnds == nil {
			return false
		}
		if len(mi.MyEnds) != len(other.MyEnds) {
			return false
		}
		for i := 0; i < len(mi.MyEnds); i++ {
			if mi.MyEnds[i] != other.MyEnds[i] {
				return false
			}
		}
	}
	// WARNING: panics without the ampersand !
	return mi.BaseNode.Equal(&other.BaseNode)
}

// SERIALIZATION ////////////////////////////////////////////////////

func (mi *MemberInfo) Strings() (ss []string) {
	ss = []string{"memberInfo {"}
	bns := mi.BaseNode.Strings()
	for i := 0; i < len(bns); i++ {
		ss = append(ss, "    "+bns[i])
	}
	ss = append(ss, fmt.Sprintf("    attrs: 0x%016x", mi.Attrs))
	ss = append(ss, "    endPoints {")
	for i := 0; i < len(mi.MyEnds); i++ {
		ss = append(ss, "        "+mi.MyEnds[i])
	}
	ss = append(ss, "    }")
	ss = append(ss, "}")
	return
}

func (mi *MemberInfo) String() string {
	return strings.Join(mi.Strings(), "\n")
}
func collectAttrs(mi *MemberInfo, ss []string) (rest []string, err error) {
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
				mi.Attrs = attrs
			}
		}
	} else {
		err = BadAttrsLine
	}
	return
}
func collectMyEnds(mi *MemberInfo, ss []string) (rest []string, err error) {
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
			mi.MyEnds = append(mi.MyEnds, line)
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
func ParseMemberInfo(s string) (
	mi *MemberInfo, rest []string, err error) {

	ss := strings.Split(s, "\n")
	return ParseMemberInfoFromStrings(ss)
}

func ParseMemberInfoFromStrings(ss []string) (
	mi *MemberInfo, rest []string, err error) {

	bn, rest, err := xn.ParseBNFromStrings(ss, "memberInfo")
	if err == nil {
		mi = &MemberInfo{BaseNode: *bn}
		rest, err = collectAttrs(mi, rest)
		if err == nil {
			rest, err = collectMyEnds(mi, rest)
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
