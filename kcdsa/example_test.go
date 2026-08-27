package kcdsa_test

import (
	"crypto/rand"
	_ "crypto/sha256" // Register the hash used by the parameter set.
	"fmt"

	"github.com/RyuaNerin/go-krypto/kcdsa"
)

func ExampleSign() {
	var key kcdsa.PrivateKey
	defer key.Zeroize()
	const sizes = kcdsa.A2048B256SHA256
	if err := kcdsa.GenerateParameters(&key.Parameters, rand.Reader, sizes); err != nil {
		panic(err)
	}
	if err := kcdsa.GenerateKey(&key, rand.Reader); err != nil {
		panic(err)
	}

	message := []byte("hello, KCDSA")
	// Sign takes the original message, rather than a precomputed digest.
	r, s, err := kcdsa.Sign(rand.Reader, &key, sizes, message)
	if err != nil {
		panic(err)
	}
	fmt.Println(kcdsa.Verify(&key.PublicKey, sizes, message, r, s))
	fmt.Println(kcdsa.Verify(&key.PublicKey, sizes, []byte("changed"), r, s))
	// Output:
	// true
	// false
}
