package kipher_test

// based on https://github.com/golang/go/blob/master/src/crypto/cipher/benchmark_test.go

import (
	"crypto/aes"
	"crypto/cipher"
	"strconv"
	"testing"

	ikipher "github.com/RyuaNerin/go-krypto/internal/kipher"
	"github.com/RyuaNerin/go-krypto/kipher"
	"github.com/RyuaNerin/go-krypto/lea"
)

func benchmarkGCMSeal(
	newCipher func([]byte) (cipher.Block, error),
	wrap func(cipher.Block) cipher.Block,
	newGCM func(cipher.Block) (cipher.AEAD, error),
	buf []byte, keySize int,
) func(*testing.B) {
	return func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(buf)))

		key := make([]byte, keySize)
		var nonce [12]byte
		var ad [13]byte
		block, _ := newCipher(key)
		gcm, _ := newGCM(wrap(block))
		var out []byte

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			out = gcm.Seal(out[:0], nonce[:], buf, ad[:])
		}
	}
}

func benchmarkGCMOpen(
	newCipher func([]byte) (cipher.Block, error),
	wrap func(cipher.Block) cipher.Block,
	newGCM func(cipher.Block) (cipher.AEAD, error),
	buf []byte, keySize int,
) func(*testing.B) {
	return func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(buf)))

		key := make([]byte, keySize)
		var nonce [12]byte
		var ad [13]byte
		block, _ := newCipher(key)
		gcm, _ := newGCM(wrap(block))
		var out []byte

		ct := gcm.Seal(nil, nonce[:], buf, ad[:])

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			out, _ = gcm.Open(out[:0], nonce[:], ct, ad[:])
		}
	}
}

func wrapKipher(b cipher.Block) cipher.Block {
	return ikipher.WrapKipher(b)
}

func BenchmarkGCMStd(b *testing.B) {
	for _, length := range []int{64, 1350, 8 * 1024} {
		b.Run("AES/Open-128-"+strconv.Itoa(length), benchmarkGCMOpen(aes.NewCipher, ikipher.WrapCipher, cipher.NewGCM, make([]byte, length), 128/8))
		b.Run("AES/Seal-128-"+strconv.Itoa(length), benchmarkGCMSeal(aes.NewCipher, ikipher.WrapCipher, cipher.NewGCM, make([]byte, length), 128/8))

		b.Run("AES/Open-256-"+strconv.Itoa(length), benchmarkGCMOpen(aes.NewCipher, ikipher.WrapCipher, cipher.NewGCM, make([]byte, length), 256/8))
		b.Run("AES/Seal-256-"+strconv.Itoa(length), benchmarkGCMSeal(aes.NewCipher, ikipher.WrapCipher, cipher.NewGCM, make([]byte, length), 256/8))
	}

	for _, length := range []int{64, 1350, 8 * 1024} {
		b.Run("LEA/Open-128-"+strconv.Itoa(length), benchmarkGCMOpen(lea.NewCipher, ikipher.WrapCipher, cipher.NewGCM, make([]byte, length), 128/8))
		b.Run("LEA/Seal-128-"+strconv.Itoa(length), benchmarkGCMSeal(lea.NewCipher, ikipher.WrapCipher, cipher.NewGCM, make([]byte, length), 128/8))

		b.Run("LEA/Open-256-"+strconv.Itoa(length), benchmarkGCMOpen(lea.NewCipher, ikipher.WrapCipher, cipher.NewGCM, make([]byte, length), 256/8))
		b.Run("LEA/Seal-256-"+strconv.Itoa(length), benchmarkGCMSeal(lea.NewCipher, ikipher.WrapCipher, cipher.NewGCM, make([]byte, length), 256/8))
	}
}

