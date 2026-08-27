package kbkdf_test

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	"github.com/RyuaNerin/go-krypto/kbkdf"
	"github.com/RyuaNerin/go-krypto/lea"
)

func ExampleCounterMode() {
	const key = "1234567890123456"

	prf1 := kbkdf.NewHMACPRF(sha256.New)
	prf2 := kbkdf.NewCMACPRF(lea.NewCipher)

	context := []byte("my context")
	key1 := kbkdf.CounterMode(nil, prf1, []byte(key), []byte("labelA"), context, 4, 32)
	key2 := kbkdf.CounterMode(nil, prf2, []byte(key), []byte("labelB"), context, 4, 32)

	fmt.Println(len(key1), len(key2))
	fmt.Println(bytes.Equal(key1, key2))
}
