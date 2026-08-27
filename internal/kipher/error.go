package kipher

import "errors"

const (
	msgNotFullBlocks           = "krypto/kipher: src not full blocks"
	msgSmallDst                = "krypto/kipher: dst smaller than src"
	msgBufferOverlap           = "krypto/kipher: invalid buffer overlap"
	msgBufferOverlapAdditional = "krypto/kipher: invalid buffer overlap of output and additional data"
	msgInvalidNonceZero        = "krypto/kipher: the nonce can't have zero length, or the security of the key will be immediately compromised"
	msgInvalidNonceSize        = "krypto/kipher: invalid nonce size"
	msgInvalidTagSizeGCM       = "krypto/kipher: incorrect tag size given"
	msgInvalidTagSizeCCM       = "krypto/kipher: tagSize must be 4, 6, 8, 10, 12, 14 or 16 in CCM"
	msgInvalidIVLength         = "krypto/kipher: IV length must equal block size"
	msgInvalidNonce            = "krypto/kipher: incorrect nonce length given"
	msgDataTooLarge            = "krypto/kipher: message too large"
	msgRequire128Bits          = "krypto/kipher: requires 128-bit block cipher"
	msgOpenFailed              = "krypto/kipher: message authentication failed"
	msgInvalidCFBParameters    = "krypto/kipher: invalid CFB parameters"
	msgInvalidBlockSize        = "krypto/kipher: invalid block size"
)

var ErrOpenFailed = errors.New(msgOpenFailed)
