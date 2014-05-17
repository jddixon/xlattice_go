package msg

// xlattice_go/msg

import (
	"fmt"
	xn "github.com/jddixon/xlNode_go"
	xt "github.com/jddixon/xlTransport_go"
	"time"
)

type MessagingNode struct {
	Acc       *xt.TcpAcceptor
	K         int
	StopCh    chan (bool)
	StoppedCh chan bool
	xn.Node
}

// Create a messaging node around a live node which has a list of peers
// and an open acceptor.

func NewMessagingNode(n *xn.Node, stopCh, stoppedCh chan bool) (
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
	if err == nil && stopCh == nil {
		err = NilControlCh
	}
	if err == nil {
		mn = &MessagingNode{
			Acc:       tcpAcc,
			K:         k,
			StopCh:    stopCh,
			StoppedCh: stoppedCh,
			Node:      *n,
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
				NewInHandler(&mn.Node, conn, mn.StopCh, mn.StoppedCh)
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

	// Open client-style connections to all K peers, sending keep-alives
	// to each until we get a signal to stop or an error occurs.
	for i := 0; i < mn.K; i++ {
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
				defer ohs[k].Cnx.Close()
				var fErr error
				errCh := make(chan error, 1)

				// hello initiates conversation ---------------------
				go func() { errCh <- ohs[k].SendHello() }()
				select {
				case <-time.After(time.Second):
					// handle timeout
				case fErr = <-errCh:
					// XXX WAIT FOR ACK OR TIME OUT
				}
				if fErr == nil {
					// keepalive/ack loop until stop received -------
					for {
						select {
						case <-stopCh[k]:
							ohs[k].Cnx.Close()
							break
						case <-time.After(time.Second):
							// KEEP_ALIVE_INTERVAL
						}
						// SEND KEEP ALIVE
						// WAIT FOR ACK
					}
					// out of loop: send Bye if possible
					// XXX STUB XXX

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
