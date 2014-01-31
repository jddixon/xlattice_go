package node

import (
	e "errors"
)

var (
	NilConnection      = e.New("nil connection argument")
	NilConnector       = e.New("nil connector argument")
	NilEndPoint        = e.New("nil endPoint argument")
	NilID              = e.New("nil ID argument")
	NilLFS             = e.New("nil LFS argument")
	NilNodeID          = e.New("nil nodeID argument")
	NilOverlay         = e.New("nil overlay argument")
	NilPeer            = e.New("nil peer argument")
	NotABaseNode       = e.New("not a serialized BaseNode - missing bits")
	NotAKnownPeer      = e.New("not a known peer")
	NotASerializedNode = e.New("not a serialized node")
	NotASerializedPeer = e.New("not a serialized peer")
	NotExpectedOpener  = e.New("not expected BaseNode serialization opener")
	NothingToSign      = e.New("nothing to sign - nil chunks")
)
