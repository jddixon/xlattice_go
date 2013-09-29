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
	xu "github.com/jddixon/xlattice_go/util"
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

type Client struct {
	serverName  string
	serverID    *xi.NodeID
	serverEnd   xt.EndPointI
	serverCK    *rsa.PublicKey
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

	// information on cluster members
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
			clusterSize: uint32(size),
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
			request, response                     *XLRegMsg
			op                                    XLRegMsg_Tag
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
			// DEBUG
			if err != nil {
				fmt.Printf("Client.Run(): err after write is %v\n", err)
			}
			// END
		}
		// Process HELLO REPLY --------------------------------------
		if err == nil {
			ciphertext2, err = mc.h.readData()
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
			// DEBUG
			fmt.Println("client AES engines set up")
			// END
		}
		// Send CLIENT MSG ==========================================
		if err == nil {
			var ckBytes, skBytes []byte
			var myEnds []string
			// XXX attrs not dealt with
			ckBytes, err = xc.RSAPubKeyToWire(mc.GetCommsPublicKey())
			if err == nil {
				skBytes, err = xc.RSAPubKeyToWire(mc.GetSigPublicKey())
				if err == nil {
					for i := 0; i < len(mc.endPoints); i++ {
						myEnds = append(myEnds, mc.endPoints[i].String())
					}
					clientName := mc.GetName()
					token := &XLRegMsg_Token{
						Name:     &clientName,
						Attrs:    &mc.proposedAttrs,
						ID:       mc.GetNodeID().Value(),
						CommsKey: ckBytes,
						SigKey:   skBytes,
						MyEnds:   myEnds,
					}

					op = XLRegMsg_Client
					request = &XLRegMsg{
						Op:          &op,
						ClientName:  &clientName, // XXX redundant
						ClientSpecs: token,
					}
					// SHOULD CHECK FOR TIMEOUT
					err = mc.writeMsg(request)
					// DEBUG
					fmt.Println("ClientMsg sent")
					// END
				}
			}
		}
		// Process CLIENT_OK ----------------------------------------
		if err == nil {
			// SHOULD CHECK FOR TIMEOUT
			response, err = mc.readMsg()
			if err == nil {
				mc.clientID = response.GetClientID()
				mc.decidedAttrs = response.GetAttrs()
				// DEBUG
				fmt.Println("client has received ClientOK")
				// END
			}
		}
		// Send CREATE MSG ==========================================
		if err == nil {
			op = XLRegMsg_Create
			wireSize := uint32(mc.clusterSize)
			request = &XLRegMsg{
				Op:          &op,
				ClusterName: &mc.clusterName,
				ClusterSize: &wireSize,
			}
			// SHOULD CHECK FOR TIMEOUT
			err = mc.writeMsg(request)
			// DEBUG
			fmt.Printf("Create sent for cluster %s, size %d\n",
				mc.clusterName, mc.clusterSize)
			// END
		}
		// Process CREATE REPLY -------------------------------------

		if err == nil {
			// SHOULD CHECK FOR TIMEOUT AND VERIFY THAT IT'S A CREATE REPLY
			response, err = mc.readMsg()
			op = response.GetOp()
			_ = op
			// DEBUG
			fmt.Printf("client has received CreateReply; err is %v\n", err)
			// END
			if err == nil {
				mc.clusterSize = response.GetClusterSize()
				mc.members = make([]*ClusterMember, mc.clusterSize)
				id := response.GetClusterID()
				mc.clusterID, err = xi.New(id)
			}
		}

		// Send JOIN MSG ============================================
		fmt.Printf("Pre-Join client-side cluster size: %d\n", mc.clusterSize)
		if err == nil {
			op = XLRegMsg_Join
			request = &XLRegMsg{
				Op:          &op,
				ClusterName: &mc.clusterName,
			}
			// SHOULD CHECK FOR TIMEOUT
			err = mc.writeMsg(request)
			// DEBUG
			fmt.Printf("Join by name sent for cluster %s\n", mc.clusterName)
			// END
		}
		// Process JOIN REPLY ---------------------------------------
		if err == nil {
			// SHOULD CHECK FOR TIMEOUT AND VERIFY THAT IT'S A JOIN REPLY
			response, err = mc.readMsg()
			op = response.GetOp()
			_ = op
			// DEBUG
			fmt.Printf("client has received JoinReply; err is %v\n", err)
			// END
			if err == nil {
				// XXX We collect this information for the second time;
				// it might be different!
				clusterSizeNow := response.GetClusterSize()
				if mc.clusterSize != clusterSizeNow {
					mc.clusterSize = clusterSizeNow
					mc.members = make([]*ClusterMember, mc.clusterSize)
				}
				id := response.GetClusterID()
				mc.clusterID, err = xi.New(id)
			}
		} // GEEP

		// COLLECT INFORMATION ON ALL CLUSTER MEMBERS ***************
		fmt.Printf("Cluster size after Join: %d\n", mc.clusterSize)
		stillToGet := xu.LowNMap(uint(mc.clusterSize))

		if err == nil {
			MAX_GET := 16
			// XXX It should be impossible for mc.members to be nil
			// at this point
			if mc.members == nil {
				mc.members = make([]*ClusterMember, mc.clusterSize)
			}
			for count := 0; count < MAX_GET && stillToGet.Any(); count++ {

				fmt.Printf("STILL TO GET: %d (bits 0x%x)\n",
					stillToGet.Count(), stillToGet.Bits)

				// Send GET MSG =========================================
				op = XLRegMsg_Get
				request = &XLRegMsg{
					Op:        &op,
					ClusterID: mc.clusterID.Value(),
					Which:     &stillToGet.Bits,
				}
				// SHOULD CHECK FOR TIMEOUT
				err = mc.writeMsg(request)

				// Process MEMBERS = GET REPLY --------------------------
				if err != nil {
					break
				}
				response, err = mc.readMsg()
				// XXX HANDLE ANY ERROR
				op = response.GetOp()
				// XXX op MUST BE XLRegMsg_Members
				_ = op

				if err == nil {
					id := response.GetClusterID()
					_ = id // XXX ignore for now
					which := xu.NewBitMap64(response.GetWhich())
					// DEBUG
					fmt.Printf("client has received %d Members\n",
						which.Count())
					// END
					tokens := response.GetTokens() // a slice
					offset := 0
					if which.Any() {
						for i := uint(0); i < uint(mc.clusterSize); i++ {
							if which.Test(i) {
								token := tokens[offset]
								offset++
								mc.members[i], err = NewClusterMemberFromToken(
									token)
								stillToGet.Clear(i)
							}
						}
					}
					if stillToGet.None() {
						fmt.Println("HAVE ALL MEMBERS") // DEBUG
						break
					}
					time.Sleep(10 * time.Millisecond)
				}
			} // FOO
		}
		// Send BYE MSG =============================================
		if err == nil {
			op = XLRegMsg_Bye
			request = &XLRegMsg{
				Op: &op,
			}
			// SHOULD CHECK FOR TIMEOUT
			err = mc.writeMsg(request)
			// DEBUG
			fmt.Println("Bye sent")
			// END
		}
		// Process ACK = BYE REPLY ----------------------------------
		if err == nil {
			// SHOULD CHECK FOR TIMEOUT AND VERIFY THAT IT'S AN ACK
			response, err = mc.readMsg()
			op = response.GetOp()
			_ = op
			// DEBUG
			fmt.Printf("client has received Ack; err is %v\n", err)
			// END
			if err == nil {

			}
		} // GEEP

		// END OF RUN ===============================================
		if cnx != nil {
			cnx.Close()
		}

		fmt.Print("CLIENT RUN COMPLETE ")
		if err != nil {
			fmt.Printf("- ERROR: %v", err)
		}
		fmt.Println("")

		mc.err = err
		mc.doneCh <- true
	}()
	return
}

// Read the next message over the connection
func (mc *Client) readMsg() (m *XLRegMsg, err error) {
	inBuf, err := mc.h.readData()
	if err == nil && inBuf != nil {
		m, err = DecryptUnpadDecode(inBuf, mc.decrypterC)
	}
	return
}

// Write a message out over the connection
func (mc *Client) writeMsg(m *XLRegMsg) (err error) {
	var data []byte
	// serialize, marshal the message
	data, err = EncodePadEncrypt(m, mc.encrypterC)
	if err == nil {
		err = mc.h.writeData(data)
	}
	return
}
