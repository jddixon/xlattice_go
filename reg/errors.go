package reg

import (
	"errors"
)

var (
	BadAttrsLine		= errors.New("badly formed attrs line") 
	ClusterMustHaveTwo	= errors.New("cluster must have at least two members")
	MissingClosingBrace	= errors.New("missing closing brace")
	WrongNumberOfBytesInAttrs = errors.New("wrong number of bytes in attrs")
)
