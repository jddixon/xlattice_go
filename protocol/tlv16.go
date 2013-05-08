package protocol

/**
 * Represents the classic Type-Length-Value field found in many
 * network protocols.
 *
 * This is an experiment; the implementation may change drastically
 * at a later date.
 * 
 * @author Jim Dixon
 */

import "errors"

type TLV16 struct {
    _type       uint16      // final
    length      uint16      // private, length in bytes
    value       *[]byte     // final
}

func (p *TLV16) Init (t uint16, v *[]byte) error {
    p._type     = t
    if v == nil { 
        return errors.New("IllegalArgument: nil data buffer" )
    }
    p.value     = v
    p.length    = uint16(len(*v)) 
    return nil
}

func (p *TLV16) Length() (uint16) {
    return p.length
}
func (p *TLV16) Type() (uint16) {
    return p._type
}
func (p *TLV16) Value() *[]byte {
    return p.value
}
func decodeUInt16(msg *[]byte, offset uint16) (n uint16) {

    n = uint16((*msg)[offset]) << 8
    offset ++
    n += uint16((*msg)[offset]) & 0xff
    return
}

func Decode (msg *[]byte, offset uint16) *TLV16 {

    if (msg == nil) {
        panic("IllegalArgument: nil msg")
    }
    _type := decodeUInt16(msg, offset)
    offset += 2
    msgLen  := decodeUInt16(msg, offset)
    offset += 2
    
//  // offset now points to beginning of value
//  if (msg.length < offset + len) {
//      throw new IllegalStateException(
//          "TLV in buffer of length " + msg.length 
//        + " but offset of value is " + offset 
//        + " and length is " + len)
//  }
    var val []byte
    val = make([]byte, msgLen)
    for i := uint16(0); i < msgLen; i++ {
        val[i] = (*msg)[offset + i]
    }
    p := new(TLV16)
    p.Init(_type, &val)
    return p
}
/* big-endian encoding of an unsigned int16 */
func encodeUInt16( n uint16, p *[]byte, offset uint16) (uint16) {
    (*p)[offset] = byte(n >> 8)
    offset++
    (*p)[offset] = byte(n)
    offset++
    return offset 
} 
    
/** 
 * Write this TLV onto the message buffer at the offset indicated.
 * 
 * @param buffer buffer to write TLV on
 * @param offset  byte offset where we start writing
 * @return offset after writing the values
 * @throws IndexOutOfBoundsException, NullPointerException
 */
 func (p *TLV16 ) Encode (buffer *[]byte, offset uint16) (uint16) {
    offset = encodeUInt16(p._type,  buffer, offset)
    offset = encodeUInt16(p.length, buffer, offset)
    for i := uint16(0); i < p.length; i++ {
        (*buffer)[offset + i] = (*p.value)[i]
    }
    return offset + p.length
}
