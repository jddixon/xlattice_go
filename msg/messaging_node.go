package msg

// xlattice_go/msg

import (
	"fmt"
	xn "github.com/jddixon/xlattice_go/node"
	xt "github.com/jddixon/xlattice_go/transport"
	"time"
)

type MessagingNode struct {
	Acc    *xt.TcpAcceptor
	K      int
	KillCh chan (bool)
	xn.Node
}

// Create a messaging node around a live node which has a list of peers
// and an open acceptor.

func NewMessagingNode(n *xn.Node, killCh chan (bool)) (
	mn *MessagingNode, err error) {

	var k int
	if n == nil {
		err = NilNode
	}
	if err == nil {
		if k = n.SizePeers(); k == 0 {
			err = NoPeers
		}
	}
	tcpAcc := n.GetAcceptor(0).(*xt.TcpAcceptor)
	if err == nil && tcpAcc == nil {
		err = AcceptorNotLive
	}
	if err == nil && killCh == nil {
		err = NilControlCh
	}
	if err == nil {
		mn = &MessagingNode{
			Acc:    tcpAcc,
			K:      k,
			KillCh: killCh,
			Node:   *n,
		}
	}
	return

}

func (mn *MessagingNode) Start() (err error) {

	// Reactive (server) side ---------------------------------------
	go func() {
		for {
			var conn xt.ConnectionI
			conn, err = mn.Acc.Accept()
			if err != nil {
				break
			}
			go func() {
				NewInHandler(&mn.Node, conn)
			}()
		}
	}()

	// Active (client) side -----------------------------------------
	// A pair of channels for each active component (where node acts
	// as client.
	stopCh := make([]chan bool, mn.K)
	stoppedCh := make([]chan bool, mn.K)
	ohs := make([]*OutHandler, mn.K)
	for i := 0; i < mn.K; i++ {
		stopCh[i] = make(chan bool, 1) // tells the outHandler to stop
		stoppedCh[i] = make(chan bool) // replies that it has done so
	}

	for i := 0; i < mn.K; i++ { // all K peers
		var cnx *xt.TcpConnection
		peer := mn.Node.GetPeer(i)
		if peer == nil {
			fmt.Printf("peer number %d is nil\n", i)
			continue
		}
		ctor := peer.GetConnector(0)
		if ctor == nil {
			fmt.Printf("connector for peer %d is nil\n", i)
			continue
		}
		conn, kErr := ctor.Connect(xt.ANY_TCP_END_POINT)
		if kErr == nil {
			cnx = conn.(*xt.TcpConnection)
			ohs[i] = &OutHandler{
				Node:       &mn.Node,
				CnxHandler: CnxHandler{Cnx: cnx, Peer: peer}}
			go func(k int) {
				defer cnx.Close()
				fErr := ohs[k].SendHello()
				// XXX WAIT FOR ACK OR TIME OUT
				if fErr == nil {
					for {
						select {
						case <-stopCh[k]:
							break
						default:
							time.Sleep(time.Second) // KEEP_ALIVE_INTERVAL
						}
						// SEND KEEP ALIVE
						// WAIT FOR ACK
					}
					// SEND BYE
					// WAIT FOR ACK BUT ON A TIMEOUT
					stoppedCh[k] <- true
				}
			}(i)
		}
	}

	// // wait for all stopped signal
	// select
	//     case timeout in 5 sec:
	//         ;
	//     default:
	//         for n := 0; n < K; n++ {
	//             <- stoppedCh[n]
	//         }
	//
	// XXX STUB XXX
	_, _ = stopCh, stoppedCh
	return
}
