package hashdrbg_test

import (
	"crypto/sha256"
	"fmt"
	"runtime"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/drbg/hashdrbg"
)

const entropySize = 32

type benchCase struct {
	pr int    // 예측 내성 사용
	ps int    // 개별화 문자열 사용
	ai int    // 추가 입력 사용
	ri uint64 // 갱신 주기 (0은 기본값)
}

func benchInput(size int, offset byte) []byte {
	buf := make([]byte, size)
	for i := range buf {
		buf[i] = byte(i) + offset
	}
	return buf
}

func (bc benchCase) GenerateFields() (entropy, ps, nonce, ai []byte) {
	entropy = benchInput(entropySize, 1)
	if bc.ps == 1 {
		ps = benchInput(entropySize, 2)
	}
	nonce = benchInput(entropySize/2, 3)
	if bc.ai == 1 {
		ai = benchInput(entropySize, 4)
	}
	return
}

func BenchmarkInstantiate(b *testing.B) {
	benchCases := []benchCase{
		{ps: 0},
		{ps: 1},
	}

	for _, bc := range benchCases {
		b.Run(fmt.Sprintf("SHA-256/ps=%d", bc.ps), func(b *testing.B) {
			entropy, ps, nonce, _ := bc.GenerateFields()

			h := sha256.New()
			var state hashdrbg.State
			b.ReportAllocs()
			for b.Loop() {
				entropy[0]++
				state.Instantiate(h, entropy, nonce, ps, false, 0)
				// Include disposal of each instance, as in the other DRBG benchmarks.
				state.Zeroize()
			}
		})
	}
}

func BenchmarkReseed(b *testing.B) {
	benchCases := []benchCase{
		{ai: 0},
		{ai: 1},
	}

	for _, bc := range benchCases {
		b.Run(fmt.Sprintf("SHA-256/ai=%d", bc.ai), func(b *testing.B) {
			entropy, ps, nonce, ai := bc.GenerateFields()

			h := sha256.New()
			var state hashdrbg.State
			state.Instantiate(h, entropy, nonce, ps, false, 0)
			b.Cleanup(state.Zeroize)

			b.ReportAllocs()
			for b.Loop() {
				entropy[0]++
				state.Reseed(entropy, ai)
			}
		})
	}
}

func BenchmarkGenerate(b *testing.B) {
	benchCases := []benchCase{
		{pr: 0, ri: 0, ai: 0},
		{pr: 0, ri: 0, ai: 1},
		{pr: 0, ri: 1, ai: 0},
		{pr: 0, ri: 1, ai: 1},
		{pr: 0, ri: 2, ai: 0},
		{pr: 0, ri: 2, ai: 1},
		{pr: 1, ri: 0, ai: 0},
		{pr: 1, ri: 0, ai: 1},
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
			b.Run(fmt.Sprintf("SHA-256/pr=%d,ri=%d,ai=%d/n=%s", bc.pr, bc.ri, bc.ai, n.name), func(b *testing.B) {
				entropy, ps, nonce, additionalInput := bc.GenerateFields()
				entropyBuf := make([]byte, len(entropy))
				fnEntropy := func() ([]byte, error) {
					// Generate clears returned entropy after reseeding.
					copy(entropyBuf, entropy)
					return entropyBuf, nil
				}

				h := sha256.New()
				var state hashdrbg.State
				state.Instantiate(h, entropy, nonce, ps, bc.pr == 1, bc.ri)
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
