// Package drbg implements deterministic random bit generators as defined in
// TTAK.KO-12.0331, TTAK.KO-12.0332, TTAK.KO-12.0189, and NIST SP 800-90A.
//
// The DRBGs returned by this package are safe for concurrent use by multiple
// goroutines.
package drbg

// Option configures an optional DRBG setting. Required cryptographic primitives
// and key sizes are supplied to the New*DRBG or New*Builder constructors.
type Option func(args *args)

type args struct {
	nonce                       []byte
	nonceIsSet                  bool
	personalizationString       []byte
	requirePredictionResistance bool
	requireDerivationFunction   bool
	reseedInterval              uint64
	strengthBytes               int

	entropy          Entropy
	initEntropy      []byte
	initEntropyIsSet bool

	// ctr
	ctrLen int
}

// WithNonce sets the nonce used to instantiate the DRBG. If omitted, Hash/HMAC
// and CTR with a derivation function generate a fresh random nonce. An explicit
// empty or short nonce requires extra entropy to supply the missing nonce bytes.
// CTR without a derivation function ignores the nonce.
func WithNonce(nonce []byte) Option {
	return func(args *args) {
		args.nonce = nonce
		args.nonceIsSet = true
	}
}

// WithPersonalizationString sets the personalization string used to instantiate the DRBG.
func WithPersonalizationString(personalizationString []byte) Option {
	return func(args *args) {
		args.personalizationString = personalizationString
	}
}

// WithPredictionResistance requires prediction resistance when the DRBG is instantiated
// and reseeded.
func WithPredictionResistance(require bool) Option {
	return func(args *args) {
		args.requirePredictionResistance = require
	}
}

// WithDerivationFunction requires the use of the derivation function when the DRBG
// is instantiated and reseeded.
func WithDerivationFunction(require bool) Option {
	return func(args *args) {
		args.requireDerivationFunction = require
	}
}

// WithStrengthBytes sets the desired strength of the generated output, in bytes.
func WithStrengthBytes(value int) Option {
	return func(args *args) {
		args.strengthBytes = value
	}
}

// WithReseedInterval sets the maximum number of Generate calls between reseed operations.
func WithReseedInterval(reseedInterval uint64) Option {
	return func(args *args) {
		args.reseedInterval = reseedInterval
	}
}

// WithCustomEntropy sets the entropy source used by the DRBG. The caller owns
// the source; closing the DRBG does not zeroize it. See Entropy for the contract
// when a source is shared between DRBGs.
func WithCustomEntropy(entropy Entropy) Option {
	return func(args *args) {
		args.entropy = entropy
	}
}

// WithInstantiateEntropy sets the entropy input used to instantiate the DRBG.
func WithInstantiateEntropy(entropyInput []byte) Option {
	return func(args *args) {
		args.initEntropy = entropyInput
		args.initEntropyIsSet = true
	}
}
