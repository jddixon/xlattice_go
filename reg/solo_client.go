package reg

// xlattice_go/reg/solo_client.go

import (
	"crypto/rsa"
	"fmt"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xt "github.com/jddixon/xlattice_go/transport"
	"io"
)

var _ = fmt.Print

// In the current implementation, a SoloClient provides a simple way to
// create an XLattice Node with a unique NodeID, private RSA keys,
// and an initialized LFS with the configuration stored in
// LFS/.xlattice/node.config with a default mode of 0400.
//
// It requires a registry to give it its NodeID.
type SoloClient struct {
	// In this implementation, SoloClient is a one-shot, launched
	// to create a single cluster

	ClientNode
}

func NewSoloClient(name, lfs string,
	serverName string, serverID *xi.NodeID, serverEnd xt.EndPointI,
	serverCK, serverSK *rsa.PublicKey,
	e []xt.EndPointI) (
	sc *SoloClient, err error) {

	cn, err := NewClientNode(name, lfs, ATTR_SOLO,
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
		clientName := cn.name
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
		err = cn.Persist()

		// DEBUG
		fmt.Printf("SoloClient %s run complete; err is %v\n", clientName, err)
		if err != nil && err != io.EOF {
			fmt.Printf("- ERROR: %v", err)
		}
		fmt.Println("")
		// END

		cn.err = err
		cn.doneCh <- true
	}()
	return
}
