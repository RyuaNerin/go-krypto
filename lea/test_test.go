package lea

import (
	"bytes"
	"encoding/hex"
	"testing"
)

type testCase struct {
	Key    []byte
	Plain  []byte
	Secure []byte
}

func Test_LEA128_Encrypt_1Block_Go(t *testing.T) { testEncryptGo(t, testCases128) }
func Test_LEA196_Encrypt_1Block_Go(t *testing.T) { testEncryptGo(t, testCases196) }
func Test_LEA256_Encrypt_1Block_Go(t *testing.T) { testEncryptGo(t, testCases256) }

func Test_LEA128_Decrypt_1Block_Go(t *testing.T) { testDecryptGo(t, testCases128) }
func Test_LEA196_Decrypt_1Block_Go(t *testing.T) { testDecryptGo(t, testCases196) }
func Test_LEA256_Decrypt_1Block_Go(t *testing.T) { testDecryptGo(t, testCases256) }

func testEncryptGo(t *testing.T, testCases []testCase) {
	dst := make([]byte, BlockSize)

	var ctx leaContextGo
	for _, tc := range testCases {
		err := ctx.initContext(tc.Key)
		if err != nil {
			t.Error(err)
		}

		leaEnc1Go(ctx.round, ctx.rk, dst, tc.Plain)
		if !bytes.Equal(dst, tc.Secure) {
			t.Errorf("encrypt failed.\nresult: %s\nanswer: %s", hex.EncodeToString(dst), hex.EncodeToString(tc.Secure))
		}
	}
}

func testDecryptGo(t *testing.T, testCases []testCase) {
	dst := make([]byte, BlockSize)

	var ctx leaContextGo
	for _, tc := range testCases {
		err := ctx.initContext(tc.Key)
		if err != nil {
			t.Error(err)
		}

		leaDec1Go(ctx.round, ctx.rk, dst, tc.Secure)
		if !bytes.Equal(dst, tc.Plain) {
			t.Errorf("encrypt failed.\nresult: %s\nanswer: %s", hex.EncodeToString(dst), hex.EncodeToString(tc.Plain))
		}
	}
}
