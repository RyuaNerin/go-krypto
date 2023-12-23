package testingutil

import (
	"bytes"
	"encoding/hex"
	"hash"
	"testing"
)

type HS func(dst, src []byte) []byte // Hash Sum

// Hash Test All
func HTA(
	t *testing.T,
	tests []CipherSize,
	do func(*testing.T, int),
	skip bool,
) {
	if skip {
		t.Skip()
		return
	}

	for _, test := range tests {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			do(t, test.Size)
		})
	}
}

// Hash Test
func HT(
	t *testing.T,
	h hash.Hash,
	testCases []HashTestCase,
	skip bool,
) {
	if skip {
		t.Skip()
		return
	}

	var dst []byte

	for _, tc := range testCases {
		tc.parse()

		h.Reset()
		h.Write(tc.MsgBytes)
		dst = h.Sum(dst[:0])

		if !bytes.Equal(dst, tc.MDBytes) {
			t.Errorf("hash failed.\nresult: %s\nanswer: %s", hex.EncodeToString(dst), hex.EncodeToString(tc.MDBytes))
			return
		}
	}
}

// Hash Test Sum
func HTS(
	t *testing.T,
	hashSize int, // in bites,
	sum1 HS,
	sum2 HS,
	skip bool,
) {
	if skip {
		t.Skip()
		return
	}

	src := make([]byte, continusHashTestIter)
	rnd.Read(src)

	dstA := make([]byte, hashSize/8)
	dstB := make([]byte, hashSize/8)

	m := continusHashTestIter - hashSize

	for srcLen := 1; srcLen < continusHashTestIter; srcLen++ {
		dstA = sum1(dstA[:0], src[:srcLen])
		dstB = sum2(dstB[:0], src[:srcLen])

		if !bytes.Equal(dstA, dstB) {
			t.Error("did not match")
			return
		}

		copy(src[srcLen%m:], dstA)
	}
}

// Hash Test Sum All
func HTSA(
	t *testing.T,
	sizes []CipherSize,
	h1New func(size int) hash.Hash,
	h2New func(size int) hash.Hash,
	skip bool,
) {
	if skip {
		t.Skip()
		return
	}

	TA(
		t,
		sizes,
		func(t *testing.T, size int) {
			h1 := h1New(size)
			h2 := h2New(size)

			HTS(
				t,
				size,
				func(dst, p []byte) []byte {
					h1.Reset()
					h1.Write(p)
					return h1.Sum(dst)
				},
				func(dst, p []byte) []byte {
					h2.Reset()
					h2.Write(p)
					return h2.Sum(dst)
				},
				false,
			)
		},
		false,
	)
}

// Hash Bench
func HB(
	b *testing.B,
	h hash.Hash,
	inputSize int,
	skip bool,
) {
	if skip {
		b.Skip()
		return
	}

	benchBuf := make([]byte, inputSize)
	rnd.Read(benchBuf)

	sum := make([]byte, h.Size())

	b.ReportAllocs()
	b.SetBytes(int64(inputSize))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Reset()
		h.Write(benchBuf)
		h.Sum(sum[:0])
	}
}

// Hash Bench All
func HBA(
	b *testing.B,
	sizes []CipherSize,
	newHash func(size int) hash.Hash,
	inputSize int,
	skip bool,
) {
	if skip {
		b.Skip()
		return
	}

	BA(b, sizes, func(b *testing.B, bitSize int) {
		HB(b, newHash(bitSize), inputSize, false)
	}, false)
}
