package transport

// xlattice_go/transport/connector.go

import (
	"strings"
)

// Parse a serialized connector such as "TcpConnector: 127.0.0.1:80",
// returning a pointer to the reconstructed connector"

func ParseConnector(str string) (ctor ConnectorI, err error) {
	parts := strings.Split(str, ": ")
	if len(parts) != 2 {
		err = NotAConnector
	} else {
		if parts[0] == "TcpConnector" {
			addr := strings.TrimSpace(parts[1])
			var ep *TcpEndPoint
			ep, err = NewTcpEndPoint(addr)
			if err == nil {
				ctor, err = NewTcpConnector(ep)
			}
		} else {
			err = NotAKnownConnector
		}
	}
	return
}
