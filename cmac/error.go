package cmac

const (
	msgUnsupportedBlockCipher = "krypto/cmac: unsupported block cipher"
	msgInvalidTagSize         = "krypto/cmac: tag size must be positive and less than or equal to the block size"
	msgInvalidInternalState   = "krypto/cmac: invalid internal state"
)
