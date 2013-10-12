package reg

// xlattice_go/reg/user_client.go

import (
	"crypto/rsa"
	"fmt"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xt "github.com/jddixon/xlattice_go/transport"
	"io"
)

var _ = fmt.Print

// The UserClient is created to enable the caller to join a cluster
// and learn information about the cluster's other members.  Once the
// client has learned that information, it is done.

// As implemented so far, this is an ephemeral client, meaning that it
// neither saves nor restores its Node; keys and such are generated for
// each instance.

// For practical use, it is essential that the UserClient create its
// Node when NewUserClient() is first called, but then save its
// configuration.  This is conventionally written to LFS/.xlattice/config.
// On subsequent the client reads its configuration file rather than
// regenerating keys, etc.

type UserClient struct {
	members []ClusterMember

	ClientNode
}

func NewUserClient(
	name, lfs string,
	serverName string, serverID *xi.NodeID, serverEnd xt.EndPointI,
	serverCK, serverSK *rsa.PublicKey,
	clusterName string, size int, 
	epCount int, e []xt.EndPointI) (ac *UserClient, err error) {

	var attrs uint64

	if lfs == "" {
		attrs |= ATTR_EPHEMERAL
	}
	cn, err := NewClientNode(name, lfs, attrs, 
		serverName, serverID, serverEnd,
		serverCK, serverSK, //  *rsa.PublicKey,
		clusterName, size,
		epCount, e)

	if err == nil {
		// Run() fills in clusterID
		ac = &UserClient{
			ClientNode:  *cn,
		}
	}
	return

}

// Start the client running in separate goroutine, so that this function
// is non-blocking.

func (uc *UserClient) Run() (err error) {

	cn := &uc.ClientNode

	go func() {
		var (
			version1 uint32
		)
		clientName := cn.GetName()
		cnx, version2, err := cn.SessionSetup(version1)
		_ = version2 // not yet used
		if err == nil {
			err = cn.ClientAndOK()
		}
		// XXX MODIFY TO SKIP THIS STEP
		if err == nil {
			err = cn.CreateAndReply()
		}
		// XXX MODIFY TO USE CLUSTER_ID PASSED TO UserClient
		if err == nil {
			err = cn.JoinAndReply()
		}
		if err == nil {
			err = cn.GetAndMembers()
		}
		if err == nil {
			err = cn.ByeAndAck()
		}

		// END OF RUN ===============================================
		if cnx != nil {
			cnx.Close()
		}
		// DEBUG
		fmt.Printf("user client %s run complete ", clientName)
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
