package node

import (
	"errors"
)

var (
	NilConnection      = errors.New("nil connection argument")
	NilConnector       = errors.New("nil connector argument")
	NilEndPoint        = errors.New("nil endPoint argument")
	NilID			   = errors.New("nil ID argument")
	NilLFS             = errors.New("nil LFS argument")
	NilNodeID          = errors.New("nil nodeID argument")
	NilOverlay         = errors.New("nil overlay argument")
	NilPeer            = errors.New("nil peer argument")
	NotABaseNode       = errors.New("not a serialized BaseNode - missing bits")
	NotAKnownPeer      = errors.New("not a known peer")
	NotASerializedNode = errors.New("not a serialized node")
	NotASerializedPeer = errors.New("not a serialized peer")
	NotExpectedOpener  = errors.New("not expected BaseNode serialization opener")
	NothingToSign      = errors.New("nothing to sign - nil chunks")
)
