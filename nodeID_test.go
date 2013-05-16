package xlattice_go

//import "fmt"
import "github.com/bmizerany/assert"
import . "github.com/jddixon/xlattice_go/rnglib"
import "testing"
import "time"

func makeRNG() *SimpleRNG {
	t := time.Now().Unix()
	rng := NewSimpleRNG(t)
	return rng
}
func TestBadNodeIDs(t *testing.T) {
	assert.Equal(t, false, IsValidID(nil))
	candidate := make([]byte, SHA1_LEN-1)
	assert.Equal(t, false, IsValidID(candidate))
	candidate = make([]byte, SHA1_LEN)
	assert.Equal(t, true, IsValidID(candidate))
	candidate = make([]byte, SHA1_LEN+1)
	assert.Equal(t, false, IsValidID(candidate))
	candidate = make([]byte, SHA3_LEN)
	assert.Equal(t, true, IsValidID(candidate))
}
func TestThisAndThat(t *testing.T) {
	rng := makeRNG()
	v1 := make([]byte, SHA1_LEN)
	rng.NextBytes(&v1)
	v2 := make([]byte, SHA1_LEN)
	rng.NextBytes(&v2)
	id1 := NewNodeID(v1)
	id2 := NewNodeID(v2)
	// XXX this should be assert.False(id1.Equal(id2))
	assert.NotEqual(t, id1, id2)

	v1a := id1.Value()
	v2a := id2.Value()

	// not identical XXX test doesn't work because assert package
	//                   can't handle byte arrays
	// assert.NotEqual(t, v1, v1a)
	// assert.NotEqual(t, v2, v2a)

	// XXX assert can't handle these tests either
	// assert.NotEqual(t, &v1, &v1a)
	// assert.NotEqual(t, &v2, &v2a)

	assert.Equal(t, SHA1_LEN, len(v1a))
	assert.Equal(t, SHA1_LEN, len(v2a))
	for i := 0; i < SHA1_LEN; i++ {
		assert.Equal(t, v1[i], v1a[i])
		assert.Equal(t, v2[i], v2a[i])
	}
	assert.Equal(t, false, id1.Equal(nil))
	assert.Equal(t, true, id1.Equal(id1))
	assert.Equal(t, false, id1.Equal(id2))
}

// func TestComparator(t *testing.T) {
//    NOT IMPLEMENTED because nodeID comparator not yet implemented
// }
