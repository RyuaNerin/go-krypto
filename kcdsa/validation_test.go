package kcdsa_test

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/RyuaNerin/go-krypto/kcdsa"
	"github.com/stretchr/testify/require"
)

func TestValidatePublicKey(t *testing.T) {
	tc := testCaseTTAK[0]
	valid := kcdsa.PublicKey{
		Parameters: kcdsa.Parameters{P: tc.P, Q: tc.Q, G: tc.G},
		Y:          tc.Y,
	}
	require.NoError(t, valid.Validate())

	testCases := []struct {
		Name   string
		Mutate func(*kcdsa.PublicKey)
	}{
		{"negative P", func(k *kcdsa.PublicKey) { k.P.Neg(k.P) }},
		{"negative Q", func(k *kcdsa.PublicKey) { k.Q.Neg(k.Q) }},
		{"composite P", func(k *kcdsa.PublicKey) { k.P.SetBit(k.P, 0, 0) }},
		{"composite Q", func(k *kcdsa.PublicKey) { k.Q.SetBit(k.Q, 0, 0) }},
		{"G outside subgroup", func(k *kcdsa.PublicKey) { k.G.Sub(k.P, k.G) }},
		{"Y outside subgroup", func(k *kcdsa.PublicKey) { k.Y.Sub(k.P, k.Y) }},
	}

	for _, tc := range testCases {
		name := tc.Name
		mutate := tc.Mutate

		t.Run(name, func(t *testing.T) {
			key := valid.Clone()
			mutate(&key)
			require.Error(t, key.Validate())
		})
	}

	var nilKey *kcdsa.PublicKey
	require.Error(t, nilKey.Validate())

	var nilParameters *kcdsa.Parameters
	require.Error(t, nilParameters.Validate())
}

func TestRegenerateParametersRejects(t *testing.T) {
	tc0 := testCaseTTAK[0]

	testCases := []struct {
		Name   string
		Mutate func(*kcdsa.GenerationParameters)
	}{
		{"short J", func(p *kcdsa.GenerationParameters) { p.J = big.NewInt(75) }},
		{"composite J", func(p *kcdsa.GenerationParameters) { p.J.SetBit(p.J, 0, 0) }},
		{"seed mismatch", func(p *kcdsa.GenerationParameters) { p.Seed[0] ^= 1 }},
		{"zero count", func(p *kcdsa.GenerationParameters) { p.Count = 0 }},
		{"negative count", func(p *kcdsa.GenerationParameters) { p.Count = -1 }},
		{"excess count", func(p *kcdsa.GenerationParameters) { p.Count = 1<<24 + 1 }},
	}

	valid := kcdsa.GenerationParameters{
		J:     tc0.J,
		Seed:  tc0.Seedb,
		Count: tc0.Count,
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			params := kcdsa.Parameters{GenParameters: valid.Clone()}
			tc.Mutate(&params.GenParameters)
			require.ErrorIs(t,
				kcdsa.RegenerateParameters(&params, rand.Reader, tc0.Sizes),
				kcdsa.ErrInvalidGenerationParameters)
			require.Nil(t, params.P)
			require.Nil(t, params.Q)
			require.Nil(t, params.G)
		})
	}
}
