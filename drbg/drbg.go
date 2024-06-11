// drbg implemented deterministic random bit generator
// as defined in TTAK.KO-12.0331, TTAK.KO-12.0332, TTAK.KO-12.0189/R1, NIST SP 800-90A.
package drbg

import (
	"io"
)

type DRBG interface {
	io.Reader
	io.Closer

	Generate(dst []byte, additionalInput []byte) (n int, err error)

	Reseed(additionalInput []byte) error
	ReseedWithEntropy(additionalInput, entropyInput []byte) error
}

type Option func(args *args)

type args struct {
	nonce                       []byte
	personalizationString       []byte
	requirePredictionResistance bool
	requireDerivationFunction   bool
	reseedInterval              uint64
	strengthBits                int

	entropy          Entropy
	initEntropy      []byte
	initEntropyIsSet bool

	// ctr
	ctrLen int
}

func WithNonce(nonce []byte) Option {
	return func(args *args) {
		args.nonce = nonce
	}
}

func WithPersonalizationString(personalizationString []byte) Option {
	return func(args *args) {
		args.personalizationString = personalizationString
	}
}

func WithPredictionResistance(require bool) Option {
	return func(args *args) {
		args.requirePredictionResistance = require
	}
}

func WithDerivationFunction(require bool) Option {
	return func(args *args) {
		args.requireDerivationFunction = require
	}
}

func WithStrengthBits(value int) Option {
	return func(args *args) {
		args.strengthBits = value
	}
}

func WithReseedInterval(reseedInterval uint64) Option {
	return func(args *args) {
		args.reseedInterval = reseedInterval
	}
}

func WithCustomEntropy(entropy Entropy) Option {
	return func(args *args) {
		args.entropy = entropy
	}
}

func WithInstantiateEntropy(entropyInput []byte) Option {
	return func(args *args) {
		args.initEntropy = entropyInput
		args.initEntropyIsSet = true
	}
}
