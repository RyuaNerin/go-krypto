package ctrdrbg_test

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"runtime"
	"testing"

	"github.com/RyuaNerin/go-krypto/hight"
	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/drbg/ctrdrbg"
	"github.com/RyuaNerin/go-krypto/lea"
)

const (
	keySize     = 32
	entropySize = 32
)

type benchCase struct {
	df int // 유도 함수 사용
	pr int // 예측 내성 사용
	ps int // 개별화 문자열 사용
	ai int // 추가 입력 사용
}

// Compare tail dispatch through the generic AES adapter and native batch ciphers.
func BenchmarkGenerateBlocks(b *testing.B) {
	ciphers := []struct {
		name      string
		newCipher func([]byte) (cipher.Block, error)
		keySize   int
		blockSize int
	}{
		{name: "AES-256", newCipher: aes.NewCipher, keySize: 32, blockSize: aes.BlockSize},
		{name: "LEA-256", newCipher: lea.NewCipher, keySize: 32, blockSize: lea.BlockSize},
		{name: "HIGHT", newCipher: hight.NewCipher, keySize: 16, blockSize: hight.BlockSize},
	}

	for _, c := range ciphers {
		sizes := []struct {
			name string
			n    int
		}{
			{name: "blocks=1", n: c.blockSize},
			{name: "blocks=2", n: 2 * c.blockSize},
			{name: "blocks=3", n: 3 * c.blockSize},
			{name: "blocks=4", n: 4 * c.blockSize},
			{name: "blocks=5", n: 5 * c.blockSize},
			{name: "blocks=6", n: 6 * c.blockSize},
			{name: "blocks=7", n: 7 * c.blockSize},
			{name: "blocks=8", n: 8 * c.blockSize},
			{name: "n=1K", n: 1024},
		}
		for _, size := range sizes {
			b.Run(c.name+"/"+size.name, func(b *testing.B) {
				entropy := benchInput(c.keySize+c.blockSize, 1)
				var state ctrdrbg.State
				if err := state.Instantiate(
					internal.ToBlock8Fn(c.newCipher), c.keySize, 0, 0,
					entropy, nil, nil, false, false,
				); err != nil {
					b.Fatal(err)
				}
				b.Cleanup(state.Zeroize)

				buf := make([]byte, size.n)
				b.ReportAllocs()
				b.SetBytes(int64(len(buf)))
				for b.Loop() {
					if err := state.Generate(buf, nil, nil); err != nil {
						b.Fatal(err)
					}
				}
				runtime.KeepAlive(buf)
			})
		}
	}
}

var newCipher = internal.ToBlock8Fn(aes.NewCipher)

func benchInput(size int, offset byte) []byte {
	buf := make([]byte, size)
	for i := range buf {
		buf[i] = byte(i) + offset
	}
	return buf
}

func (bc benchCase) GenerateFields() (entropy, ps, nonce, ai []byte) {
	size := entropySize
	if bc.df == 0 {
		// Without a derivation function, entropy must have seedlen bytes.
		size = keySize + aes.BlockSize
	}
	entropy = benchInput(size, 1)
	if bc.ps == 1 {
		ps = benchInput(entropySize, 2)
	}
	if bc.df == 1 {
		nonce = benchInput(entropySize/2, 3)
	}
	if bc.ai == 1 {
		ai = benchInput(entropySize, 4)
	}
	return
}

func BenchmarkInstantiate(b *testing.B) {
	benchCases := []benchCase{
		{df: 1, ps: 0},
		{df: 1, ps: 1},
		{df: 0, ps: 0},
		{df: 0, ps: 1},
	}

	for _, bc := range benchCases {
		b.Run(fmt.Sprintf("AES-256/df=%d,ps=%d", bc.df, bc.ps), func(b *testing.B) {
			entropy, ps, nonce, _ := bc.GenerateFields()

			var state ctrdrbg.State
			b.ReportAllocs()
			for b.Loop() {
				entropy[0]++
				if err := state.Instantiate(
					newCipher, keySize, 0, 0, entropy, nonce, ps, bc.df == 1, false,
				); err != nil {
					b.Fatal(err)
				}
				// Include disposal of each instance, as in the other DRBG benchmarks.
				state.Zeroize()
			}
		})
	}
}

func BenchmarkReseed(b *testing.B) {
	benchCases := []benchCase{
		{df: 1, ai: 0},
		{df: 1, ai: 1},
		{df: 0, ai: 0},
		{df: 0, ai: 1},
	}

	for _, bc := range benchCases {
		b.Run(fmt.Sprintf("AES-256/df=%d,ai=%d", bc.df, bc.ai), func(b *testing.B) {
			entropy, ps, nonce, ai := bc.GenerateFields()

			var state ctrdrbg.State
			if err := state.Instantiate(
				newCipher, keySize, 0, 0, entropy, nonce, ps, bc.df == 1, false,
			); err != nil {
				b.Fatal(err)
			}
			b.Cleanup(state.Zeroize)

			b.ReportAllocs()
			for b.Loop() {
				entropy[0]++
				if err := state.Reseed(entropy, ai); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkGenerate(b *testing.B) {
	benchCases := []benchCase{
		{df: 1, pr: 0, ai: 0},
		{df: 1, pr: 0, ai: 1},
		{df: 1, pr: 1, ai: 0},
		{df: 1, pr: 1, ai: 1},
		{df: 0, pr: 0, ai: 0},
		{df: 0, pr: 0, ai: 1},
		{df: 0, pr: 1, ai: 0},
		{df: 0, pr: 1, ai: 1},
	}
	nSizes := []struct {
		name string
		n    int
	}{
		{name: "1", n: 1},
		{name: "8", n: 8},
		{name: "1K", n: 1024},
		{name: "8K", n: 8 * 1024},
		{name: "64K", n: 64 * 1024},
	}

	for _, bc := range benchCases {
		for _, n := range nSizes {
			b.Run(fmt.Sprintf("AES-256/df=%d,pr=%d,ai=%d/n=%s", bc.df, bc.pr, bc.ai, n.name), func(b *testing.B) {
				entropy, ps, nonce, additionalInput := bc.GenerateFields()
				entropyBuf := make([]byte, len(entropy))
				fnEntropy := func() ([]byte, error) {
					// Generate clears returned entropy after reseeding.
					copy(entropyBuf, entropy)
					return entropyBuf, nil
				}

				var state ctrdrbg.State
				if err := state.Instantiate(
					newCipher, keySize, 0, 0, entropy, nonce, ps, bc.df == 1, bc.pr == 1,
				); err != nil {
					b.Fatal(err)
				}
				b.Cleanup(state.Zeroize)

				buf := make([]byte, n.n)
				b.ReportAllocs()
				b.SetBytes(int64(len(buf)))
				for b.Loop() {
					if err := state.Generate(buf, fnEntropy, additionalInput); err != nil {
						b.Fatal(err)
					}
				}
				runtime.KeepAlive(buf)
			})
		}
	}
}
