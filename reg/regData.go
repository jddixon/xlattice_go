package reg

// xlattice_go/reg/regData.go

// This file contains functions and structures used to describe
// and manage the cluster data managed by the registry.

import (
	"crypto/rsa"
	"fmt"
	//xm "github.com/jddixon/xlattice_go/msg"
	xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	//xo "github.com/jddixon/xlattice_go/overlay"
	//xt "github.com/jddixon/xlattice_go/transport"
)

var _ = fmt.Print

// bit flags
const (
	EPHEMERAL = 1 << iota
	FOO
	BAR
)

type RegCluster struct {
	Name          string // must be unique
	ID            []byte // must be unique
	Size          int    // a maximum > 1
	Members       []*ClusterMember
	MembersByName map[string]*ClusterMember
	MembersByID   *xn.BaseNodeMap
}

func NewRegCluster(name string, id *xi.NodeID, size int) (
	rc *RegCluster, err error) {

	if name == "" {
		name = "xlCluster"
	}
	if size < 2 {
		err = ClusterMustHaveTwo
	} else {
		var bnm xn.BaseNodeMap // empty map
		rc = &RegCluster{
			Name:        name,
			ID:          id.Value(),
			Size:        size,
			MembersByID: &bnm,
		}
	}
	return
}

func (rc *RegCluster) AddToCluster(name string, id *xi.NodeID,
	commsPubKey, sigPubKey *rsa.PublicKey, attrs uint64) (err error) {

	if _, ok := rc.MembersByName[name]; ok {
		// XXX surely something more complicated is called for!
		return
	}
	member, err := NewClusterMember(name, id, commsPubKey, sigPubKey, attrs)
	if err == nil {
		rc.MembersByName[name] = member

		// XXX add to MembersByID

	}
	return
}

type ClusterMember struct {
	Attrs       uint64
	xn.BaseNode // name and ID must be unique
}

func NewClusterMember(name string, id *xi.NodeID,
	commsPubKey, sigPubKey *rsa.PublicKey, attrs uint64) (
	member *ClusterMember, err error) {

	base, err := xn.NewBaseNode(name, id, commsPubKey, sigPubKey, nil)
	if err == nil {
		member = &ClusterMember{Attrs: attrs, BaseNode: *base}
	}
	return
}

type RegData struct {
	clusters []*RegCluster
}
