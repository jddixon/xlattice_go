package filters

import (
)

const (
	MIN_M	= uint(2)
	MAX_M	= uint(24)		// XXX arguments for limit?
	MIN_K	= 1

	// ostensibly "too many hash functions for filter size"
	MAX_MK_PRODUCT	= 256

	SIZEOF_UINT64 = 8		// bytes
)
