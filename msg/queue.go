package msg

// xlattice_go/msg/queue.go

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	xn "github.com/jddixon/xlNode_go"
	xt "github.com/jddixon/xlTransport_go"
	"sync"
)

var _ = fmt.Print

// A node maintains a queue of message carriers for each outbound
// connection to a peer.  Each message carrier has a pointer to a message.
// The message carrier contains information about the message queued.
// This may include its sequence number and a send count for messages
// which may need to be retransmitted.  It will include a scheduled send
// time and a flag field.  The low-level bit in the flag field is set if
// the message is marshaled to wire format and otherwise clear.
//
// A unique, ascending sequence number is assigned to each message sent
// on a connection, except that where a message is retransmitted it
// retains the same sequence number.  A message expecting an acknowlegement
// will be retransmitted up to maxSend times if it is not acknowledged.
// (XXX We probably need a callback to handle messages that have been
// resent too many times.)
//
// When a message requiring an acknowlegement is sent, a copy is
// added to the message queue with its sendCount set to 1.  When
// the acknowlegement is received, it will contain the same sequence
// number, making it possible to search the message queue and delete
// the copy scheduled for retransmission.

const (
	MSG_BUF_LEN = 16 * 1024
)

var ONE = uint64(1)

type MsgQueue struct {
	first    *MsgCarrier
	nextMsgN uint64
	t        int64 // ns
	mu       sync.Mutex
}

// a field of 64 bit flags
const (
	WIRE_FORM uint64 = 1 << iota // set if message has been marshaled
	FOO_FOO
)

type MsgCarrier struct {
	next      *MsgCarrier
	seqN      uint64
	t         int64 // ns since epoch
	flags     uint64
	sendCount int // incremented when msg resent
	maxSend   int
	msg       *interface{}
}

type CnxHandler struct {
	Cnx   *xt.TcpConnection
	Peer  *xn.Peer
	MsgN  uint64
	State int
}

// Read the next message over the connection
func (h *CnxHandler) readMsg() (m *XLatticeMsg, err error) {
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
func (h *CnxHandler) writeMsg(m *XLatticeMsg) (err error) {
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
} // GEEP

func DecodePacket(data []byte) (*XLatticeMsg, error) {
	var m XLatticeMsg
	err := proto.Unmarshal(data, &m)
	// XXX do some filtering, eg for nil op
	return &m, err
}

func EncodePacket(msg *XLatticeMsg) (data []byte, err error) {
	return proto.Marshal(msg)
}
