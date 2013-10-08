package merkletree

import (
	"errors"
)

var (
	CantParseFirstLine   = errors.New("can't parse first line")
	CantParseOtherLine   = errors.New("can't parse other line")
	DirectoryNotFound    = errors.New("directory not found")
	EmptyName            = errors.New("empty name argument")
	EmptyPath            = errors.New("empty path argument")
	EmptySerialization   = errors.New("empty serialization")
	FileNotFound         = errors.New("file not found")
	InitialIndent        = errors.New("indented first line")
	InvalidHashLength    = errors.New("invalid hash length")
	NilMerkleNode        = errors.New("nil MerkleNode")
	NilNode              = errors.New("nil node argument")
	NilTreeButNotBinding = errors.New("nil tree but not binding")
)
