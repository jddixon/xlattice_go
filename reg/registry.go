package reg

// xlattice_go/reg/registry.go

// This file contains functions and structures used to describe
// and manage the cluster data managed by the registry.

import (
	"crypto/rsa"
	"fmt"
	xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xo "github.com/jddixon/xlattice_go/overlay"
	xt "github.com/jddixon/xlattice_go/transport"
)

var _ = fmt.Print

type Registry struct {
	// registry data
	Clusters       []*RegCluster
	ClustersByName map[string]*RegCluster
	ClustersByID   *xn.BNIMap
	MembersByID    *xn.BNIMap

	// the extended XLattice node, so files, communications, and keys
	RegNode
}

func NewRegistry(clusters []*RegCluster, name string, id *xi.NodeID,
	lfs string, ckPriv, skPriv *rsa.PrivateKey,
	overlay xo.OverlayI, endPoint xt.EndPointI) (reg *Registry, err error) {

	rn, err := NewRegNode(name, id, lfs, ckPriv, skPriv, overlay, endPoint)
	if err == nil {
		reg = &Registry{
			Clusters:       clusters,
			ClustersByName: make(map[string]*RegCluster),
			RegNode:        *rn,
		}
		if clusters != nil {
			// XXX need to populate the indexes here
		}
	}
	return
}
