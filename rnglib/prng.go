// xlattice_go/prng.go
package rnglib

import (
	xf "github.com/jddixon/xlattice_go/util/lfs"
	"math/rand"
	"os"
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
//	p := new(PRNG) // allocates
//	src := NewMTSource(seed)
//	p.rng = rand.New(src)
//	p.Seed(seed)
//	return s
//}

func (p *PRNG) Seed(seed int64) {
	p.rng.Seed(seed)
}

// expose the rand.Random interface
func (p *PRNG) Int63() int64         { return p.rng.Int63() }
func (p *PRNG) Uint32() uint32       { return p.rng.Uint32() }
func (p *PRNG) Int31() int32         { return p.rng.Int31() }
func (p *PRNG) Int() int             { return p.rng.Int() }
func (p *PRNG) Int63n(n int64) int64 { return p.rng.Int63n(n) }
func (p *PRNG) Int31n(n int32) int32 { return p.rng.Int31n(n) }
func (p *PRNG) Intn(n int) int       { return p.rng.Intn(n) }
func (p *PRNG) Float64() float64     { return p.rng.Float64() }
func (p *PRNG) Float32() float32     { return p.rng.Float32() }
func (p *PRNG) Perm(n int) []int     { return p.rng.Perm(n) }

// PRNG functions
func (p *PRNG) NextBoolean() bool { return p.rng.Int63()&1 == 0 }
func (p *PRNG) NextByte() byte    { return byte(p.rng.Int63n(256)) }

// miraculously inefficient
func (p *PRNG) NextBytes(buffer []byte) {
	for n := range buffer {
		buffer[n] = p.NextByte()
	}
}
func (p *PRNG) NextInt32(n uint32) uint32 { return uint32(float32(n) * p.rng.Float32()) }
func (p *PRNG) NextInt64(n uint64) uint64 { return uint64(float64(n) * p.rng.Float64()) }
func (p *PRNG) NextFloat32() float32      { return p.rng.Float32() }
func (p *PRNG) NextFloat64() float64      { return p.rng.Float64() }

// These produce strings which are acceptable POSIX file namep.
// and also advance a cursor by a multiple of 64 bitp.  All strings
// are at least one byte and less than maxLen bytes n length.  We
// arbitrarily limit file names to less than 256 characterp.
func (p *PRNG) _nextFileName(nameLen int) string {
	/* always returns at least one character */
	maxStarterNdx := uint32(len(_FILE_NAME_STARTERS))
	ndx := p.NextInt32(maxStarterNdx)
	var chars []string
	chars = append(chars, _FILE_NAME_STARTERS[ndx])
	maxCharNdx := uint32(len(_FILE_NAME_CHARS))
	for n := 0; n < nameLen; n++ {
		ndx := p.NextInt32(maxCharNdx)
		chars = append(chars, _FILE_NAME_CHARS[ndx])
	}
	return strings.Join(chars, "")
}

func (p *PRNG) NextFileName(maxLen int) string {
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
		nameLen = p.NextInt32(_maxLen) // so len < 256
	}
	for {
		name = p._nextFileName(int(nameLen))
		if (len(name) > 0) && !strings.Contains(name, "..") {
			break
		}
	}
	return name
}

// Return a string which is an acceptable POSIX path with at least
// minParts parts and fewer than maxParts parts.  If the parameters
// supplied are unreasonable, they are adjusted.
//
// If uniqueness or some other constraint needs to be met, call this
// function then test compliance, then as necessary repeat.
func (p *PRNG) NextPosixPath(maxPartLen, minParts, maxParts int) (pp string) {

	if maxParts < 2 {
		maxParts = 2 // so at least one
	}
	if minParts <= 0 || maxParts <= minParts {
		minParts = maxParts - 1
	}
	count := minParts + p.Intn(maxParts-minParts)
	var parts []string
	for i := 0; i < count; i++ {
		parts = append(parts, p.NextFileName(maxPartLen))
	}
	return strings.Join(parts, "/")
}

// Return a string which is a valid fully qualified domain name.  Only
// names in .com, .net, and .org are returned by this version.  The
// parameter is a suggestion as to the complexity of the result, with
// zero returning the simplest response and values greater than zero
// possibly returning a more complex result.
func (p *PRNG) NextFQDN(complexity int) string {

	if complexity == 0 {
		complexity = 0
	} else if complexity > 3 {
		complexity = 3
	}
	top := p.Intn(3)
	var topLevel string
	switch top {
	case 0:
		topLevel = "com"
	case 1:
		topLevel = "net"
	default:
		topLevel = "org"
	}
	var parts []string
	count := 1 + complexity
	for i := 0; i < count; i++ {
		parts = append(parts, p.NextFileName(8))
	}
	parts = append(parts, topLevel)
	return strings.Join(parts, ".")
}

// Return a string which is a well-formed email address.  Complexity
// is a suggestion as to how complex the result should be.
func (p *PRNG) NextEmailAddress(complexity int) (e string) {
	if complexity == 0 {
		complexity = 0
	} else if complexity > 3 {
		complexity = 3
	}
	fqdn := p.NextFQDN(complexity)
	name := p.NextFileName(5 + complexity)
	return name + "@" + fqdn
}

// OPERATIONS ON THE FILE SYSTEM ------------------------------------

// These are operations on the file system.  Directory depth is at least 1
// and no more than 'depth'.  Likewise for width, the number of
// files in a directory, where a file is either a data file or a subdirectory.
// The number of bytes in a file is at least minLen and less than maxLen.
// Subdirectory names may be random
//
// XXX Changed to return an int64 length.
//
func (p *PRNG) NextDataFile(dirName string, maxLen int, minLen int) (
	int64, string) {

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
	pathToFile := dirName + "/" + p.NextFileName(16)
	pathExists, err := xf.PathExists(pathToFile)
	if err != nil {
		panic(err)
	}
	for pathExists {
		pathToFile := dirName + "/" + p.NextFileName(16)
		pathExists, err = xf.PathExists(pathToFile)
		if err != nil {
			panic(err)
		}
	}
	count := minLen + int(p.NextFloat32()*float32((maxLen-minLen)))
	data := make([]byte, count)
	p.NextBytes(data) // fill with random bytes
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

// NextDataDir creates a directory tree populated with data filep.
//
// BUGS
// * on at least one occasion with width = 4 only 3 files/directories
//   were created at the top level (2 were subdirs)
// DEFICIENCIES:
// * no control over percentage of directories
// * no guarantee that depth will be reached
//
func (p *PRNG) NextDataDir(pathToDir string, depth int, width int,
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
			if (p.NextFloat32() > 0.25) &&
				((i < width-1) || (subdirSoFar > 0)) {
				// 25% are subdirs
				// data file i
				// SPECIFICATION ERROR: file name may not be unique
				// count, pathToFile
				p.NextDataFile(pathToDir,
					maxLen, minLen)
			} else {
				// directory
				subdirSoFar += 1
				// create unique name
				fileName := p.NextFileName(16)
				pathToSubdir := pathToDir + "/" + fileName
				p.NextDataDir(pathToSubdir, depth-1, width,
					maxLen, minLen)
			}
		} else {
			// data file
			// XXX SPECIFICATION ERROR: file name may not be unique
			// count, pathToFile
			p.NextDataFile(pathToDir, maxLen, minLen)
		}
	}
}
