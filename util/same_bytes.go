package util

// xlattice_go/util/same_bytes.go

import (
	"unsafe"
)

// Return true if the two byte slices have the same length and no
// byte differ.

func SameBytes(a, b []byte) bool {
	var aInt, bInt int
	if len(a) != len(b) {
		return false
	}
	length := len(a)
	if length == 0 {
		return true
	}
	sizeInt := int(unsafe.Sizeof(aInt))
	var i int
	for i := 0; i < length-sizeInt; i += sizeInt {
		aInt = *(*int)(unsafe.Pointer(&a[i]))
		bInt = *(*int)(unsafe.Pointer(&b[i]))
		if aInt != bInt {
			return false
		}
	}
	// check for leftover bytes
	if i > length {
		for i -= sizeInt; i < len(a); i++ {
			if a[i] != b[i] {
				return false
			}
		}
	}
	return true
}
