package cavp

import (
	"hash"
	"testing"

	"github.com/stretchr/testify/require"
)

type DrbgHmacKAT struct {
	PredictionResistance     bool `cavp:"PredictionResistance,header"`
	EntropyInputLen          int  `cavp:"EntropyInputLen,header"`
	NonceLen                 int  `cavp:"NonceLen,header"`
	PersonalizationStringLen int  `cavp:"PersonalizationStringLen,header"`
	AdditionalInputLen       int  `cavp:"AdditionalInputLen,header"`
	ReturnedBitsLen          int  `cavp:"ReturnedBitsLen,header"`

	Count                 int       `cavp:"COUNT"`
	EntropyInput          []byte    `cavp:"EntropyInput,hex"`
	Nonce                 []byte    `cavp:"Nonce,hex"`
	PersonalizationString []byte    `cavp:"PersonalizationString,hex"`
	EntropyInputReseed    []byte    `cavp:"EntropyInputReseed,hex,optional"`
	AdditionalInputReseed []byte    `cavp:"AdditionalInputReseed,hex,optional"`
	EntropyInputPR        [2][]byte `cavp:"EntropyInputPR,hex,optional"`
	AdditionalInput       [2][]byte `cavp:"AdditionalInput,hex"`
	ReturnedBits          []byte    `cavp:"ReturnedBits,hex,optional"`
}

func ProcessDrbgHmacKAT(input, output string, newHash func() hash.Hash) error {
	return processDrbgKATFile(input, output, func(v DrbgHashKAT) ([]byte, error) {
		return drbgHMAC(DrbgHmacKAT(v).vector(), newHash)
	})
}

func TestDrbgHmacKAT(
	t testing.TB,
	input string,
	newHash func() hash.Hash,
) {
	t.Helper()
	Test(t, input, func(v DrbgHmacKAT) {
		got, err := drbgHMAC(v.vector(), newHash)
		require.NoError(t, err)
		require.Equal(t, v.ReturnedBits, got, "COUNT = %d", v.Count)
	})
}

func (v DrbgHmacKAT) vector() drbgVector {
	return drbgVector{
		EntropyInput:          v.EntropyInput,
		Nonce:                 v.Nonce,
		PersonalizationString: v.PersonalizationString,
		EntropyInputReseed:    v.EntropyInputReseed,
		AdditionalInputReseed: v.AdditionalInputReseed,
		EntropyInputPR:        v.EntropyInputPR,
		AdditionalInput:       v.AdditionalInput,
		ReturnedBitsLen:       v.ReturnedBitsLen,
		PredictionResistance:  v.PredictionResistance,
	}
}
