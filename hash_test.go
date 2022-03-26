package krypto

import (
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"io/fs"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/RyuaNerin/go-krypto/lsh256"
	"github.com/RyuaNerin/go-krypto/lsh512"
	"github.com/RyuaNerin/go-krypto/test"
)

func Test_SHA2_224(t *testing.T) {
	testHash(t, `SHA2\(224\)(?P<test>[^\.]+).txt`, sha256.New224)
}
func Test_SHA2_256(t *testing.T) {
	testHash(t, `SHA2\(256\)(?P<test>[^\.]+).txt`, sha256.New)
}
func Test_SHA2_384(t *testing.T) {
	testHash(t, `SHA2\(384\)(?P<test>[^\.]+).txt`, sha512.New384)
}
func Test_SHA2_512(t *testing.T) {
	testHash(t, `SHA2\(512\)(?P<test>[^\.]+).txt`, sha512.New)
}

func Test_LSH256_224(t *testing.T) {
	testHash(t, `LSH\(256-224\)(?P<test>[^\.]+).txt`, lsh256.New224)
}
func Test_LSH256(t *testing.T) {
	testHash(t, `LSH\(256-256\)(?P<test>[^\.]+).txt`, lsh256.New)
}

func Test_LSH512_224(t *testing.T) {
	testHash(t, `LSH\(512-224\)(?P<test>[^\.]+).txt`, lsh512.New224)
}

func Test_LSH512_256(t *testing.T) {
	testHash(t, `LSH\(512-256\)(?P<test>[^\.]+).txt`, lsh512.New256)
}

func Test_LSH512_384(t *testing.T) {
	testHash(t, `LSH\(512-384\)(?P<test>[^\.]+).txt`, lsh512.New384)
}

func Test_LSH512_512(t *testing.T) {
	testHash(t, `LSH\(512-512\)(?P<test>[^\.]+).txt`, lsh512.New)
}

func testHash(t *testing.T, regexStr string, newHash func() hash.Hash) {
	re := regexp.MustCompile(regexStr)

	filepath.Walk(
		"test",
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			var fn test.HashTestFunc

			m := re.FindStringSubmatch(filepath.Base(path))
			if m == nil {
				return nil
			}
			for i, name := range re.SubexpNames() {
				switch name {
				case "test":
					switch m[i] {
					case "LongMsg":
						fn = test.HashTest
					case "ShortMsg":
						fn = test.HashTest
					}
				}
			}

			if fn != nil {
				fn(t, path, newHash)
			}

			return nil
		},
	)
}
