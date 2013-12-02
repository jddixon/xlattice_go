package reg

// xlattice_go/reg/aes_cnx.go

// Assume that a generator for this code is parameterized by
//	package name	- defaults to using local directory name
//  protocol name	- defaults to using whatever precedes "Msg" in *.proto
//  MSG_BUF_LEN		- defaults to 16 (K assumed)
//  file name		- defaults to protocol name + "_aes_cnx.go"
//  struct name		- defaults to protocol name + "AesCnxHandler"
//
// Generator is tested by generating the text for xlattice_go/reg
// and then comparing it to this file, with this comment block dropped.

import (
	"code.google.com/p/goprotobuf/proto"
	"crypto/aes"
	"crypto/cipher"
	xc "github.com/jddixon/xlattice_go/crypto"
	xt "github.com/jddixon/xlattice_go/transport"
	// "sync"
)

const (
	MSG_BUF_LEN = 16 * 1024
)

type AesCnxHandler struct {
	State                              int
	Cnx                                *xt.TcpConnection
	engine                             cipher.Block
	encrypter                          cipher.BlockMode
	decrypter                          cipher.BlockMode
	iv1, key1, iv2, key2, salt1, salt2 []byte
}

func (a *AesCnxHandler) SetupSessionKey() (err error) {
	a.engine, err = aes.NewCipher(a.key2)
	if err == nil {
		a.encrypter = cipher.NewCBCEncrypter(a.engine, a.iv2)
		a.decrypter = cipher.NewCBCDecrypter(a.engine, a.iv2)
	}
	return
}

// Read data from the connection.  XXX This will not handle partial
// reads correctly.
//
func (a *AesCnxHandler) ReadData() (data []byte, err error) {
	data = make([]byte, MSG_BUF_LEN)
	count, err := a.Cnx.Read(data)
	if err == nil && count > 0 {
		data = data[:count]
		return
	}
	return nil, err
}

func (a *AesCnxHandler) WriteData(data []byte) (err error) {
	count, err := a.Cnx.Write(data)

	// XXX handle cases where not all bytes written

	_ = count
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

func EncodePadEncrypt(msg *XLRegMsg, engine cipher.BlockMode) (
	ciphertext []byte, err error) {

	var paddedData []byte

	cData, err := EncodePacket(msg)
	if err == nil {
		paddedData, err = xc.AddPKCS7Padding(cData, aes.BlockSize)
	}
	if err == nil {
		msgLen := len(paddedData)
		nBlocks := (msgLen + aes.BlockSize - 2) / aes.BlockSize
		ciphertext = make([]byte, nBlocks*aes.BlockSize)
		engine.CryptBlocks(ciphertext, paddedData) // dest <- src
	}
	return
}

func DecryptUnpadDecode(ciphertext []byte, engine cipher.BlockMode) (
	msg *XLRegMsg, err error) {

	plaintext := make([]byte, len(ciphertext))
	engine.CryptBlocks(plaintext, ciphertext) // dest <- src

	unpaddedCData, err := xc.StripPKCS7Padding(plaintext, aes.BlockSize)
	if err == nil {
		msg, err = DecodePacket(unpaddedCData)
	}
	return
}
