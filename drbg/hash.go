package drbg

import (
	"errors"
	"hash"
	"io"

	"github.com/RyuaNerin/go-krypto/internal/drbg/hashdrbg"
)

type hashDRGB struct {
	state                       *hashdrbg.State
	entropy                     *entropy
	closed                      bool
	requirePredictionResistance bool
}

// New Deterministic Random Bit Generator based on Hash Function
//
// TTAK.KO-12.0331
func NewHashDRGB(
	rand io.Reader,
	h hash.Hash,
	strengthBits int,
	nonce []byte,
	seed []byte,
	requirePredictionResistance bool,
) (DRBG, error) {
	if strengthBits > h.BlockSize()*8 {
		return nil, errors.New("krypto/drbg: invalid strength")
	}
	if len(seed) > hashdrbg.MaxPersonalizationStringLength {
		return nil, errors.New("krypto/drbg: personalization_string too long")
	}

	security_strength := hashdrbg.GetSecurityStrength(strengthBits)

	entropy := newEntropy(rand, security_strength, security_strength*3)
	entropyInput, err := entropy.Get()
	if err != nil {
		return nil, err
	}

	state := hashdrbg.Instantiate_Hash_DRBG(
		h,
		strengthBits,
		requirePredictionResistance,
		entropyInput,
		nonce,
		seed,
	)

	return &hashDRGB{
		state:                       state,
		entropy:                     entropy,
		requirePredictionResistance: requirePredictionResistance,
	}, nil
}

func (h *hashDRGB) Read(dst []byte) (n int, err error) {
	_, err = h.Generate(dst, len(dst), nil)
	if err != nil {
		return 0, err
	}
	return len(dst), nil
}

func (h *hashDRGB) Generate(dst []byte, length int, additionalInput []byte) ([]byte, error) {
	if h.closed {
		return nil, ErrUninstantiated
	}

	if length > hashdrbg.MaxNoOfBitsPerRequest/8 {
		return nil, errors.New("krypto/drbg: too many bits requested")
	}

	if len(additionalInput) > hashdrbg.MaxAdditionalInputLength {
		return nil, errors.New("krypto/drbg: additional_input too long")
	}

	return h.state.Generate_Hash_DRBG(dst, length*8, h.entropy.Get, additionalInput)
}

func (h *hashDRGB) Reseed(additionalInput []byte) error {
	if h.closed {
		return ErrUninstantiated
	}

	entropyInput, err := h.entropy.Get()
	if err != nil {
		return err
	}

	h.state.Reseed_Hash_DRBG(entropyInput, additionalInput)
	return nil
}

func (h *hashDRGB) Close() error {
	h.closed = true
	h.entropy.buf = nil
	h.state.Uninstantiate_Hash_DRBG()
	return nil
}
