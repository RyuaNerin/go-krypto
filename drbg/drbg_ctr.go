package drbg

import (
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/drbg/ctrdrbg"
)

func WithCTRLength(ctrLen int) Option {
	return func(args *args) {
		args.ctrLen = ctrLen
	}
}

type ctrDRGB struct {
	state   *ctrdrbg.State
	entropy Entropy
	closed  bool

	securityStrength          int
	requireDerivationFunction bool
}

// NewCTRDRBG Deterministic Random Bit Generator based on Block Cipher
//
// default entropy: if requireDerivationFunction is true, NewEntropy(rand.Reader, strengthBits), otherwise NewEntropy(rand.Reader, blocksize + keysize)
// If strengthBits is unset or smaller than keySize, strengthBits is set to keySize.
//
// TTAK.KO-12.0189/R1
func NewCTRDRBG(
	newCipher func(key []byte) (cipher.Block, error),
	keySize int, // Bytes
	options ...Option,
) (DRBG, error) {
	var args args
	for _, v := range options {
		v(&args)
	}

	b, err := newCipher(make([]byte, keySize))
	if err != nil {
		return nil, err
	}
	blockSize := b.BlockSize()
	seedLen := blockSize + keySize

	if args.requireDerivationFunction {
		if len(args.personalizationString) > seedLen {
			return nil, fmt.Errorf(msgAdditionalInputIsTooLongFormat, seedLen)
		}
	} else {
		if uint64(len(args.personalizationString)) > ctrdrbg.MaxAdditionalInputLength {
			return nil, fmt.Errorf(msgAdditionalInputIsTooLongFormat, ctrdrbg.MaxAdditionalInputLength)
		}
	}

	if args.ctrLen == 0 {
		args.ctrLen = blockSize
	}
	if args.ctrLen < 4 || blockSize < args.ctrLen {
		return nil, fmt.Errorf(msgInvalidCtrLengthFormat2, 4, blockSize)
	}

	if args.strengthBits < keySize {
		args.strengthBits = keySize
	}

	if args.entropy == nil {
		if args.requireDerivationFunction {
			args.entropy = newEntropy(rand.Reader, internal.BitsToBytes(args.strengthBits))
		} else {
			args.entropy = newEntropy(rand.Reader, seedLen)
		}
	}

	entropyInput := args.initEntropy
	if !args.initEntropyIsSet {
		entropyInput, err = args.entropy.Get()
		if err != nil {
			return nil, err
		}
	}

	state := ctrdrbg.Instantiate(
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
		state:            state,
		entropy:          args.entropy,
		securityStrength: args.strengthBits,

		requireDerivationFunction: args.requireDerivationFunction,
	}

	return drbg, nil
}

func (h *ctrDRGB) Read(dst []byte) (n int, err error) {
	remain := len(dst)
	for remain > 0 {
		toRead := remain
		if toRead > h.state.MaxNoOfBytesPerRequest {
			toRead = h.state.MaxNoOfBytesPerRequest
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
		return 0, errors.New(msgErrorUninstantiated)
	}

	if len(dst) > h.state.MaxNoOfBytesPerRequest {
		return 0, fmt.Errorf(msgTooManyBytesRequestedFormat, h.state.MaxNoOfBytesPerRequest)
	}

	if h.requireDerivationFunction {
		if uint64(len(additionalInput)) > ctrdrbg.MaxAdditionalInputLength {
			return 0, fmt.Errorf(msgAdditionalInputIsTooLongFormat, ctrdrbg.MaxAdditionalInputLength)
		}
	} else {
		if len(additionalInput) > h.state.SeedLenByte {
			return 0, fmt.Errorf(msgAdditionalInputIsTooLongFormat, h.state.SeedLenByte)
		}
	}

	err = h.state.Generate(dst, h.entropy.Get, additionalInput)
	if err != nil {
		return 0, err
	}
	return len(dst), nil
}

func (h *ctrDRGB) ReseedWithEntropy(additionalInput, entropyInput []byte) error {
	if h.closed {
		return errors.New(msgErrorUninstantiated)
	}

	if h.requireDerivationFunction {
		if uint64(len(additionalInput)) > ctrdrbg.MaxAdditionalInputLength {
			return fmt.Errorf(msgAdditionalInputIsTooLongFormat, ctrdrbg.MaxAdditionalInputLength)
		}
	} else {
		if len(additionalInput) > h.state.SeedLenByte {
			return fmt.Errorf(msgAdditionalInputIsTooLongFormat, h.state.SeedLenByte)
		}
	}

	if len(entropyInput) < h.securityStrength && ctrdrbg.MaxLength < uint64(len(entropyInput)) {
		return fmt.Errorf(msgInvalidEntropyFormat, h.securityStrength, ctrdrbg.MaxLength)
	}

	h.state.Reseed(entropyInput, additionalInput)
	return nil
}

func (h *ctrDRGB) Reseed(additionalInput []byte) error {
	if h.closed {
		return errors.New(msgErrorUninstantiated)
	}

	if h.requireDerivationFunction {
		if uint64(len(additionalInput)) > ctrdrbg.MaxAdditionalInputLength {
			return fmt.Errorf(msgAdditionalInputIsTooLongFormat, ctrdrbg.MaxAdditionalInputLength)
		}
	} else {
		if len(additionalInput) > h.state.SeedLenByte {
			return fmt.Errorf(msgAdditionalInputIsTooLongFormat, h.state.SeedLenByte)
		}
	}

	entropyInput, err := h.entropy.Get()
	if err != nil {
		return err
	}

	h.state.Reseed(entropyInput, additionalInput)
	return nil
}

func (h *ctrDRGB) Close() error {
	if h.closed {
		return errors.New(msgErrorUninstantiated)
	}
	h.closed = true

	h.state.Uninstantiate()

	h.state = nil
	h.entropy = nil
	return nil
}
