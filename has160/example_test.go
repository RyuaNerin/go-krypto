package has160_test

import (
	"encoding/hex"
	"fmt"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/has160"
)

func ExampleNew() {
	h := has160.New()
	defer krypto.Zeroize(h)

	for _, part := range []string{"hello", ", world"} {
		if _, err := h.Write([]byte(part)); err != nil {
			panic("failed to write to has160 hash")
		}
	}
	fmt.Println("Hashing 'hello, world' with has160:")
	fmt.Println(hex.Dump(h.Sum(nil)))

	h.Reset()
	h.Write([]byte("hello, world"))
	fmt.Println("Hashing 'hello, world' after reset:")
	fmt.Println(hex.Dump(h.Sum(nil)))
}

func ExampleSum() {
	sum := has160.Sum([]byte("abc"))

	fmt.Println("Sum Hash160:")
	fmt.Println(hex.Dump(sum[:]))
}
