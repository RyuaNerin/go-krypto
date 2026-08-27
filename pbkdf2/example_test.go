package pbkdf2_test

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/RyuaNerin/go-krypto/pbkdf2"
)

func ExampleKey() {
	const password = "example password"

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		panic(err)
	}

	key := pbkdf2.Key([]byte(password), salt, 100_000, 32, sha256.New)
	defer clear(key)

	fmt.Println(hex.Dump(key))
}
