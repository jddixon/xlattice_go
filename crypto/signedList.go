package crypto

// xlatttice_go/crypto/signedList.go

import (
	"bufio"
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	xu "github.com/jddixon/xlattice_go/util"
	"io"
	"strings"
	"time"
)

var (
	CONTENT_START = []byte("# BEGIN CONTENT ")
	CONTENT_END   = []byte("# END CONTENT ")
)

/**
 * In its serialized form a SignedList consists of a public key line,
 * a title line, a timestamp line, a number of content lines, and a
 * digital signature.  Each of the lines ends with a CR-LF sequence.
 * A blank line follows the last content line.  The timestamp (in
 * CCYY-MM-DD HH:MM:SS form) represents the time at which the list
 * was signed using the RSA private key corresponding to the key in
 * the public key line.
 *
 * The SHA1withRSA digital signature is on the entire SignedList excluding
 * the digital signature line.  All line endings are converted to
 * CRLF before taking the digital signature.
 *
 * The SignedList itself has a 20-byte extended hash, the 20-byte SHA1
 * digest of a function of the public key and the title.  This means
 * that the owner of the RSA key can create any number of documents
 * with the same hash but different timestamps with the intention
 * being that users can choose to regard the document with the most
 * recent timestamp as authentic.
 *
 * What the content line contains varies between subclasses.
 */

type SignedList struct {
	pubKey    *rsa.PublicKey
	title     string
	timestamp xu.Timestamp // set when signed
	digSig    []byte
}

func NewSignedList(pubKey *rsa.PublicKey, title string) (
	sl *SignedList, err error) {

	if pubKey == nil {
		err = NilPublicKey
	} else if title == "" {
		err = EmptyTitle
	} else {
		sl = &SignedList{
			pubKey: pubKey,
			title:  title,
		}
	}
	return
}

// PROPERTIES ///////////////////////////////////////////////////

func (sl *SignedList) GetPublicKey() *rsa.PublicKey {
	return sl.pubKey
}
func (sl *SignedList) GetTitle() string {
	return sl.title
}

func (sl *SignedList) IsSigned() bool {
	return sl.digSig != nil
}

/**
 * Return this SignedList's SHA1 hash, a byte array 20 bytes
 * long.  The hash is over first the public key in its 'wire' form
 * and then over the title.
 *
 * XXX This is completely different from the current Java implementation;
 * the two must be reconciled.
 */
func (sl *SignedList) GetHash() []byte {

	d := sha1.New()

	// public key in PKCS1 format
	pk, _ := RSAPubKeyToWire(sl.pubKey)
	d.Write(pk)

	d.Write([]byte(sl.title))
	return d.Sum(nil)
}

// Return the SHA1 hash of the SignedList, excluding the digital
// signature but expecting the timestamp to have been set.
func (sl *SignedList) HashBody() (hash []byte, err error) {
	d := sha1.New()

	// public key in SSH format ---------------------------
	pk, _ := RSAPubKeyToDisk(sl.pubKey)
	d.Write(pk)

	// title ----------------------------------------------
	d.Write([]byte(sl.title))

	// timestamp ------------------------------------------
	d.Write([]byte(sl.timestamp.String()))

	// content lines --------------------------------------
	for i := 0; err == nil && i < sl.Size(); i++ {
		var line string
		line, err = sl.Get(i)
		if err == nil || err == io.EOF {
			d.Write([]byte(line))
			if err == io.EOF {
				err = nil
				break
			}
		}
	}
	if err == nil {
		hash = d.Sum(nil)
	}
	return
}

/**
 * Set a timestamp and calculate a digital signature.  First
 * calculate the SHA1 hash of the pubKey, title, timestamp,
 * and content lines, excluding the terminating CRLF in each
 * case, then encrypt that using the RSA private key supplied.
 *
 * @param key RSAKey whose secret materials are used to sign
 */
func (sl *SignedList) Sign(skPriv *rsa.PrivateKey) (err error) {

	var (
		digSig, hash []byte
	)

	if sl.digSig != nil {
		err = ListAlreadySigned
	} else if skPriv == nil {
		err = NilPrivateKey
	} else {
		sl.timestamp = xu.Timestamp(time.Now().UnixNano())
		hash, err = sl.HashBody()
		if err == nil {
			digSig, err = rsa.SignPKCS1v15(
				rand.Reader, skPriv, crypto.SHA1, hash)
			if err == nil {
				sl.digSig = digSig
			}
		}
		if err != nil {
			sl.timestamp = 0 // restore to default
		}
	}
	return
}

