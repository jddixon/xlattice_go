package transport

// xlattice_go/transport/endpoint.go

import (
	"strings"
)

// Parse a serialized endPoint such as "TcpEndPoint: 127.0.0.1:80",
// returning a pointer to the reconstructed endPoint"

func ParseEndPoint(str string) (ep EndPointI, err error) {
	parts := strings.Split(str, ": ")
	if len(parts) != 2 {
		err = NotAnEndPoint
	} else {
		if parts[0] == "TcpEndPoint" {
			addr := strings.TrimSpace(parts[1])
			ep, err = NewTcpEndPoint(addr)
		} else {
			err = NotAKnownEndPoint
		}
	}
	return
}
