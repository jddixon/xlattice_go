package reg

// xlattice_go/reg/old_client_run.go

//////////////////////////
// THIS IS BEING REPLACED.
//////////////////////////

import (
	"crypto/aes"
	"crypto/cipher"
	//"crypto/rsa"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	xm "github.com/jddixon/xlattice_go/msg"
	// xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xt "github.com/jddixon/xlattice_go/transport"
	xu "github.com/jddixon/xlattice_go/util"
	"io"
	"time"
)

func (mc *OldClient) SessionSetup(version1 uint32) (
	cnx *xt.TcpConnection, version2 uint32, err error) {
	var (
		ciphertext1, iv1, key1, salt1, salt1c []byte
		ciphertext2, iv2, key2, salt2         []byte
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
		fmt.Printf("client %s AES engines set up\n", mc.GetName())
		// END
	}
	return
}

func (mc *OldClient) ClientAndOK() (err error) {

	var (
		ckBytes, skBytes []byte
		myEnds           []string
	)
	clientName := mc.GetName()
	// XXX attrs not dealt with

	// Send CLIENT MSG ==========================================
	ckBytes, err = xc.RSAPubKeyToWire(mc.GetCommsPublicKey())
	if err == nil {
		skBytes, err = xc.RSAPubKeyToWire(mc.GetSigPublicKey())
		if err == nil {
			for i := 0; i < len(mc.endPoints); i++ {
				myEnds = append(myEnds, mc.endPoints[i].String())
			}
			token := &XLRegMsg_Token{
				Name:     &clientName,
				Attrs:    &mc.proposedAttrs,
				ID:       mc.GetNodeID().Value(),
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
			err = mc.writeMsg(request)
			// DEBUG
			fmt.Printf("CLIENT_MSG for %s sent\n", clientName)
			// END
		}
	}
	// Process CLIENT_OK --------------------------------------------
	// SHOULD CHECK FOR TIMEOUT
	response, err := mc.readMsg()
	if err == nil {
		mc.clientID = response.GetClientID()
		mc.decidedAttrs = response.GetClientAttrs()
		// DEBUG
		fmt.Printf("    client %s has received ClientOK\n",
			mc.GetName())
		// END
	}
	return
} // GEEP2

func (mc *OldClient) CreateAndReply() (err error) {

	var response *XLRegMsg
	clientName := mc.GetName()

	// Send CREATE MSG ==========================================
	op := XLRegMsg_Create
	wireSize := uint32(mc.clusterSize)
	request := &XLRegMsg{
		Op:          &op,
		ClusterName: &mc.clusterName,
		ClusterSize: &wireSize,
	}
	// SHOULD CHECK FOR TIMEOUT
	err = mc.writeMsg(request)
	// DEBUG
	fmt.Printf("client %s sends CREATE for cluster %s, size %d\n",
		clientName, mc.clusterName, mc.clusterSize)
	// END

	if err == nil {
		// Process CREATE REPLY -------------------------------------
		// SHOULD CHECK FOR TIMEOUT AND VERIFY THAT IT'S A CREATE REPLY
		response, err = mc.readMsg()
		op = response.GetOp()
		_ = op
		// DEBUG
		fmt.Printf("    client has received CreateReply; err is %v\n", err)
		// END
		if err == nil {
			mc.clusterSize = response.GetClusterSize()
			mc.members = make([]*ClusterMember, mc.clusterSize)
			id := response.GetClusterID()
			mc.clusterID, err = xi.New(id)
		}
	}
	return
} // GEEPGEEP

func (mc *OldClient) JoinAndReply() (err error) {

	clientName := mc.GetName() // DEBUG

	// Send JOIN MSG ============================================
	fmt.Printf("Pre-Join client-side cluster size: %d\n", mc.clusterSize)
	op := XLRegMsg_Join
	request := &XLRegMsg{
		Op:          &op,
		ClusterName: &mc.clusterName,
	}
	// SHOULD CHECK FOR TIMEOUT
	err = mc.writeMsg(request)
	// DEBUG
	fmt.Printf("Client %s sends JOIN by name sent for cluster %s\n",
		clientName, mc.clusterName)
	// END

	// Process JOIN REPLY ---------------------------------------
	if err == nil {
		var response *XLRegMsg

		// SHOULD CHECK FOR TIMEOUT AND VERIFY THAT IT'S A JOIN REPLY
		response, err = mc.readMsg()
		op := response.GetOp()
		_ = op
		// DEBUG
		fmt.Printf("    client has received JoinReply; err is %v\n", err)
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
	} // GEEP3
	return
}

// Collect information on all cluster members
func (mc *OldClient) GetAndMembers() (err error) {

	clientName := mc.GetName() // DEBUG

	MAX_GET := 16
	// XXX It should be impossible for mc.members to be nil
	// at this point
	if mc.members == nil {
		mc.members = make([]*ClusterMember, mc.clusterSize)
		// DEBUG
		fmt.Println("Client.Run after Join: UNEXPECTED MAKE mc.members")
	} else {
		fmt.Println("Client.Run after Join: NO NEED TO MAKE mc.members")
		// END
	}
	stillToGet := xu.LowNMap(uint(mc.clusterSize))
	for count := 0; count < MAX_GET && stillToGet.Any(); count++ {
		var response *XLRegMsg

		for i := uint(0); i < uint(mc.clusterSize); i++ {
			if mc.members[i] != nil {
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
			fmt.Printf("    client has received %d MEMBERS\n",
				which.Count())
			// END
			tokens := response.GetTokens() // a slice
			if which.Any() {
				offset := 0
				for i := uint(0); i < uint(mc.clusterSize); i++ {
					if which.Test(i) {
						token := tokens[offset]
						offset++
						mc.members[i], err = NewClusterMemberFromToken(
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

func (mc *OldClient) ByeAndAck() (err error) {

	clientName := mc.GetName() // DEBUG

	op := XLRegMsg_Bye
	request := &XLRegMsg{
		Op: &op,
	}
	// SHOULD CHECK FOR TIMEOUT
	err = mc.writeMsg(request)
	// DEBUG
	fmt.Printf("client %s BYE sent\n", clientName)
	// END

	// Process ACK = BYE REPLY ----------------------------------
	if err == nil {
		var response *XLRegMsg

		// SHOULD CHECK FOR TIMEOUT AND VERIFY THAT IT'S AN ACK
		response, err = mc.readMsg()
		op := response.GetOp()
		_ = op
		// DEBUG
		fmt.Printf("    client %s has received ACK; err is %v\n",
			clientName, err)
		// END
	}
	return
} // GEEP6

// Start the client running in separate goroutine, so that this function
// is non-blocking.

func (mc *OldClient) Run() (err error) {
	go func() {
		var (
			version1 uint32
		)
		clientName := mc.GetName()
		cnx, version2, err := mc.SessionSetup(version1)
		_ = version2 // not yet used
		if err == nil {
			err = mc.ClientAndOK()
		}
		if err == nil {
			err = mc.CreateAndReply()
		}
		if err == nil {
			err = mc.JoinAndReply()
		}
		if err == nil {
			err = mc.GetAndMembers()
		}
		if err == nil {
			err = mc.ByeAndAck()
		}

		// END OF RUN ===============================================
		if cnx != nil {
			cnx.Close()
		}
		// DEBUG
		fmt.Printf("client %s run complete ", clientName)
		if err != nil && err != io.EOF {
			fmt.Printf("- ERROR: %v", err)
		}
		fmt.Println("")
		// END

		mc.err = err
		mc.doneCh <- true
	}()
	return
}
