package node

import (
	"sync"
)

// XXX This almost certainly belongs in its own package.
//
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

type MsgQueue struct {
	first    *MsgCarrier
	nextSeqN uint64
	t        int64 // ns
	mu       sync.Mutex
}

const (
	WIRE_FORM uint64 = 2 * iota // set if message has been marshaled
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
