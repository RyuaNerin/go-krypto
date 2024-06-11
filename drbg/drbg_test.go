package drbg_test

import (
	"bufio"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"hash"
	"io"
	"testing"

	"github.com/RyuaNerin/go-krypto/drbg"
	"github.com/RyuaNerin/go-krypto/lea"
)

var rnd = bufio.NewReaderSize(rand.Reader, 1<<10)

const benchDstLength = 256

func BenchmarkCTRDRBG(b *testing.B) {
	b.Run("LEA-128/DF=0/PR=0/PS=0/AI=0", benchCTRDRBG(lea.NewCipher, 128/8, 0, 0, 0, 0))
	b.Run("LEA-128/DF=0/PR=0/PS=0/AI=1", benchCTRDRBG(lea.NewCipher, 128/8, 0, 0, 0, 1))
	b.Run("LEA-128/DF=0/PR=0/PS=1/AI=0", benchCTRDRBG(lea.NewCipher, 128/8, 0, 0, 1, 0))
	b.Run("LEA-128/DF=0/PR=0/PS=1/AI=1", benchCTRDRBG(lea.NewCipher, 128/8, 0, 0, 1, 1))
	b.Run("LEA-128/DF=0/PR=1/PS=0/AI=0", benchCTRDRBG(lea.NewCipher, 128/8, 0, 1, 0, 0))
	b.Run("LEA-128/DF=0/PR=1/PS=0/AI=1", benchCTRDRBG(lea.NewCipher, 128/8, 0, 1, 0, 1))
	b.Run("LEA-128/DF=0/PR=1/PS=1/AI=0", benchCTRDRBG(lea.NewCipher, 128/8, 0, 1, 1, 0))
	b.Run("LEA-128/DF=0/PR=1/PS=1/AI=1", benchCTRDRBG(lea.NewCipher, 128/8, 0, 1, 1, 1))
	b.Run("LEA-128/DF=1/PR=0/PS=0/AI=0", benchCTRDRBG(lea.NewCipher, 128/8, 1, 0, 0, 0))
	b.Run("LEA-128/DF=1/PR=0/PS=0/AI=1", benchCTRDRBG(lea.NewCipher, 128/8, 1, 0, 0, 1))
	b.Run("LEA-128/DF=1/PR=0/PS=1/AI=0", benchCTRDRBG(lea.NewCipher, 128/8, 1, 0, 1, 0))
	b.Run("LEA-128/DF=1/PR=0/PS=1/AI=1", benchCTRDRBG(lea.NewCipher, 128/8, 1, 0, 1, 1))
	b.Run("LEA-128/DF=1/PR=1/PS=0/AI=0", benchCTRDRBG(lea.NewCipher, 128/8, 1, 1, 0, 0))
	b.Run("LEA-128/DF=1/PR=1/PS=0/AI=1", benchCTRDRBG(lea.NewCipher, 128/8, 1, 1, 0, 1))
	b.Run("LEA-128/DF=1/PR=1/PS=1/AI=0", benchCTRDRBG(lea.NewCipher, 128/8, 1, 1, 1, 0))
	b.Run("LEA-128/DF=1/PR=1/PS=1/AI=1", benchCTRDRBG(lea.NewCipher, 128/8, 1, 1, 1, 1))
}

