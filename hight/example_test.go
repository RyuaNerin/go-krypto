package hight_test

import (
	"encoding/hex"
	"fmt"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/hight"
)

func ExampleNewCipher() {
	// Define a 16-byte key.
	key := []byte("1234567890123456")
	defer krypto.Zeroize(key)

	// Create a new cipher block with the given key.
	block, err := hight.NewCipher(key)
	if err != nil {
		panic(err)
	}
	defer krypto.Zeroize(block)

	//                   12345678
	plaintext := []byte("hello!!!") // 8 bytes = 1 Block

	// Encrypt plaintext
	ciphertext := make([]byte, block.BlockSize())
	block.Encrypt(ciphertext, plaintext)
	fmt.Println("Ciphertext:")
	fmt.Println(hex.Dump(ciphertext))

	// Decrypt ciphertext
	decrypted := make([]byte, block.BlockSize())
	block.Decrypt(decrypted, ciphertext)
	fmt.Println("Decrypted:")
	fmt.Println(hex.Dump(decrypted))
}
