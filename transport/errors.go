package transport

// xlattice_go/transport/errors.go

import (
	"errors"
)

// Go won't accept these as constants
var (
	AlreadyBound       = errors.New("cnx has already been bound")
	AlreadyConnected   = errors.New("cnx has already been connected")
	NotAConnector      = errors.New("Not a connector")
	NotAKnownConnector = errors.New("Not a known connector type")
	NotAKnownEndPoint  = errors.New("Not a known endPoint type")
	NilConnection      = errors.New("nil connection")
	NilEndPoint        = errors.New("nil endpoint argument")
	NotBound           = errors.New("connection has not been bound")
	NotAMockEndPoint   = errors.New("Not a mock endPoint")
	NotAnEndPoint      = errors.New("Not an endPoint")
	NotImplemented     = errors.New("not implemented")
	NotMockEndPoint    = errors.New("not a Mock endpoint")
	NotTcpEndPoint     = errors.New("not a Tcp endpoint")
)
