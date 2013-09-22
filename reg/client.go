package reg

// xlattice_go/reg/client.go

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

type Client struct {
	serverName  string
	serverID    *xi.NodeID
	serverEnd   xt.EndPointI
	serverCK    *rsa.PublicKey
	clusterName string
	clusterID   *xi.NodeID
	size        int

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

	// information on other cluster members
	others []*ClusterMember

	// Information on this cluster member.  By convention endPoints[0]
	// is used for member-member communications and [1] for comms with
	// cluster clients, should they exist.

	endPoints []xt.EndPointI
	xn.BaseNodeI
}

// Given contact information for a registry and the name of a cluster,
// the client joins the cluster, collects information on the other members,
// and terminates when it has info on the entire membership.

func NewClient(
	serverName string, serverID *xi.NodeID, serverEnd xt.EndPointI,
	serverCK *rsa.PublicKey,
	clusterName string, clusterID *xi.NodeID, size int,
	e []xt.EndPointI, bni xn.BaseNodeI) (
	mc *Client, err error) {

	// sanity checks on parameter list
	if serverName == "" || serverID == nil || serverEnd == nil || serverCK == nil {
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
		mc = &Client{
			doneCh:      make(chan bool, 1),
			serverName:  serverName,
			serverID:    serverID,
			serverEnd:   serverEnd,
			serverCK:    serverCK,
			clusterName: clusterName,
			clusterID:   clusterID,
			size:        size,
			h:           cnxHandler,
			endPoints:   e,
			BaseNodeI:   bni,
		}
	}
	return
}

// Start the client running in separate goroutine, so that this function
// is non-blocking.

func (mc *Client) Run() (err error) {
	go func() {
		var (
			err                                   error
			cnx                                   *xt.TcpConnection
			ciphertext1, iv1, key1, salt1, salt1c []byte
			ciphertext2, iv2, key2, salt2         []byte
			version1, version2                    uint32
		)
		// Set up connection to server. -----------------------------
		ctor, err := xt.NewTcpConnector(mc.serverEnd)
		if err == nil {
			var conn xt.ConnectionI
			conn, err = ctor.Connect(nil)
			if err == nil {
				cnx = conn.(*xt.TcpConnection)
			}
		}
		// Send HELLO -----------------------------------------------
		if err == nil {
			mc.h.Cnx = cnx
			ciphertext1, iv1, key1, salt1,
				err = xm.ClientEncodeHello(version1, mc.serverCK)
		}
		if err == nil {
			err = mc.h.writeData(ciphertext1)
		}
		// Process HELLO REPLY --------------------------------------
		if err == nil {
			ciphertext2, err = mc.h.readData()
		}
		if err == nil {
			iv2, key2, salt2, salt1c, version2,
				err = xm.ClientDecodeHelloReply(ciphertext2, iv1, key1)
			_ = salt1c // XXX
		}
		// Set up AES engine ----------------------------------------
		if err == nil {
			mc.salt1 = salt1
			mc.iv2 = iv2
			mc.key2 = key2
			mc.salt2 = salt2
			mc.version2 = version2
			mc.engineC, err = aes.NewCipher(key2)
			if err == nil {
				mc.encrypterC = cipher.NewCBCEncrypter(mc.engineC, iv2)
				mc.decrypterC = cipher.NewCBCDecrypter(mc.engineC, iv2)
			}
		}
		// Send CLIENT MSG ------------------------------------------
		if err == nil {
			var ckBytes, skBytes, ciphertext []byte
			var myEnds []string
			// XXX attrs not dealt with
			ckBytes, err = xc.RSAPubKeyToWire(mc.GetCommsPublicKey())
			if err == nil {
				skBytes, err = xc.RSAPubKeyToWire(mc.GetSigPublicKey())
				if err == nil {
					for i := 0; i < len(mc.endPoints); i++ {
						myEnds = append(myEnds, mc.endPoints[i].String())
					}
					token := &XLRegMsg_Token{
						Attrs:    &mc.proposedAttrs,
						ID:       mc.GetNodeID().Value(),
						CommsKey: ckBytes,
						SigKey:   skBytes,
						MyEnds:   myEnds,
					}

					op := XLRegMsg_Client
					clientName := mc.GetName()
					clientMsg := XLRegMsg{
						Op:          &op,
						ClientName:  &clientName,
						ClientSpecs: token,
					}
					ciphertext, err = EncodePadEncrypt(&clientMsg, mc.encrypterC)
					err = mc.h.writeData(ciphertext)
				}
			}
		}
		// Process CLIENT_OK ----------------------------------------

		// END OF RUN -----------------------------------------------
		if cnx != nil {
			cnx.Close()
		}

		fmt.Println("CLIENT RUN COMPLETE")

		mc.err = err
		mc.doneCh <- true
	}()
	return
}
