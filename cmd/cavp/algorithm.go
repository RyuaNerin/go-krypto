package main

import (
	"crypto/cipher"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"strings"

	"github.com/RyuaNerin/elliptic2/nist"
	"github.com/RyuaNerin/go-krypto/aria"
	"github.com/RyuaNerin/go-krypto/hight"
	"github.com/RyuaNerin/go-krypto/kcdsa"
	"github.com/RyuaNerin/go-krypto/lea"
	"github.com/RyuaNerin/go-krypto/lsh256"
	"github.com/RyuaNerin/go-krypto/lsh512"
	"github.com/RyuaNerin/go-krypto/seed"
)

type HashInfo struct {
	w   int
	New func() hash.Hash
}

type (
	funcNewBlockCipher func(key []byte) (cipher.Block, error)
	funcNewBlockMode   func(cipher cipher.Block, iv []byte) funcProcessBlock
	funcProcessBlock   func(dst []byte, src []byte, encrypt bool)
	funcMCT            func(cavp *cavpProcessor, fnCipher funcNewBlockCipher, key, pt, ct, iv [][]byte, blockSize int)
)

var (
	algHash = map[string]*HashInfo{
		"LSH-256-224": {32, lsh256.New224},
		"LSH-256-256": {32, lsh256.New},
		"LSH-512-224": {64, lsh512.New224},
		"LSH-512-256": {64, lsh512.New256},
		"LSH-512-384": {64, lsh512.New384},
		"LSH-512-512": {64, lsh512.New},

		"SHA-224": {32, sha256.New224},
		"SHA-256": {32, sha256.New},
		"SHA-384": {64, sha512.New384},
		"SHA-512": {64, sha512.New},

		"SHA224": {32, sha256.New224},
		"SHA256": {32, sha256.New},
		"SHA384": {64, sha512.New384},
		"SHA512": {64, sha512.New},
	}

	algBlocks = map[string]funcNewBlockCipher{
		"ARIA-128": aria.NewCipher,
		"ARIA-192": aria.NewCipher,
		"ARIA-256": aria.NewCipher,
		"ARIA128":  aria.NewCipher,
		"ARIA192":  aria.NewCipher,
		"ARIA256":  aria.NewCipher,

		"HIGHT": hight.NewCipher,

		"LEA-128": lea.NewCipher,
		"LEA-192": lea.NewCipher,
		"LEA-256": lea.NewCipher,
		"LEA128":  lea.NewCipher,
		"LEA192":  lea.NewCipher,
		"LEA256":  lea.NewCipher,

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

		"HIGHT": 128 / 8,

		"LEA-128": 128 / 8,
		"LEA-192": 192 / 8,
		"LEA-256": 256 / 8,
		"LEA128":  128 / 8,
		"LEA192":  192 / 8,
		"LEA256":  256 / 8,

		"SEED-128": 128 / 8,
		"SEED128":  128 / 8,
	}

	algKCDSA = map[string]kcdsa.ParameterSizes{
		"(2048)(224)_SHA-224": kcdsa.L2048N224SHA224,
		"(2048)(224)_SHA-256": kcdsa.L2048N224SHA256,
		"(2048)(256)_SHA-256": kcdsa.L2048N256SHA256,
		"(3072)(256)_SHA-256": kcdsa.L3072N256SHA256,
	}

	algECKCDSACurve = map[string]elliptic.Curve{
		"(B-233)": nist.B233(),
		"(B-283)": nist.B283(),
		"(K-233)": nist.K233(),
		"(K-283)": nist.K283(),
		"(P-224)": elliptic.P224(),
		"(P-256)": elliptic.P256(),
	}

	algBlockMode = map[string]funcNewBlockMode{
		"(ECB)":    blockECB,
		"(CBC)":    blockCBC,
		"(CFB8)":   blockCFB(8),
		"(CFB32)":  blockCFB(32),
		"(CFB64)":  blockCFB(64),
		"(CFB128)": blockCFB(128),
		"(OFB)":    blockOFB,
		"(CTR)":    blockCTR,
	}
	algBlockModeMCT = map[string]funcMCT{
		"(ECB)":    processBlock_MCT_ECB,
		"(CBC)":    processBlock_MCT_CBC,
		"(CFB8)":   processBlock_MCT_CFB(8),
		"(CFB32)":  processBlock_MCT_CFB(32),
		"(CFB64)":  processBlock_MCT_CFB(64),
		"(CFB128)": processBlock_MCT_CFB(128),
		"(OFB)":    processBlock_MCT_OFB,
		"(CTR)":    processBlock_MCT_CTR,
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

func getBlock(filename string) funcNewBlockCipher {
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

func getKCDSA(filename string) kcdsa.ParameterSizes {
	for substr, alg := range algKCDSA {
		if strings.Contains(filename, substr) {
			return alg
		}
	}
	return 0
}

func getECKCDSACurve(filename string) elliptic.Curve {
	for substr, alg := range algECKCDSACurve {
		if strings.Contains(filename, substr) {
			return alg
		}
	}
	return nil
}

func getBlockMode(filename string) funcNewBlockMode {
	for substr, alg := range algBlockMode {
		if strings.Contains(filename, substr) {
			return alg
		}
	}
	return nil
}

func getBlockModeMCT(filename string) funcMCT {
	for substr, alg := range algBlockModeMCT {
		if strings.Contains(filename, substr) {
			return alg
		}
	}
	return nil
}
