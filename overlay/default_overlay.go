package overlay

import (
	"errors"
	"fmt"
	xt "github.com/jddixon/xlattice_go/transport"
)

var _ = fmt.Print

func DefaultOverlay(e xt.EndPointI) (o OverlayI, err error) {
	t := e.Transport()
	switch t {
	case "ip":
		fallthrough
	case "udp":
		fallthrough
	case "tcp":
		t = "ip"
	default:
		return nil, errors.New("not implemented")
	}

	// KLUDGE
	tcpE := e.(*xt.TcpEndPoint)
	tcpA := tcpE.GetTcpAddr() // IP, Port, Zone fields
	v4Addr := tcpA.IP[12:]

	// 127/8 or 10/8 or 172.16/12 or 192.168/16
	if v4Addr[0] == 127 {
		aRange, err := NewCIDRAddrRange("127.0.0.0/8")
		if err == nil {
			o, err = NewIPOverlay("localhost", aRange, "ip", 1.0)
		}
	} else if v4Addr[0] == 10 {
		aRange, err := NewCIDRAddrRange("10.0.0.0/8")
		if err == nil {
			o, err = NewIPOverlay("privateA", aRange, "ip", 1.0)
		}
	} else if v4Addr[0] == 172 && v4Addr[1] >= 16 && v4Addr[1] < 32 {
		aRange, err := NewCIDRAddrRange("172.16.0.0/12")
		if err == nil {
			o, err = NewIPOverlay("privateB", aRange, "ip", 1.0)
		}
	} else if v4Addr[0] == 192 && v4Addr[1] == 168 {
		aRange, err := NewCIDRAddrRange("192.168.0.0/16")
		if err == nil {
			o, err = NewIPOverlay("privateC", aRange, "ip", 1.0)
		}
	} else {
		aRange, err := NewCIDRAddrRange("0.0.0.0/0")
		if err == nil {
			o, err = NewIPOverlay("globalV4", aRange, "ip", 1.0)
		}
	}
	return
}
