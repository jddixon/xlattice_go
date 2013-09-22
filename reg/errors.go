package reg

import (
	"errors"
)

var (
	BadAttrsLine              = errors.New("badly formed attrs line")
	ClientMustHaveEndPoint    = errors.New("client must have at least one endPoint")
	ClusterMustHaveTwo        = errors.New("cluster must have at least two members")
	IllFormedCluster          = errors.New("ill-formed cluster serialization")
	MissingClosingBrace       = errors.New("missing closing brace")
	MissingClusterNameOrID    = errors.New("missing cluster name or ID")
	MissingEndPointsSection   = errors.New("missing endPoints section")
	MissingMembersList        = errors.New("missing members list")
	MissingServerInfo         = errors.New("missing server info")
	NilRegistry               = errors.New("nil registry argument")
	UnexpectedMsgType         = errors.New("unexpected message type")
	WrongNumberOfBytesInAttrs = errors.New("wrong number of bytes in attrs")
)
