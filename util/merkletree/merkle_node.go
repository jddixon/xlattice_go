package merkletree

// xlattice_go/util/merkletree/merkle_node.go

import (
	//"code.google.com/p/go.crypto/sha3"
	//"crypto/sha1"
	//"encoding/hex"
	"fmt"
	xu "github.com/jddixon/xlattice_go/util"
	//"hash"
	//"os"
	//"path"
	//re "regexp"
	//"strings"
)

var _ = fmt.Print

type MerkleNodeI interface {
	Name() string
	GetHash() []byte
	SetHash([]byte) error
	UsingSHA1() bool
	IsLeaf() bool

	Equal(any interface{}) bool
	ToString(indent string) (string, error)
	ToStrings(indent string, ss *[]string) error
	// XXX DELAY THESE FOR A WHILE
	// GetPath()		        string
	// SetPath(value	string) error
}

type MerkleNode struct {
	name      string
	hash      []byte
	usingSHA1 bool
}

func NewMerkleNode(name string, hash []byte, usingSHA1 bool) (
	mn *MerkleNode, err error) {

	if name == "" {
		err = EmptyName
	}
	if err == nil {
		length := len(hash)
		if length != 0 && length != SHA1_LEN && length != SHA3_LEN {
			err = InvalidHashLength
		}
	}
	if err == nil {
		mn = &MerkleNode{
			name:      name,
			hash:      hash,
			usingSHA1: usingSHA1,
		}
	}
	return
}
func (mn *MerkleNode) Name() string {
	return mn.name
}

// XXX THIS IS A MAJOR CHANGE FROM THE PYTHON, where the hash is a
// hex value
func (mn *MerkleNode) GetHash() []byte {
	return mn.hash
}
func (mn *MerkleNode) SetHash(value []byte) (err error) {
	// XXX SOME VALIDATION NEEDED
	mn.hash = value
	return
}
func (mn *MerkleNode) UsingSHA1() bool {
	return mn.usingSHA1
}
func (mn *MerkleNode) Equal(any interface{}) bool {
	if any == mn {
		return true
	}
	if any == nil {
		return false
	}
	switch v := any.(type) {
	case *MerkleNode:
		_ = v
	default:
		return false
	}
	other := any.(*MerkleNode) // type assertion

	return mn.name == other.name && xu.SameBytes(mn.hash, other.hash) &&
		mn.usingSHA1 == other.usingSHA1
}