func BenchmarkGCMKipher(b *testing.B) {
	for _, length := range []int{64, 1350, 8 * 1024} {
		b.Run("AES/Open-128-"+strconv.Itoa(length), benchmarkGCMOpen(aes.NewCipher, wrapKipher, kipher.NewGCM, make([]byte, length), 128/8))
		b.Run("AES/Seal-128-"+strconv.Itoa(length), benchmarkGCMSeal(aes.NewCipher, wrapKipher, kipher.NewGCM, make([]byte, length), 128/8))

		b.Run("AES/Open-256-"+strconv.Itoa(length), benchmarkGCMOpen(aes.NewCipher, wrapKipher, kipher.NewGCM, make([]byte, length), 256/8))
		b.Run("AES/Seal-256-"+strconv.Itoa(length), benchmarkGCMSeal(aes.NewCipher, wrapKipher, kipher.NewGCM, make([]byte, length), 256/8))
	}

	for _, length := range []int{64, 1350, 8 * 1024} {
		b.Run("LEA/Open-128-"+strconv.Itoa(length), benchmarkGCMOpen(lea.NewCipher, wrapKipher, kipher.NewGCM, make([]byte, length), 128/8))
		b.Run("LEA/Seal-128-"+strconv.Itoa(length), benchmarkGCMSeal(lea.NewCipher, wrapKipher, kipher.NewGCM, make([]byte, length), 128/8))

		b.Run("LEA/Open-256-"+strconv.Itoa(length), benchmarkGCMOpen(lea.NewCipher, wrapKipher, kipher.NewGCM, make([]byte, length), 256/8))
		b.Run("LEA/Seal-256-"+strconv.Itoa(length), benchmarkGCMSeal(lea.NewCipher, wrapKipher, kipher.NewGCM, make([]byte, length), 256/8))
	}
}

func benchmarkStream(
	newCipher func([]byte) (cipher.Block, error),
	wrap func(cipher.Block) cipher.Block,
	mode func(cipher.Block, []byte) cipher.Stream,
	buf []byte,
) func(*testing.B) {
	return func(b *testing.B) {
		b.SetBytes(int64(len(buf)))

		var key [16]byte
		var iv [16]byte
		block, _ := newCipher(key[:])
		stream := mode(wrap(block), iv[:])

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			stream.XORKeyStream(buf, buf)
		}
	}
}

// If we test exactly 1K blocks, we would generate exact multiples of
// the cipher's block size, and the cipher stream fragments would
// always be wordsize aligned, whereas non-aligned is a more typical
// use-case.
const (
	almost1K = 1024 - 5
	almost8K = 8*1024 - 5
)

func BenchmarkCFBEncrypt1K(b *testing.B) {
	b.Run("AES/Std", benchmarkStream(aes.NewCipher, ikipher.WrapCipher, cipher.NewCFBEncrypter, make([]byte, almost1K)))
	b.Run("AES/Kipher", benchmarkStream(aes.NewCipher, wrapKipher, kipher.NewCFBEncrypter, make([]byte, almost1K)))
	b.Run("LEA/Std", benchmarkStream(lea.NewCipher, ikipher.WrapCipher, cipher.NewCFBEncrypter, make([]byte, almost1K)))
	b.Run("LEA/Kipher", benchmarkStream(lea.NewCipher, wrapKipher, kipher.NewCFBEncrypter, make([]byte, almost1K)))
}

func BenchmarkCFBDecrypt1K(b *testing.B) {
	b.Run("AES/Std", benchmarkStream(aes.NewCipher, ikipher.WrapCipher, cipher.NewCFBDecrypter, make([]byte, almost1K)))
	b.Run("AES/Kipher", benchmarkStream(aes.NewCipher, wrapKipher, kipher.NewCFBDecrypter, make([]byte, almost1K)))
	b.Run("LEA/Std", benchmarkStream(lea.NewCipher, ikipher.WrapCipher, cipher.NewCFBDecrypter, make([]byte, almost1K)))
	b.Run("LEA/Kipher", benchmarkStream(lea.NewCipher, wrapKipher, kipher.NewCFBDecrypter, make([]byte, almost1K)))
}

func BenchmarkCFBDecrypt8K(b *testing.B) {
	b.Run("AES/Std", benchmarkStream(aes.NewCipher, ikipher.WrapCipher, cipher.NewCFBDecrypter, make([]byte, almost8K)))
	b.Run("AES/Kipher", benchmarkStream(aes.NewCipher, wrapKipher, kipher.NewCFBDecrypter, make([]byte, almost8K)))
	b.Run("LEA/Std", benchmarkStream(lea.NewCipher, ikipher.WrapCipher, cipher.NewCFBDecrypter, make([]byte, almost8K)))
	b.Run("LEA/Kipher", benchmarkStream(lea.NewCipher, wrapKipher, kipher.NewCFBDecrypter, make([]byte, almost8K)))
}

