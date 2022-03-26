package krypto

import (
	"crypto/cipher"
	"io/fs"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/RyuaNerin/go-krypto/aria"
	"github.com/RyuaNerin/go-krypto/hight"
	"github.com/RyuaNerin/go-krypto/lea"
	"github.com/RyuaNerin/go-krypto/seed"
	"github.com/RyuaNerin/go-krypto/test"
)

func Test_ARIA128(t *testing.T) {
	testBlock(t, `ARIA128\((?P<block>[^\)]+)\)(?P<test>[^\.]+)\.txt`, aria.NewCipher)
}
func Test_ARIA192(t *testing.T) {
	testBlock(t, `ARIA192\((?P<block>[^\)]+)\)(?P<test>[^\.]+)\.txt`, aria.NewCipher)
}
func Test_ARIA256(t *testing.T) {
	testBlock(t, `ARIA256\((?P<block>[^\)]+)\)(?P<test>[^\.]+)\.txt`, aria.NewCipher)
}

func Test_LEA128(t *testing.T) {
	testBlock(t, `LEA128\((?P<block>[^\)]+)\)(?P<test>[^\.]+)\.txt`, lea.NewCipher)
}
func Test_LEA192(t *testing.T) {
	testBlock(t, `LEA192\((?P<block>[^\)]+)\)(?P<test>[^\.]+)\.txt`, lea.NewCipher)
}
func Test_LEA256(t *testing.T) {
	testBlock(t, `LEA256\((?P<block>[^\)]+)\)(?P<test>[^\.]+)\.txt`, lea.NewCipher)
}
func Test_SEED128(t *testing.T) {
	testBlock(t, `SEED128\((?P<block>[^\)]+)\)(?P<test>[^\.]+)\.txt`, seed.NewCipher)
}

/*
func Test_SEED256(t *testing.T) {
	testBlock(t, "test/SEED256\.txt", seed.NewCipher)
}
*/

func Test_HIGHT(t *testing.T) {
	testBlock(t, `HIGHT\((?P<block>[^\)]+)\)(?P<test>[^\.]+)\.txt`, hight.NewCipher)
}

func testBlock(t *testing.T, regexStr string, newBlock func(key []byte) (cipher.Block, error)) {
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

			var blockMode test.CipherMode
			var fn test.BlockTestFunc

			m := re.FindStringSubmatch(filepath.Base(path))
			if m == nil {
				return nil
			}
			for i, name := range re.SubexpNames() {
				switch name {
				case "block":
					switch m[i] {
					case "ECB":
						blockMode = test.CipherModeECB
					case "CBC":
						blockMode = test.CipherModeCBC
					case "OFB":
						blockMode = test.CipherModeOFB
					case "CTR":
						blockMode = test.CipherModeCTR
					}
				case "test":
					switch m[i] {
					case "KAT":
						fn = test.BlockTest
					case "MMT":
						fn = test.BlockTest
					}
				}
			}

			if blockMode != 0 && fn != nil {
				fn(t, path, blockMode, newBlock)
			}

			return nil
		},
	)
}
