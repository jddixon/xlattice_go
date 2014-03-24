package chunks

// xlattice_go/protocol/chunks/digiList.go

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	xu "github.com/jddixon/xlattice_go/util"
	"strings"
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
	timestamp xu.Timestamp
	digSig    []byte
}

func NewDigiList(sk *rsa.PublicKey, title string, timestamp xu.Timestamp) (
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

func (dl *DigiList) Timestamp() xu.Timestamp {
	return dl.timestamp
}

// Return a copy if the digital signature, or an error if the chunkList
// is unsigned.
func (dl *DigiList) GetDigSig() (digSig []byte, err error) {
	if dl.digSig == nil || len(dl.digSig) == 0 {
		err = NoDigSig
	} else {
		digSig = make([]byte, len(dl.digSig))
		copy(digSig, dl.digSig)
	}
	return
}

// DigiListI INTERFACE //////////////////////////////////////////////

// If dl.sk is not nil, return an error if it does not match the
// public part of private key.
//
// If there are any items in the DigiList, sign it.  If this succeeds,
// any existing signature is overwritten and the public part of the key
// is written to the data structure, to dl.sk.
//
func (dl *DigiList) Sign(skPriv *rsa.PrivateKey, subClass DigiListI) (
	err error) {

	var (
		hash   []byte
		n      uint
		sk     *rsa.PublicKey
		skWire []byte
	)
	if skPriv == nil {
		err = NilRSAPrivKey
	} else if subClass == nil {
		err = NilSubClass
	} else {
		n = subClass.Size()
		// DEBUG
		// fmt.Printf("DigiList.Sign(): subclass has %d chunks\n", n)
		// END
		sk = &skPriv.PublicKey

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
			itemHash, err = subClass.HashItem(i)
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
			rand.Reader, skPriv, crypto.SHA1, hash)
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
			itemHash, err = subClass.HashItem(i)
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

func ParseDigiList(s string) (dl *DigiList, rest []string, err error) {

	ss := strings.Split(s, "\r\n")
	lineCount := len(ss)
	if lineCount > 0 && ss[lineCount-1] == "" {
		ss = ss[:lineCount-1]
	}
	return ParseDigiListFromStrings(ss)
}

func ParseDigiListFromStrings(ss []string) (
	dl *DigiList, rest []string, err error) {

	var (
		sk     *rsa.PublicKey
		title  string
		t      xu.Timestamp
		digSig []byte
	)
	if len(ss) < 4 {
		err = TooShortForDigiList
	}
	if err == nil {
		sk, err = xc.RSAPubKeyFromDisk([]byte(ss[0]))
	}
	if err == nil {
		title = ss[1]
		t, err = xu.ParseTimestamp(ss[2])
	}
	if err == nil {
		digSig, err = base64.StdEncoding.DecodeString(ss[3])
	}
	if err == nil {
		dl = &DigiList{sk, title, t, digSig}
		rest = ss[4:]
	}
	return
}

// Serialize the DigiList, terminating each field and each item
// with a CRLF.
func (dl *DigiList) String() (s string) {
	ss := dl.Strings()
	return strings.Join(ss, "\r\n") + "\r\n"
}

func (dl *DigiList) Strings() (ss []string) {

	skSSH, err := xc.RSAPubKeyToDisk(dl.sk)
	if err != nil {
		panic(err)
	}
	ss = append(ss, fmt.Sprintf("sk: %s", skSSH))
	ss = append(ss, fmt.Sprintf("title: %s", dl.title))
	ss = append(ss, fmt.Sprintf("timestamp: %s", dl.timestamp.String()))
	ss = append(ss, fmt.Sprintf("digSig: %s",
		base64.StdEncoding.EncodeToString(dl.digSig)))

	return
}
