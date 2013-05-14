package xlattice_go

import "github.com/bmizerany/assert"
import . "github.com/jddixon/xlattice_go/rnglib"
import "testing"

func makeNodeID(rng *SimpleRNG) *NodeID {
    var buffer []byte
    if rng.NextBoolean() {
        buffer = make([]byte, SHA1_LEN)
    } else {
        buffer = make([]byte, SHA3_LEN)
    }
    rng.NextBytes(&buffer)
    return NewNodeID(&buffer)
}

func doKeyTests(t *testing.T, node *Node, rng *SimpleRNG) {
    // XXX STUB

}
func TestNewNewCtor(t *testing.T) {
    rng := makeRNG()
    _, err := NewNewNode(nil)
    assert.NotEqual(t, nil, err)

    id := makeNodeID(rng)
    assert.NotEqual(t, nil, id)
    n, err2 := NewNewNode(id)
    assert.NotEqual(t, nil, n)
    assert.Equal(t, nil, err2)
    actualID := n.GetNodeID()
    assert.Equal(t, true, id.Equal(actualID) )
    doKeyTests(t, n, rng) 
    assert.Equal(t, 0, n.SizePeers())
    assert.Equal(t, 0, n.SizeOverlays())
    assert.Equal(t, 0, n.SizeConnections())
}

func TestNewCtor(t *testing.T) {
    // rng := makeRNG()
    
    // if  constructor assigns a nil NodeID, we should get an
    // IllegalArgument panic
    // XXX STUB

    // if assigned a nil key, the NewNode constructor should panic 
    // with an IllegalArgument string
    // XXX STUB


}

