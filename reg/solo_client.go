package reg

// xlattice_go/reg/solo_client.go

import (
	"crypto/rsa"
	"fmt"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xt "github.com/jddixon/xlattice_go/transport"
)

var _ = fmt.Print

// In the current implementation, a SoloClient provides a simple way to
// create an XLattice Node with a unique NodeID, private RSA keys,
// and an initialized LFS with the configuration stored in
// LFS/.xlattice/node.config with a default mode of 0400.
//
// It requires a registry to give it its NodeID.
//
// In other words, use a SoloClient to create and persist an XLattice
// node that will NOT be a member of a cluster but WILL have its
// configuration saved to permanent storage, to its local file system
// (LFS).

type SoloClient struct {
	// In this implementation, SoloClient is a one-shot, launched
	// to create a solitary node.

	ClientNode
}

func NewSoloClient(name, lfs string,
	serverName string, serverID *xi.NodeID, serverEnd xt.EndPointI,
	serverCK, serverSK *rsa.PublicKey,
	e []xt.EndPointI) (
	sc *SoloClient, err error) {

	cn, err := NewClientNode(name, lfs, nil, nil, ATTR_SOLO,
		serverName, serverID, serverEnd, serverCK, serverSK,
		"", uint64(0), nil, 0, // no cluster
		len(e), e)

	if err == nil {
		// Run() fills in clusterID
		sc = &SoloClient{
			ClientNode: *cn,
		}
	}
	return
}

// Start the client running in separate goroutine, so that this function
// is non-blocking.

func (sc *SoloClient) Run() (err error) {

	cn := &sc.ClientNode

	go func() {
		var (
			version1 uint32
		)
		cnx, version2, err := cn.SessionSetup(version1)
		_ = version2 // not yet used
		if err == nil {
			err = cn.ClientAndOK()
		}
		if err == nil {
			err = cn.ByeAndAck()
		}
		// END OF RUN ===============================================

		if cnx != nil {
			cnx.Close()
		}
		// Create the Node and write its configuration to the usual place
		// in the file system: LFS/.xlattice/node.config.
		err = cn.PersistNode()
		cn.Err = err
		cn.DoneCh <- true
	}()
	return
}
