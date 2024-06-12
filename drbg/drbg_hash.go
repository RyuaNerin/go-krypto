package drbg

import (
	"crypto/rand"
	"errors"
	"fmt"
	"hash"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/drbg/ctrdrbg"
	"github.com/RyuaNerin/go-krypto/internal/drbg/hashdrbg"
	"github.com/RyuaNerin/go-krypto/internal/memory"
)

type hashDRGB struct {
	state   *hashdrbg.State
	entropy Entropy
	closed  bool

	securityStrength int
}

// New Deterministic Random Bit Generator based on Hash Function
//
// default entropy: if requireDerivationFunction is true, NewEntropy(rand.Reader, strengthBits, strengthBits*2), otherwise NewEntropy(rand.Reader, blocksize + keysize, blocksize + keysize)
// If strengthBits is unset or smaller than keySize, strengthBits is set to keySize.
//
// TTAK.KO-12.0331
func NewHashDRGB(
	h hash.Hash,
	options ...Option,
) (drbg DRBG, err error) {
	var args args
	for _, v := range options {
		v(&args)
	}

	hashOutLen := h.Size()

	if args.strengthBits > hashOutLen*8 {
		return nil, fmt.Errorf(msgInvalidStrengthFormat, hashOutLen*8)
	}
	if args.strengthBits == 0 {
		args.strengthBits = hashOutLen * 8
	}

	if uint64(len(args.personalizationString)) > hashdrbg.MaxPersonalizationStringLength {
		return nil, fmt.Errorf(msgPersonalizationStringIsTooLongFormat, hashdrbg.MaxPersonalizationStringLength)
	}

	securityStrengthBits := hashdrbg.GetSecurityStrengthBits(args.strengthBits)

	if args.entropy == nil {
		args.entropy = newEntropy(rand.Reader, internal.BitsToBytes(securityStrengthBits))
	}

	entropyInput := args.initEntropy
	if !args.initEntropyIsSet {
		entropyInput, err = args.entropy.Get()
		if err != nil {
			return nil, err
		}
	}

	state := hashdrbg.Instantiate(
		h,
		entropyInput,
		args.nonce,
		args.personalizationString,
		args.requirePredictionResistance,
		args.reseedInterval,
	)

	drbg = &hashDRGB{
		state:            state,
		entropy:          args.entropy,
		securityStrength: internal.BitsToBytes(securityStrengthBits),
	}

	return drbg, nil
}

func (h *hashDRGB) Read(dst []byte) (n int, err error) {
	remain := len(dst)
	for remain > 0 {
		toRead := remain
		if toRead > hashdrbg.MaxNoOfBytesPerRequest {
			toRead = hashdrbg.MaxNoOfBytesPerRequest
		}

		_, err = h.Generate(dst[:toRead], nil)
		if err != nil {
			return 0, err
		}

		remain -= toRead
	}

	return len(dst), nil
}

func (h *hashDRGB) Generate(dst []byte, additionalInput []byte) (n int, err error) {
	if h.closed {
		return 0, errors.New(msgErrorUninstantiated)
	}

	if len(dst) > hashdrbg.MaxNoOfBytesPerRequest {
		return 0, fmt.Errorf(msgTooManyBytesRequestedFormat, hashdrbg.MaxNoOfBytesPerRequest)
	}

	if uint64(len(additionalInput)) > hashdrbg.MaxAdditionalInputLength {
		return 0, fmt.Errorf(msgAdditionalInputIsTooLongFormat, hashdrbg.MaxAdditionalInputLength)
	}

	return len(dst), h.state.Generate(dst, h.entropy.Get, additionalInput)
}

func (h *hashDRGB) ReseedWithEntropy(additionalInput, entropyInput []byte) error {
	if h.closed {
		return errors.New(msgErrorUninstantiated)
	}

	if uint64(len(additionalInput)) > hashdrbg.MaxAdditionalInputLength {
		return fmt.Errorf(msgAdditionalInputIsTooLongFormat, hashdrbg.MaxAdditionalInputLength)
	}

	if len(entropyInput) < h.securityStrength && ctrdrbg.MaxLength < uint64(len(entropyInput)) {
		return fmt.Errorf(msgInvalidEntropyFormat, h.securityStrength, ctrdrbg.MaxLength)
	}

	h.state.Reseed(entropyInput, additionalInput)
	return nil
}

func (h *hashDRGB) Reseed(additionalInput []byte) error {
	if h.closed {
		return errors.New(msgErrorUninstantiated)
	}

	if uint64(len(additionalInput)) > hashdrbg.MaxAdditionalInputLength {
		return fmt.Errorf(msgAdditionalInputIsTooLongFormat, hashdrbg.MaxAdditionalInputLength)
	}

	entropyInput, err := h.entropy.Get()
	if err != nil {
		return err
	}

	h.state.Reseed(entropyInput, additionalInput)
	return nil
}

func (h *hashDRGB) Close() error {
	if h.closed {
		return errors.New(msgErrorUninstantiated)
	}
	h.closed = true

	memory.MemclrI(h.state)

	h.state = nil
	h.entropy = nil
	return nil
}
