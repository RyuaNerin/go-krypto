package hmac

const (
	msgUnsupportedHashCipher   = "krypto/hmac: unsupported hash cipher"
	msgInvalidTagSize          = "krypto/hmac: tag size must be positive and less than or equal to the hash size"
	msgHashGenerationNotUnique = "krypto/hmac: hash generation function does not produce unique values"
)
