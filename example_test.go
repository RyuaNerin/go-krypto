package krypto_test

import (
	"fmt"

	"github.com/RyuaNerin/go-krypto"
	_ "github.com/RyuaNerin/go-krypto/lsh256" // Register LSH-256.
)

func ExampleKryptoHash_New() {
	h := krypto.LSH256.New()
	if _, err := h.Write([]byte("hello")); err != nil {
		panic(err)
	}
	fmt.Println(krypto.LSH256, len(h.Sum(nil)))
	// Output: LSH-256 32
}

func ExampleZeroize() {
	secret := []byte{1, 2, 3, 4}
	krypto.Zeroize(secret)
	fmt.Println(secret)
	// Output: [0 0 0 0]
}

func ExampleZeroizeT() {
	secret := [4]byte{1, 2, 3, 4}
	krypto.ZeroizeT(&secret)
	fmt.Println(secret)
	// Output: [0 0 0 0]
}
