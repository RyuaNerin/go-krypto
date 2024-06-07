package drbg

import (
	"errors"
	"hash"
	"io"
	"runtime"

	"github.com/RyuaNerin/go-krypto/internal/drbg/hashdrbg"
)

type hashDRGB struct {
	state   *hashdrbg.State
	entropy *entropy
	closed  bool
}

// New Deterministic Random Bit Generator based on Hash Function
//
// TTAK.KO-12.0331
func NewHashDRGB(
	rand io.Reader,
	h hash.Hash,
	strengthBits int,
	options ...Option,
) (DRBG, error) {
	var args args
	for _, v := range options {
		v(&args)
	}

	if strengthBits > h.BlockSize()*8 {
		return nil, errors.New(msgInvalidStrength)
	}
	if uint64(len(args.personalizationString)) > hashdrbg.MaxPersonalizationStringLength {
		return nil, errors.New(msgPersonalizationStringIsTooLong)
	}

	securityStrength := hashdrbg.GetSecurityStrength(strengthBits)

	entropy := newEntropy(rand, securityStrength/8, securityStrength/8*3)
	entropyInput, err := entropy.Get()
	if err != nil {
		return nil, err
	}

	state := hashdrbg.Instantiate_Hash_DRBG(
		h,
		strengthBits,
		entropyInput,
		args.nonce,
		args.personalizationString,
		args.requirePredictionResistance,
	)

	drbg := &hashDRGB{
		state:   state,
		entropy: entropy,
	}

	runtime.SetFinalizer(drbg, drbg.Close)
	return drbg, nil
}

func (h *hashDRGB) Read(dst []byte) (n int, err error) {
	remain := len(dst)
	for remain > 0 {
		toRead := remain
		if toRead > hashdrbg.MaxNoOfBitsPerRequest/8 {
			toRead = hashdrbg.MaxNoOfBitsPerRequest / 8
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

	if len(dst) > hashdrbg.MaxNoOfBitsPerRequest/8 {
		return 0, errors.New(msgTooManyBitsRequested)
	}

	if uint64(len(additionalInput)) > hashdrbg.MaxAdditionalInputLength {
		return 0, errors.New(msgAdditionalInputIsTooLong)
	}

	return len(dst), h.state.Generate_Hash_DRBG(dst, h.entropy.Get, additionalInput)
}

func (h *hashDRGB) Reseed(additionalInput []byte) error {
	if h.closed {
		return errors.New(msgErrorUninstantiated)
	}

	entropyInput, err := h.entropy.Get()
	if err != nil {
		return err
	}

	h.state.Reseed_Hash_DRBG(entropyInput, additionalInput)
	return nil
}

func (h *hashDRGB) Close() error {
	if h.closed {
		return errors.New(msgErrorUninstantiated)
	}
	h.closed = true

	h.state.Uninstantiate_Hash_DRBG()

	h.state = nil
	h.entropy = nil
	return nil
}
