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
	V4_ADDR_PAT     = `^` + QUAD_PAT + `\.` + QUAD_PAT + `\.` + QUAD_PAT + `\.` + QUAD_PAT + `$`
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
	Address
}

// Expect an IPV4 address in the form A.B.C.D:P, where P is the
// port number and the :P is optional.
func NewV4Address(val string) (addr *V4Address, err error) {
	if v4AddrRE == nil {
		if err = makeRE(); err != nil {
			panic(err)
		}
	}
	var validAddr bool // false by default
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
			}
		}
		if err == nil {
			validAddr = v4AddrRE.MatchString(parts[0])
		}
	}
	if validAddr {
		addr = &V4Address{Address{val}}
	}
	// DEBUG
	if !validAddr {
		fmt.Printf("    validAddr flag NOT set for %s\n", val)
	}
	// END
	return
}
