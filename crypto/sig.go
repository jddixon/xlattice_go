// xlattice_go/crypto/sig.go

package crypto

import (
	cr "crypto"
	"crypto/rsa"
	"crypto/sha1"
)

// XXX POSSIBLE CHANGE IN SPEC: Rather than panicking, we could just
// return err, and then interpret a nil value as meaning "OK".
func SigVerify(pubkey *rsa.PublicKey, msg []byte, sig []byte) bool {
	// presumably a rare error, so let's just complain
	if pubkey == nil || msg == nil || sig == nil {
		panic("IllegalArgument: nil parameter")
	}
	d := sha1.New()
	d.Write(msg)
	hash := d.Sum(nil)

	err := rsa.VerifyPKCS1v15(pubkey, cr.SHA1, hash, sig)
	return err == nil
}
