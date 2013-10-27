package u

// xlattice_go/u/const.go

const (
	SHA1_LEN = 40 // length of hex version
	SHA3_LEN = 64

	//           ....x....1....x....2....x....3....x....4
	SHA1_NONE = "0000000000000000000000000000000000000000"

	//          ....x....1....x....2....x....3....x....4....x....5....x....6....
	SHA3_NONE = "0000000000000000000000000000000000000000000000000000000000000000"
	DEFAULT_BUFFER_SIZE = 256 * 256
)

type DirStruc int
const (
	DIR_FLAT	DirStruc = iota
	DIR16x16
	DIR256x256
)

