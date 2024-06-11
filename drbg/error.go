package drbg

const (
	msgErrorUninstantiated                  = "krypto/drbg: uninstantiated"
	msgInvalidEntropyFormat                 = "krypto/drbg: invalid entropy size. must be between %d and %d bytes"
	msgInvalidStrengthFormat                = "krypto/drbg: invalid strength. maximum length is %d bytes"
	msgPersonalizationStringIsTooLongFormat = "krypto/drbg: personalization_string is too long. maximum length is %d bytes"
	msgTooManyBytesRequestedFormat          = "krypto/drbg: too many bytes requested. maximum length is %d bytes"
	msgAdditionalInputIsTooLongFormat       = "krypto/drbg: additionalInput is too long. miximum length is %d bytes"
	msgInvalidCtrLengthFormat2              = "krypto/drbg: invalid counter length. must be between %d and %d"
)
