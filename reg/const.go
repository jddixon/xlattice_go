package reg

// xlattice_go/reg/const.go

const (
	SHA1_LEN = 20 // length in bytes of binary hash
	SHA3_LEN = 32

	// The version MUST consist of three parts separated by dots,
	// with each part being one or two digits.  It is converted
	// into a uint32 in in_handler.go init()
	VERSION      = "0.1.1"
	VERSION_DATE = "2013-10-10"
)

// client attrs bits, also used for member attrs
const (
	ATTR_EPHEMERAL = 1 << iota
	ATTR_ADMIN
	ATTR_SOLO // no related cluster, persists config to LFS
)
