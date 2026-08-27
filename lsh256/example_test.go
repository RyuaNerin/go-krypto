package lsh256_test

import (
	"encoding/hex"
	"fmt"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/lsh256"
)

func ExampleNew() {
	h := lsh256.New()
	defer krypto.Zeroize(h)

	for _, part := range []string{"hello", ", world"} {
		if _, err := h.Write([]byte(part)); err != nil {
			panic("failed to write to lsh256 hash")
		}
	}
	fmt.Println("Hashing 'hello, world' with lsh256:")
	fmt.Println(hex.Dump(h.Sum(nil)))

	h.Reset()
	if _, err := h.Write([]byte("hello, world")); err != nil {
		panic("failed to write to lsh256 hash")
	}
	fmt.Println("Hashing 'hello, world' after reset:")
	fmt.Println(hex.Dump(h.Sum(nil)))
}

func ExampleSum224() {
	sum := lsh256.Sum224([]byte("abc"))

	fmt.Println("Sum LSH-256-224:")
	fmt.Println(hex.Dump(sum[:]))
}

func ExampleSum256() {
	sum := lsh256.Sum256([]byte("abc"))

	fmt.Println("Sum LSH-256-256:")
	fmt.Println(hex.Dump(sum[:]))
}
