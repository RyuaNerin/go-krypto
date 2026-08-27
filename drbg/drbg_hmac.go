package drbg

import (
	"crypto/rand"
	"hash"
	"sync"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/drbg/hmacdrbg"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

type hmacDRGB struct {
	mu sync.Mutex // guards everything below; a DRBG is safe for concurrent use

	z           zeroize.ZeroizedFlag
	state       hmacdrbg.State
	entropy     Entropy
	ownsEntropy bool

	securityStrength int
}

var _ krypto.Zeroizable = (*hmacDRGB)(nil)

func (h *hmacDRGB) Close() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if err := h.z.ReturnIfZeroized(); err != nil {
		return err
	}

	h.zeroize()

	return nil
}

func (h *hmacDRGB) Zeroize() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.zeroize()
}

func (h *hmacDRGB) zeroize() {
	h.z.SetInvalid()

	h.state.Zeroize()
	if h.ownsEntropy {
		zeroize.Zeroize(h.entropy)
	}
	h.entropy = nil
	h.ownsEntropy = false

	h.securityStrength = 0
}

func (h *hmacDRGB) getEntropy() ([]byte, error) {
	return readEntropy(h.entropy, h.securityStrength, hmacdrbg.MaxLength)
}

// NewHMACDRBG is a deterministic random bit generator (DRBG) seeded by HMAC.
//
// TTAK.KO-12.0332
func NewHMACDRBG(
	h func() hash.Hash,
	options ...Option,
) (krypto.DRBG, error) {
	var args args
	for _, v := range options {
		v(&args)
	}
	if args.ctrLen != 0 {
		return nil, ErrCTROptionUnsupported
	}
	return newHMAC(h, args)
}

// NewHMACDRGB is a misspelled alias of [NewHMACDRBG].
//
// Deprecated: Use [NewHMACDRBG] instead.
func NewHMACDRGB(h func() hash.Hash, options ...Option) (krypto.DRBG, error) {
	return NewHMACDRBG(h, options...)
}

func newHMAC(h func() hash.Hash, args args) (krypto.DRBG, error) {
	hashOutLen := h().Size()

	maxStrength := min(hashOutLen, maxHashSecurityStrength(hashOutLen))
	if args.strengthBytes < 0 || args.strengthBytes > maxStrength {
		return nil, InvalidStrengthError(maxStrength)
	}
	if args.strengthBytes == 0 {
		args.strengthBytes = maxStrength
	}

	if args.reseedInterval > hmacdrbg.ReseedInterval {
		return nil, InvalidReseedIntervalError(hmacdrbg.ReseedInterval)
	}

	if uint64(len(args.personalizationString)) > hmacdrbg.MaxPersonalizationStringLength {
		return nil, PersonalizationStringIsTooLongError(hmacdrbg.MaxPersonalizationStringLength)
	}

	securityStrengthBytes := min(
		hmacdrbg.GetSecurityStrengthBytes(args.strengthBytes),
		maxHashSecurityStrength(hashOutLen),
	)

	ownsEntropy := args.entropy == nil
	if ownsEntropy {
		args.entropy = newEntropy(rand.Reader, securityStrengthBytes)
	}

	var err error
	entropyInput := args.initEntropy
	if args.initEntropyIsSet {
		err = checkEntropy(entropyInput, securityStrengthBytes, hmacdrbg.MaxLength)
	} else {
		entropyInput, err = readEntropy(args.entropy, securityStrengthBytes, hmacdrbg.MaxLength)
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

	drbg := &hmacDRGB{
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

func (h *hmacDRGB) Read(dst []byte) (n int, err error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if err = h.z.ReturnIfZeroized(); err != nil {
		return 0, err
	}

	remain := len(dst)
	for remain > 0 {
		toRead := min(remain, hmacdrbg.MaxNoOfBytesPerRequest)

		_, err = h.generate(dst[len(dst)-remain:len(dst)-remain+toRead], nil)
		if err != nil {
			return len(dst) - remain, err
		}

		remain -= toRead
	}

	return len(dst), nil
}

func (h *hmacDRGB) Generate(dst []byte, additionalInput []byte) (n int, err error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.generate(dst, additionalInput)
}

func (h *hmacDRGB) generate(dst []byte, additionalInput []byte) (n int, err error) {
	if err = h.z.ReturnIfZeroized(); err != nil {
		return 0, err
	}

	if len(dst) > hmacdrbg.MaxNoOfBytesPerRequest {
		return 0, TooManyBytesRequestedError(hmacdrbg.MaxNoOfBytesPerRequest)
	}

	if uint64(len(additionalInput)) > hmacdrbg.MaxAdditionalInputLength {
		return 0, AdditionalInputIsTooLongError(hmacdrbg.MaxAdditionalInputLength)
	}

	// The state consumes reseed entropy synchronously. Retain the owned buffer
	// only until Generate returns, without changing custom source buffers.
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

func (h *hmacDRGB) ReseedWithEntropy(additionalInput, entropyInput []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if err := h.z.ReturnIfZeroized(); err != nil {
		return err
	}

	if uint64(len(additionalInput)) > hmacdrbg.MaxAdditionalInputLength {
		return AdditionalInputIsTooLongError(hmacdrbg.MaxAdditionalInputLength)
	}

	if err := checkEntropy(entropyInput, h.securityStrength, hmacdrbg.MaxLength); err != nil {
		return err
	}

	h.state.Reseed(entropyInput, additionalInput)
	return nil
}

func (h *hmacDRGB) Reseed(additionalInput []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if err := h.z.ReturnIfZeroized(); err != nil {
		return err
	}

	if uint64(len(additionalInput)) > hmacdrbg.MaxAdditionalInputLength {
		return AdditionalInputIsTooLongError(hmacdrbg.MaxAdditionalInputLength)
	}

	entropyInput, err := h.getEntropy()
	if err != nil {
		return err
	}
	defer clear(entropyInput)

	h.state.Reseed(entropyInput, additionalInput)
	return nil
}
