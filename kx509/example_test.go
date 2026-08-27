package kx509_test

import (
	"crypto/elliptic"
	"crypto/rand"
	"fmt"

	"github.com/RyuaNerin/go-krypto/eckcdsa"
	"github.com/RyuaNerin/go-krypto/kx509"
)

func ExampleMarshalPKCS8PrivateKey() {
	key, err := eckcdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	defer key.Zeroize()

	der, err := kx509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		panic(err)
	}
	defer clear(der)
	parsed, err := kx509.ParsePKCS8PrivateKey(der)
	if err != nil {
		panic(err)
	}
	decoded, ok := parsed.(*eckcdsa.PrivateKey)
	if !ok {
		panic("unexpected private key type")
	}
	defer decoded.Zeroize()
	fmt.Println(key.Equal(decoded))
	// Output: true
}

func ExampleMarshalPKIXPublicKey() {
	key, err := eckcdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	defer key.Zeroize()

	der, err := kx509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		panic(err)
	}
	parsed, err := kx509.ParsePKIXPublicKey(der)
	if err != nil {
		panic(err)
	}
	fmt.Println(key.PublicKey.Equal(parsed))
	// Output: true
}
