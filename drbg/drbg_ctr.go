package drbg

import (
	"crypto/cipher"
	"crypto/rand"
	"sync"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/drbg/ctrdrbg"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

func WithCTRLength(ctrLen int) Option {
	return func(args *args) {
		args.ctrLen = ctrLen
	}
}

type ctrDRGB struct {
	mu sync.Mutex // guards everything below; a DRBG is safe for concurrent use

	z           zeroize.ZeroizedFlag
	state       ctrdrbg.State
	entropy     Entropy
	ownsEntropy bool

	securityStrength          int
	requireDerivationFunction bool
}

var _ krypto.DRBG = (*ctrDRGB)(nil)

func (h *ctrDRGB) Close() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if err := h.z.ReturnIfZeroized(); err != nil {
		return err
	}

	h.zeroize()

	return nil
}

func (h *ctrDRGB) Zeroize() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.zeroize()
}

func (h *ctrDRGB) zeroize() {
	h.z.SetInvalid()

	h.state.Zeroize()
	if h.ownsEntropy {
		zeroize.Zeroize(h.entropy)
	}
	h.entropy = nil
	h.ownsEntropy = false

	h.securityStrength = 0
	h.requireDerivationFunction = false
}

func ctrEntropyLen(requireDerivationFunction bool, securityStrength, seedLen int) (minLen int, maxLen uint64) {
	if requireDerivationFunction {
		return securityStrength, ctrdrbg.MaxLength
	}
	return seedLen, uint64(seedLen)
}

func (h *ctrDRGB) getEntropy() ([]byte, error) {
	minLen, maxLen := ctrEntropyLen(h.requireDerivationFunction, h.securityStrength, h.state.SeedLen)
	return readEntropy(h.entropy, minLen, maxLen)
}

// NewCTRDRBG is a deterministic random bit generator (DRBG) based on a block cipher.
//
// TTAK.KO-12.0189/R1
func NewCTRDRBG(
	newCipher func(key []byte) (cipher.Block, error),
	keySize int, // Bytes
	options ...Option,
) (krypto.DRBG, error) {
	var args args
	for _, v := range options {
		v(&args)
	}
	return newCTR(newCipher, keySize, args)
}

func newCTR(
	newCipher func(key []byte) (cipher.Block, error),
	keySize int, // Bytes
	args args,
) (krypto.DRBG, error) {
	if keySize <= 0 {
		return nil, InvalidKeySizeError(keySize)
	}
	if args.strengthBytes < 0 || args.strengthBytes > keySize {
		return nil, InvalidStrengthError(keySize)
	}
	b, err := newCipher(make([]byte, keySize))
	if err != nil {
		return nil, err
	}
	defer zeroize.Zeroize(b)
	blockSize := b.BlockSize()
	if blockSize != 8 && blockSize != 16 {
		return nil, InvalidBlockSizeError(blockSize)
	}
	seedLen := blockSize + keySize
	if args.reseedInterval > ctrdrbg.MaxReseedInterval(blockSize) {
		return nil, InvalidReseedIntervalError(ctrdrbg.MaxReseedInterval(blockSize))
	}

	if args.requireDerivationFunction {
		if uint64(len(args.personalizationString)) > ctrdrbg.MaxPersonalizationStringLength {
			return nil, PersonalizationStringIsTooLongError(ctrdrbg.MaxPersonalizationStringLength)
		}
	} else {
		// Without a derivation function the personalization string is XORed into
		// the seed, so it can not be longer than the seed.
		if len(args.personalizationString) > seedLen {
			return nil, PersonalizationStringIsTooLongError(seedLen)
		}
	}

	if args.ctrLen == 0 {
		args.ctrLen = blockSize
	}
	if args.ctrLen < 1 || blockSize < args.ctrLen {
		return nil, InvalidCtrLengthError{Min: 1, Max: blockSize}
	}

	if args.strengthBytes == 0 {
		args.strengthBytes = keySize
	}

	ownsEntropy := args.entropy == nil
	if ownsEntropy {
		if args.requireDerivationFunction {
			args.entropy = newEntropy(rand.Reader, args.strengthBytes)
		} else {
			args.entropy = newEntropy(rand.Reader, seedLen)
		}
	}

	minLen, maxLen := ctrEntropyLen(args.requireDerivationFunction, args.strengthBytes, seedLen)
	entropyInput := args.initEntropy
	if args.initEntropyIsSet {
		err = checkEntropy(entropyInput, minLen, maxLen)
	} else {
		entropyInput, err = readEntropy(args.entropy, minLen, maxLen)
	}
	if err != nil {
		if ownsEntropy {
			zeroize.Zeroize(args.entropy)
		}
		return nil, err
	}
	if !args.initEntropyIsSet {
		defer clear(entropyInput)
	}

	nonce := args.nonce
	if args.requireDerivationFunction {
		nonce, err = instantiateNonce(args, args.strengthBytes, len(entropyInput))
		if err != nil {
			if ownsEntropy {
				zeroize.Zeroize(args.entropy)
			}
			return nil, err
		}
		if !args.nonceIsSet {
			defer clear(nonce)
		}
	}

	drbg := &ctrDRGB{
		entropy:          args.entropy,
		ownsEntropy:      ownsEntropy,
		securityStrength: args.strengthBytes,

		requireDerivationFunction: args.requireDerivationFunction,
	}

	err = drbg.state.Instantiate(
		internal.ToBlock8Fn(newCipher),
		keySize,
		args.reseedInterval,
		args.ctrLen,
		entropyInput,
		nonce,
		args.personalizationString,
		args.requireDerivationFunction,
		args.requirePredictionResistance,
	)
	if err != nil {
		drbg.zeroize()
		return nil, err
	}

	return drbg, nil
}

