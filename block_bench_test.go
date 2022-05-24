package krypto

import (
	"crypto/aes"
	"crypto/cipher"
	"testing"

	"github.com/RyuaNerin/go-krypto/aria"
	"github.com/RyuaNerin/go-krypto/hight"
	"github.com/RyuaNerin/go-krypto/lea"
	"github.com/RyuaNerin/go-krypto/seed"
)

func Benchmark_CBC_Encrypt_1K_AES(b *testing.B) {
	c, _ := aes.NewCipher(make([]byte, 16))
	benchBlock1k(b, true, c)
}
func Benchmark_CBC_Decrypt_1K_AES(b *testing.B) {
	c, _ := aes.NewCipher(make([]byte, 16))
	benchBlock1k(b, false, c)
}

func Benchmark_CBC_Encrypt_1K_SEED128(b *testing.B) {
	c, _ := seed.NewCipher(make([]byte, 16))
	benchBlock1k(b, true, c)
}
func Benchmark_CBC_Decrypt_1K_SEED128(b *testing.B) {
	c, _ := seed.NewCipher(make([]byte, 16))
	benchBlock1k(b, false, c)
}

func Benchmark_CBC_Encrypt_1K_SEED256(b *testing.B) {
	c, _ := seed.NewCipher(make([]byte, 32))
	benchBlock1k(b, true, c)
}
func Benchmark_CBC_Decrypt_1K_SEED256(b *testing.B) {
	c, _ := seed.NewCipher(make([]byte, 32))
	benchBlock1k(b, false, c)
}

func Benchmark_CBC_Encrypt_1K_HIGHT(b *testing.B) {
	c, _ := hight.NewCipher(make([]byte, 16))
	benchBlock1k(b, true, c)
}
func Benchmark_CBC_Decrypt_1K_HIGHT(b *testing.B) {
	c, _ := hight.NewCipher(make([]byte, 16))
	benchBlock1k(b, false, c)
}

func Benchmark_CBC_Encrypt_1K_ARIA_16(b *testing.B) {
	c, _ := aria.NewCipher(make([]byte, 16))
	benchBlock1k(b, true, c)
}
func Benchmark_CBC_Decrypt_1K_ARIA_16(b *testing.B) {
	c, _ := aria.NewCipher(make([]byte, 16))
	benchBlock1k(b, false, c)
}

func Benchmark_CBC_Encrypt_1K_ARIA_24(b *testing.B) {
	c, _ := aria.NewCipher(make([]byte, 24))
	benchBlock1k(b, true, c)
}
func Benchmark_CBC_Decrypt_1K_ARIA_24(b *testing.B) {
	c, _ := aria.NewCipher(make([]byte, 24))
	benchBlock1k(b, false, c)
}

func Benchmark_CBC_Encrypt_1K_ARIA_32(b *testing.B) {
	c, _ := aria.NewCipher(make([]byte, 32))
	benchBlock1k(b, true, c)
}
func Benchmark_CBC_Decrypt_1K_ARIA_32(b *testing.B) {
	c, _ := aria.NewCipher(make([]byte, 32))
	benchBlock1k(b, false, c)
}

func Benchmark_CBC_Encrypt_1K_LEA_16(b *testing.B) {
	c, _ := lea.NewCipher(make([]byte, 16))
	benchBlock1k(b, true, c)
}
func Benchmark_CBC_Decrypt_1K_LEA_16(b *testing.B) {
	c, _ := lea.NewCipher(make([]byte, 16))
	benchBlock1k(b, false, c)
}

func Benchmark_CBC_Encrypt_1K_LEA_24(b *testing.B) {
	c, _ := lea.NewCipher(make([]byte, 24))
	benchBlock1k(b, true, c)
}
func Benchmark_CBC_Decrypt_1K_LEA_24(b *testing.B) {
	c, _ := lea.NewCipher(make([]byte, 24))
	benchBlock1k(b, false, c)
}

func Benchmark_CBC_Encrypt_1K_LEA_32(b *testing.B) {
	c, _ := lea.NewCipher(make([]byte, 32))
	benchBlock1k(b, true, c)
}
func Benchmark_CBC_Decrypt_1K_LEA_32(b *testing.B) {
	c, _ := lea.NewCipher(make([]byte, 32))
	benchBlock1k(b, false, c)
}

func benchBlock1k(b *testing.B, encryptMode bool, block cipher.Block) {
	buf := make([]byte, 1024)
	b.SetBytes(int64(len(buf)))

	var bm cipher.BlockMode

	iv := make([]byte, block.BlockSize())
	if encryptMode {
		bm = cipher.NewCBCEncrypter(block, iv)
	} else {
		bm = cipher.NewCBCDecrypter(block, iv)
	}

	b.SetBytes(int64(len(buf)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bm.CryptBlocks(buf, buf)
	}
}
