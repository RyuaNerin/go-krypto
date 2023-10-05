package testingutil

import (
	"bytes"
	"testing"
)

const modifiedSrcIterations = 16 * 1024

// Block Test Src changing
func BTSC(
	t *testing.T,
	keySize int, // in bites,
	additionalSize int, // in bytes, iv, nonce ...
	srcSize int,
	init BI,
	do BD,
	skip bool,
) {
	key := make([]byte, keySize/8)
	additional := make([]byte, additionalSize)

	src := make([]byte, srcSize)
	src2 := make([]byte, srcSize)
	dst := make([]byte, srcSize)

	for i := 0; i < modifiedSrcIterations; i++ {
		rnd.Read(key)
		rnd.Read(additional)
		rnd.Read(src)
		copy(src2, src)

		c, err := init(key, additional)
		if err != nil {
			t.Error(err)
			return
		}

		do(c, dst, src)

		if !bytes.Equal(src, src2) {
			t.Fail()
			return
		}
	}
}

// Block Test Src changing All
func BTSCA(
	t *testing.T,
	sizes []CipherSize,
	additionalSize int, // in bytes, iv, nonce ...
	srcSize int,
	init BI,
	do BD,
	skip bool,
) {
	if skip {
		t.Skip()
		return
	}

	TA(
		t,
		sizes,
		func(t *testing.T, bitSize int) {
			BTSC(
				t,
				bitSize,
				additionalSize,
				srcSize,
				init,
				do,
				false,
			)
		},
		false,
	)
}
