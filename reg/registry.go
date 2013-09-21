package reg

// xlattice_go/reg/registry.go

// This file contains functions and structures used to describe
// and manage the cluster data managed by the registry.

import (
	"fmt"
	xn "github.com/jddixon/xlattice_go/node"
)

var _ = fmt.Print

type Registry struct {
	// registry data
	Clusters       []*RegCluster
	ClustersByName map[string]*RegCluster
	ClustersByID   *xn.BNIMap
	MembersByID    *xn.BNIMap

	// the extended XLattice node, so files, communications, and keys
	Node *RegNode
}
