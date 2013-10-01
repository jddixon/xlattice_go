package overlay

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	//"strings"
	xt "github.com/jddixon/xlattice_go/transport"
	xu "github.com/jddixon/xlattice_go/util"
)

var _ = fmt.Print

const (
	NAME       = `[` + xu.NAME_STARTERS + `][` + xu.NAME_CHARS + `]*`
	QUAD       = xt.QUAD_PAT
	ADDR_RANGE = QUAD + `.` + QUAD + `.` + QUAD + `.` + QUAD + `/\d+`
	IP_OVERLAY = `^overlay:\s*(` + NAME + `),\s*(` + NAME + `),\s*(` + ADDR_RANGE + `),\s*(\d+\.\d*)$`
)

var (
	overlayRE *regexp.Regexp

	BadOverlayRE      = errors.New("bad overlay regexp")
	NotAnOverlay      = errors.New("not an overlay")
	NotAnAddressRange = errors.New("not a valid address range")
)

func makeRE() (err error) {
	overlayRE, err = regexp.Compile(IP_OVERLAY)
	return
}

// Parse a string like
//   overlay: localHost, tcp, 127.0.0.1/55, 1.5
// where the parameters are name, transport, address range, and cost

func Parse(s string) (o OverlayI, err error) {
	if overlayRE == nil {
		err = makeRE()
	}
	if err == nil {
		groups := overlayRE.FindStringSubmatch(s)
		if groups == nil {
			// fmt.Printf("NO MATCH for '%s'\n", s)
			err = NotAnOverlay
		} else {
			// fmt.Println("MATCH")
			name := groups[1]
			transport := groups[2]
			aRange := groups[3]
			cost := groups[4] // a string

			// DEBUG
			// fmt.Printf("name:      %s\n", name)
			// fmt.Printf("transport: %s\n", transport)
			// fmt.Printf("aRange:    %s\n", aRange)
			// fmt.Printf("cost:      %s\n", cost)
			// END

			var ar *AddrRange
			var rCost float32
			if transport != "ip" && transport != "tcp" && transport != "udp" {
				err = errors.New("unknown transport " + transport)
			} else {
				ar, err = NewCIDRAddrRange(aRange)
			}
			if err == nil {
				var _rCost float64
				_rCost, err = strconv.ParseFloat(cost, 32)
				rCost = float32(_rCost)
			}
			if err == nil {
				o, err = NewIPOverlay(name, ar, transport, rCost)
			}
		}
	}
	return
}
