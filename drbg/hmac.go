package drbg

import (
	"errors"
	"hash"
	"io"

	"github.com/RyuaNerin/go-krypto/internal/drbg/hashdrbg"
	"github.com/RyuaNerin/go-krypto/internal/drbg/hmacdrbg"
)

type hmacDRGB struct {
	state                       *hmacdrbg.State
	entropy                     *entropy
	closed                      bool
	requirePredictionResistance bool
}

// New Deterministic Random Bit Generator based on HMAC
//
// TTAK.KO-12.0332
func NewHMACDRGB(
	rand io.Reader,
	h func() hash.Hash,
	strengthBits int,
	nonce []byte,
	seed []byte,
	requirePredictionResistance bool,
) (DRBG, error) {
	outlen := h().Size()

	if strengthBits > outlen*8 {
		return nil, errors.New("krypto/drbg: invalid strength")
	}
	if len(seed) > hmacdrbg.MaxPersonalizationStringLength {
		return nil, errors.New("krypto/drbg: personalization_string too long")
	}

	security_strength := hmacdrbg.GetSecurityStrength(strengthBits)

	entropy := newEntropy(rand, security_strength, security_strength*3)
	entropyInput, err := entropy.Get()
	if err != nil {
		return nil, err
	}

	state := hmacdrbg.Instantiate_HMAC_DRBG(
		h,
		strengthBits,
		requirePredictionResistance,
		entropyInput,
		nonce,
		seed,
	)

	return &hmacDRGB{
		state:                       state,
		entropy:                     entropy,
		requirePredictionResistance: requirePredictionResistance,
	}, nil
}

func (h *hmacDRGB) Read(dst []byte) (n int, err error) {
	_, err = h.Generate(dst, len(dst), nil)
	if err != nil {
		return 0, err
	}
	return len(dst), nil
}

func (h *hmacDRGB) Generate(dst []byte, length int, additionalInput []byte) ([]byte, error) {
	if h.closed {
		return nil, ErrUninstantiated
	}

	if length > hashdrbg.MaxNoOfBitsPerRequest/8 {
		return nil, errors.New("krypto/drbg: too many bits requested")
	}

	if len(additionalInput) > hashdrbg.MaxAdditionalInputLength {
		return nil, errors.New("krypto/drbg: additional_input too long")
	}

	return h.state.Generate_HMAC_DRBG(dst, length*8, h.entropy.Get, additionalInput)
}

func (h *hmacDRGB) Reseed(additionalInput []byte) error {
	if h.closed {
		return ErrUninstantiated
	}

	entropyInput, err := h.entropy.Get()
	if err != nil {
		return err
	}

	h.state.Reseed_HMAC_DRBG(entropyInput, additionalInput)
	return nil
}

func (h *hmacDRGB) Close() error {
	h.closed = true
	h.entropy.buf = nil
	h.state.Uninstantiate_HMAC_DRBG()
	return nil
}
