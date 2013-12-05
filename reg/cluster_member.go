package reg

import (
	"encoding/hex"
	"fmt"
	xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	"strconv"
	"strings"
)

var (
	INDENT = xn.INDENT // in a better go, this would be const
)

type ClusterMember struct {
	Attrs        uint64 // negotiated with/decreed by reg server
	ClusterName  string
	ClusterAttrs uint64
	ClusterID    *xi.NodeID
	ClusterSize  uint32 // this is a FIXED size, aka MaxSize, including self
	SelfIndex    uint32 // which member we are in the Members slice

	Members []*MemberInfo // information on (other) cluster members

	// EpCount is the number of endPoints dedicated to use for cluster-
	// related purposes.  By convention endPoints[0] is used for
	// member-member communications and [1] for comms with cluster clients,
	// should they exist. The first EpCount endPoints are passed
	// to other cluster members via the registry.
	EpCount uint32

	xn.Node
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

	if cm.Attrs != other.Attrs || cm.ClusterName != other.ClusterName ||
		cm.ClusterAttrs != other.ClusterAttrs ||
		cm.ClusterSize != other.ClusterSize || cm.EpCount != other.EpCount {
		return false
	}
	if !cm.ClusterID.Equal(other.ClusterID) {
		return false
	}
	for i := 0; i < len(cm.Members); i++ {
		if !cm.Members[i].Equal(other.Members[i]) {
			return false
		}
	}
	return true
}

// SERIALIZATION ////////////////////////////////////////////////////

func (cm *ClusterMember) Strings() (ss []string) {
	ss = []string{"clusterMember {"}
	ns := cm.Node.Strings()
	for i := 0; i < len(ns); i++ {
		ss = append(ss, INDENT+ns[i])
	}
	ss = append(ss, fmt.Sprintf("%sattrs: %d", INDENT, cm.Attrs))
	ss = append(ss, fmt.Sprintf("%sclusterName: %s", INDENT, cm.ClusterName))
	ss = append(ss, fmt.Sprintf("%sclusterAttrs: %d", INDENT, cm.ClusterAttrs))
	ss = append(ss, fmt.Sprintf("%sclusterID: %s", INDENT,
		hex.EncodeToString(cm.ClusterID.Value())))
	ss = append(ss, fmt.Sprintf("%sclusterSize: %d", INDENT, cm.ClusterSize))
	ss = append(ss, fmt.Sprintf("%sselfIndex: %d", INDENT, cm.SelfIndex))
	ss = append(ss, fmt.Sprintf("%smembers {", INDENT))
	for i := 0; i < len(cm.Members); i++ {
		// DEBUG
		//fmt.Printf("serializing member %d\n", i)
		// END
		miss := cm.Members[i].Strings()
		for j := 0; j < len(miss); j++ {
			ss = append(ss, fmt.Sprintf("%s%s", INDENT, miss[j]))
		}
	}
	ss = append(ss, fmt.Sprintf("%s}", INDENT))
	ss = append(ss, fmt.Sprintf("%sepCount: %d", INDENT, cm.EpCount))
	ss = append(ss, "}")
	return
}
func (cm *ClusterMember) String() string {
	return strings.Join(cm.Strings(), "\n")
}
func ParseClusterMember(s string) (
	cm *ClusterMember, rest []string, err error) {

	ss := strings.Split(s, "\n")
	return ParseClusterMemberFromStrings(ss)
}

