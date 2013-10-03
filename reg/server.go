package reg

// xlattice_go/reg/server.go

import (
	xt "github.com/jddixon/xlattice_go/transport"
)

type RegServer struct {
	acc       xt.AcceptorI // volatile, not serialized
	Testing   bool         // serialized
	Verbosity int          // serialized
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
			// As each client connects its connection is passed to a
			// handler running in a separate goroutine.
			cnx, err := rs.acc.Accept()
			if err != nil {
				// Any I/O error shuts down the server.
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
				// XXX notice the error has no effect
			}()
		}
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
