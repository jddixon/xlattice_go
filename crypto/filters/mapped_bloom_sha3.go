package filters

import (
	"fmt" // DEBUG
	xf "github.com/jddixon/xlattice_go/util/lfs"
	gm "launchpad.net/gommap"
	"os"
	"reflect"
	"unsafe"
)

type MappedBloomSHA3 struct {
	backingFile string
	f           *os.File
	BloomSHA3
}

func NewMappedBloomSHA3(m, k uint, backingFile string) (
	mb3 *MappedBloomSHA3, err error) {

	var (
		f           *os.File
		filterBits  uint
		filterWords uint
		Filter      []uint64
		inCore      []byte
		b3          *BloomSHA3
	)
	if m < MIN_M || m > MAX_M {
		err = MOutOfRange
	}
	if err == nil && (k < MIN_K || (k*m > MAX_MK_PRODUCT)) {
		err = TooManyHashFunctions
	}
	if err == nil {
		var found bool
		found, err = xf.PathExists(backingFile)
		if err == nil {
			if found {
				f, err = os.Open(backingFile)
				if err == nil {
					var fi os.FileInfo
					fi, err = f.Stat()
					if err == nil {
						// XXX could be huge ...
						size := uint(fi.Size())
						if size < uint(1)<<m {
							err = MappingFileTooSmall
						}
					}
				}
			} else {
				// ! found {
				f, err = os.OpenFile(backingFile,
					os.O_CREATE|os.O_RDWR, 0640)
			}
		}
	}
	if err == nil {
		filterBits = uint(1) << m
		filterWords = filterBits / 64
		size := int64(filterBits / 8) // bytes
		// DEBUG
		fmt.Printf("m = %d, mapping %d (0x%0x) bytes as %d words\n",
			m, size, size, filterWords)
		// END
		inCore, err = gm.MapAt(0, f.Fd(), 0, size,
			gm.PROT_READ|gm.PROT_WRITE, gm.MAP_SHARED)
	}
	if err == nil {
		ih := (*reflect.SliceHeader)(unsafe.Pointer(&inCore))
		fh := (*reflect.SliceHeader)(unsafe.Pointer(&Filter))
		fh.Data = ih.Data               // Filter slice points at same data
		fh.Len = ih.Len / SIZEOF_UINT64 // length suitably modified
		fh.Cap = ih.Cap / SIZEOF_UINT64 // likewise for capacity

		// DEBUG
		fmt.Printf("inCore: addr %v, len %d, cap %d\n", ih.Data, ih.Len, ih.Cap)
		fmt.Printf("Filter: addr %v, len %d, cap %d\n", fh.Data, fh.Len, fh.Cap)
		// END
		var ks *KeySelector

		b3 = &BloomSHA3{
			m:          m,
			k:          k,
			Filter:     Filter,
			wordOffset: make([]uint, k),
			bitOffset:  make([]byte, k),

			// comments say these are convenience variables but they
			// are actually used
			filterBits:  filterBits,
			filterWords: filterWords,
		}
		// DEBUG
		fmt.Printf("MappedBloom.doClear; %d filterWords\n", filterWords)
		// END
		b3.doClear() // no lock
		// offsets into the filter
		ks, err = NewKeySelector(m, k, b3.bitOffset, b3.wordOffset)
		if err == nil {
			b3.ks = ks
		} else {
			b3 = nil
		}
	}
	if err == nil {
		mb3 = &MappedBloomSHA3{
			backingFile: backingFile,
			f:           f,
			BloomSHA3:   *b3,
		}
	}
	return
}
