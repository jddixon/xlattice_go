// xlattice_go/crypto/sig.go

package crypto

import (
	cr "crypto"
	"crypto/rsa"
	"crypto/sha1"
	"errors"
)

// XXX CHANGE IN SPEC: Rather than panicking, we just
// return err, and then interpret a nil value as meaning "OK".

func SigVerify(pubkey *rsa.PublicKey, msg []byte, sig []byte) error {
	// presumably a rare error, so let's just complain
	if pubkey == nil || msg == nil || sig == nil {
		return errors.New("IllegalArgument: nil parameter")
	}
	d := sha1.New()
	d.Write(msg)
	hash := d.Sum(nil)

	return rsa.VerifyPKCS1v15(pubkey, cr.SHA1, hash, sig)
}
