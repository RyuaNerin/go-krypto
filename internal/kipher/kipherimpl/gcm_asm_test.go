//go:build gc && (arm64 || amd64) && !noasm && !noasm_krypto && !noasm_krypto_kipher

package kipherimpl_test

import (
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/kipher/kipherimpl"
	"github.com/RyuaNerin/go-krypto/internal/randutil"
	"github.com/RyuaNerin/go-krypto/lea"

	"github.com/stretchr/testify/require"
)

func TestGCMAuthAsm(t *testing.T) {
	if !kipherimpl.SupportsGFMUL {
		t.Skip("PCLMULQDQ or PMULL not available")
		return
	}

	rnd, _ := randutil.NewTestRand()

	var tagMaskA, tagMaskB [kipherimpl.GCMBlockSize]byte
	outA, outG := make([]byte, 16), make([]byte, 16)

	key := make([]byte, 16)
	ciphertext := make([]byte, 16*16+15)
	additionalData := make([]byte, 16*16+1)

	_, err := rnd.Read(key)
	require.NoError(t, err)
	_, err = rnd.Read(ciphertext)
	require.NoError(t, err)
	_, err = rnd.Read(additionalData)
	require.NoError(t, err)

	block, err := lea.NewCipher(key)
	require.NoError(t, err)
	kb := internal.ToBlock8(block)

	var gcmA kipherimpl.GCM
	var gcmG kipherimpl.GCM

	kipherimpl.GCMInitGo(&gcmA, kb)
	kipherimpl.GCMInitAsm(&gcmG, kb)

	kipherimpl.GCMAuthGo(&gcmA, outA, ciphertext, additionalData, &tagMaskA)
	kipherimpl.GCMAuthAsm(&gcmG, outG, ciphertext, additionalData, &tagMaskB)

	require.Equal(t, outA, outG)
}

func BenchmarkGCMAuthAsm(b *testing.B) {
	if !kipherimpl.SupportsGFMUL {
		b.Skip("PCLMULQDQ or PMULL not available")
		return
	}

	bench := func(
		blocks int,
		gcmAuth func(g *kipherimpl.GCM, out, ciphertext, additionalData []byte, tagMask *[kipherimpl.GCMTagSize]byte),
	) func(b *testing.B) {
		return func(b *testing.B) {
			rnd, _ := randutil.NewBenchRand()

			var g kipherimpl.GCM

			block, _ := lea.NewCipher(make([]byte, 16))
			kb := internal.ToBlock8(block)

			kipherimpl.GCMInitGo(&g, kb)

			var tagMask [kipherimpl.GCMBlockSize]byte
			out := make([]byte, 16)

			ciphertext := make([]byte, kipherimpl.GCMBlockSize*blocks)

			_, err := rnd.Read(ciphertext)
			require.NoError(b, err)

			for b.Loop() {
				gcmAuth(&g, out, ciphertext, nil, &tagMask)
				copy(ciphertext, out)
			}
		}
	}

	b.Run("generic-8", bench(8, kipherimpl.GCMAuthGo))
	b.Run("generic-1K", bench(1*1024-5, kipherimpl.GCMAuthGo))
	b.Run("generic-8K", bench(8*1024-5, kipherimpl.GCMAuthGo))

	b.Run("assembly-8", bench(8, kipherimpl.GCMAuthAsm))
	b.Run("assembly-1K", bench(1*1024-5, kipherimpl.GCMAuthAsm))
	b.Run("assembly-8K", bench(8*1024-5, kipherimpl.GCMAuthAsm))
}
