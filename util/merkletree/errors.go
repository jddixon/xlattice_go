package merkletree

import (
	"errors"
)

var (
	DirectoryNotFound    = errors.New("directory not found")
	EmptyPath            = errors.New("empty path argument")
	FileNotFound         = errors.New("file not found")
	NilTreeButNotBinding = errors.New("nil tree but not binding")
)
