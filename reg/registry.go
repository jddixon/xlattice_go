package reg

// xlattice_go/reg/registry.go

// This file contains functions and structures used to describe
// and manage the cluster data managed by the registry.

import (
	"crypto/rsa"
	"fmt"
	xn "github.com/jddixon/xlattice_go/node"
	"log"
	"os"
	"sync"
)

var _ = fmt.Print

type Registry struct {
	// registry data
	Clusters       []*RegCluster          // serialized
	ClustersByName map[string]*RegCluster // volatile, not serialized
	ClustersByID   *xn.BNIMap             // -ditto-
	RegMembersByID *xn.BNIMap             // -ditto-
	Logger         *log.Logger            // -ditto-
	mu             sync.RWMutex           // -ditto-

	// the extended XLattice node, so files, communications, and keys
	RegNode
}

func NewRegistry(clusters []*RegCluster, node *xn.Node, 
	ckPriv, skPriv *rsa.PrivateKey, logger *log.Logger) (
	reg *Registry, err error) {

	rn, err := NewRegNode(node, ckPriv, skPriv)
	if err == nil {
		var bniMap xn.BNIMap
		if logger == nil {
			logger = log.New(os.Stderr, "", log.Ldate|log.Ltime)
		}
		reg = &Registry{
			Clusters:       clusters,
			ClustersByName: make(map[string]*RegCluster),
			ClustersByID:   &bniMap,
			Logger:         logger,
			RegNode:        *rn,
		}
		if clusters != nil {
			// XXX need to populate the indexes here
		}
	}
	return
}

// XXX RegMembersByID is not being updated!  This is the redundant and so
// possibly inconsistent index of members of registry clusters

func (reg *Registry) AddCluster(cluster *RegCluster) (index int, err error) {

	if cluster == nil {
		err = NilCluster
	} else {
		name := cluster.Name
		id := cluster.ID // []byte

		reg.mu.Lock()
		defer reg.mu.Unlock()

		if _, ok := reg.ClustersByName[name]; ok {
			err = NameAlreadyInUse
		} else if reg.ClustersByID.FindBNI(id) != nil {
			err = IDAlreadyInUse
		}
		if err == nil {
			index = len(reg.Clusters)
			reg.Clusters = append(reg.Clusters, cluster)
			reg.ClustersByName[name] = cluster
			err = reg.ClustersByID.AddToBNIMap(cluster)
		}
	}
	if err != nil {
		index = -1
	}
	return
}

// SERIALIZATION ====================================================

func ParseRegistry(s string) (reg *Registry, rest []string, err error) {


	// XXX STUB
	return
}

func (reg *Registry) String() (s string) {

	// STUB XXX
	return
}

func (reg *Registry) Strings() (s []string) {

	// STUB XXX
	return
}
