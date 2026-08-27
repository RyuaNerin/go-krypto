package cmac_test

import (
	"fmt"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/cmac"
	"github.com/RyuaNerin/go-krypto/lea"
)

func ExampleNew() {
	block, err := lea.NewCipher([]byte("1234567980123456")) // 16-byte key for LEA-128
	if err != nil {
		panic(err)
	}
	defer krypto.Zeroize(block)

	// Create a CMAC instance
	mac := cmac.New(block)
	defer krypto.Zeroize(mac)

	const message = "hello world!!!"

	// Write messages
	if _, err := mac.Write([]byte(message)); err != nil {
		panic(err)
	}
	tag := mac.Sum(nil)

	// Verify the message with the same key
	verifier := cmac.New(block)
	defer krypto.Zeroize(verifier)

	if _, err := verifier.Write([]byte(message)); err != nil {
		panic(err)
	}
	fmt.Println(cmac.Equal(tag, verifier.Sum(nil)))
	// Output: true
}
