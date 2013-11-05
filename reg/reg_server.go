package reg

// xlattice_go/reg/reg_server.go

import (
	xt "github.com/jddixon/xlattice_go/transport"
	"net"
)

type RegServer struct {
	acc       xt.AcceptorI // volatile, not serialized
	Testing   bool         // serialized
	Verbosity int          // serialized
	DoneCh    chan (bool)
	Registry
}

func NewRegServer(reg *Registry, testing bool, verbosity int) (
	rs *RegServer, err error) {

	if reg == nil {
		err = NilRegistry
	} else {
		acc := reg.GetAcceptor(0) // by convention
		rs = &RegServer{
			acc:       acc,
			Testing:   testing,
			Verbosity: verbosity,
			Registry:  *reg,
			DoneCh:    make(chan bool, 1),
		}
	}
	return
}

func (rs *RegServer) Close() {
	if rs.acc != nil {
		rs.acc.Close()
	}
}
func (rs *RegServer) GetAcceptor() xt.AcceptorI {
	return rs.acc
}

// Starts the server running in a goroutine.  Does not block.
func (rs *RegServer) Run() (err error) {

	go func() {
		for {
			logger := rs.Registry.Logger

			// As each client connects its connection is passed to a
			// handler running in a separate goroutine.
			cnx, err := rs.acc.Accept()
			if err != nil {
				// SHOULD NOT CONTINUE IF 'use of closed network connection";
				// this yields an infinite loop if the listening socket has
				// been closed to shut down the server.
				netOpError, ok := err.(*net.OpError)
				if ok && netOpError.Err.Error() == "use of closed network connection" {
					err = nil
				} else {
					logger.Printf(
						"fatal I/O error %v, shutting down the server\n",
						err)
				}
				break
			}
			go func() {
				var (
					h *InHandler
				)
				h, err = NewInHandler(&rs.Registry, cnx)
				if err == nil {
					err = h.Run()
				}
				if err != nil {
					logger.Printf("I/O error %v, closing client connection\n",
						err)
					cnx := h.Cnx
					if cnx != nil {
						cnx.Close()
					}
				}
			}()
		}
		rs.DoneCh <- true
	}()
	return
}

// SERIALIZATION ====================================================

func ParseRegServer(s string) (rs *RegServer, rest []string, err error) {

	// XXX STUB
	return
}

func (rs *RegServer) String() (s string) {

	// STUB XXX
	return
}

func (rs *RegServer) Strings() (s []string) {

	// STUB XXX
	return
}
