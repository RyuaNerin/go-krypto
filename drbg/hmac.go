package drbg

import (
	"errors"
	"hash"
	"io"
	"runtime"

	"github.com/RyuaNerin/go-krypto/internal/drbg/hmacdrbg"
)

type hmacDRGB struct {
	state   *hmacdrbg.State
	entropy *entropy
	closed  bool
}

// New Deterministic Random Bit Generator based on HMAC
//
// TTAK.KO-12.0332
func NewHMACDRGB(
	rand io.Reader,
	h func() hash.Hash,
	strengthBits int,
	options ...DRBGOption,
) (DRBG, error) {
	var args args
	for _, v := range options {
		v(&args)
	}

	outlen := h().Size()

	if strengthBits > outlen*8 {
		return nil, errors.New("krypto/drbg: invalid strength")
	}
	if uint64(len(args.personalizationString)) > hmacdrbg.MaxPersonalizationStringLength {
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
		entropyInput,
		args.nonce,
		args.personalizationString,
		args.requirePredictionResistance,
	)

	drbg := &hmacDRGB{
		state:   state,
		entropy: entropy,
	}
	runtime.SetFinalizer(drbg, drbg.Close)

	return drbg, nil
}

func (h *hmacDRGB) Read(dst []byte) (n int, err error) {
	remain := len(dst)
	for remain > 0 {
		toRead := remain
		if toRead > hmacdrbg.MaxNoOfBitsPerRequest/8 {
			toRead = hmacdrbg.MaxNoOfBitsPerRequest / 8
		}

		_, err = h.Generate(dst[:toRead], nil)
		if err != nil {
			return 0, err
		}

		remain -= toRead
	}

	return len(dst), nil
}

func (h *hmacDRGB) Generate(dst []byte, additionalInput []byte) (n int, err error) {
	if h.closed {
		return 0, ErrUninstantiated
	}

	if len(dst) > hmacdrbg.MaxNoOfBitsPerRequest/8 {
		return 0, errors.New("krypto/drbg: too many bits requested")
	}

	if uint64(len(additionalInput)) > hmacdrbg.MaxAdditionalInputLength {
		return 0, errors.New("krypto/drbg: additional_input too long")
	}

	return len(dst), h.state.Generate_HMAC_DRBG(dst, h.entropy.Get, additionalInput)
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
	if h.closed {
		return ErrUninstantiated
	}
	h.closed = true

	h.state.Uninstantiate_HMAC_DRBG()

	h.state = nil
	h.entropy = nil
	return nil
}
