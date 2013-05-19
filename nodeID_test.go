package xlattice_go

import "github.com/bmizerany/assert"
import "testing"

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
	rng := MakeRNG()
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

func TestComparator(t *testing.T) {
	rng := MakeRNG()
	v1 := make([]byte, SHA1_LEN)
	rng.NextBytes(&v1)
	v3 := make([]byte, SHA3_LEN)
	rng.NextBytes(&v3)
	id1 := NewNodeID(v1)				// SHA1
	id3 := NewNodeID(v3)				// SHA3

	// make clones which will sort before and after v1
	v1a := make([]byte, SHA1_LEN)			// sorts AFTER v1
	for i := 0; i < SHA1_LEN; i++ {
		if v1[i] == 255 {
			v1a[i] = 255
		} else {
			v1a[i] = v1[i] + 1
			break
		}
	}
	v1b := make([]byte, SHA1_LEN)			// sorts BEFORE v1
	for i := 0; i < SHA1_LEN; i++ {
		if v1[i] == 0 {
			v1b[i] = 0
		} else {
			v1b[i] = v1[i] - 1
			break
		}
	}
	id1a := NewNodeID(v1a)
	id1b := NewNodeID(v1b)

	result, err := id1.Compare(id1)				// self
	assert.Equal(t, 0, result)
	assert.Equal(t, err, nil)

	result, err = id1.Compare(id1.Clone())		// identical copy
	assert.Equal(t, 0, result)
	assert.Equal(t, err, nil)

	result, err = id1.Compare(nil)		// nil comparand
	assert.NotEqual(t, err, nil)

	result, err = id1.Compare(id3)
	assert.NotEqual(t, err, nil)		// different lengths

	result, err = id1.Compare(id1a)
	assert.Equal(t, -1, result)			// id1 < id1a
	assert.Equal(t, err, nil)

	result, err = id1.Compare(id1b)		// id1 > id1b
	assert.Equal(t, 1, result)
	assert.Equal(t, err, nil)

	result, err = id1a.Compare(id1b)	// id1a > id1b
	assert.Equal(t, 1, result)
	assert.Equal(t, err, nil)

}
