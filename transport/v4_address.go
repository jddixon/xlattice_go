package transport

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var _ = fmt.Print

const (
	QUAD_PAT = `\d|\d\d|(?:[01]\d\d|2(?:[01234]\d|5[0-5]))`
	// XXX If you add the dollar sign, the match fails
	V4_ADDR_PAT = `^` + QUAD_PAT + `\.` + QUAD_PAT + `\.` + QUAD_PAT + `\.` + QUAD_PAT + `$`

	// 2013-07-18 results using this pattern appear identical
	V4_ADDR_PAT2 = "^\\d|\\d\\d|(?:[01]\\d\\d|2(?:[01234]\\d|5[0-5]))\\.\\d|\\d\\d|(?:[01]\\d\\d|2(?:[01234]\\d|5[0-5]))\\.\\d|\\d\\d|(?:[01]\\d\\d|2(?:[01234]\\d|5[0-5]))\\.\\d|\\d\\d|(?:[01]\\d\\d|2(?:[01234]\\d|5[0-5]))$"

	bad_port_number = "not a valid IPv4 port number: "
	bad_ipv4_addr   = "not a valid IPv4 address: "
)

var v4AddrRE *regexp.Regexp

func makeRE() (err error) {
	v4AddrRE, err = regexp.Compile(V4_ADDR_PAT)
	return err
}

// An IPv4 address
type V4Address struct {
	host	string
	port	string		// if it's an int, the default is zero
}

// Expect an IPV4 address in the form A.B.C.D:P, where P is the
// port number and the :P is optional.
//
// NOTE that in Go usage ":8080" is a valid address, with an
// implicit "127.0.0.1" host part.

func NewV4Address(val string) (addr *V4Address, err error) {
	if v4AddrRE == nil {
		if err = makeRE(); err != nil {
			panic(err)
		}
	}
	var validAddr bool // false by default
	var portPart	string

	parts := strings.Split(val, `:`)
	partsCount := len(parts)
	if partsCount == 0 || partsCount > 2 {
		err = errors.New(bad_ipv4_addr + val)
	} else if partsCount == 1 {
		// no colon
		validAddr = v4AddrRE.MatchString(val)
	} else {
		// we have a colon
		var port int
		if port, err = strconv.Atoi(parts[1]); err == nil {
			if port >= 256*256 {
				err = errors.New(bad_port_number + parts[1])
			} else {
				portPart = parts[1]
			}
		}
		if err == nil {
			validAddr = v4AddrRE.MatchString(parts[0])
		}
	}
	if validAddr {
		addr = &V4Address{parts[0], portPart}
	}
	// DEBUG
	if !validAddr {
		fmt.Printf("    validAddr flag NOT set for %s\n", val)
	}
	// END
	return
}
func (a *V4Address) Clone() (AddressI, error) {
	return NewV4Address(a.String())			// .(AddressI)
}
func (a *V4Address) Equal(any interface{}) bool {
	if any == nil {
		return false
	}
	if any == a {
		return true
	}
	switch v := any.(type) {
	case *V4Address:
		_ = v
	default:
		return false
	}
	other := any.(*V4Address)
	return a.host == other.host && a.port == other.port
}
func (a *V4Address) String() string {
	if a.port == "" {
		return a.host
	} else {
		return a.host + ":" + a.port
	}
}
