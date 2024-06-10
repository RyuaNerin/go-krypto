package main

import (
	"bytes"
	"crypto/cipher"
	"encoding/hex"
	"hash"
	"log"

	"github.com/RyuaNerin/go-krypto/aria"
	"github.com/RyuaNerin/go-krypto/has160" //nolint:staticcheck
	"github.com/RyuaNerin/go-krypto/hight"
	"github.com/RyuaNerin/go-krypto/lea"
	"github.com/RyuaNerin/go-krypto/lsh256"
	"github.com/RyuaNerin/go-krypto/lsh512"
	"github.com/RyuaNerin/go-krypto/seed"
)

var (
	blockKey = []byte("password")
	input    = []byte("Hello, World!")
)

func main() {
	printBlocks()
	printHashs()
}

func printBlocks() {
	blocks := []struct {
		name    string
		keyBits int
		f       func(key []byte) (cipher.Block, error)
	}{
		{"ARIA-128", 128, aria.NewCipher},
		{"ARIA-196", 196, aria.NewCipher},
		{"ARIA-256", 256, aria.NewCipher},
		{"HIGHT-128", 128, hight.NewCipher},
		{"LEA-256", 128, lea.NewCipher},
		{"LEA-256", 196, lea.NewCipher},
		{"LEA-256", 256, lea.NewCipher},
		{"SEED-128", 128, seed.NewCipher},
	}

	dst := make([]byte, 100)
	dst2 := make([]byte, 100)

	for _, b := range blocks {
		key := pkcs7pad(blockKey, b.keyBits/8)
		c, err := b.f(key)
		if err != nil {
			log.Fatal(err)
			return
		}

		blockSize := c.BlockSize()

		input := pkcs7pad(input, blockSize)

		for i := 0; i < len(input); i += blockSize {
			c.Encrypt(dst[i:], input[i:])
		}
		for i := 0; i < len(input); i += blockSize {
			c.Decrypt(dst2[i:], dst[i:])
		}

		log.Println(b.name, hex.EncodeToString(input), "->", hex.EncodeToString(dst[:len(input)]), "->", hex.EncodeToString(dst2[:len(input)]))
	}
}

func printHashs() {
	Hashes := []struct {
		name string
		h    hash.Hash
	}{
		{"HAS-165", has160.New()},
		{"LSH-256-224", lsh256.New224()},
		{"LSH-256", lsh256.New()},
		{"LSH-512-224", lsh512.New224()},
		{"LSH-512-256", lsh512.New256()},
		{"LSH-512-384", lsh512.New384()},
		{"LSH-512", lsh512.New()},
	}

	dst := make([]byte, 100)

	for _, h := range Hashes {
		h.h.Write(input)
		dst = h.h.Sum(dst[:0])

		log.Println(h.name, hex.EncodeToString(input), "->", hex.EncodeToString(dst))
	}
}

func pkcs7pad(data []byte, blockSize int) []byte {
	padLen := blockSize - len(data)%blockSize
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)

	out := append(make([]byte, 0, blockSize), data...)
	return append(out, padding...)
}
