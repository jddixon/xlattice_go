package reg

import (
	"fmt"
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestTestDir(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_TEST_DIR")
	}

	// XXX JUST A TEMPLATE AT THE MOMENT ////////////////////////////

	/////////////////////////////////////////////////////////////////
	// HELLO - REPLY TESTS 
	/////////////////////////////////////////////////////////////////

	// 1. Read key_rsa as key *rsa.PrivateKey

	// 2. Extract public key as pubkey *rsa.PublicKey

	// 3. Read key_rsa.pub as pubkey2 *rsa.PublicKey

	// 4. Verify pubkey == pubkey2

	// 5. Read version1.str as v1Str

	// 6. Read version1 as []byte

	// 7. Convert to dv1 DecimalVersion

	// 8. Verify v1Str == dv1.String()

	// 9, 10, 11, 12 same as 5-8 for version2

	// 13, 14, 15, 16 read iv1, key1, salt1, hello-data as []byte

	// 17. helloPlain = iv1 + key1 + salt1 + version1

	// 18. Verify helloPlain == helloData

	// 19. Read hello-encrypted as []byte

	// 20. Decrypt helloEncrypted using key => helloDecrypted

	// 21. Verify helloDecrypted == helloData

	// 22, 23, 24, 25, 26 read iv2, key2, salt2, padding, reply-data as []byte
	
	// 27. helloReply = concat iv2, key2, salt2, salt1, padding

	// 28. Verify helloReply == replyData

	// 29. Create aesEngineS1 from iv1, key1

	// 30. helloReplyMsg = aesEngineS1.encrypt(helloReply)

	// 31. Read reply-encrypted as replyEncrypted []byte

	// 32. Verify helloReplyMsg == replyEncrypted

	// 33. Create aesEngineC1 from iv1, key1

	// 34. Use aesEngineC1.decrypt(replyEncrypted) => replyDecrypted

	// 35. Verify replyDecrypted == replyData
}
