package lsh256

import (
	"github.com/RyuaNerin/go-krypto"
)

func init() {
	krypto.RegisterHash(krypto.LSH256, New)
	krypto.RegisterHash(krypto.LSH256_224, New224)
}