func BenchmarkOFB1K(b *testing.B) {
	b.Run("AES/Std", benchmarkStream(aes.NewCipher, ikipher.WrapCipher, cipher.NewOFB, make([]byte, almost1K)))
	b.Run("AES/Kipher", benchmarkStream(aes.NewCipher, wrapKipher, kipher.NewOFB, make([]byte, almost1K)))
	b.Run("LEA/Std", benchmarkStream(lea.NewCipher, ikipher.WrapCipher, cipher.NewOFB, make([]byte, almost1K)))
	b.Run("LEA/Kipher", benchmarkStream(lea.NewCipher, wrapKipher, kipher.NewOFB, make([]byte, almost1K)))
}

func BenchmarkCTR1K(b *testing.B) {
	b.Run("AES/Std", benchmarkStream(aes.NewCipher, ikipher.WrapCipher, cipher.NewCTR, make([]byte, almost8K)))
	b.Run("AES/Kipher", benchmarkStream(aes.NewCipher, wrapKipher, kipher.NewCTR, make([]byte, almost8K)))
	b.Run("LEA/Std", benchmarkStream(lea.NewCipher, ikipher.WrapCipher, cipher.NewCTR, make([]byte, almost8K)))
	b.Run("LEA/Kipher", benchmarkStream(lea.NewCipher, wrapKipher, kipher.NewCTR, make([]byte, almost8K)))
}

func BenchmarkCTR8K(b *testing.B) {
	b.Run("AES/Std", benchmarkStream(aes.NewCipher, ikipher.WrapCipher, cipher.NewCTR, make([]byte, almost8K)))
	b.Run("AES/Kipher", benchmarkStream(aes.NewCipher, wrapKipher, kipher.NewCTR, make([]byte, almost8K)))
	b.Run("LEA/Std", benchmarkStream(lea.NewCipher, ikipher.WrapCipher, cipher.NewCTR, make([]byte, almost8K)))
	b.Run("LEA/Kipher", benchmarkStream(lea.NewCipher, wrapKipher, kipher.NewCTR, make([]byte, almost8K)))
}

func benchmarkBlock(
	newCipher func([]byte) (cipher.Block, error),
	wrap func(cipher.Block) cipher.Block,
	mode func(cipher.Block, []byte) cipher.BlockMode,
	buf []byte,
) func(*testing.B) {
	return func(b *testing.B) {
		b.SetBytes(int64(len(buf)))

		var key [16]byte
		var iv [16]byte
		block, _ := newCipher(key[:])
		stream := mode(wrap(block), iv[:])

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			stream.CryptBlocks(buf, buf)
		}
	}
}

func BenchmarkCBCEncrypt1K(b *testing.B) {
	b.Run("AES/Std", benchmarkBlock(aes.NewCipher, ikipher.WrapCipher, cipher.NewCBCEncrypter, make([]byte, 1*1024)))
	b.Run("AES/Kipher", benchmarkBlock(aes.NewCipher, wrapKipher, kipher.NewCBCEncrypter, make([]byte, 1*1024)))
	b.Run("LEA/Std", benchmarkBlock(lea.NewCipher, ikipher.WrapCipher, cipher.NewCBCEncrypter, make([]byte, 1*1024)))
	b.Run("LEA/Kipher", benchmarkBlock(lea.NewCipher, wrapKipher, kipher.NewCBCEncrypter, make([]byte, 1*1024)))
}

func BenchmarkCBCDecrypt1K(b *testing.B) {
	b.Run("AES/Std", benchmarkBlock(aes.NewCipher, ikipher.WrapCipher, cipher.NewCBCDecrypter, make([]byte, 1*1024)))
	b.Run("AES/Kipher", benchmarkBlock(aes.NewCipher, wrapKipher, kipher.NewCBCDecrypter, make([]byte, 1*1024)))
	b.Run("LEA/Std", benchmarkBlock(lea.NewCipher, ikipher.WrapCipher, cipher.NewCBCDecrypter, make([]byte, 1*1024)))
	b.Run("LEA/Kipher", benchmarkBlock(lea.NewCipher, wrapKipher, kipher.NewCBCDecrypter, make([]byte, 1*1024)))
}
