package reg

//xlattice_go/reg/xlreg_errors.go

import (
	e "errors"
)

var (
	BadAttrsLine                   = e.New("badly formed attrs line")
	BadVersion                     = e.New("badly formated VERSION")
	CantFindClusterByID            = e.New("cannot find cluster with this ID")
	CantFindClusterByName          = e.New("cannot find cluster with this name")
	ClientMustHaveEndPoint         = e.New("client must have at least one endPoint")
	ClusterMembersMustHaveEndPoint = e.New("cluster members must have at least one endPoint")
	ClusterMustHaveMember          = e.New("cluster must have a member!")
	ClusterMustHaveTwo             = e.New("cluster must have at least two members")
	IDAlreadyInUse                 = e.New("ID already in use")
	IllFormedCluster               = e.New("ill-formed cluster serialization")
	IllFormedClusterMember         = e.New("ill-formed cluster member serialization")
	IllFormedRegCred               = e.New("ill-formed regCred serialization")
	MissingClosingBrace            = e.New("missing closing brace")
	MissingClusterNameOrID         = e.New("missing cluster name or ID")
	MissingEndPointsSection        = e.New("missing endPoints section")
	MissingMembersList             = e.New("missing members list")
	MissingPrivateKey              = e.New("missing private key line")
	MissingRegNodeLine             = e.New("missing regNode line")
	MissingServerInfo              = e.New("missing server info")
	NameAlreadyInUse               = e.New("name already in use")
	NilCluster                     = e.New("nil cluster argument")
	NilNode                        = e.New("nil node argument")
	NilPrivateKey                  = e.New("nil private key argument")
	NilRegistry                    = e.New("nil registry argument")
	NilRegNode                     = e.New("nil RegNode argument")
	NilToken                       = e.New("nil XLRegMsg_Token argument")
	NoNodeNoKeys                   = e.New("no node and no keys to build one")
	RcvdInvalidMsgForState         = e.New("invalid msg type for current state")
	TagOutOfRange                  = e.New("message tag of of range")
	UnexpectedMsgType              = e.New("unexpected message type")
	UnknownClient                  = e.New("client unknown, not in registry")
	WrongNumberOfBytesInAttrs      = e.New("wrong number of bytes in attrs")
)
