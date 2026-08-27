package kipher_test

import (
	"crypto/cipher"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/testing/tblockmode"
	"github.com/RyuaNerin/go-krypto/kipher"
)

var ecbInfo = tblockmode.BlockModeInfo{BlockSize: blockSize}

func newECBE(_ []byte) cipher.BlockMode { return kipher.NewECBEncryptor(newLEABlock()) }
func newECBD(_ []byte) cipher.BlockMode { return kipher.NewECBDecryptor(newLEABlock()) }

func Test_ECB_New(t *testing.T)       { ecbInfo.TestNew(t, newECBE, newECBD, false) }
func Test_ECB_BlockMode(t *testing.T) { ecbInfo.TestBlockMode(t, newECBE, newECBD, false) }

func Fuzz_ECB_New(f *testing.F)         { ecbInfo.FuzzNew(f, newECBE, newECBD, false) }
func Fuzz_ECB_BlockMode(f *testing.F)   { ecbInfo.FuzzBlockMode(f, newECBE, newECBD, false) }
func Fuzz_ECB_ShortBuffer(f *testing.F) { ecbInfo.FuzzBlockModeShortBuffer(f, newECBE, newECBD, false) }

func Benchmark_ECB_New(b *testing.B)       { ecbInfo.BenchNew(b, newECBE, newECBD, false) }
func Benchmark_ECB_BlockMode(b *testing.B) { ecbInfo.BenchBlockMode(b, newECBE, newECBD, false) }
