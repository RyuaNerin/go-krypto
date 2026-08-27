package main

import (
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"strings"

	"github.com/RyuaNerin/go-krypto/aria"
	"github.com/RyuaNerin/go-krypto/hight"
	"github.com/RyuaNerin/go-krypto/internal/cavp"
	"github.com/RyuaNerin/go-krypto/kcdsa"
	"github.com/RyuaNerin/go-krypto/lea"
	"github.com/RyuaNerin/go-krypto/lsh256"
	"github.com/RyuaNerin/go-krypto/lsh512"
	"github.com/RyuaNerin/go-krypto/seed"

	"github.com/RyuaNerin/elliptic2/nist"
)

type HashInfo struct {
	New func() hash.Hash
}

var (
	algHash = map[string]*HashInfo{
		"LSH-256-224": {lsh256.New224},
		"LSH-256-256": {lsh256.New},
		"LSH-512-224": {lsh512.New224},
		"LSH-512-256": {lsh512.New256},
		"LSH-512-384": {lsh512.New384},
		"LSH-512-512": {lsh512.New},

		"SHA-224":  {sha256.New224},
		"SHA-256":  {sha256.New},
		"SHA-384":  {sha512.New384},
		"SHA-512":  {sha512.New},
		"SHA2-224": {sha256.New224},
		"SHA2-256": {sha256.New},
		"SHA2-384": {sha512.New384},
		"SHA2-512": {sha512.New},

		"SHA224": {sha256.New224},
		"SHA256": {sha256.New},
		"SHA384": {sha512.New384},
		"SHA512": {sha512.New},
	}

	algBlocks = map[string]cavp.NewCipherFunc{
		"ARIA-128": aria.NewCipher,
		"ARIA-192": aria.NewCipher,
		"ARIA-256": aria.NewCipher,
		"ARIA128":  aria.NewCipher,
		"ARIA192":  aria.NewCipher,
		"ARIA256":  aria.NewCipher,
		"HIGHT":    hight.NewCipher,
		"LEA-128":  lea.NewCipher,
		"LEA-192":  lea.NewCipher,
		"LEA-256":  lea.NewCipher,
		"LEA128":   lea.NewCipher,
		"LEA192":   lea.NewCipher,
		"LEA256":   lea.NewCipher,
		"SEED-128": seed.NewCipher,
		"SEED128":  seed.NewCipher,
		"SEED":     seed.NewCipher,
	}

	algBlockKeySizes = map[string]int{
		"ARIA-128": 128 / 8,
		"ARIA-192": 192 / 8,
		"ARIA-256": 256 / 8,
		"ARIA128":  128 / 8,
		"ARIA192":  192 / 8,
		"ARIA256":  256 / 8,
		"HIGHT":    128 / 8,
		"LEA-128":  128 / 8,
		"LEA-192":  192 / 8,
		"LEA-256":  256 / 8,
		"LEA128":   128 / 8,
		"LEA192":   192 / 8,
		"LEA256":   256 / 8,
		"SEED-128": 128 / 8,
		"SEED128":  128 / 8,
		"SEED":     128 / 8,
	}

	algKCDSA = map[string]kcdsa.ParameterSizes{
		"(2048)(224)_SHA-224":  kcdsa.A2048B224SHA224,
		"(2048)(224)_SHA-256":  kcdsa.A2048B224SHA256,
		"(2048)(256)_SHA-256":  kcdsa.A2048B256SHA256,
		"(3072)(256)_SHA-256":  kcdsa.A3072B256SHA256,
		"(2048)(224)(SHA-224)": kcdsa.A2048B224SHA224,
		"(2048)(224)(SHA-256)": kcdsa.A2048B224SHA256,
		"(2048)(256)(SHA-256)": kcdsa.A2048B256SHA256,
		"(3072)(256)(SHA-256)": kcdsa.A3072B256SHA256,
	}

	algECKCDSACurve = map[string]elliptic.Curve{
		"(B-233)": nist.B233(),
		"(B-283)": nist.B283(),
		"(K-233)": nist.K233(),
		"(K-283)": nist.K283(),
		"(P-224)": elliptic.P224(),
		"(P-256)": elliptic.P256(),
	}

	algBlockMode = map[string]cavp.BlockMode{
		"(ECB)":    cavp.ECB,
		"(CBC)":    cavp.CBC,
		"(CFB8)":   cavp.CFB8,
		"(CFB32)":  cavp.CFB32,
		"(CFB64)":  cavp.CFB64,
		"(CFB128)": cavp.CFB128,
		"(OFB)":    cavp.OFB,
		"(CTR)":    cavp.CTR,
	}
)

func getHash(filename string) *HashInfo {
	for substr, alg := range algHash {
		if strings.Contains(filename, substr) {
			return alg
		}
	}
	return nil
}

func getBlock(filename string) cavp.NewCipherFunc {
	for substr, alg := range algBlocks {
		if strings.Contains(filename, substr) {
			return alg
		}
	}
	return nil
}

func getBlockKeySize(filename string) int {
	for substr, size := range algBlockKeySizes {
		if strings.Contains(filename, substr) {
			return size
		}
	}
	return 0
}

func getKCDSA(filename string) (kcdsa.ParameterSizes, bool) {
	for substr, alg := range algKCDSA {
		if strings.Contains(filename, substr) {
			return alg, true
		}
	}
	return 0, false
}

func getECKCDSACurve(filename string) elliptic.Curve {
	for substr, alg := range algECKCDSACurve {
		if strings.Contains(filename, substr) {
			return alg
		}
	}
	return nil
}

func getBlockMode(filename string) (cavp.BlockMode, bool) {
	for substr, alg := range algBlockMode {
		if strings.Contains(filename, substr) {
			return alg, true
		}
	}
	return 0, false
}
