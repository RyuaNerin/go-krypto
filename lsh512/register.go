package lsh512

import "github.com/RyuaNerin/go-krypto"

func init() {
	krypto.RegisterHash(krypto.LSH512, New)
	krypto.RegisterHash(krypto.LSH512_384, New384)
	krypto.RegisterHash(krypto.LSH512_256, New256)
	krypto.RegisterHash(krypto.LSH256_224, New224)
}