func ParseClusterMemberFromStrings(ss []string) (
	cm *ClusterMember, rest []string, err error) {

	var (
		node         *xn.Node
		attrs        uint64
		clusterName  string
		clusterAttrs uint64
		clusterID    *xi.NodeID
		clusterSize  uint32
		selfIndex    uint32
		memberInfos  []*MemberInfo
		epCount      uint32
	)
	line := xn.NextNBLine(&ss)
	if line != "clusterMember {" {
		err = IllFormedClusterMember
	} else {
		node, rest, err = xn.ParseFromStrings(ss)
	}
	if err == nil {
		line = xn.NextNBLine(&rest)
		parts := strings.Split(line, ": ")
		if len(parts) == 2 && parts[0] == "attrs" {
			var n int
			raw := strings.TrimSpace(parts[1])
			n, err = strconv.Atoi(raw)
			if err == nil {
				attrs = uint64(n)
			}
		} else {
			err = IllFormedClusterMember
		}
	} // FOO
	if err == nil {
		line = xn.NextNBLine(&rest)
		parts := strings.Split(line, ": ")
		if len(parts) == 2 && parts[0] == "clusterName" {
			clusterName = strings.TrimLeft(parts[1], " \t")
		} else {
			err = IllFormedClusterMember
		}
	}
	if err == nil {
		line = xn.NextNBLine(&rest)
		parts := strings.Split(line, ": ")
		if len(parts) == 2 && parts[0] == "clusterAttrs" {
			var n int
			raw := strings.TrimSpace(parts[1])
			n, err = strconv.Atoi(raw)
			if err == nil {
				clusterAttrs = uint64(n)
			}
		} else {
			err = IllFormedClusterMember
		}
	}
	if err == nil {
		line = xn.NextNBLine(&rest)
		parts := strings.Split(line, ": ")
		if len(parts) == 2 && parts[0] == "clusterID" {
			var h []byte
			raw := strings.TrimSpace(parts[1])
			h, err = hex.DecodeString(raw)
			if err == nil {
				clusterID, err = xi.New(h)
			}
		} else {
			err = IllFormedClusterMember
		}
	}
	if err == nil {
		line = xn.NextNBLine(&rest)
		parts := strings.Split(line, ": ")
		if len(parts) == 2 && parts[0] == "clusterSize" {
			var n int
			raw := strings.TrimSpace(parts[1])
			n, err = strconv.Atoi(raw)
			if err == nil {
				clusterSize = uint32(n)
			}
		} else {
			err = IllFormedClusterMember
		}
	}
	if err == nil {
		line = xn.NextNBLine(&rest)
		parts := strings.Split(line, ": ")
		if len(parts) == 2 && parts[0] == "selfIndex" {
			var n int
			raw := strings.TrimSpace(parts[1])
			n, err = strconv.Atoi(raw)
			if err == nil {
				selfIndex = uint32(n)
			}
		} else {
			err = IllFormedClusterMember
		}
	} // GEEP

	if err == nil {
		line = xn.NextNBLine(&rest)
		if line == "members {" {
			line = strings.TrimSpace(rest[0]) // a peek
			for line == "memberInfo {" {
				var mi *MemberInfo
				mi, rest, err = ParseMemberInfoFromStrings(rest)
				if err != nil {
					break
				} else {
					memberInfos = append(memberInfos, mi)
					line = strings.TrimSpace(rest[0]) // a peek
				}
			}
			// we need a closing brace at this point
			line = xn.NextNBLine(&rest)
			if line != "}" {
				err = IllFormedClusterMember
			}
		} else {
			err = IllFormedClusterMember
		}
	}
	if err == nil {
		line = xn.NextNBLine(&rest)
		parts := strings.Split(line, ": ")
		if len(parts) == 2 && parts[0] == "epCount" {
			var n int
			raw := strings.TrimSpace(parts[1])
			n, err = strconv.Atoi(raw)
			if err == nil {
				epCount = uint32(n)
			}
		} else {
			err = IllFormedClusterMember
		}
	}
	if err == nil {
		line = xn.NextNBLine(&rest)
		if line != "}" {
			err = IllFormedClusterMember
		}
	}
	if err == nil {
		cm = &ClusterMember{
			Attrs:        attrs,
			ClusterName:  clusterName,
			ClusterAttrs: clusterAttrs,
			ClusterID:    clusterID,
			ClusterSize:  clusterSize,
			SelfIndex:    selfIndex,
			Members:      memberInfos,
			EpCount:      epCount,
			Node:         *node,
		}
	}
	return
}
