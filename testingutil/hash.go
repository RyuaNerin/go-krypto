package testingutil

import (
	"bytes"
	"encoding/hex"
	"hash"
	"log"
	"testing"
)

type HS func(dst, src []byte) []byte // Hash Sum

// Hash Test All
func HTA(t *testing.T, tests []CipherSize, do func(*testing.T, int)) {
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
) {
	var dst []byte

	for _, tc := range testCases {
		tc.parse()

		log.Println(len(tc.MsgBytes))

		h.Reset()
		h.Write(tc.MsgBytes)
		dst = h.Sum(dst[:0])

		if !bytes.Equal(dst, tc.MDBytes) {
			t.Errorf("hash failed.\nresult: %s\nanswer: %s", hex.EncodeToString(dst), hex.EncodeToString(tc.MDBytes))
			return
		}
	}
}

// Hash Test Two Cipher
func HTTC(
	t *testing.T,
	hashSize int, // in bites,
	DoA HS,
	DoB HS,
) {
	src := make([]byte, continusHashTestIter)
	rnd.Read(src)

	dstA := make([]byte, hashSize/8)
	dstB := make([]byte, hashSize/8)

	m := continusHashTestIter - hashSize

	for srcLen := 1; srcLen < continusHashTestIter; srcLen++ {
		dstA = DoA(dstA[:0], src[:srcLen])
		dstB = DoB(dstB[:0], src[:srcLen])

		if !bytes.Equal(dstA, dstB) {
			t.Error("did not match")
			return
		}

		copy(src[srcLen%m:], dstA)
	}
}

// Hash Bench
func HB(b *testing.B, h hash.Hash, inputSize int) {
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
