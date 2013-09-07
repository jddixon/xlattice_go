package reg

import (
	xn "github.com/jddixon/xlattice_go/node"
)

type Registry struct {
}

type RegCluster struct {
	name    string // must be unique
	id      []byte // must be unique
	size    int    // a maximum
	members []*ClusterMember
}

type ClusterMember struct {
	xn.Peer
}

type RegUser struct {
	name string // must be unique
	id   []byte // must be unique
	attr uint64 // should be cluster-specific
}
