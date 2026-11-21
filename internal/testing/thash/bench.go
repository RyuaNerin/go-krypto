package thash

import (
	"hash"
	"io"
	"strconv"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/randutil"

	"github.com/stretchr/testify/require"
)

type benchmarkConfig struct {
	Name string
	Size int
}

func newBenchmarkConfigs(blockSize int) []benchmarkConfig {
	return []benchmarkConfig{
		{Name: "8", Size: 8},
		{Name: strconv.Itoa(blockSize), Size: blockSize},
		{Name: strconv.Itoa(blockSize * 2), Size: blockSize * 2},
		{Name: "1K", Size: 1 * 1024},
		{Name: "8K", Size: 8 * 1024},
		{Name: "64K", Size: 64 * 1024},
		{Name: "1M", Size: 1024 * 1024},
	}
}

func (ht *HashTest) BenchNew(b *testing.B, skip bool) {
	b.Helper()
	if skip {
		b.Skip("skipping benchmark due to missing CPU feature")
	}

	for _, sz := range ht.HashInfo.Sizes {
		fnNew := ht.getNew(b, sz.DigestSize)

		b.Run(
			sz.Name,
			func(b *testing.B) {
				for b.Loop() {
					h := fnNew()
					require.NotNil(b, h)
				}
			},
		)
	}
}

type HashFunc interface {
	func(size int) hash.Hash | func() hash.Hash
}

func BenchHash[TNewHash HashFunc](b *testing.B, f TNewHash, size int) {
	b.Helper()

	var newHash func() hash.Hash
	switch fn := any(f).(type) {
	case func(size int) hash.Hash:
		newHash = func() hash.Hash { return fn(size) }
	case func() hash.Hash:
		newHash = fn
	default:
		b.Fatalf("unsupported hash function type: %T", f)
	}

	r, _ := randutil.NewBenchRand()

	benchConfigs := newBenchmarkConfigs(newHash().BlockSize())

	for _, bc := range benchConfigs {
		b.Run(
			bc.Name,
			func(b *testing.B) {
				input := make([]byte, bc.Size)
				_, err := io.ReadFull(r, input)
				require.NoError(b, err, "failed to fill buffer with random data")

				h := newHash()
				require.NotNil(b, h, "failed to create LSH hash")

				sum := make([]byte, h.Size())

				b.ReportAllocs()
				b.SetBytes(int64(bc.Size))
				for b.Loop() {
					h.Reset()
					_, err := h.Write(input)
					require.NoError(b, err, "failed to write to hash")
					sum = h.Sum(sum[:0])
					require.NoError(b, err, "failed to compute hash sum")
				}
			},
		)
	}
}

func (ht *HashTest) BenchHash(b *testing.B, skip bool) {
	b.Helper()
	if skip {
		b.Skip("skipping benchmark due to missing CPU feature")
	}

	r, _ := randutil.NewBenchRand()

	benchConfigs := newBenchmarkConfigs(ht.HashInfo.BlockSize)

	for _, sz := range ht.HashInfo.Sizes {
		fnNew := ht.getNew(b, sz.DigestSize)

		b.Run(
			sz.Name,
			func(b *testing.B) {
				h := fnNew()
				require.NotNil(b, h)

				for _, bc := range benchConfigs {
					b.Run(bc.Name, func(b *testing.B) {
						buf := make([]byte, bc.Size)
						_, err := r.Read(buf)
						require.NoError(b, err)

						sum := make([]byte, h.Size())

						b.ReportAllocs()
						b.SetBytes(int64(bc.Size))
						for b.Loop() {
							h.Reset()
							h.Write(buf)
							sum = h.Sum(sum[:0])
						}
					})
				}
			},
		)
	}
}

func (ht *HashTest) BenchSum(b *testing.B, skip bool) {
	b.Helper()
	if skip {
		b.Skip("skipping benchmark due to missing CPU feature")
	}
	r, _ := randutil.NewBenchRand()

	for _, sz := range ht.HashInfo.Sizes {
		sum := ht.getSum(b, sz.DigestSize)

		var fnBench func(b *testing.B)
		switch fn := any(sum).(type) {
		case func([]byte) []byte:
			fnBench = benchSum(ht, r, fn)
		case func(size int, data []byte) []byte:
			fnBench = benchSumSize(ht, r, sz.DigestSize, fn)

		case func([]byte) [20]byte:
			fnBench = benchSum(ht, r, fn)
		case func(size int, data []byte) [20]byte:
			fnBench = benchSumSize(ht, r, sz.DigestSize, fn)

		case func([]byte) [28]byte:
			fnBench = benchSum(ht, r, fn)
		case func(size int, data []byte) [28]byte:
			fnBench = benchSumSize(ht, r, sz.DigestSize, fn)

		case func([]byte) [32]byte:
			fnBench = benchSum(ht, r, fn)
		case func(size int, data []byte) [32]byte:
			fnBench = benchSumSize(ht, r, sz.DigestSize, fn)

		case func([]byte) [48]byte:
			fnBench = benchSum(ht, r, fn)
		case func(size int, data []byte) [48]byte:
			fnBench = benchSumSize(ht, r, sz.DigestSize, fn)

		case func([]byte) [64]byte:
			fnBench = benchSum(ht, r, fn)
		case func(size int, data []byte) [64]byte:
			fnBench = benchSumSize(ht, r, sz.DigestSize, fn)

		default:
			b.Fatalf("unsupported sum function type: %T", sum)
		}

		b.Run(sz.Name, fnBench)
	}
}

func benchSumSize[T Digest](
	ht *HashTest,
	r io.Reader,
	digestSize int,
	fn func(blockSize int, b []byte) T,
) func(b *testing.B) {
	benchConfigs := newBenchmarkConfigs(ht.HashInfo.BlockSize)

	return func(b *testing.B) {
		b.Helper()

		for _, bc := range benchConfigs {
			b.Run(bc.Name, func(b *testing.B) {
				buf := make([]byte, bc.Size)
				_, err := r.Read(buf)
				require.NoError(b, err)

				b.SetBytes(int64(bc.Size))
				for b.Loop() {
					fn(digestSize, buf)
				}
			})
		}
	}
}

func benchSum[T Digest](
	ht *HashTest,
	r io.Reader,
	fn func([]byte) T,
) func(b *testing.B) {
	benchmarkInputsizes := newBenchmarkConfigs(ht.HashInfo.BlockSize)

	return func(b *testing.B) {
		b.Helper()

		for _, bs := range benchmarkInputsizes {
			b.Run(bs.Name, func(b *testing.B) {
				buf := make([]byte, bs.Size)
				_, err := r.Read(buf)
				require.NoError(b, err)

				b.SetBytes(int64(bs.Size))
				for b.Loop() {
					fn(buf)
				}
			})
		}
	}
}
