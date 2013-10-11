package reg

// xlattice_go/reg/client_node.go

import (
	//"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"fmt"
	//xc "github.com/jddixon/xlattice_go/crypto"
	//xm "github.com/jddixon/xlattice_go/msg"
	xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xt "github.com/jddixon/xlattice_go/transport"
	//xu "github.com/jddixon/xlattice_go/util"
	//"io"
	//"time"
)

var _ = fmt.Print

const (
	CLIENT_START = iota
	HELLO_SENT
	CLIENT_SENT
	CLUSTER_SENT
	JOIN_SENT
	GET_SENT
	BYE_SENT
	CLIENT_CLOSED
)

type ClientNode struct {
	serverName string
	serverID   *xi.NodeID
	serverEnd  xt.EndPointI
	serverCK   *rsa.PublicKey

	// MOVE TO UserClient
	//clusterName string
	//clusterID   *xi.NodeID
	//clusterSize uint32 // this is a FIXED size, aka MaxSize

	// run information
	doneCh chan bool
	err    error
	h      *CnxHandler

	proposedAttrs, decidedAttrs uint64

	iv1, key1    []byte // one-shot
	version1     uint32 // proposed by client
	iv2, key2    []byte // session
	version2     uint32 // decreed by server
	salt1, salt2 []byte // not currently used
	engineC      cipher.Block
	encrypterC   cipher.BlockMode
	decrypterC   cipher.BlockMode

	// information on cluster members		// MOVED TO UserClient
	members []*ClusterMember

	// Information on this cluster member.  By convention endPoints[0]
	// is used for member-member communications and [1] for comms with
	// cluster clients, should they exist.
	clientID  []byte // XXX or *xi.NodeID
	endPoints []xt.EndPointI
	xn.BaseNodeI
}

// Given contact information for a registry and the name of a cluster,
// the client joins the cluster, collects information on the other members,
// and terminates when it has info on the entire membership.

func NewClientNode(
	serverName string, serverID *xi.NodeID, serverEnd xt.EndPointI,
	serverCK *rsa.PublicKey,
	clusterName string, clusterID *xi.NodeID, size int,
	e []xt.EndPointI, bni xn.BaseNodeI) (
	cn *ClientNode, err error) {

	// sanity checks on parameter list
	if serverName == "" || serverID == nil || serverEnd == nil ||
		serverCK == nil {

		err = MissingServerInfo
	} else if clusterName == "" || clusterID == nil {
		err = MissingClusterNameOrID
	} else if size < 2 {
		err = ClusterMustHaveTwo
	} else if len(e) < 1 {
		err = ClientMustHaveEndPoint
	}
	if err == nil {
		cnxHandler := &CnxHandler{State: CLIENT_START}
		cn = &ClientNode{
			doneCh:     make(chan bool, 1),
			serverName: serverName,
			serverID:   serverID,
			serverEnd:  serverEnd,
			serverCK:   serverCK,
			h:          cnxHandler,
			endPoints:  e,

			// THIS BECOMES A NODE XXX, so Node: node,
			BaseNodeI: bni,
		}
	}
	return
}

// Read the next message over the connection
func (cn *ClientNode) readMsg() (m *XLRegMsg, err error) {
	inBuf, err := cn.h.readData()
	if err == nil && inBuf != nil {
		m, err = DecryptUnpadDecode(inBuf, cn.decrypterC)
	}
	return
}

// Write a message out over the connection
func (cn *ClientNode) writeMsg(m *XLRegMsg) (err error) {
	var data []byte
	// serialize, marshal the message
	data, err = EncodePadEncrypt(m, cn.encrypterC)
	if err == nil {
		err = cn.h.writeData(data)
	}
	return
}
