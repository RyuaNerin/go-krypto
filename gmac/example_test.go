package gmac_test

import (
	"crypto/rand"
	"fmt"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/gmac"
	"github.com/RyuaNerin/go-krypto/lea"
)

func ExampleNewGMAC() {
	block, err := lea.NewCipher([]byte("1234567980123456")) // 16-byte key for LEA-128
	if err != nil {
		panic(err)
	}
	defer krypto.Zeroize(block)

	// generate a nonce
	nonce := make([]byte, gmac.StandardNonceSize)
	if _, err := rand.Read(nonce); err != nil {
		panic(err)
	}

	// Create a GMAC instance
	mac, err := gmac.NewGMAC(block, nonce)
	if err != nil {
		panic(err)
	}
	defer krypto.Zeroize(mac)

	const message = "hello world!!!"

	// Write messages
	if _, err := mac.Write([]byte(message)); err != nil {
		panic(err)
	}
	tag := mac.Sum(nil)

	// Verify the message with the same nonce and key
	verifier, err := gmac.NewGMAC(block, nonce)
	if err != nil {
		panic(err)
	}
	defer krypto.Zeroize(verifier)

	if _, err := verifier.Write([]byte(message)); err != nil {
		panic(err)
	}
	fmt.Println(gmac.Equal(tag, verifier.Sum(nil)))
}
