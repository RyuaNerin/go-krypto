package seed

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

func Test_SEED128_Encrypt(t *testing.T) { testEncrypt(t, testCases128) }

func Test_SEED128_Decrypt(t *testing.T) { testDecrypt(t, testCases128) }

func testEncrypt(t *testing.T, testCases []testCase) {
	dst := make([]byte, BlockSize)

	for _, tc := range testCases128 {
		c, err := NewCipher(tc.Key)
		if err != nil {
			t.Error(err)
		}

		c.Encrypt(dst, tc.Plain)
		if !bytes.Equal(dst, tc.Secure) {
			t.Errorf("encrypt failed.\nresult: %s\nanswer: %s", hex.EncodeToString(dst), hex.EncodeToString(tc.Secure))
		}
	}
}
func testDecrypt(t *testing.T, testCases []testCase) {
	dst := make([]byte, BlockSize)

	for _, tc := range testCases128 {
		c, err := NewCipher(tc.Key)
		if err != nil {
			t.Error(err)
		}

		c.Decrypt(dst, tc.Secure)
		if !bytes.Equal(dst, tc.Plain) {
			t.Errorf("decrypt failed.\nresult: %s\nanswer: %s", hex.EncodeToString(dst), hex.EncodeToString(tc.Plain))
		}
	}
}
