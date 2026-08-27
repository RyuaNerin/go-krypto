package kipher_test

import (
	"crypto/rand"
	"fmt"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/kipher"
	"github.com/RyuaNerin/go-krypto/lea"
)

func ExampleNewGCM() {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	defer clear(key)
	block, err := lea.NewCipher(key)
	if err != nil {
		panic(err)
	}
	defer krypto.Zeroize(block)
	aead, err := kipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	defer krypto.Zeroize(aead)

	// Generate a fresh nonce for every message encrypted with this key.
	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		panic(err)
	}
	additionalData := []byte("message header")
	sealed := aead.Seal(nil, nonce, []byte("hello, GCM"), additionalData)
	plaintext, err := aead.Open(nil, nonce, sealed, additionalData)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(plaintext))
	// Output: hello, GCM
}

func ExampleNewCCM() {
	key := make([]byte, 16)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	defer clear(key)
	block, err := lea.NewCipher(key)
	if err != nil {
		panic(err)
	}
	defer krypto.Zeroize(block)
	aead, err := kipher.NewCCM(block, 12, 16)
	if err != nil {
		panic(err)
	}
	defer krypto.Zeroize(aead)

	// Generate a fresh nonce for every message encrypted with this key.
	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		panic(err)
	}
	sealed := aead.Seal(nil, nonce, []byte("hello, CCM"), nil)
	plaintext, err := aead.Open(nil, nonce, sealed, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(plaintext))
	// Output: hello, CCM
}
