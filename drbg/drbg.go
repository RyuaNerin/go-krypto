// drbg implemented deterministic random bit generator
// as defined in TTAK.KO-12.0331, TTAK.KO-12.0332, TTAK.KO-12.0189/R1
package drbg

import (
	"errors"
	"io"
)

type DRBG interface {
	io.Reader
	io.Closer

	Generate(dst []byte, additionalInput []byte) (n int, err error)
	Reseed(additionalInput []byte) error
}

var ErrUninstantiated = errors.New("krypto/drbg: state is uninstantiated")

type DRBGOption func(args *args)

type args struct {
	nonce                       []byte
	personalizationString       []byte
	requirePredictionResistance bool
	requireDerivationFunction   bool

	// ctr
	ctrLen         int
	reseedInterval int
}

func WithNonce(nonce []byte) DRBGOption {
	return func(args *args) {
		args.nonce = nonce
	}
}
func WithPersonalizationString(personalizationString []byte) DRBGOption {
	return func(args *args) {
		args.personalizationString = personalizationString
	}
}
func WithPredictionResistance(require bool) DRBGOption {
	return func(args *args) {
		args.requirePredictionResistance = require
	}
}
func WithDerivationFunction(require bool) DRBGOption {
	return func(args *args) {
		args.requireDerivationFunction = require
	}
}