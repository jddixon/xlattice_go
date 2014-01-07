package crypto

// xlatttice_go/crypto/signedListI.go

import (
	"crypto/rsa"
)

/**
 * In its serialized form a SignedList consists of a public key line,
 * a title line, a timestamp line, a number of content lines, and a
 * digital signature.  Each of the lines ends with a CR-LF sequence.
 * A blank line follows the last content line.  The timestamp (in
 * CCYY-MM-DD HH:MM:SS form) represents the time at which the list
 * was signed using the RSA private key corresponding to the key in
 * the public key line.  The public key itself is base-64 encoded.
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

type SignedListI interface {

	/** @return a clone of the public key */
	GetPublicKey() *rsa.PublicKey

	GetTitle() string

	IsSigned() bool

	/**
	 * Return this SignedList's SHA1 hash, a byte array 20 bytes
	 * long.
	 */
	GetHash() []byte

	/**
	 * Set a timestamp and calculate a digital signature.  First
	 * calculate the SHA1 hash of the pubkey, title, timestamp,
	 * and content lines, excluding the terminating CRLF in each
	 * case, then encrypt that using the RSA private key supplied.
	 *
	 * @param key RSAKey whose secret materials are used to sign
	 */
	Sign(key *rsa.PrivateKey) error

	/**
	 * The number of items in the list, excluding the header lines
	 * (public key, title, timestamp) and the footer lines (blank
	 * line, digital signature).
	 *
	 * @return the number of content items
	 */
	Size() int

	/**
	 * Verify that the SignedList agrees with its digital signature.
	 * Returns nil if the digital signature is correct.
	 */
	Verify() error

	// SERIALIZATION ////////////////////////////////////////////////

	/**
	 * Serialize the entire document.  All lines are CRLF-terminated.
	 * Subclasses are responsible for formatting their content lines,
	 * without any termination.
	 */
	String() string

	/**
	 * Nth content item in String form, without any terminating
	 * CRLF.
	 *
	 * @param n index of content item to be serialized
	 * @return  serialized content item
	 */
	Get(n int) (string, error)
}
