package drbg

import (
	"crypto/cipher"
	"hash"

	"github.com/RyuaNerin/go-krypto"
)

type CTRBuilder struct {
	newCipher func([]byte) (cipher.Block, error)
	keySize   int
	args      args
}

// NewCTRBuilder returns a CTR DRBG builder with default settings.
func NewCTRBuilder(newCipher func([]byte) (cipher.Block, error), keySize int) *CTRBuilder {
	return &CTRBuilder{newCipher: newCipher, keySize: keySize}
}

// WithNonce sets the nonce used to instantiate the DRBG.
func (b *CTRBuilder) WithNonce(nonce []byte) *CTRBuilder {
	WithNonce(nonce)(&b.args)
	return b
}

// WithPersonalizationString sets the personalization string used to instantiate the DRBG.
func (b *CTRBuilder) WithPersonalizationString(ps []byte) *CTRBuilder {
	WithPersonalizationString(ps)(&b.args)
	return b
}

// WithPredictionResistance sets whether prediction resistance is required.
func (b *CTRBuilder) WithPredictionResistance(required bool) *CTRBuilder {
	WithPredictionResistance(required)(&b.args)
	return b
}

// WithStrengthBytes sets the requested security strength in bytes; zero uses the default.
func (b *CTRBuilder) WithStrengthBytes(size int) *CTRBuilder {
	WithStrengthBytes(size)(&b.args)
	return b
}

// WithReseedInterval sets the maximum number of Generate calls between reseeds; zero uses the default.
func (b *CTRBuilder) WithReseedInterval(interval uint64) *CTRBuilder {
	WithReseedInterval(interval)(&b.args)
	return b
}

// WithEntropy sets the entropy source used to instantiate and reseed the DRBG.
func (b *CTRBuilder) WithEntropy(entropy Entropy) *CTRBuilder {
	WithCustomEntropy(entropy)(&b.args)
	return b
}

// WithInstantiateEntropy sets the initial entropy input, bypassing the entropy source during instantiation.
func (b *CTRBuilder) WithInstantiateEntropy(entropyInput []byte) *CTRBuilder {
	WithInstantiateEntropy(entropyInput)(&b.args)
	return b
}

// WithDerivationFunction sets whether to use the CTR derivation function.
func (b *CTRBuilder) WithDerivationFunction(required bool) *CTRBuilder {
	WithDerivationFunction(required)(&b.args)
	return b
}

// WithCTRLength sets the counter length in bytes; zero uses the block size.
func (b *CTRBuilder) WithCTRLength(length int) *CTRBuilder {
	WithCTRLength(length)(&b.args)
	return b
}

// New instantiates a DRBG using the configured cipher, key size, and options.
func (b *CTRBuilder) New() (krypto.DRBG, error) {
	return newCTR(b.newCipher, b.keySize, b.args)
}

type HashBuilder struct {
	h    hash.Hash
	args args
}

// NewHashBuilder returns a Hash DRBG builder with default settings.
//
// The hash must not be shared with another DRBG.
func NewHashBuilder(h hash.Hash) *HashBuilder {
	return &HashBuilder{h: h}
}

// WithNonce sets the nonce used to instantiate the DRBG.
func (b *HashBuilder) WithNonce(nonce []byte) *HashBuilder {
	WithNonce(nonce)(&b.args)
	return b
}

// WithPersonalizationString sets the personalization string used to instantiate the DRBG.
func (b *HashBuilder) WithPersonalizationString(ps []byte) *HashBuilder {
	WithPersonalizationString(ps)(&b.args)
	return b
}

// WithPredictionResistance sets whether prediction resistance is required.
func (b *HashBuilder) WithPredictionResistance(required bool) *HashBuilder {
	WithPredictionResistance(required)(&b.args)
	return b
}

// WithStrengthBytes sets the requested security strength in bytes; zero uses the default.
func (b *HashBuilder) WithStrengthBytes(size int) *HashBuilder {
	WithStrengthBytes(size)(&b.args)
	return b
}

// WithReseedInterval sets the maximum number of Generate calls between reseeds; zero uses the default.
func (b *HashBuilder) WithReseedInterval(interval uint64) *HashBuilder {
	WithReseedInterval(interval)(&b.args)
	return b
}

// WithEntropy sets the entropy source used to instantiate and reseed the DRBG.
func (b *HashBuilder) WithEntropy(entropy Entropy) *HashBuilder {
	WithCustomEntropy(entropy)(&b.args)
	return b
}

// WithInstantiateEntropy sets the initial entropy input, bypassing the entropy source during instantiation.
func (b *HashBuilder) WithInstantiateEntropy(entropyInput []byte) *HashBuilder {
	WithInstantiateEntropy(entropyInput)(&b.args)
	return b
}

// New instantiates a DRBG using the configured hash and options.
// Call New only once successfully per builder, since the DRBG retains the hash.
func (b *HashBuilder) New() (krypto.DRBG, error) {
	return newHash(b.h, b.args)
}

type HMACBuilder struct {
	h    func() hash.Hash
	args args
}

// NewHMACBuilder returns a HMAC DRBG builder with default settings.
func NewHMACBuilder(h func() hash.Hash) *HMACBuilder {
	return &HMACBuilder{h: h}
}

// WithNonce sets the nonce used to instantiate the DRBG.
func (b *HMACBuilder) WithNonce(nonce []byte) *HMACBuilder {
	WithNonce(nonce)(&b.args)
	return b
}

// WithPersonalizationString sets the personalization string used to instantiate the DRBG.
func (b *HMACBuilder) WithPersonalizationString(ps []byte) *HMACBuilder {
	WithPersonalizationString(ps)(&b.args)
	return b
}

// WithPredictionResistance sets whether prediction resistance is required.
func (b *HMACBuilder) WithPredictionResistance(required bool) *HMACBuilder {
	WithPredictionResistance(required)(&b.args)
	return b
}

// WithStrengthBytes sets the requested security strength in bytes; zero uses the default.
func (b *HMACBuilder) WithStrengthBytes(size int) *HMACBuilder {
	WithStrengthBytes(size)(&b.args)
	return b
}

// WithReseedInterval sets the maximum number of Generate calls between reseeds; zero uses the default.
func (b *HMACBuilder) WithReseedInterval(interval uint64) *HMACBuilder {
	WithReseedInterval(interval)(&b.args)
	return b
}

// WithEntropy sets the entropy source used to instantiate and reseed the DRBG.
func (b *HMACBuilder) WithEntropy(entropy Entropy) *HMACBuilder {
	WithCustomEntropy(entropy)(&b.args)
	return b
}

// WithInstantiateEntropy sets the initial entropy input, bypassing the entropy source during instantiation.
func (b *HMACBuilder) WithInstantiateEntropy(entropyInput []byte) *HMACBuilder {
	WithInstantiateEntropy(entropyInput)(&b.args)
	return b
}

// New instantiates a DRBG using the configured hash constructor and options.
func (b *HMACBuilder) New() (krypto.DRBG, error) {
	return newHMAC(b.h, b.args)
}
