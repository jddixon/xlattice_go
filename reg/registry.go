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
	"strings"
	"sync"
)

var _ = fmt.Print

type Registry struct {
	LogFile        string
	Logger         *log.Logger            // volatile, not serialized
	
	// registry data
	Clusters       []*RegCluster          // serialized
	ClustersByName map[string]*RegCluster // volatile, not serialized
	ClustersByID   *xn.BNIMap             // -ditto-
	RegMembersByID *xn.BNIMap             // -ditto-
	mu             sync.RWMutex			  // -ditto-

	// the extended XLattice node, so id, lfs, keys, etc
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

// Tentatively the registry is serialized separately from the regNode
// and so consists of a sequence of serialized clusters

func (reg *Registry) String() (s string) {
	return strings.Join(reg.Strings(), "\n")
}

// If we change the serialization so that there is no closing brace,
// it will be possible to simply append cluster serializations to the
// registry configuration file while the registry is running.

func (reg *Registry) Strings() (ss []string) {
	ss = []string {"registry {"}
	ss = append(ss, fmt.Sprintf("    LogFile: %s", reg.LogFile))
	ss = append(ss, "}")

	for i := 0 ; i < len(reg.Clusters); i++ {
		cs := reg.Clusters[i].Strings()
		for j := 0; j < len(cs); j++ {
			ss = append(ss, cs[j])
		}
	}
	return
}

func ParseRegistry(s string) (reg *Registry, rest []string, err error) {

	// XXX STUB
	return
}
