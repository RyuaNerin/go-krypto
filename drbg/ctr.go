package drbg

import (
	"crypto/cipher"
	"errors"
	"io"
	"runtime"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/drbg/ctrdrbg"
)

type CTRDRBGOption DRBGOption

func WithCTRLength(ctrLen int) CTRDRBGOption {
	return func(args *args) {
		args.ctrLen = ctrLen
	}
}
func WithCTRReseedInterval(reseedInterval int) CTRDRBGOption {
	return func(args *args) {
		args.reseedInterval = reseedInterval
	}
}

type ctrDRGB struct {
	state   *ctrdrbg.State
	entropy *entropy
	closed  bool

	requireDerivationFunction bool
}

// New Deterministic Random Bit Generator based on Block Cipher
//
// TTAK.KO-12.0189/R1
func New(
	rand io.Reader,
	newCipher func(key []byte) (cipher.Block, error),
	keySize int, // Bytes
	strengthBits int,
	options ...CTRDRBGOption,
) (DRBG, error) {
	var args args
	for _, v := range options {
		v(&args)
	}

	b, err := newCipher(make([]byte, keySize))
	if err != nil {
		return nil, err
	}
	seedLen := b.BlockSize() + keySize

	if args.ctrLen == 0 {
		args.ctrLen = b.BlockSize()
	}

	if args.requireDerivationFunction {
		if len(args.personalizationString) > seedLen {
			return nil, errors.New("krypto/drbg: additionalInput too long")
		}
	} else {
		if len(args.personalizationString) > ctrdrbg.MaxAdditionalInputLength {
			return nil, errors.New("krypto/drbg: additionalInput too long")
		}
	}
	if args.ctrLen < 4 || b.BlockSize() < args.ctrLen {
		return nil, errors.New("krypto/drbg: invalid ctrLen")
	}

	var entropy *entropy
	if args.requireDerivationFunction {
		entropy = newEntropy(rand, internal.Bytes(strengthBits), internal.Bytes(strengthBits)*3)
	} else {
		entropy = newEntropy(rand, seedLen, seedLen)
	}

	entropyInput, err := entropy.Get()
	if err != nil {
		return nil, err
	}

	state := ctrdrbg.Instantiate_CTR_DRBG(
		newCipher,
		keySize,
		args.reseedInterval,
		args.ctrLen,
		entropyInput,
		args.nonce,
		args.personalizationString,
		args.requireDerivationFunction,
		args.requirePredictionResistance,
	)

	drbg := &ctrDRGB{
		state:   state,
		entropy: entropy,

		requireDerivationFunction: args.requireDerivationFunction,
	}
	runtime.SetFinalizer(drbg, drbg.Close)

	return drbg, nil
}

func (h *ctrDRGB) Read(dst []byte) (n int, err error) {
	remain := len(dst)
	for remain > 0 {
		toRead := remain
		if toRead > h.state.MaxNoOfBitsPerRequest {
			toRead = h.state.MaxNoOfBitsPerRequest
		}

		_, err = h.Generate(dst[:toRead], nil)
		if err != nil {
			return 0, err
		}

		remain -= toRead
	}

	return len(dst), nil
}

func (h *ctrDRGB) Generate(dst []byte, additionalInput []byte) (n int, err error) {
	if h.closed {
		return 0, ErrUninstantiated
	}

	if len(dst) > h.state.MaxNoOfBitsPerRequest {
		return 0, errors.New("krypto/drbg: too many bits requested")
	}

	if h.requireDerivationFunction {
		if len(additionalInput) > ctrdrbg.MaxAdditionalInputLength {
			return 0, errors.New("krypto/drbg: additionalInput too long")
		}
	} else {
		if len(additionalInput) > h.state.SeedLenByte {
			return 0, errors.New("krypto/drbg: additionalInput too long")
		}
	}

	return len(dst), h.state.Generate_CTR_DRBG(dst, h.entropy.Get, additionalInput)
}

func (h *ctrDRGB) Reseed(additionalInput []byte) error {
	if h.closed {
		return ErrUninstantiated
	}

	if h.requireDerivationFunction {
		if len(additionalInput) > ctrdrbg.MaxAdditionalInputLength {
			return errors.New("krypto/drbg: additionalInput too long")
		}
	} else {
		if len(additionalInput) > h.state.SeedLenByte {
			return errors.New("krypto/drbg: additionalInput too long")
		}
	}

	entropyInput, err := h.entropy.Get()
	if err != nil {
		return err
	}

	h.state.Reseed_CTR_DRBG(entropyInput, additionalInput)
	return nil
}

func (h *ctrDRGB) Close() error {
	if h.closed {
		return ErrUninstantiated
	}
	h.closed = true

	h.state.Uninstantiate_CTR_DRBG()

	h.state = nil
	h.entropy = nil
	return nil
}
