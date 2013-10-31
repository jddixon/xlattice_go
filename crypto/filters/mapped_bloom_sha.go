package filters

import (
	"fmt" // DEBUG
	xf "github.com/jddixon/xlattice_go/util/lfs"
	"io/ioutil"
	gm "launchpad.net/gommap"
	"os"
	"reflect"
	"unsafe"
)

var _ = fmt.Print

type MappedBloomSHA struct {
	backingFile string
	f           *os.File
	inCore      gm.MMap
	BloomSHA
}

func NewMappedBloomSHA(m, k uint, backingFile string) (
	mb3 *MappedBloomSHA, err error) {

	var (
		f           *os.File
		filterBits  uint
		filterWords uint
		size        int64
		Filter      []uint64
		inCore      gm.MMap
		b3          *BloomSHA
	)
	if m < MIN_M || m > MAX_M {
		err = MOutOfRange
	} else {
		filterBits = uint(1) << m
		filterWords = filterBits / 64
		size = int64(filterBits / 8) // bytes
	}
	if err == nil && (k < MIN_K || (k*m > MAX_MK_PRODUCT)) {
		err = TooManyHashFunctions
	}
	if err == nil {
		var found bool
		found, err = xf.PathExists(backingFile)
		if err == nil {
			if found {
				f, err = os.OpenFile(backingFile,
					os.O_RDWR, 0640)
				if err == nil {
					var fi os.FileInfo
					fi, err = f.Stat()
					if err == nil {
						// XXX could be huge ...
						fileSize := int64(fi.Size())
						if fileSize < size {
							err = MappingFileTooSmall
						}
					}
				}
			} else {
				// ! found {
				f, err = os.OpenFile(backingFile,
					os.O_CREATE|os.O_RDWR, 0640)

				// XXX should write blocks in a loop
				zeroes := make([]byte, size)
				err = ioutil.WriteFile(backingFile, zeroes, 0640)
			}
		}
	}
	if err == nil {
		// 2013-11-01: was getting EPERM errors here because I used
		// os.Open to open the backing file.  Switched to OpenFile
		// with explicit permissions and everything worked.
		inCore, err = gm.MapAt(0, f.Fd(), 0, size,
			gm.PROT_READ|gm.PROT_WRITE,
			gm.MAP_SHARED)
	}
	if err == nil {
		ih := (*reflect.SliceHeader)(unsafe.Pointer(&inCore))
		fh := (*reflect.SliceHeader)(unsafe.Pointer(&Filter))
		fh.Data = ih.Data               // Filter slice points at same data
		fh.Len = ih.Len / SIZEOF_UINT64 // length suitably modified
		fh.Cap = ih.Cap / SIZEOF_UINT64 // likewise for capacity

		b3 = &BloomSHA{
			m:      m,
			k:      k,
			Filter: Filter,

			// comments say these are convenience variables but they
			// are actually used
			filterBits:  filterBits,
			filterWords: filterWords,
		}
		b3.doClear() // no lock
	}
	if err == nil {
		mb3 = &MappedBloomSHA{
			backingFile: backingFile,
			f:           f,
			inCore:      inCore,
			BloomSHA:    *b3,
		}
	}
	return
}
func (mb3 *MappedBloomSHA) Close() {
	if mb3.f != nil {

		mb3.inCore.Sync(gm.MS_SYNC)
		mb3.f.Close()
	}
}
