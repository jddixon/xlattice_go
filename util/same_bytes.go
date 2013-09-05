package util

// xlattice_go/util/same_bytes.go

import (
	"unsafe"
)

func SameBytes(a, b []byte) bool {
	var aInt, bInt int
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	length := len(a)
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
