package msg

func SameBytes(a, b []byte) bool {
	if a == nil || b == nil {
		return false
	} 
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
} 
