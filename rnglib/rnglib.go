package rnglib

// import "fmt"
import "os"
import "math/rand"
import "strings"

func Version() (string, string) {
	return "0.1.0", "2013-04-15"
}

// a crude attempt at properties
var _FILE_NAME_STARTERS = strings.Split(
	"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_", "")

func FILE_NAME_STARTERS() []string {
	return _FILE_NAME_STARTERS
}

var _FILE_NAME_CHARS = strings.Split(
	"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-.", "")

func FILE_NAME_CHARS() []string {
	return _FILE_NAME_CHARS
}

type RNG interface {
	Seed(seed int64)
	NextBool()
	NextByte()
	NextBytes([]byte)
	NextInt32(uint32)
	NextInt64(uint64)
	NextFloat32()
	NextFloat64()
	NextFileName(int)
}

type SimpleRNG struct {
	rng *rand.Rand
}

func NewSimpleRNG(seed int64) *SimpleRNG {
	s := new(SimpleRNG) // allocates
	src := rand.NewSource(seed)
	s.rng = rand.New(src)
	s.Seed(seed)
	return s
}

func (s *SimpleRNG) Seed(seed int64) {
	s.rng.Seed(seed)
}
func (s *SimpleRNG) NextBoolean() bool {
	return s.rng.Int63()%2 == 0
}
func (s *SimpleRNG) NextByte() byte {
	return byte(s.rng.Int63n(256))
}
func (s *SimpleRNG) NextBytes(buffer *[]byte) {
	for n := range *buffer {
		(*buffer)[n] = s.NextByte()
	}
}
func (s *SimpleRNG) NextInt32(n uint32) uint32 {
	return uint32(float32(n) * s.rng.Float32())
}
func (s *SimpleRNG) NextInt64(n uint64) uint64 {
	return uint64(float64(n) * s.rng.Float64())
}
func (s *SimpleRNG) NextFloat32() float32 {
	return s.rng.Float32()
}
func (s *SimpleRNG) NextFloat64() float64 {
	return s.rng.Float64()
}

// These produce strings which are acceptable POSIX file names.
// and also advance a cursor by a multiple of 64 bits.  All strings
// are at least one byte and less than maxLen bytes n length.  We
// arbitrarily limit file names to less than 256 characters.

func (s *SimpleRNG) _nextFileName(nameLen int) string {
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

func (s *SimpleRNG) NextFileName(maxLen int) string {
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

func (s *SimpleRNG) NextDataFile(dirName string, maxLen int, minLen int) (int, string) {
	// silently convert paramaters to reasonable values
	if minLen < 0 {
		minLen = 0
	}
	if maxLen < minLen+1 {
		maxLen = minLen + 1
	}

	// create the data directory if it does not exist
	dirExists, err := PathExists(dirName)
	if err != nil {
		panic(err)
	}
	if !dirExists {
		os.MkdirAll(dirName, 0755)
	}

	// loop until the file does not exist
	pathToFile := dirName + "/" + s.NextFileName(16)
	pathExists, err := PathExists(pathToFile)
	if err != nil {
		panic(err)
	}
	for pathExists {
		pathToFile := dirName + "/" + s.NextFileName(16)
		pathExists, err = PathExists(pathToFile)
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
	return count, pathToFile
}

/* XXx this should be in another packge */
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// NextDataDir creates a directory tree populated with data files.
//
// BUGS
// * on at least one occasion with width = 4 only 3 files/directories
//   were created at the top level (2 were subdirs)
// DEFICIENCIES:
// * no control over percentage of directories
// * no guarantee that depth will be reached

func (s *SimpleRNG) NextDataDir(pathToDir string, depth int, width int,
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
	pathExists, err := PathExists(pathToDir)
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
