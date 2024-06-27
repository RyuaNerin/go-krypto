package gmac_test

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"hash"
	"testing"

	"github.com/RyuaNerin/go-krypto/gmac"
	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/lea"

	. "github.com/RyuaNerin/testingutil"
)

func Test_GMAC_ShortWrite(t *testing.T) {
	c, err := aes.NewCipher(make([]byte, 16))
	if err != nil {
		t.Fatal(err)
	}

	h, err := gmac.NewGMAC(c, nil)
	if err != nil {
		t.Fatal(err)
	}

	HTSW(t, h, false)
}

func Test_GMAC(t *testing.T) {
	var (
		/**
		https://cryptopp.com/wiki/GMAC

		$ ./test.exe
		Key: 00000000000000000000000000000000
		IV:  00000000000000000000000000000000
		Message: Yoda said, do or do not. There is no try.
		GMAC: E7EE2C63B4DC328EED4A86B3FB3490AF
		*/
		key    = internal.HB(`00000000000000000000000000000000`)
		iv     = internal.HB(`00000000000000000000000000000000`)
		msg    = []byte(`Yoda said, do or do not. There is no try.`)
		expect = internal.HB(`E7EE2C63B4DC328EED4A86B3FB3490AF`)
	)

	c, err := aes.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}

	h, err := gmac.NewGMAC(c, iv)
	if err != nil {
		t.Fatal(err)
	}

	h.Write(msg)
	actual := h.Sum(nil)
	if !bytes.Equal(actual, expect) {
		t.Errorf("expect: %x, actual: %x", expect, actual)
		return
	}

	h.Reset()
	h.Write(msg)
	actual = h.Sum(actual[:0])
	if !bytes.Equal(actual, expect) {
		t.Errorf("expect: %x, actual: %x", expect, actual)
		return
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

var (
	aes128, _ = aes.NewCipher(make([]byte, 16))
	aes192, _ = aes.NewCipher(make([]byte, 24))
	aes256, _ = aes.NewCipher(make([]byte, 32))

	lea128, _ = lea.NewCipher(make([]byte, 16))
	lea192, _ = lea.NewCipher(make([]byte, 24))
	lea256, _ = lea.NewCipher(make([]byte, 32))
)

func newHash(b cipher.Block) hash.Hash {
	h, _ := gmac.NewGMAC(b, make([]byte, 16))
	return h
}

func Benchmark_Hash_8(b *testing.B) {
	b.Run("AES-128", func(b *testing.B) { HB(b, newHash(aes128), 8, false) })
	b.Run("AES-192", func(b *testing.B) { HB(b, newHash(aes192), 8, false) })
	b.Run("AES-256", func(b *testing.B) { HB(b, newHash(aes256), 8, false) })

	b.Run("LEA-128", func(b *testing.B) { HB(b, newHash(lea128), 8, false) })
	b.Run("LEA-192", func(b *testing.B) { HB(b, newHash(lea192), 8, false) })
	b.Run("LEA-256", func(b *testing.B) { HB(b, newHash(lea256), 8, false) })
}

func Benchmark_Hash_1K(b *testing.B) {
	b.Run("AES-128", func(b *testing.B) { HB(b, newHash(aes128), 1024, false) })
	b.Run("AES-192", func(b *testing.B) { HB(b, newHash(aes192), 1024, false) })
	b.Run("AES-256", func(b *testing.B) { HB(b, newHash(aes256), 1024, false) })

	b.Run("LEA-128", func(b *testing.B) { HB(b, newHash(lea128), 1024, false) })
	b.Run("LEA-192", func(b *testing.B) { HB(b, newHash(lea192), 1024, false) })
	b.Run("LEA-256", func(b *testing.B) { HB(b, newHash(lea256), 1024, false) })
}

func Benchmark_Hash_8K(b *testing.B) {
	b.Run("AES-128", func(b *testing.B) { HB(b, newHash(aes128), 8192, false) })
	b.Run("AES-192", func(b *testing.B) { HB(b, newHash(aes192), 8192, false) })
	b.Run("AES-256", func(b *testing.B) { HB(b, newHash(aes256), 8192, false) })

	b.Run("LEA-128", func(b *testing.B) { HB(b, newHash(lea128), 8192, false) })
	b.Run("LEA-192", func(b *testing.B) { HB(b, newHash(lea192), 8192, false) })
	b.Run("LEA-256", func(b *testing.B) { HB(b, newHash(lea256), 8192, false) })
}
