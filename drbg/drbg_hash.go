package drbg

import (
	"crypto/rand"
	"hash"
	"sync"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/drbg/hashdrbg"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

type hashDRGB struct {
	mu sync.Mutex // guards everything below; a DRBG is safe for concurrent use

	z           zeroize.ZeroizedFlag
	state       hashdrbg.State
	entropy     Entropy
	ownsEntropy bool

	securityStrength int
}

var _ krypto.Zeroizable = (*hashDRGB)(nil)

func (h *hashDRGB) Close() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if err := h.z.ReturnIfZeroized(); err != nil {
		return err
	}

	h.zeroize()

	return nil
}

func (h *hashDRGB) Zeroize() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.zeroize()
}

func (h *hashDRGB) zeroize() {
	h.z.SetInvalid()

	h.state.Zeroize()
	if h.ownsEntropy {
		zeroize.Zeroize(h.entropy)
	}
	h.entropy = nil
	h.ownsEntropy = false

	h.securityStrength = 0
}

func (h *hashDRGB) getEntropy() ([]byte, error) {
	return readEntropy(h.entropy, h.securityStrength, hashdrbg.MaxLength)
}

// NewHashDRBG is a deterministic random bit generator (DRBG) seeded by a hash function.
//
// TTAK.KO-12.0331
func NewHashDRBG(
	h hash.Hash,
	options ...Option,
) (krypto.DRBG, error) {
	var args args
	for _, v := range options {
		v(&args)
	}
	if args.ctrLen != 0 {
		return nil, ErrCTROptionUnsupported
	}

	return newHash(h, args)
}

// NewHashDRGB is a misspelled alias of [NewHashDRBG].
//
// Deprecated: Use [NewHashDRBG] instead.
func NewHashDRGB(h hash.Hash, options ...Option) (krypto.DRBG, error) {
	return NewHashDRBG(h, options...)
}

func newHash(h hash.Hash, args args) (krypto.DRBG, error) {
	hashOutLen := h.Size()

	maxStrength := min(hashOutLen, maxHashSecurityStrength(hashOutLen))
	if args.strengthBytes < 0 || args.strengthBytes > maxStrength {
		return nil, InvalidStrengthError(maxStrength)
	}
	if args.strengthBytes == 0 {
		args.strengthBytes = maxStrength
	}

	if args.reseedInterval > hashdrbg.ReseedInterval {
		return nil, InvalidReseedIntervalError(hashdrbg.ReseedInterval)
	}

	if uint64(len(args.personalizationString)) > hashdrbg.MaxPersonalizationStringLength {
		return nil, PersonalizationStringIsTooLongError(hashdrbg.MaxPersonalizationStringLength)
	}

	securityStrengthBytes := min(
		hashdrbg.GetSecurityStrengthBytes(args.strengthBytes),
		maxHashSecurityStrength(hashOutLen),
	)

	ownsEntropy := args.entropy == nil
	if ownsEntropy {
		args.entropy = newEntropy(rand.Reader, securityStrengthBytes)
	}

	var err error
	entropyInput := args.initEntropy
	if args.initEntropyIsSet {
		err = checkEntropy(entropyInput, securityStrengthBytes, hashdrbg.MaxLength)
	} else {
		entropyInput, err = readEntropy(args.entropy, securityStrengthBytes, hashdrbg.MaxLength)
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

	nonce, err := instantiateNonce(args, securityStrengthBytes, len(entropyInput))
	if err != nil {
		if ownsEntropy {
			zeroize.Zeroize(args.entropy)
		}
		return nil, err
	}
	if !args.nonceIsSet {
		defer clear(nonce)
	}

	drbg := &hashDRGB{
		entropy:          args.entropy,
		ownsEntropy:      ownsEntropy,
		securityStrength: securityStrengthBytes,
	}
	drbg.state.Instantiate(
		h,
		entropyInput,
		nonce,
		args.personalizationString,
		args.requirePredictionResistance,
		args.reseedInterval,
	)

	return drbg, nil
}

func (h *hashDRGB) Read(dst []byte) (n int, err error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if err := h.z.ReturnIfZeroized(); err != nil {
		return 0, err
	}

	remain := len(dst)
	for remain > 0 {
		toRead := min(remain, hashdrbg.MaxNoOfBytesPerRequest)

		_, err = h.generate(dst[len(dst)-remain:len(dst)-remain+toRead], nil)
		if err != nil {
			return len(dst) - remain, err
		}

		remain -= toRead
	}

	return len(dst), nil
}

func (h *hashDRGB) Generate(dst []byte, additionalInput []byte) (n int, err error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.generate(dst, additionalInput)
}

func (h *hashDRGB) generate(dst []byte, additionalInput []byte) (n int, err error) {
	if err := h.z.ReturnIfZeroized(); err != nil {
		return 0, err
	}

	if len(dst) > hashdrbg.MaxNoOfBytesPerRequest {
		return 0, TooManyBytesRequestedError(hashdrbg.MaxNoOfBytesPerRequest)
	}

	if uint64(len(additionalInput)) > hashdrbg.MaxAdditionalInputLength {
		return 0, AdditionalInputIsTooLongError(hashdrbg.MaxAdditionalInputLength)
	}

	var entropyInput []byte
	defer func() { clear(entropyInput) }()
	getEntropy := func() (input []byte, err error) {
		entropyInput, err = h.getEntropy()
		return entropyInput, err
	}
	if err = h.state.Generate(dst, getEntropy, additionalInput); err != nil {
		return 0, err
	}
	return len(dst), nil
}

func (h *hashDRGB) ReseedWithEntropy(additionalInput, entropyInput []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if err := h.z.ReturnIfZeroized(); err != nil {
		return err
	}

	if uint64(len(additionalInput)) > hashdrbg.MaxAdditionalInputLength {
		return AdditionalInputIsTooLongError(hashdrbg.MaxAdditionalInputLength)
	}

	if err := checkEntropy(entropyInput, h.securityStrength, hashdrbg.MaxLength); err != nil {
		return err
	}

	h.state.Reseed(entropyInput, additionalInput)
	return nil
}

func (h *hashDRGB) Reseed(additionalInput []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if err := h.z.ReturnIfZeroized(); err != nil {
		return err
	}

	if uint64(len(additionalInput)) > hashdrbg.MaxAdditionalInputLength {
		return AdditionalInputIsTooLongError(hashdrbg.MaxAdditionalInputLength)
	}

	entropyInput, err := h.getEntropy()
	if err != nil {
		return err
	}
	defer clear(entropyInput)

	h.state.Reseed(entropyInput, additionalInput)
	return nil
}
