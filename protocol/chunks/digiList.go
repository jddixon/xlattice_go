package chunks

// xlattice_go/protocol/chunks/digiList.go

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/binary"
	xc "github.com/jddixon/xlattice_go/crypto"
)

// A DigiList is an abstract class.  Such a list is associated with an
// RSA key.  The public part of that key is part of this structure.
// The list has a title and a timestamp.  The owner of the key is
// responsible for making these meaningful.
//
// Items are added to the list using the DigiListI interface's
// HashItem (n int, interface{}) function.  When all items have been
// added, DigiList.Sign() is used to add the digital signature to
// the data structure.
type DigiList struct {
	sk        *rsa.PublicKey
	title     string
	timestamp int64
	digSig    []byte
}

func NewDigiList(sk *rsa.PublicKey, title string, timestamp int64) (
	dl *DigiList, err error) {

	// nil public key is acceptable
	if title == "" {
		err = EmptyTitle
	} else {
		dl = &DigiList{
			sk:        sk,
			title:     title,
			timestamp: timestamp,
		}
	}
	return
}

func (dl *DigiList) PublicKey() *rsa.PublicKey {
	return dl.sk // may be nil
}

func (dl *DigiList) Title() string {
	return dl.title
}

func (dl *DigiList) Timestamp() int64 {
	return dl.timestamp
}

func (dl *DigiList) DigSig() []byte {
	return dl.digSig
}

// DigiListI INTERFACE //////////////////////////////////////////////

// If dl.sk is not nil, return an error if it does not match the
// public part of the key.
//
// If there are any items in the DigiList, sign it.  If this succeeds,
// any existing signature is overwritten and the public part of the key
// is written to the data structure, to dl.sk.
//
func (dl *DigiList) Sign(key *rsa.PrivateKey, subClass DigiListI) (
	err error) {

	var (
		hash   []byte
		n      uint
		sk     *rsa.PublicKey
		skWire []byte
	)
	if key == nil {
		err = NilRSAPrivKey
	} else if subClass == nil {
		err = NilSubClass
	} else {
		n = subClass.Size()
		sk = &key.PublicKey

		// DEVIATION FROM SPEC - we ignore any existing dl.sk

		skWire, err = xc.RSAPubKeyToWire(sk)
	}
	if err == nil {
		d := sha1.New()
		d.Write(skWire)           // public key to hash
		d.Write([]byte(dl.title)) // title
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(dl.timestamp))
		d.Write(b) // timestamp to hash
		for i := uint(0); i < n; i++ {
			var itemHash []byte
			itemHash, err = subClass.HashItem(n)
			if err != nil {
				break
			}
			d.Write(itemHash)
		}
		if err == nil {
			hash = d.Sum(nil)
		}
	}
	if err == nil {
		dl.sk = sk
		dl.digSig, err = rsa.SignPKCS1v15(
			rand.Reader, key, crypto.SHA1, hash)
	}
	return
}

// If the DigiList has been signed, verify the digital signature,
// returning an error if it does not validate.
func (dl *DigiList) Verify(subClass DigiListI) (err error) {

	var (
		hash   []byte
		n      uint
		skWire []byte
	)
	if subClass == nil {
		err = NilSubClass
	} else if dl.digSig == nil {
		err = NoDigSig
	} else {
		n = subClass.Size()
		skWire, err = xc.RSAPubKeyToWire(dl.sk)
	}
	if err == nil {
		d := sha1.New()
		d.Write(skWire)           // public key to hash
		d.Write([]byte(dl.title)) // title
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(dl.timestamp))
		d.Write(b) // timestamp to hash
		for i := uint(0); i < n; i++ {
			var itemHash []byte
			itemHash, err = subClass.HashItem(n)
			if err != nil {
				break
			}
			d.Write(itemHash)
		}
		if err == nil {
			hash = d.Sum(nil)
		}
	}
	if err == nil {
		err = rsa.VerifyPKCS1v15(dl.sk, crypto.SHA1, hash, dl.digSig)
	}
	return
}

// SERIALIZATION ////////////////////////////////////////////////////

func ParseDigiList(str string) (dl *DigiList, err error) {

	// XXX STUB
	return
}

// Serialize the DigiList, terminating each field and each item
// with a CRLF.  This is a default implementation; subclasses
// satisfying DigiListI can override.
func (dl *DigiList) String() (str string) {
	// XXX STUB
	return
}
