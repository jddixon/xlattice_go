package reg

// xlattice_go/reg/reg_cluster.go

// This file contains functions and structures used to describe
// and manage the clusters managed by the registry.

import (
	"crypto/rsa"
	"fmt"
	xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	"strings"
)

var _ = fmt.Print

// cluster bit flags 
const (
	CLUSTER_DELETED = 1 << iota
)

type RegCluster struct {
	attrs			uint64		// a field of bit flags
	Name			string // must be unique
	ID				[]byte // must be unique
	Size			int    // a maximum > 1
	Members			[]*ClusterMember
	MembersByName	map[string]*ClusterMember
	MembersByID		*xn.BaseNodeMap
}

func NewRegCluster(name string, id *xi.NodeID, size int) (
	rc *RegCluster, err error) {

	// all attrs bits are zero by default

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

// SERIALIZATION ////////////////////////////////////////////////////
func (rc *RegCluster) Strings() (s []string) {
	// XXX STUB 

	return
}

func (rc *RegCluster) String() string {
	return strings.Join(rc.Strings(), "\n")
}
func ParseRegCluster(rc *RegCluster, rest []string, err error) {

	// XXX STUB

	return
}
