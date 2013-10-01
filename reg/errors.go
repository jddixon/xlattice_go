package reg

import (
	"errors"
)

var (
	BadAttrsLine              = errors.New("badly formed attrs line")
	BadVersion				  = errors.New("badly formated VERSION")
	CantFindClusterByID       = errors.New("cannot find cluster with this ID")
	CantFindClusterByName     = errors.New("cannot find cluster with this name")
	ClientMustHaveEndPoint    = errors.New("client must have at least one endPoint")
	ClusterMustHaveTwo        = errors.New("cluster must have at least two members")
	IDAlreadyInUse            = errors.New("ID already in use")
	IllFormedCluster          = errors.New("ill-formed cluster serialization")
	MissingClosingBrace       = errors.New("missing closing brace")
	MissingClusterNameOrID    = errors.New("missing cluster name or ID")
	MissingEndPointsSection   = errors.New("missing endPoints section")
	MissingMembersList        = errors.New("missing members list")
	MissingServerInfo         = errors.New("missing server info")
	NameAlreadyInUse          = errors.New("name already in use")
	NilCluster                = errors.New("nil cluster argument")
	NilPrivateKey             = errors.New("nil private key argument")
	NilRegistry               = errors.New("nil registry argument")
	NilToken                  = errors.New("nil XLRegMsg_Token argument")
	RcvdInvalidMsgForState    = errors.New("invalid msg type for current state")
	UnexpectedMsgType         = errors.New("unexpected message type")
	WrongNumberOfBytesInAttrs = errors.New("wrong number of bytes in attrs")
)
