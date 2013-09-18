package reg

import (
	"errors"
)

var (
	BadAttrsLine              = errors.New("badly formed attrs line")
	ClusterMustHaveTwo        = errors.New("cluster must have at least two members")
	IllFormedCluster          = errors.New("ill-formed cluster serialization")
	MissingClosingBrace       = errors.New("missing closing brace")
	MissingMembersList        = errors.New("missing members list")
	WrongNumberOfBytesInAttrs = errors.New("wrong number of bytes in attrs")
)
