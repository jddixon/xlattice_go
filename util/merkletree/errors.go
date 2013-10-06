package merkletree

import (
	"errors"
)

var (
	CantParseFirstLine	 = errors.New("can't parse first line")
	CantParseOtherLine	 = errors.New("can't parse other line")
	DirectoryNotFound    = errors.New("directory not found")
	EmptyName            = errors.New("empty name argument")
	EmptyPath            = errors.New("empty path argument")
	FileNotFound         = errors.New("file not found")
	InvalidHashLength	 = errors.New("invalid hash length")
	NilMerkleNode		 = errors.New("nil MerkleNode")
	NilTreeButNotBinding = errors.New("nil tree but not binding")
)