func (h *ctrDRGB) Read(dst []byte) (n int, err error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if err := h.z.ReturnIfZeroized(); err != nil {
		return 0, err
	}

	remain := len(dst)
	for remain > 0 {
		toRead := min(remain, h.state.MaxNoOfBytesPerRequest)

		_, err = h.generate(dst[len(dst)-remain:len(dst)-remain+toRead], nil)
		if err != nil {
			return len(dst) - remain, err
		}

		remain -= toRead
	}

	return len(dst), nil
}

func (h *ctrDRGB) Generate(dst []byte, additionalInput []byte) (n int, err error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.generate(dst, additionalInput)
}

func (h *ctrDRGB) generate(dst []byte, additionalInput []byte) (n int, err error) {
	if err := h.z.ReturnIfZeroized(); err != nil {
		return 0, err
	}

	if len(dst) > h.state.MaxNoOfBytesPerRequest {
		return 0, TooManyBytesRequestedError(uint64(h.state.MaxNoOfBytesPerRequest))
	}

	if h.requireDerivationFunction {
		if uint64(len(additionalInput)) > ctrdrbg.MaxAdditionalInputLength {
			return 0, AdditionalInputIsTooLongError(ctrdrbg.MaxAdditionalInputLength)
		}
	} else {
		if len(additionalInput) > h.state.SeedLen {
			return 0, AdditionalInputIsTooLongError(uint64(h.state.SeedLen))
		}
	}

	// The state consumes reseed entropy synchronously. Clear the owned copy
	// when generation finishes; custom source buffers remain untouched.
	var entropyInput []byte
	defer func() { clear(entropyInput) }()
	getEntropy := func() (input []byte, err error) {
		entropyInput, err = h.getEntropy()
		return entropyInput, err
	}
	err = h.state.Generate(dst, getEntropy, additionalInput)
	if err != nil {
		return 0, err
	}
	return len(dst), nil
}

func (h *ctrDRGB) ReseedWithEntropy(additionalInput, entropyInput []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if err := h.z.ReturnIfZeroized(); err != nil {
		return err
	}

	if h.requireDerivationFunction {
		if uint64(len(additionalInput)) > ctrdrbg.MaxAdditionalInputLength {
			return AdditionalInputIsTooLongError(ctrdrbg.MaxAdditionalInputLength)
		}
	} else {
		if len(additionalInput) > h.state.SeedLen {
			return AdditionalInputIsTooLongError(uint64(h.state.SeedLen))
		}
	}

	minLen, maxLen := ctrEntropyLen(h.requireDerivationFunction, h.securityStrength, h.state.SeedLen)
	if err := checkEntropy(entropyInput, minLen, maxLen); err != nil {
		return err
	}

	return h.state.Reseed(entropyInput, additionalInput)
}

func (h *ctrDRGB) Reseed(additionalInput []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if err := h.z.ReturnIfZeroized(); err != nil {
		return err
	}

	if h.requireDerivationFunction {
		if uint64(len(additionalInput)) > ctrdrbg.MaxAdditionalInputLength {
			return AdditionalInputIsTooLongError(ctrdrbg.MaxAdditionalInputLength)
		}
	} else {
		if len(additionalInput) > h.state.SeedLen {
			return AdditionalInputIsTooLongError(uint64(h.state.SeedLen))
		}
	}

	entropyInput, err := h.getEntropy()
	if err != nil {
		return err
	}
	defer clear(entropyInput)

	return h.state.Reseed(entropyInput, additionalInput)
}
