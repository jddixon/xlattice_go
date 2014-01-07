package crypto

// xlattice_go/crypto/pkcs7.go

import ()

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
