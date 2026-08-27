package eckcdsa_test

import (
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/RyuaNerin/go-krypto/eckcdsa"
)

func ExampleSign() {
	key, err := eckcdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	defer key.Zeroize()

	message := []byte("hello, EC-KCDSA")
	// Sign takes the original message, rather than a precomputed digest.
	r, s, err := eckcdsa.Sign(rand.Reader, key, sha256.New(), message)
	if err != nil {
		panic(err)
	}
	fmt.Println(eckcdsa.Verify(&key.PublicKey, sha256.New(), message, r, s))
	fmt.Println(eckcdsa.Verify(&key.PublicKey, sha256.New(), []byte("changed"), r, s))
	// Output:
	// true
	// false
}
