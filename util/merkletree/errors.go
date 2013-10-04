package merkletree

import (
	"errors"
)

var (
	DirectoryNotFound		= errors.New("directory not found")
	EmptyPath				= errors.New("empty path argument")
	NilTreeButNotBinding	= errors.New("nil tree but not binding")

)