func BenchmarkHashDRBG(b *testing.B) {
	b.Run("SHA-256/PR=0/RI=1/PS=1/AI=1", benchHashDRBG(sha256.New, 0, 1, 1, 1))
	b.Run("SHA-256/PR=0/RI=1/PS=0/AI=1", benchHashDRBG(sha256.New, 0, 1, 0, 1))
	b.Run("SHA-256/PR=0/RI=1/PS=1/AI=0", benchHashDRBG(sha256.New, 0, 1, 1, 0))
	b.Run("SHA-256/PR=0/RI=1/PS=0/AI=0", benchHashDRBG(sha256.New, 0, 1, 0, 0))
	b.Run("SHA-256/PR=0/RI=2/PS=1/AI=1", benchHashDRBG(sha256.New, 0, 2, 1, 1))
	b.Run("SHA-256/PR=0/RI=2/PS=0/AI=1", benchHashDRBG(sha256.New, 0, 2, 0, 1))
	b.Run("SHA-256/PR=0/RI=2/PS=1/AI=0", benchHashDRBG(sha256.New, 0, 2, 1, 0))
	b.Run("SHA-256/PR=0/RI=2/PS=0/AI=0", benchHashDRBG(sha256.New, 0, 2, 0, 0))
	b.Run("SHA-256/PR=1/RI=0/PS=1/AI=1", benchHashDRBG(sha256.New, 1, 0, 1, 1))
	b.Run("SHA-256/PR=1/RI=0/PS=0/AI=1", benchHashDRBG(sha256.New, 1, 0, 0, 1))
	b.Run("SHA-256/PR=1/RI=0/PS=1/AI=0", benchHashDRBG(sha256.New, 1, 0, 1, 0))
	b.Run("SHA-256/PR=1/RI=0/PS=0/AI=0", benchHashDRBG(sha256.New, 1, 0, 0, 0))
}

func BenchmarkHMACDRBG(b *testing.B) {
	b.Run("SHA-256/PR=0/RI=1/PS=1/AI=1", benchHmacDRBG(sha256.New, 0, 1, 1, 1))
	b.Run("SHA-256/PR=0/RI=1/PS=0/AI=1", benchHmacDRBG(sha256.New, 0, 1, 0, 1))
	b.Run("SHA-256/PR=0/RI=1/PS=1/AI=0", benchHmacDRBG(sha256.New, 0, 1, 1, 0))
	b.Run("SHA-256/PR=0/RI=1/PS=0/AI=0", benchHmacDRBG(sha256.New, 0, 1, 0, 0))
	b.Run("SHA-256/PR=0/RI=2/PS=1/AI=1", benchHmacDRBG(sha256.New, 0, 2, 1, 1))
	b.Run("SHA-256/PR=0/RI=2/PS=0/AI=1", benchHmacDRBG(sha256.New, 0, 2, 0, 1))
	b.Run("SHA-256/PR=0/RI=2/PS=1/AI=0", benchHmacDRBG(sha256.New, 0, 2, 1, 0))
	b.Run("SHA-256/PR=0/RI=2/PS=0/AI=0", benchHmacDRBG(sha256.New, 0, 2, 0, 0))
	b.Run("SHA-256/PR=1/RI=0/PS=1/AI=1", benchHmacDRBG(sha256.New, 1, 0, 1, 1))
	b.Run("SHA-256/PR=1/RI=0/PS=0/AI=1", benchHmacDRBG(sha256.New, 1, 0, 0, 1))
	b.Run("SHA-256/PR=1/RI=0/PS=1/AI=0", benchHmacDRBG(sha256.New, 1, 0, 1, 0))
	b.Run("SHA-256/PR=1/RI=0/PS=0/AI=0", benchHmacDRBG(sha256.New, 1, 0, 0, 0))
}

func benchCTRDRBG(
	cipher func(key []byte) (cipher.Block, error),
	keySize int,
	useDF, usePR, usePS, useAI int,
) func(b *testing.B) {
	c, _ := cipher(make([]byte, keySize))
	blockSize := c.BlockSize()

	return func(b *testing.B) {
		dst := make([]byte, benchDstLength)

		var opts []drbg.Option
		{
			nonce := make([]byte, blockSize)
			io.ReadFull(rnd, nonce)
			opts = append(opts, drbg.WithNonce(nonce))
		}
		if useDF == 1 {
			opts = append(opts, drbg.WithDerivationFunction(true))
		}
		if usePR == 1 {
			opts = append(opts, drbg.WithPredictionResistance(true))
		}
		if usePS == 1 {
			personal := make([]byte, blockSize)
			io.ReadFull(rnd, personal)
			opts = append(opts, drbg.WithPersonalizationString(personal))
		}

		state, err := drbg.NewCTRDRBG(
			cipher,
			keySize,
			opts...,
		)
		if err != nil {
			b.Fatalf("failed to create CTR_DRBG: %v", err)
			return
		}

		b.ResetTimer()
		b.ReportAllocs()
		b.SetBytes(benchDstLength)
		if useAI == 1 {
			for i := 0; i < b.N; i++ {
				_, err = state.Read(dst)
				if err != nil {
					b.Fatalf("failed to generate random: %v", err)
					return
				}
			}
		} else {
			for i := 0; i < b.N; i++ {
				_, err = state.Generate(dst[:0], dst[:blockSize])
				if err != nil {
					b.Fatalf("failed to generate random: %v", err)
					return
				}
			}
		}
	}
}

