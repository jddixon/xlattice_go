package reg

// xlattice_go/reg/client_node.go

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	xm "github.com/jddixon/xlattice_go/msg"
	xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xt "github.com/jddixon/xlattice_go/transport"
	xu "github.com/jddixon/xlattice_go/util"
	// "io"
	"time"
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

	// The significance of these fields is different in different
	// subclasses.  To the user, this is the cluster being joined.
	// To an admin client, it is a cluster being created.
	clusterName string
	clusterID   *xi.NodeID
	clusterSize uint32 // this is a FIXED size, aka MaxSize

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

// RUN CODE =========================================================

func (cn *ClientNode) SessionSetup(version1 uint32) (
	cnx *xt.TcpConnection, version2 uint32, err error) {
	var (
		ciphertext1, iv1, key1, salt1, salt1c []byte
		ciphertext2, iv2, key2, salt2         []byte
	)
	// Set up connection to server. -----------------------------
	ctor, err := xt.NewTcpConnector(cn.serverEnd)
	if err == nil {
		var conn xt.ConnectionI
		conn, err = ctor.Connect(nil)
		if err == nil {
			cnx = conn.(*xt.TcpConnection)
		}
	}
	// Send HELLO -----------------------------------------------
	if err == nil {
		cn.h.Cnx = cnx
		ciphertext1, iv1, key1, salt1,
			err = xm.ClientEncodeHello(version1, cn.serverCK)
	}
	if err == nil {
		err = cn.h.writeData(ciphertext1)
		// DEBUG
		if err != nil {
			fmt.Printf("Client.Run(): err after write is %v\n", err)
		}
		// END
	}
	// Process HELLO REPLY --------------------------------------
	if err == nil {
		ciphertext2, err = cn.h.readData()
		// DEBUG
		if err != nil {
			fmt.Printf("Client.Run(): err after read is %v\n", err)
		}
		// END
	}
	if err == nil {
		iv2, key2, salt2, salt1c, version2,
			err = xm.ClientDecodeHelloReply(ciphertext2, iv1, key1)
		_ = salt1c // XXX
	}
	// Set up AES engines ---------------------------------------
	if err == nil {
		cn.salt1 = salt1
		cn.iv2 = iv2
		cn.key2 = key2
		cn.salt2 = salt2
		cn.version2 = version2
		cn.engineC, err = aes.NewCipher(key2)
		if err == nil {
			cn.encrypterC = cipher.NewCBCEncrypter(cn.engineC, iv2)
			cn.decrypterC = cipher.NewCBCDecrypter(cn.engineC, iv2)
		}
		// DEBUG
		fmt.Printf("client %s AES engines set up\n", cn.GetName())
		// END
	}
	return
}

func (cn *ClientNode) ClientAndOK() (err error) {

	var (
		ckBytes, skBytes []byte
		myEnds           []string
	)
	clientName := cn.GetName()
	// XXX attrs not dealt with

	// Send CLIENT MSG ==========================================
	ckBytes, err = xc.RSAPubKeyToWire(cn.GetCommsPublicKey())
	if err == nil {
		skBytes, err = xc.RSAPubKeyToWire(cn.GetSigPublicKey())
		if err == nil {
			for i := 0; i < len(cn.endPoints); i++ {
				myEnds = append(myEnds, cn.endPoints[i].String())
			}
			token := &XLRegMsg_Token{
				Name:     &clientName,
				Attrs:    &cn.proposedAttrs,
				ID:       cn.GetNodeID().Value(),
				CommsKey: ckBytes,
				SigKey:   skBytes,
				MyEnds:   myEnds,
			}

			op := XLRegMsg_Client
			request := &XLRegMsg{
				Op:          &op,
				ClientName:  &clientName, // XXX redundant
				ClientSpecs: token,
			}
			// SHOULD CHECK FOR TIMEOUT
			err = cn.writeMsg(request)
			// DEBUG
			fmt.Printf("CLIENT_MSG for %s sent\n", clientName)
			// END
		}
	}
	// Process CLIENT_OK --------------------------------------------
	// SHOULD CHECK FOR TIMEOUT
	response, err := cn.readMsg()
	if err == nil {
		cn.clientID = response.GetClientID()
		cn.decidedAttrs = response.GetAttrs()
		// DEBUG
		fmt.Printf("    client %s has received ClientOK\n",
			cn.GetName())
		// END
	}
	return
} // GEEP2

func (cn *ClientNode) CreateAndReply() (err error) {

	var response *XLRegMsg
	clientName := cn.GetName()

	// Send CREATE MSG ==========================================
	op := XLRegMsg_Create
	wireSize := uint32(cn.clusterSize)
	request := &XLRegMsg{
		Op:          &op,
		ClusterName: &cn.clusterName,
		ClusterSize: &wireSize,
	}
	// SHOULD CHECK FOR TIMEOUT
	err = cn.writeMsg(request)
	// DEBUG
	fmt.Printf("client %s sends CREATE for cluster %s, size %d\n",
		clientName, cn.clusterName, cn.clusterSize)
	// END

	if err == nil {
		// Process CREATE REPLY -------------------------------------
		// SHOULD CHECK FOR TIMEOUT AND VERIFY THAT IT'S A CREATE REPLY
		response, err = cn.readMsg()
		op = response.GetOp()
		_ = op
		// DEBUG
		fmt.Printf("    client has received CreateReply; err is %v\n", err)
		// END
		if err == nil {
			cn.clusterSize = response.GetClusterSize()
			cn.members = make([]*ClusterMember, cn.clusterSize)
			id := response.GetClusterID()
			cn.clusterID, err = xi.New(id)
		}
	}
	return
} // GEEPGEEP

