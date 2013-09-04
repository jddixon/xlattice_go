package reg

// xlattice_go/reg/packets.go

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	// xn "github.com/jddixon/xlattice_go/node"
	xt "github.com/jddixon/xlattice_go/transport"
	// "sync"
)

var _ = fmt.Print

const (
	MSG_BUF_LEN = 16 * 1024
)

type CnxHandler struct {
	State int // as yet unused
	Cnx   *xt.TcpConnection
}

// Read the next message over the connection
func (h *CnxHandler) readMsg() (m *XLRegMsg, err error) {
	inBuf := make([]byte, MSG_BUF_LEN)
	count, err := h.Cnx.Read(inBuf)
	if err == nil && count > 0 {
		inBuf = inBuf[:count]
		// parse = deserialize, unmarshal the message
		m, err = DecodePacket(inBuf)
	}
	return
}

// Write a message out over the connection
func (h *CnxHandler) writeMsg(m *XLRegMsg) (err error) {
	var count int
	var data []byte
	// serialize, marshal the message
	data, err = EncodePacket(m)
	if err == nil {
		count, err = h.Cnx.Write(data)
		// XXX handle cases where not all bytes written
		_ = count
	}
	return
}

func DecodePacket(data []byte) (*XLRegMsg, error) {
	var m XLRegMsg
	err := proto.Unmarshal(data, &m)
	// XXX do some filtering, eg for nil op
	return &m, err
}

func EncodePacket(msg *XLRegMsg) (data []byte, err error) {
	return proto.Marshal(msg)
}
