package node

import (
	"bytes"
	"encoding/hex"
	"errors"
)

// these SHOULD be in a crypto package
const SHA1_LEN = 20			// in bytes; hex SHA1_LEN is twice this
const SHA3_LEN = 32

// END SHOULD

// CONSTRUCTORS /////////////////////////////////////////////////////
type NodeID struct {
	_nodeID []byte
}

func NewNodeID(id []byte) *NodeID {
	if id == nil {
		panic("IllegalArgument: nil nodeID")
	}
	q := new(NodeID)
	// deep copy the slice
	size := len(id)
	myID := make([]byte, size)
	for i := 0; i < size; i++ {
		myID[i] = id[i]
	}
	q._nodeID = myID
	if !IsValidID(id) {
		panic("IllegalArgument: invalid id length")
	}
	return q
}

// func NewNodeIDFromString(id string) *NodeID {
//     ...
// }

func (n *NodeID) Clone() *NodeID {
	v := n.Value()
	return NewNodeID(v)
}

// OTHER METHODS ////////////////////////////////////////////////////
func (n *NodeID) Compare(any interface{}) (int, error) {
	result := 0
	err := error(nil)
	if any == nil {
		err = errors.New("IllegalArgument: nil comparand")
	} else if any == n {
		return result, err // defaults to 0, nil
	} else {
		switch v := any.(type) {
		case *NodeID:
			_ = v
		default:
			err = errors.New("IllegalArgument: not pointer to NodeID")
		}
	}
	if err != nil {
		return result, err
	}
	other := any.(*NodeID)
	if n.Length() != other.Length() {
		return 0, errors.New("IllegalArgument: NodeIDs of different length")
	}
	return bytes.Compare(n.Value(), other.Value()), nil
}

func (n *NodeID) Equal(any interface{}) bool {
	if any == n {
		return true
	}
	if any == nil {
		return false
	}
	switch v := any.(type) {
	case *NodeID:
		_ = v
	default:
		return false
	}
	other := any.(*NodeID) // type assertion
	if n.Length() != other.Length() {
		return false
	}
	for i := 0; i < n.Length(); i++ {
		if (*n)._nodeID[i] != (*other)._nodeID[i] {
			return false
		}
	}
	return true
}

func IsValidID(value []byte) bool {
	if value == nil {
		return false
	}
	// XXX check type?
	x := len(value)
	return x == 20 || x == 32 // SHA1 or SHA3
}

func (n *NodeID) Length() int {
	return len(n._nodeID)
}

// Returns a deep copy of the slice.
func (n *NodeID) Value() []byte {
	size := len(n._nodeID)
	v := make([]byte, size)
	for i := 0; i < size; i++ {
		v[i] = n._nodeID[i]
	}
	return v
}

// SERIALIZATION ////////////////////////////////////////////////////
func (n *NodeID) String() string {
	return hex.EncodeToString(n._nodeID)
}
