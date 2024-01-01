package lea

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"testing"
)

const (
	testKey    = 1024 * 1024
	testBlocks = 16 * 1024 * 1024
)

func Benchmark_Encrypt_1Blocks(b *testing.B) { benchAll(b, bb(1, leaEnc1Go, false)) }

func Benchmark_Decrypt_1Blocks(b *testing.B) { benchAll(b, bb(1, leaDec1Go, false)) }

var rnd = bufio.NewReaderSize(rand.Reader, 1<<15)

type testCase struct {
	Key    []byte
	Plain  []byte
	Secure []byte
}

func testAll(t *testing.T, f func(*testing.T, int)) {
	tests := []struct {
		name    string
		keySize int
	}{
		{"128", 128},
		{"196", 196},
		{"256", 256},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			f(t, test.keySize)
		})
	}
}

func testEncryptGo(t *testing.T, testCases []testCase) {
	dst := make([]byte, BlockSize)

	var ctx leaContext
	for _, tc := range testCases {
		err := ctx.initContext(tc.Key)
		if err != nil {
			t.Error(err)
		}

		leaEnc1Go(&ctx, dst, tc.Plain)
		if !bytes.Equal(dst, tc.Secure) {
			t.Errorf("encrypt failed.\nresult: %s\nanswer: %s", hex.EncodeToString(dst), hex.EncodeToString(tc.Secure))
		}
	}
}

func testDecryptGo(t *testing.T, testCases []testCase) {
	dst := make([]byte, BlockSize)

	var ctx leaContext
	for _, tc := range testCases {
		err := ctx.initContext(tc.Key)
		if err != nil {
			t.Error(err)
		}

		leaDec1Go(&ctx, dst, tc.Secure)
		if !bytes.Equal(dst, tc.Plain) {
			t.Errorf("encrypt failed.\nresult: %s\nanswer: %s", hex.EncodeToString(dst), hex.EncodeToString(tc.Plain))
		}
	}
}

func tb(blocks int, funcGo, funcAsm funcBlock, skip bool) func(t *testing.T, keySize int) {
	return func(t *testing.T, keySize int) {
		if skip {
			t.Skip()
			return
		}

		k := make([]byte, keySize/8)
		rnd.Read(k)

		srcGo := make([]byte, BlockSize*blocks)
		dstGo := make([]byte, BlockSize*blocks)

		srcAsm := make([]byte, BlockSize*blocks)
		dstAsm := make([]byte, BlockSize*blocks)

		rnd.Read(srcGo)
		copy(srcAsm, srcGo)

		var ctx leaContextAsm
		err := ctx.g.initContext(k)
		if err != nil {
			t.Error(err)
			return
		}

		for i := 0; i < testBlocks/blocks; i++ {
			funcGo(&ctx.g, dstGo, srcGo)
			funcAsm(&ctx.g, dstAsm, srcAsm)

			if !bytes.Equal(dstGo, dstAsm) {
				t.Error("did not match")
				t.FailNow()
				return
			}

			copy(srcGo, dstGo)
			copy(srcAsm, dstAsm)
		}
	}
}

func benchAll(b *testing.B, f func(*testing.B, int)) {
	tests := []struct {
		name    string
		keySize int
	}{
		{"128", 128},
		{"196", 196},
		{"256", 256},
	}
	for _, test := range tests {
		test := test
		b.Run(test.name, func(b *testing.B) {
			f(b, test.keySize)
		})
	}
}

func bb(blocks int, f funcBlock, skip bool) func(b *testing.B, keySize int) {
	return func(b *testing.B, keySize int) {
		if skip {
			b.Skip()
			return
		}

		key := make([]byte, keySize/8)
		rand.Read(key)

		src := make([]byte, BlockSize*blocks)
		dst := make([]byte, BlockSize*blocks)
		rand.Read(src)

		var ctx leaContext
		err := ctx.initContext(key)
		if err != nil {
			b.Error(err)
			return
		}

		b.ReportAllocs()
		b.SetBytes(int64(len(src)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			f(&ctx, dst, src)
			copy(src, dst)
		}
	}
}
