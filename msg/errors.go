package msg

import (
	"errors"
)

var (
	AcceptorNotLive       = errors.New("acceptor is not live")
	BadSig                = errors.New("bad digital signature") // never used?
	CannotSendSecondHello = errors.New("can't send second hello")
	ExpectedMsgOne        = errors.New("expected msg number to be 1")
	MissingHandlerField   = errors.New("missing CnxHandler field")
	MissingHello          = errors.New("expected a Hello msg")
	NilConnection         = errors.New("nil connection")
	NilControlCh          = errors.New("nil control chan argument")
	NilNode               = errors.New("nil node argument")
	NoPeers               = errors.New("node has no peers")
	NotExpectedCommsKey   = errors.New("not peer's expected comms public key")
	NotExpectedNodeID     = errors.New("not peer's expected NodeID")
	NotExpectedSigKey     = errors.New("not peer's expected sig public key")
	UnexpectedMsgType     = errors.New("unexpected message type")
	WrongMsgNbr           = errors.New("wrong message number")
)
