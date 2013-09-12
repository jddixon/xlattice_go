package reg

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

// TODO: MOVE THIS TO crypto/ =======================================

// PKCS7 padding (RFC 5652) pads a message out to a whole multiple
// of the block size, with the value of each byte being the number
// of bytes of padding.  If the data passed is nil, the function
// returns a full block of padding.

func PKCS7Padding(data []byte, blockSize int) (padding []byte) {
	var length int
	if data == nil {
		length = 0
	} else {
		length = len(data)
	}
	// we want from 1 to blockSize bytes of padding
	nBlocks := (length + blockSize - 1) / blockSize
	rem := nBlocks*blockSize - length
	if rem == 0 {
		rem = blockSize
	}
	padding = make([]byte, rem)
	for i := 0; i < rem; i++ {
		padding[i] = byte(rem)
	}
	return
}

var (
	ImpossibleBlockSize   = errors.New("impossible block size")
	IncorrectPKCS7Padding = errors.New("incorrectly padded data")
	NilData               = errors.New("nil data argument")
)

func AddPKCS7Padding(data []byte, blockSize int) (out []byte, err error) {
	if blockSize <= 1 {
		err = ImpossibleBlockSize
	} else {
		padding := PKCS7Padding(data, blockSize)
		if data == nil {
			out = padding
		} else {
			out = append(data, padding...)
		}
	}
	return
}

// The data passed is presumed to have PKCS7 padding.  If possible, return
// a copy of the data without the padding.  Return an error if the padding
// is incorrect.

func StripPKCS7Padding(data []byte, blockSize int) (out []byte, err error) {
	if blockSize <= 1 {
		err = ImpossibleBlockSize
	} else if data == nil {
		err = NilData
	}
	if err == nil {
		lenData := len(data)
		if lenData < blockSize {
			err = IncorrectPKCS7Padding
		} else {
			// examine the very last byte: it must be padding and must
			// contain the number of padding bytes added
			lenPadding := int(data[lenData-1])
			if lenPadding < 1 || lenData < lenPadding {
				err = IncorrectPKCS7Padding
			} else {
				out = data[:lenData-lenPadding]
			}
		}
	}
	return
}

func EncodePadEncrypt(msg *XLRegMsg, engine cipher.BlockMode) (ciphertext []byte, err error) {
	var paddedData []byte

	cData, err := EncodePacket(msg)
	if err == nil {
		paddedData, err = AddPKCS7Padding(cData, aes.BlockSize)
	}
	if err == nil {
		msgLen := len(paddedData)
		nBlocks := (msgLen + aes.BlockSize - 2) / aes.BlockSize
		ciphertext = make([]byte, nBlocks*aes.BlockSize)
		engine.CryptBlocks(ciphertext, paddedData) // dest <- src
	}
	return
}

func DecryptUnpadDecode(ciphertext []byte, engine cipher.BlockMode) (msg *XLRegMsg, err error) {

	plaintext := make([]byte, len(ciphertext))
	engine.CryptBlocks(plaintext, ciphertext) // dest <- src

	unpaddedCData, err := StripPKCS7Padding(plaintext, aes.BlockSize)
	if err == nil {
		msg, err = DecodePacket(unpaddedCData)
	}
	return
}
