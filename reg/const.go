package reg

const (
	SHA1_LEN     = 20 // length in bytes of binary hash
	SHA3_LEN     = 32

	// The version MUST consist of three parts separated by dots,
	// with each part being one or two digits.  It is converted
	// into a uint32 in in_handler.go init()
	VERSION      = "0.1.0"
	VERSION_DATE = "2013-09-30"
)
