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
	options ...DRBGOption,
) (DRBG, error) {
	var args args
	for _, v := range options {
		v(&args)
	}

	if strengthBits > h.BlockSize()*8 {
		return nil, errors.New("krypto/drbg: invalid strength")
	}
	if len(args.personalizationString) > hashdrbg.MaxPersonalizationStringLength {
		return nil, errors.New("krypto/drbg: personalization_string too long")
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
	return h.Generate(dst, nil)
}

func (h *hashDRGB) Generate(dst []byte, additionalInput []byte) (n int, err error) {
	if h.closed {
		return 0, ErrUninstantiated
	}

	if len(dst) > hashdrbg.MaxNoOfBitsPerRequest/8 {
		return 0, errors.New("krypto/drbg: too many bits requested")
	}

	if len(additionalInput) > hashdrbg.MaxAdditionalInputLength {
		return 0, errors.New("krypto/drbg: additional_input too long")
	}

	return len(dst), h.state.Generate_Hash_DRBG(dst, h.entropy.Get, additionalInput)
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
	if h.closed {
		return ErrUninstantiated
	}
	h.closed = true

	h.state.Uninstantiate_Hash_DRBG()

	h.state = nil
	h.entropy = nil
	return nil
}