/**
 * The number of items in the list, excluding the header lines
 * (public key, title, timestamp) and the footer lines (blank
 * line, digital signature).
 *
 * @return the number of content items
 */
func (sl *SignedList) Size() int {
	// SUBCLASS MUST IMPLEMENT
	return -1
}

/**
 * Verify that the SignedList agrees with its digital signature,
 * returning nil if it is correct and an appropriate error otherwise.
 */
func (sl *SignedList) Verify() (err error) {

	var hash []byte

	if sl.digSig == nil {
		err = UnsignedList
	} else {
		hash, err = sl.HashBody()
		if err == nil {
			err = rsa.VerifyPKCS1v15(sl.pubKey, crypto.SHA1, hash, sl.digSig)
		}
	}
	return
}

// SERIALIZATION ////////////////////////////////////////////////

/**
 * Serialize the entire document.  All lines are CRLF-terminated.
 * Subclasses are responsible for formatting their content lines,
 * without any termination.  If any error is encountered, this
 * function silently returns an empty string.
 */
func (sl *SignedList) String() (s string) {

	var (
		err error
		ss  []string
	)
	// public key to SSH format -----------------------
	pk, _ := RSAPubKeyToDisk(sl.pubKey)
	ss = append(ss, string(pk))

	// title ------------------------------------------
	ss = append(ss, sl.title)

	// timestamp --------------------------------------
	ss = append(ss, sl.timestamp.String())

	// content lines ----------------------------------
	for i := 0; err == nil && i < sl.Size(); i++ {
		var line string
		line, err = sl.Get(i)
		if err == nil || err == io.EOF {
			ss = append(ss, line)
			if err == io.EOF {
				err = nil
				break
			}
		}
	}
	if err == nil {
		s = strings.Join(ss, CRLF)
	}
	return
}

/**
 * Return the Nth content item in String form, without any terminating CRLF.
 * Using code should permit the implementation to return io.EOF either
 * with the last valid line or an empty string and io.EOF on subsequent
 * calls.
 */
func (sl *SignedList) Get(n int) (s string, err error) {

	/* SUBCLASSES MUST IMPLEMENT */

	err = NotImplemented
	return
}

/**
 * Reads in content lines, stripping off line endings, storing the
 * line in a subclass-defined internal buffer (conventionally "content").
 */
func (sl *SignedList) ReadContents(*bufio.Reader) (err error) {

	/* SUBCLASSES MUST IMPLEMENT */

	err = NotImplemented
	return
}

// SERIALIZATION ////////////////////////////////////////////////////

func NextLineWithoutCRLF(in *bufio.Reader) (line []byte, err error) {
	line, err = in.ReadBytes('\n')
	if err == nil || err == io.EOF {
		line = line[:len(line)-1] // drop the \n
		lineLen := len(line)
		if lineLen > 0 && line[lineLen-1] == '\r' {
			line = line[:len(line)-1] // drop any \r
		}
	}
	return
}

// Read the header part of a signed list that has been serialized in disk
// format, returning a pointer to the deserialized object or an error.
// Subclasses should call this to get a pointer to the SignedList part
// of the subclass struct.  If the subclass is an XXXList, then expect
// the calling routine to be ParseXXXList()
//
func ParseSignedList(in *bufio.Reader) (sl *SignedList, err error) {

	var (
		line   []byte
		pubKey *rsa.PublicKey
		title  string
		t      xu.Timestamp // binary form
	)

	line, err = NextLineWithoutCRLF(in)
	if err == nil {
		pubKey, err = RSAPubKeyFromDisk(line)
		if err == nil {
			line, err = NextLineWithoutCRLF(in)
			if err == nil {
				title = string(line)
				line, err = NextLineWithoutCRLF(in)
				if err == nil {
					t, err = xu.ParseTimestamp(string(line))
					if err == nil {
						line, err = NextLineWithoutCRLF(in)
						if err == nil {
							if !bytes.Equal(line, CONTENT_START) {
								err = MissingContentStart
							}
						}
					}
				}
			}
		}
	}
	if err == nil {
		sl = &SignedList{
			pubKey:    pubKey,
			title:     title,
			timestamp: t,
		}
	}
	return
}