func (cn *ClientNode) JoinAndReply() (err error) {

	clientName := cn.GetName() // DEBUG

	// Send JOIN MSG ============================================
	fmt.Printf("Pre-Join client-side cluster size: %d\n", cn.clusterSize)
	op := XLRegMsg_Join
	request := &XLRegMsg{
		Op:          &op,
		ClusterName: &cn.clusterName,
	}
	// SHOULD CHECK FOR TIMEOUT
	err = cn.writeMsg(request)
	// DEBUG
	fmt.Printf("Client %s sends JOIN by name sent for cluster %s\n",
		clientName, cn.clusterName)
	// END

	// Process JOIN REPLY ---------------------------------------
	if err == nil {
		var response *XLRegMsg

		// SHOULD CHECK FOR TIMEOUT AND VERIFY THAT IT'S A JOIN REPLY
		response, err = cn.readMsg()
		op := response.GetOp()
		_ = op
		// DEBUG
		fmt.Printf("    client has received JoinReply; err is %v\n", err)
		// END
		if err == nil {
			// XXX We collect this information for the second time;
			// it might be different!
			clusterSizeNow := response.GetClusterSize()
			if cn.clusterSize != clusterSizeNow {
				cn.clusterSize = clusterSizeNow
				cn.members = make([]*ClusterMember, cn.clusterSize)
			}
			id := response.GetClusterID()
			cn.clusterID, err = xi.New(id)
		}
	} // GEEP3
	return
}

// Collect information on all cluster members
func (cn *ClientNode) GetAndMembers() (err error) {

	clientName := cn.GetName() // DEBUG

	MAX_GET := 16
	// XXX It should be impossible for cn.members to be nil
	// at this point
	if cn.members == nil {
		cn.members = make([]*ClusterMember, cn.clusterSize)
		// DEBUG
		fmt.Println("Client.Run after Join: UNEXPECTED MAKE cn.members")
	} else {
		fmt.Println("Client.Run after Join: NO NEED TO MAKE cn.members")
		// END
	}
	stillToGet := xu.LowNMap(uint(cn.clusterSize))
	for count := 0; count < MAX_GET && stillToGet.Any(); count++ {
		var response *XLRegMsg

		for i := uint(0); i < uint(cn.clusterSize); i++ {
			if cn.members[i] != nil {
				stillToGet = stillToGet.Clear(i)
			}
		}
		// DEBUG
		fmt.Printf("Client %s sends GET for %d members (bits 0x%x)\n",
			clientName, stillToGet.Count(), stillToGet.Bits)
		// END

		// Send GET MSG =========================================
		op := XLRegMsg_GetCluster
		request := &XLRegMsg{
			Op:        &op,
			ClusterID: cn.clusterID.Value(),
			Which:     &stillToGet.Bits,
		}
		// SHOULD CHECK FOR TIMEOUT
		err = cn.writeMsg(request)

		// Process MEMBERS = GET REPLY --------------------------
		if err != nil {
			break
		}
		response, err = cn.readMsg()
		// XXX HANDLE ANY ERROR
		op = response.GetOp()
		// XXX op MUST BE XLRegMsg_Members
		_ = op

		if err == nil {
			id := response.GetClusterID()
			_ = id // XXX ignore for now
			which := xu.NewBitMap64(response.GetWhich())
			// DEBUG
			fmt.Printf("    client has received %d MEMBERS\n",
				which.Count())
			// END
			tokens := response.GetTokens() // a slice
			if which.Any() {
				offset := 0
				for i := uint(0); i < uint(cn.clusterSize); i++ {
					if which.Test(i) {
						token := tokens[offset]
						offset++
						cn.members[i], err = NewClusterMemberFromToken(
							token)
						stillToGet = stillToGet.Clear(i)
					}
				}
			}
			if stillToGet.None() {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
	return
} // GEEP4

// Send Bye, wait for and process Ack.

func (cn *ClientNode) ByeAndAck() (err error) {

	clientName := cn.GetName() // DEBUG

	op := XLRegMsg_Bye
	request := &XLRegMsg{
		Op: &op,
	}
	// SHOULD CHECK FOR TIMEOUT
	err = cn.writeMsg(request)
	// DEBUG
	fmt.Printf("client %s BYE sent\n", clientName)
	// END

	// Process ACK = BYE REPLY ----------------------------------
	if err == nil {
		var response *XLRegMsg

		// SHOULD CHECK FOR TIMEOUT AND VERIFY THAT IT'S AN ACK
		response, err = cn.readMsg()
		op := response.GetOp()
		_ = op
		// DEBUG
		fmt.Printf("    client %s has received ACK; err is %v\n",
			clientName, err)
		// END
	}
	return
} // GEEP6

//// Start the client running in separate goroutine, so that this function
//// is non-blocking.
//
//func (cn *ClientNode) Run() (err error) {
//	go func() {
//		var (
//			version1 uint32
//		)
//		clientName := cn.GetName()
//		cnx, version2, err := cn.SessionSetup(version1)
//		_ = version2 // not yet used
//		if err == nil {
//			err = cn.ClientAndOK()
//		}
//		if err == nil {
//			err = cn.CreateAndReply()
//		}
//		if err == nil {
//			err = cn.JoinAndReply()
//		}
//		if err == nil {
//			err = cn.GetAndMembers()
//		}
//		if err == nil {
//			err = cn.ByeAndAck()
//		}
//
//		// END OF RUN ===============================================
//		if cnx != nil {
//			cnx.Close()
//		}
//		// DEBUG
//		fmt.Printf("client %s run complete ", clientName)
//		if err != nil && err != io.EOF {
//			fmt.Printf("- ERROR: %v", err)
//		}
//		fmt.Println("")
//		// END
//
//		cn.err = err
//		cn.doneCh <- true
//	}()
//	return
//}