func benchHashDRBG(
	hash func() hash.Hash,
	usePR, useRI, usePS, useAI int,
) func(b *testing.B) {
	hashSize := hash().Size()

	return func(b *testing.B) {
		dst := make([]byte, hashSize)

		var opts []drbg.Option
		{
			nonce := make([]byte, hashSize)
			io.ReadFull(rnd, nonce)
			opts = append(opts, drbg.WithNonce(nonce))
		}
		if usePR == 1 {
			opts = append(opts, drbg.WithPredictionResistance(true))
		}
		if useRI > 0 {
			opts = append(opts, drbg.WithReseedInterval(uint64(useRI)))
		}
		if usePS == 1 {
			personal := make([]byte, hashSize)
			io.ReadFull(rnd, personal)
			opts = append(opts, drbg.WithPersonalizationString(personal))
		}

		state, err := drbg.NewHashDRGB(
			hash(),
			opts...,
		)
		if err != nil {
			b.Fatalf("failed to create CTR_DRBG: %v", err)
			return
		}

		b.ResetTimer()
		b.ReportAllocs()
		b.SetBytes(benchDstLength)
		if useAI == 1 {
			for i := 0; i < b.N; i++ {
				_, err = state.Read(dst)
				if err != nil {
					b.Fatalf("failed to generate random: %v", err)
					return
				}
			}
		} else {
			for i := 0; i < b.N; i++ {
				_, err = state.Generate(dst[:0], dst[:hashSize])
				if err != nil {
					b.Fatalf("failed to generate random: %v", err)
					return
				}
			}
		}
	}
}

func benchHmacDRBG(
	hash func() hash.Hash,
	usePR, useRI, usePS, useAI int,
) func(b *testing.B) {
	hashSize := hash().Size()

	return func(b *testing.B) {
		dst := make([]byte, hashSize)

		var opts []drbg.Option
		{
			nonce := make([]byte, hashSize)
			io.ReadFull(rnd, nonce)
			opts = append(opts, drbg.WithNonce(nonce))
		}
		if usePR == 1 {
			opts = append(opts, drbg.WithPredictionResistance(true))
		}
		if useRI > 0 {
			opts = append(opts, drbg.WithReseedInterval(uint64(useRI)))
		}
		if usePS == 1 {
			personal := make([]byte, hashSize)
			io.ReadFull(rnd, personal)
			opts = append(opts, drbg.WithPersonalizationString(personal))
		}

		state, err := drbg.NewHMACDRGB(
			hash,
			opts...,
		)
		if err != nil {
			b.Fatalf("failed to create CTR_DRBG: %v", err)
			return
		}

		b.ResetTimer()
		b.ReportAllocs()
		b.SetBytes(benchDstLength)
		if useAI == 1 {
			for i := 0; i < b.N; i++ {
				_, err = state.Read(dst)
				if err != nil {
					b.Fatalf("failed to generate random: %v", err)
					return
				}
			}
		} else {
			for i := 0; i < b.N; i++ {
				_, err = state.Generate(dst[:0], dst[:hashSize])
				if err != nil {
					b.Fatalf("failed to generate random: %v", err)
					return
				}
			}
		}
	}
}
