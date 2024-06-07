package pbkdf

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/lsh256"
)

func TestSHA256(t *testing.T) {
	dst := Generate(
		nil,
		[]byte(`TTAKO!HellowWorld!SHA2256`),
		internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
		2048,
		768/8,
		sha256.New,
	)
	expect := internal.HB(`
		934cf2b39078d94ab9c7fbf7b98241ed9a1db71a4a80465e7925b27846e699a6
		c878d69473beac7ab6b2e458f22dc85052463704600597aff02022c3cd35d72f
		e833e46bcd7bcfc3c30b2abc8a141b551f4488497d3f6b17b4d14c4fcc3b448e`)

	if !bytes.Equal(dst, expect) {
		t.Fail()
	}
}

func TestLSH256(t *testing.T) {
	dst := Generate(
		nil,
		[]byte(`TTAKO!HellowWorld!LSH256256`),
		internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
		2048,
		768/8,
		lsh256.New,
	)
	expect := internal.HB(`
		053b1dd10a3a97d584e187bfe3659571af75aaae9547ccb883598726b598713b
		1028a8edf065ad71387ced631304c42fca127aa6fc4e9bdcb7b1d3d11e8a43cd
		9b91f7e603765e0e3bb53aec8ea4dd6398d38e8b46f5b30df6e0ee89ae723f1d`)

	if !bytes.Equal(dst, expect) {
		t.Fail()
	}
}

func BenchmarkGenerate(b *testing.B) {
	rnd := bufio.NewReaderSize(rand.Reader, 1<<10)

	bench := func(h func() hash.Hash, iteration, keyLen, passwordLen, saltLen int) func(b *testing.B) {
		return func(b *testing.B) {
			password := make([]byte, passwordLen)
			salt := make([]byte, saltLen)
			out := make([]byte, keyLen)

			rnd.Read(salt)
			rnd.Read(out)

			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				copy(password, salt)
				copy(salt, out)
				out = Generate(out[:0], password, salt, iteration, keyLen, h)
			}
		}
	}

	// https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html#pbkdf2
	b.Run("SHA1: ITER 1_300_000", bench(sha1.New, 1_300_000, 128, 128, 128))
	b.Run("SHA256: ITER 600_000", bench(sha256.New, 600_000, 128, 128, 128))
	b.Run("SHA512: ITER 210_000", bench(sha512.New, 210_000, 128, 128, 128))
}
