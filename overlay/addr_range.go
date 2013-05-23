package overlay

// xlattice_go/overlay/addr_range.go

import (
	"errors"
)

// An address range as the term is used in XLattice and on the global
// internet.  The range is defined by a prefix, the number of significant
// bits in that prefix, and the number of bits in a valid address.  So
// an ipv4 10/8 address range would be represented as [10], 8, 32.
type AddrRange struct {
	Prefix    []byte // all addresses in the range match this prefix
	PrefixLen uint   // number of leading bits, the 8 in 10/8
	AddrLen   uint   // in bits: 32 for ipv4, 64 for ipv6
}

func NewAddrRange(prefix []byte, prefixLen uint, addrLen uint) (*AddrRange, error) {
	if prefix == nil {
		return nil, errors.New("IllegalArgument: nil prefix")
	}
	byteLen := uint(len(prefix))
	if prefixLen > 8*byteLen {
		return nil, errors.New("IllegalArgument: too few bits in prefix")
	}
	if prefixLen > addrLen {
		return nil, errors.New("IllegalArgument: prefix too long for addr len")
	}

	return &AddrRange{prefix, prefixLen, addrLen}, nil

}
