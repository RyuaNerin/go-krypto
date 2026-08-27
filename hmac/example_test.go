package hmac_test

import (
	"crypto/sha256"
	"fmt"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/hmac"
)

func ExampleNew() {
	// Create an HMAC instance
	mac := hmac.New(sha256.New, []byte("1234567980123456"))
	defer krypto.Zeroize(mac)

	const message = "hello world!!!"

	// Write messages
	if _, err := mac.Write([]byte(message)); err != nil {
		panic(err)
	}
	tag := mac.Sum(nil)

	// Verify the message with the same key
	verifier := hmac.New(sha256.New, []byte("1234567980123456"))
	defer krypto.Zeroize(verifier)

	if _, err := verifier.Write([]byte(message)); err != nil {
		panic(err)
	}
	fmt.Println(hmac.Equal(tag, verifier.Sum(nil)))
	// Output: true
}
