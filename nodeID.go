package xlattice_go

// these SHOULD be in a crypto package
const SHA1_LEN = 20
const SHA3_LEN = 32

// END SHOULD

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
	if !q.IsValid() {
		panic("IllegalArgument: invalid id length")
	}
	return q
}

// func NewNodeIDFromString(id string) *NodeID {
//     ...
// }

func (n *NodeID) Length() int {
	return len(n._nodeID)
}

func (n *NodeID) Value() []byte {
	// deep copy the slice
	size := len(n._nodeID)
	v := make([]byte, size)
	for i := 0; i < size; i++ {
		v[i] = n._nodeID[i]
	}
	return v
}
func (n *NodeID) Clone() *NodeID {
	v := n.Value()
	return NewNodeID(v)
}

// XXX DEPRECATED
func (n *NodeID) IsValid() bool {
	x := n.Length()
	return x == 20 || x == 32 // SHA1 or SHA3
}
func IsValidID(value []byte) bool {
	if value == nil {
		return false
	}
	// XXX check type?
	x := len(value)
	return x == 20 || x == 32 // SHA1 or SHA3
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

// func (n *NodeID) ToString() string {
//
