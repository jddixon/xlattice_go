package merkletree

// xlattice_go/util/merkletree/merkletree.go

import (
	"code.google.com/p/go.crypto/sha3"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	//xu "github.com/jddixon/xlattice_go/util"
	//"hash"
	"io/ioutil"
	//"os"
	//"path"
	// re "regexp"
	//"strings"
)

var _ = fmt.Print

type MerkleLeaf struct {
	MerkleNode
}

// Creates a MerkleTree leaf node.

func NewMerkleLeaf(name string, hash []byte, usingSHA1 bool) (
	ml *MerkleLeaf, err error) {

	mn, err := NewMerkleNode(name, hash, usingSHA1)
	if err == nil {
		ml = &MerkleLeaf{*mn}
	}
	return
}

// Create a MerkleTree leaf node corresponding to a file in the file
// system.  To simplify programming, the base name of the file, which is
// part of the path, is also passed as a separate argument.

func CreateMerkleLeafFromFileSystem(pathToFile, name string, usingSHA1 bool) (
	ml *MerkleLeaf, err error) {

	var hash []byte
	if usingSHA1 {
		hash, err = SHA1File(pathToFile)
	} else {
		hash, err = SHA3File(pathToFile)
	}
	if err == nil {
		ml, err = NewMerkleLeaf(name, hash, usingSHA1)
	}
	return
}
func (ml *MerkleLeaf) IsLeaf() bool {
	return true
}
func (ml *MerkleLeaf) Equal(any interface{}) bool {
	if any == ml {
		return true
	}
	if any == nil {
		return false
	}
	switch v := any.(type) {
	case *MerkleLeaf:
		_ = v
	default:
		return false
	}
	other := any.(*MerkleLeaf) // type assertion
	return ml.MerkleNode.Equal(other.MerkleNode)
}

// Serialize the leaf node, prefixing it with 'indent', which should
// conventionally be a number of spaces.

func (ml *MerkleLeaf) ToString(indent string) string {
	var shash string
	hash := ml.hash
	if len(hash) == 0 {
		if ml.usingSHA1 {
			shash = SHA1_NONE
		} else {
			shash = SHA3_NONE
		}
	} else {
		shash = hex.EncodeToString(hash)
	}
	return fmt.Sprintf("%s%s %s\r\n", indent, shash, ml.name)
}

// Return the SHA1 hash of a file.  This is a sequence of 20 bytes.

func SHA1File(pathToFile string) (hash []byte, err error) {
	var data []byte
	found, err := PathExists(pathToFile)
	if err == nil && !found {
		err = FileNotFound
	}
	if err == nil {
		data, err = ioutil.ReadFile(pathToFile)
		if err == nil {
			digest := sha1.New()
			digest.Write(data)
			hash = digest.Sum(nil)
		}
	}
	return
}

// Return the SHA3-256 hash of a file.  This is a sequence of 32 bytes.

func SHA3File(pathToFile string) (hash []byte, err error) {
	var data []byte
	found, err := PathExists(pathToFile)
	if err == nil && !found {
		err = FileNotFound
	}
	if err == nil {
		data, err = ioutil.ReadFile(pathToFile)
		if err == nil {
			digest := sha3.NewKeccak256()
			digest.Write(data)
			hash = digest.Sum(nil)
		}
	}
	return
}
