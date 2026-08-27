package lsh512_test

import (
	"encoding/hex"
	"fmt"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/lsh512"
)

func ExampleNew() {
	h := lsh512.New()
	defer krypto.Zeroize(h)

	for _, part := range []string{"hello", ", world"} {
		if _, err := h.Write([]byte(part)); err != nil {
			panic("failed to write to lsh512 hash")
		}
	}
	fmt.Println("Hashing 'hello, world' with lsh512:")
	fmt.Println(hex.Dump(h.Sum(nil)))

	h.Reset()
	if _, err := h.Write([]byte("hello, world")); err != nil {
		panic("failed to write to lsh512 hash")
	}
	fmt.Println("Hashing 'hello, world' after reset:")
	fmt.Println(hex.Dump(h.Sum(nil)))
}

func ExampleSum224() {
	sum := lsh512.Sum224([]byte("abc"))

	fmt.Println("Sum LSH-512-224:")
	fmt.Println(hex.Dump(sum[:]))
}

func ExampleSum256() {
	sum := lsh512.Sum256([]byte("abc"))

	fmt.Println("Sum LSH-512-256:")
	fmt.Println(hex.Dump(sum[:]))
}

func ExampleSum384() {
	sum := lsh512.Sum384([]byte("abc"))

	fmt.Println("Sum LSH-512-384:")
	fmt.Println(hex.Dump(sum[:]))
}

func ExampleSum512() {
	sum := lsh512.Sum512([]byte("abc"))

	fmt.Println("Sum LSH-512-512:")
	fmt.Println(hex.Dump(sum[:]))
}
