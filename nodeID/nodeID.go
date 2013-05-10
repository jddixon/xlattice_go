package nodeID

type NodeID struct {
    nodeID []byte
}

func NewNodeID (id []byte) *NodeID {
    if id == nil {
        panic( "IllegalArgument: nil nodeID")
    }
    q := new(NodeID)
    q.nodeID = id
    if !q.IsValid() {
        panic ("IllegalArgument: invalid id length")
    }
    return q
}

// func NewNodeIDFromString(id string) *NodeID {
//     ...
// }

func (n *NodeID) Length() int {
    return len(n.nodeID)
}

func (n *NodeID) Value() []byte {
    v := n.nodeID
    // XXX need assurance that this is a copy
    return v
}
func (n *NodeID) Clone() *NodeID {
    v := n.Value()
    return NewNodeID(v)
}
func (n *NodeID) IsValid() bool {
    x := n.Length()
    return x == 20  || x == 32      // SHA1 or SHA3
}
// func (n *NodeID) Equal (any interface{}) bool {
// }

// func (n *NodeID) ToString() string {
// 
