package reg

// xlattice_go/reg/registry.go

// This file contains functions and structures used to describe
// and manage the cluster data managed by the registry.

import (
	"crypto/rsa"
	"fmt"
	xf "github.com/jddixon/xlattice_go/crypto/filters"
	xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xt "github.com/jddixon/xlattice_go/transport"
	xu "github.com/jddixon/xlattice_go/util"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var _ = fmt.Print

type Registry struct {
	LogFile string
	Logger  *log.Logger // volatile, not serialized

	// registry data
	m, k     uint          // serialized
	Clusters []*RegCluster // serialized

	idFilter xf.BloomSHAI

	ClustersByName map[string]*RegCluster // volatile, not serialized
	ClustersByID   *xn.BNIMap             // -ditto-
	RegMembersByID *xn.BNIMap             // -ditto-
	mu             sync.RWMutex           // -ditto-

	// the extended XLattice node, so id, lfs, keys, etc
	RegNode
}

func NewRegistry(clusters []*RegCluster, node *xn.Node,
	ckPriv, skPriv *rsa.PrivateKey, opt *RegOptions) (
	reg *Registry, err error) {

	var (
		idFilter      xf.BloomSHAI
		rn            *RegNode
		serverVersion xu.DecimalVersion
	)
	serverVersion, err = xu.ParseDecimalVersion(VERSION)
	// DEBUG
	fmt.Printf("NewRegistry: server version is %s\n", serverVersion.String())
	// END
	if err == nil {
		rn, err = NewRegNode(node, ckPriv, skPriv)
	}
	if err == nil {
		if opt.BackingFile == "" {
			idFilter, err = xf.NewBloomSHA(opt.M, opt.K)
		} else {
			idFilter, err = xf.NewMappedBloomSHA(opt.M, opt.K, opt.BackingFile)
		}
	}
	if err == nil {
		// registry's own ID added to Bloom filter
		idFilter.Insert(node.GetNodeID().Value())

		var bniMap xn.BNIMap
		logger := opt.Logger
		if logger == nil {
			logger = log.New(os.Stderr, "", log.Ldate|log.Ltime)
		}
		reg = &Registry{
			idFilter:       idFilter,
			Clusters:       clusters,
			ClustersByName: make(map[string]*RegCluster),
			ClustersByID:   &bniMap,
			Logger:         logger,
			RegNode:        *rn,
		}
		if clusters != nil {
			// XXX need to populate the indexes here
		}
		myLFS := rn.GetLFS()
		if myLFS != "" {
			var ep []xt.EndPointI
			for i := 0; i < rn.SizeEndPoints(); i++ {
				ep = append(ep, rn.GetEndPoint(i))
			}
			regCred := &RegCred{
				Name:        rn.GetName(),
				ID:          rn.GetNodeID(),
				CommsPubKey: rn.GetCommsPublicKey(),
				SigPubKey:   rn.GetSigPublicKey(),
				EndPoints:   ep,
				Version:     serverVersion,
			}
			serialized := regCred.String() // shd have terminating CRLF
			pathToFile := filepath.Join(myLFS, "regCred.dat")
			err = ioutil.WriteFile(pathToFile, []byte(serialized), 0640)
		}
	}
	return
}

func (reg *Registry) ContainsID(n *xi.NodeID) (bool, error) {
	return reg.idFilter.Member(n.Value())
}
func (reg *Registry) InsertID(n *xi.NodeID) (err error) {
	b := n.Value()
	found, err := reg.idFilter.Member(b)
	if err == nil && found {
		err = IDAlreadyInUse
	}
	if err == nil {
		reg.idFilter.Insert(b)
	}
	return
}
func (reg *Registry) IDCount() uint {
	return reg.idFilter.Size()
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

// This function generates a good-quality random NodeID (a 32-byte
// value) that is not already known to the registry and then adds
// the new NodeID to the registry's Bloom filter.
func (reg *Registry) UniqueNodeID() (nodeID *xi.NodeID, err error) {

	nodeID, err = xi.New(nil)
	found, err := reg.ContainsID(nodeID)
	for err == nil && found {
		nodeID, err = xi.New(nil)
		found, err = reg.ContainsID(nodeID)
	}
	if err == nil {
		err = reg.idFilter.Insert(nodeID.Value())
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
	ss = []string{"registry {"}
	ss = append(ss, fmt.Sprintf("    LogFile: %s", reg.LogFile))
	ss = append(ss, "}")

	for i := 0; i < len(reg.Clusters); i++ {
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
