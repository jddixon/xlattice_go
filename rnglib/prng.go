// xlattice_go/prng.go
package rnglib

import (
	xf "github.com/jddixon/xlattice_go/util/lfs"
	"os"
	"math/rand"
	"strings"
)

type PRNG struct {
	rng *rand.Rand
}

//func NewMTSource(seed int64) rand.Source {
//	var mt64 MT64
//	mt64.Seed(seed)
//	return &mt64
//}
//
//func NewSimpleRNG(seed int64) *PRNG {
//	s := new(PRNG) // allocates
//	src := NewMTSource(seed)
//	s.rng = rand.New(src)
//	s.Seed(seed)
//	return s
//} // GEEP

func (s *PRNG) Seed(seed int64) {
	s.rng.Seed(seed)
}

// expose the rand.Random interface /////////////////////////////////
func (s *PRNG) Int63() int64         { return s.rng.Int63() }
func (s *PRNG) Uint32() uint32       { return s.rng.Uint32() }
func (s *PRNG) Int31() int32         { return s.rng.Int31() }
func (s *PRNG) Int() int             { return s.rng.Int() }
func (s *PRNG) Int63n(n int64) int64 { return s.rng.Int63n(n) }
func (s *PRNG) Int31n(n int32) int32 { return s.rng.Int31n(n) }
func (s *PRNG) Intn(n int) int       { return s.rng.Intn(n) }
func (s *PRNG) Float64() float64     { return s.rng.Float64() }
func (s *PRNG) Float32() float32     { return s.rng.Float32() }
func (s *PRNG) Perm(n int) []int     { return s.rng.Perm(n) }

// PRNG functions //////////////////////////////////////////////
func (s *PRNG) NextBoolean() bool { return s.rng.Int63()&1 == 0 }
func (s *PRNG) NextByte() byte    { return byte(s.rng.Int63n(256)) }

// miraculously inefficient
func (s *PRNG) NextBytes(buffer *[]byte) {
	for n := range *buffer {
		(*buffer)[n] = s.NextByte()
	}
}
func (s *PRNG) NextInt32(n uint32) uint32 { return uint32(float32(n) * s.rng.Float32()) }
func (s *PRNG) NextInt64(n uint64) uint64 { return uint64(float64(n) * s.rng.Float64()) }
func (s *PRNG) NextFloat32() float32      { return s.rng.Float32() }
func (s *PRNG) NextFloat64() float64      { return s.rng.Float64() }

// These produce strings which are acceptable POSIX file names.
// and also advance a cursor by a multiple of 64 bits.  All strings
// are at least one byte and less than maxLen bytes n length.  We
// arbitrarily limit file names to less than 256 characters.

func (s *PRNG) _nextFileName(nameLen int) string {
	/* always returns at least one character */
	maxStarterNdx := uint32(len(_FILE_NAME_STARTERS))
	ndx := s.NextInt32(maxStarterNdx)
	var chars []string
	chars = append(chars, _FILE_NAME_STARTERS[ndx])
	maxCharNdx := uint32(len(_FILE_NAME_CHARS))
	for n := 0; n < nameLen; n++ {
		ndx := s.NextInt32(maxCharNdx)
		chars = append(chars, _FILE_NAME_CHARS[ndx])
	}
	return strings.Join(chars, "")
}

func (s *PRNG) NextFileName(maxLen int) string {
	_maxLen := uint32(maxLen)
	if _maxLen < 2 {
		_maxLen = 2 // this is a ceiling which cannot be reached
	}
	if _maxLen > 256 {
		_maxLen = 256
	}
	var name string
	nameLen := uint32(0)
	for nameLen == 0 {
		nameLen = s.NextInt32(_maxLen) // so len < 256
	}
	for {
		name = s._nextFileName(int(nameLen))
		if (len(name) > 0) && !strings.Contains(name, "..") {
			break
		}
	}
	return name
}

// These are operations on the file system.  Directory depth is at least 1
// and no more than 'depth'.  Likewise for width, the number of
// files in a directory, where a file is either a data file or a subdirectory.
// The number of bytes in a file is at least minLen and less than maxLen.
// Subdirectory names may be random
//
// XXX CHANGED TO RETURN AN INT64 LENGTH

func (s *PRNG) NextDataFile(dirName string, maxLen int, minLen int) (int64, string) {
	// silently convert parameters to reasonable values
	if minLen < 0 {
		minLen = 0
	}
	if maxLen < minLen+1 {
		maxLen = minLen + 1
	}

	// create the data directory if it does not exist
	dirExists, err := xf.PathExists(dirName)
	if err != nil {
		panic(err)
	}
	if !dirExists {
		os.MkdirAll(dirName, 0755)
	}

	// loop until the file does not exist
	pathToFile := dirName + "/" + s.NextFileName(16)
	pathExists, err := xf.PathExists(pathToFile)
	if err != nil {
		panic(err)
	}
	for pathExists {
		pathToFile := dirName + "/" + s.NextFileName(16)
		pathExists, err = xf.PathExists(pathToFile)
		if err != nil {
			panic(err)
		}
	}
	count := minLen + int(s.NextFloat32()*float32((maxLen-minLen)))
	data := make([]byte, count)
	s.NextBytes(&data) // fill with random bytes
	// XXX may cause panic
	fo, err := os.Create(pathToFile)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}() // XXX wakaranai
	// XXX this should be chunked
	// XXX data should be slice
	if _, err := fo.Write(data); err != nil {
		panic(err)
	}
	// XXX respec to also return err
	return int64(count), pathToFile
}

// NextDataDir creates a directory tree populated with data files.
//
// BUGS
// * on at least one occasion with width = 4 only 3 files/directories
//   were created at the top level (2 were subdirs)
// DEFICIENCIES:
// * no control over percentage of directories
// * no guarantee that depth will be reached

func (s *PRNG) NextDataDir(pathToDir string, depth int, width int,
	maxLen int, minLen int) {
	// number of directory levels; 1 means no subdirectories
	if depth < 1 {
		depth = 1
	}
	// number of members (files, subdirectories) at each level
	if width < 1 {
		width = 1
	}
	// XXX may panic
	pathExists, err := xf.PathExists(pathToDir)
	if err != nil {
		panic(err)
	}
	if !pathExists {
		os.MkdirAll(pathToDir, 0755)
	}
	subdirSoFar := 0
	for i := 0; i < width; i++ {
		if depth > 1 {
			if (s.NextFloat32() > 0.25) &&
				((i < width-1) || (subdirSoFar > 0)) {
				// 25% are subdirs
				// data file i
				// SPECIFICATION ERROR: file name may not be unique
				// count, pathToFile
				s.NextDataFile(pathToDir,
					maxLen, minLen)
			} else {
				// directory
				subdirSoFar += 1
				// create unique name
				fileName := s.NextFileName(16)
				pathToSubdir := pathToDir + "/" + fileName
				s.NextDataDir(pathToSubdir, depth-1, width,
					maxLen, minLen)
			}
		} else {
			// data file
			// XXX SPECIFICATION ERROR: file name may not be unique
			// count, pathToFile
			s.NextDataFile(pathToDir, maxLen, minLen)
		}
	} // end for
}
