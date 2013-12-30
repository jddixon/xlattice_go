package nodeID

import (
	"bytes"
	"encoding/hex"
	"errors"
	xr "github.com/jddixon/xlattice_go/rnglib"
)

// TAKE CARE: these in bytes; hex values are twice these
const SHA1_LEN = 20
const SHA3_LEN = 32

// CONSTRUCTORS /////////////////////////////////////////////////////
type NodeID struct {
	_nodeID []byte
}

var (
	BadNodeIDLen = errors.New("bad length for nodeID")
	NilNodeID    = errors.New("nil byte array for nodeID")
)

func New(id []byte) (q *NodeID, err error) {
	q = new(NodeID)
	if id == nil {
		id = make([]byte, SHA3_LEN)
		rng := xr.MakeSystemRNG()
		rng.NextBytes(id)
		q._nodeID = id
	} else {
		// deep copy the slice
		size := len(id)
		myID := make([]byte, size)
		for i := 0; i < size; i++ {
			myID[i] = id[i]
		}
		q._nodeID = myID
	}
	if !IsValidID(id) {
		err = BadNodeIDLen
	}
	return
}

// XXX CONSIDER ME DEPRECATED
func NewNodeID(id []byte) (q *NodeID, err error) {
	return New(id)
}

// func NewNodeIDFromString(id string) *NodeID {
//     ...
// }

func (n *NodeID) Clone() (*NodeID, error) {
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

func SameNodeID(a, b *NodeID) (same bool) {
	if a == nil || b == nil {
		return false
	}
	aVal, bVal := a.Value(), b.Value()

	return bytes.Equal(aVal, bVal)

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
	//for i := 0; i < n.Length(); i++ {
	//	if (*n)._nodeID[i] != (*other)._nodeID[i] {
	//		return false
	//	}
	//}
	// return true
	return SameNodeID(n, other)
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
