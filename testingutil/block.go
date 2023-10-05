package testingutil

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"testing"
)

type BI func(key, additional []byte) (interface{}, error) // Block Init
type BD func(c interface{}, dst, src []byte)              // Block Do

func CE(data interface{}, dst, src []byte) { data.(cipher.Block).Encrypt(dst, src) } // Cipher.Encrypt
func CD(data interface{}, dst, src []byte) { data.(cipher.Block).Decrypt(dst, src) } // Cipher.Decrypt

// Block Init Wrap
func BIW(f func(key []byte) (cipher.Block, error)) BI {
	return func(key, additional []byte) (interface{}, error) {
		return f(key)
	}
}

// Block Test Encrypt
func BTE(
	t *testing.T,
	init BI,
	do BD,
	testCases []BlockTestCase,
	skip bool,
) {
	if skip {
		t.Skip()
		return
	}

	var dst []byte

	for _, tc := range testCases {
		tc.parse()

		c, err := init(tc.KeyBytes, tc.IVBytes)
		if err != nil {
			t.Error(err)
			return
		}

		if len(dst) < len(tc.PlainBytes) {
			dst = make([]byte, len(tc.PlainBytes))
		}

		do(c, dst, tc.PlainBytes)
		if !bytes.Equal(dst, tc.SecureBytes) {
			t.Errorf("encrypt failed.\nresult: %s\nanswer: %s", hex.EncodeToString(dst), hex.EncodeToString(tc.SecureBytes))
			return
		}
	}
}

// Block Test Decrypt
func BTD(
	t *testing.T,
	init BI,
	do BD,
	testCases []BlockTestCase,
	skip bool,
) {
	if skip {
		t.Skip()
		return
	}

	var dst []byte

	for _, tc := range testCases {
		tc.parse()

		c, err := init(tc.KeyBytes, tc.IVBytes)
		if err != nil {
			t.Error(err)
			return
		}

		if len(dst) < len(tc.SecureBytes) {
			dst = make([]byte, len(tc.SecureBytes))
		}

		do(c, dst, tc.SecureBytes)
		if !bytes.Equal(dst, tc.PlainBytes) {
			t.Errorf("decrypt failed.\nresult: %s\nanswer: %s", hex.EncodeToString(dst), hex.EncodeToString(tc.PlainBytes))
			return
		}
	}
}

// Block Test Two Cipher
func BTTC(
	t *testing.T,
	keySize int, // in bites,
	additionalSize int, // in bytes, iv, nonce ...
	srcSize int,
	rndFitSize int, // 0 = ignore
	init BI,
	DoA BD,
	DoB BD,
	skip bool,
) {
	if skip {
		t.Skip()
		return
	}

	key := make([]byte, keySize/8)
	additional := make([]byte, additionalSize)
	rnd.Read(key)
	rnd.Read(additional)

	c, err := init(key, additional)
	if err != nil {
		t.Error(err)
		return
	}

	src := make([]byte, srcSize)
	rnd.Read(src)

	dstA := make([]byte, srcSize)
	dstB := make([]byte, srcSize)

	l := srcSize
	lRaw := make([]byte, 4)

	for i := 0; i < continusBlockTestIter; i++ {
		if rndFitSize != 0 {
			rnd.Read(lRaw)
			l = int(binary.LittleEndian.Uint32(lRaw))
			l = l % (srcSize/rndFitSize - 1)
			l = (1 + l) * rndFitSize
		}

		DoA(c, dstA[:l], src[:l])
		DoB(c, dstB[:l], src[:l])

		if !bytes.Equal(dstA[:l], dstB[:l]) {
			t.Error("did not match")
			return
		}

		copy(src[len(src)%l:], dstA[:l])
	}
}

// Block Bench New
func BBN(
	b *testing.B,
	keySize int, // in bites
	additionalSize int, // in bytes, iv, nonce ...
	init BI,
	skip bool,
) {
	if skip {
		b.Skip()
		return
	}

	key := make([]byte, keySize/8)
	additional := make([]byte, additionalSize)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rnd.Read(key)
		rnd.Read(additional)
		_, err := init(key, additional)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

// Block Bench New All
func BBNA(
	b *testing.B,
	sizes []CipherSize,
	additionalSize int, // in bytes, iv, nonce ...
	init BI,
	skip bool,
) {
	if skip {
		b.Skip()
		return
	}

	BA(b, sizes, func(b *testing.B, keySize int) {
		BBN(b, keySize, additionalSize, init, false)
	}, false)
}

// Block Bench Do
func BBD(
	b *testing.B,
	keySize int, // in bites
	additionalSize int, // in bytes, iv, nonce ...
	srcSize int,
	init BI,
	do BD,
	skip bool,
) {
	if skip {
		b.Skip()
		return
	}

	key := make([]byte, keySize/8)
	additional := make([]byte, additionalSize)

	rand.Read(key)
	rand.Read(additional)

	c, err := init(key, additional)
	if err != nil {
		b.Error(err)
		return
	}

	dst := make([]byte, srcSize)
	src := make([]byte, srcSize)
	rand.Read(src)

	b.ReportAllocs()
	b.SetBytes(int64(len(src)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		do(c, dst, src)
		copy(src, dst)
	}
}

// Block Bench Do All
func BBDA(
	b *testing.B,
	sizes []CipherSize,
	additionalSize int, // in bytes, iv, nonce ...
	srcSize int,
	init BI,
	do BD,
	skip bool,
) {
	if skip {
		b.Skip()
		return
	}

	BA(b, sizes, func(b *testing.B, keySize int) {
		BBD(b, keySize, 0, srcSize, init, do, false)
	}, false)
}
