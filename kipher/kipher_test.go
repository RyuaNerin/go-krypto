package kipher_test

import (
	"crypto/cipher"
	"testing"

	. "github.com/RyuaNerin/go-krypto/internal/cavp"
	"github.com/RyuaNerin/go-krypto/lea"
)

const blockSize = lea.BlockSize

func newLEABlock() cipher.Block {
	b, err := lea.NewCipher(make([]byte, blockSize))
	if err != nil {
		panic(err)
	}
	return b
}

const (
	mmtLEA128CBC    = "03. LEA/LEA-128_(CBC)_MMT"
	mmtLEA128CFB8   = "03. LEA/LEA-128_(CFB8)_MMT"
	mmtLEA128CFB64  = "03. LEA/LEA-128_(CFB64)_MMT"
	mmtLEA128CFB128 = "03. LEA/LEA-128_(CFB128)_MMT"
	mmtLEA128CTR    = "03. LEA/LEA-128_(CTR)_MMT"
	mmtLEA128ECB    = "03. LEA/LEA-128_(ECB)_MMT"
	mmtLEA128OFB    = "03. LEA/LEA-128_(OFB)_MMT"
)

func Test_CBC_MMT(t *testing.T)    { TestBlockKAT(t, mmtLEA128CBC, CBC, lea.NewCipher) }
func Test_CFB8_MMT(t *testing.T)   { TestBlockKAT(t, mmtLEA128CFB8, CFB8, lea.NewCipher) }
func Test_CFB64_MMT(t *testing.T)  { TestBlockKAT(t, mmtLEA128CFB64, CFB64, lea.NewCipher) }
func Test_CFB128_MMT(t *testing.T) { TestBlockKAT(t, mmtLEA128CFB128, CFB128, lea.NewCipher) }
func Test_CTR_MMT(t *testing.T)    { TestBlockKAT(t, mmtLEA128CTR, CTR, lea.NewCipher) }
func Test_ECB_MMT(t *testing.T)    { TestBlockKAT(t, mmtLEA128ECB, ECB, lea.NewCipher) }
func Test_OFB_MMT(t *testing.T)    { TestBlockKAT(t, mmtLEA128OFB, OFB, lea.NewCipher) }
