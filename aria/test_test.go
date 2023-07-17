package aria

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

func TestEncrypt128(t *testing.T) { testEncrypt(t, testCases128) }
func TestEncrypt196(t *testing.T) { testEncrypt(t, testCases196) }
func TestEncrypt256(t *testing.T) { testEncrypt(t, testCases256) }

func TestDecrypt128(t *testing.T) { testDecrypt(t, testCases128) }
func TestDecrypt196(t *testing.T) { testDecrypt(t, testCases196) }
func TestDecrypt256(t *testing.T) { testDecrypt(t, testCases256) }

func testEncryptDecrypt(t *testing.T, testCases []testCase) {
	plain := make([]byte, BlockSize)
	secure := make([]byte, BlockSize)

	for _, tc := range testCases {
		c, err := NewCipher(tc.Key)
		if err != nil {
			t.Error(err)
		}

		c.Encrypt(secure, tc.Plain)
		c.Decrypt(plain, secure)
		if !bytes.Equal(plain, tc.Plain) {
			t.Errorf("encrypt failed.\nresult: %s\nanswer: %s", hex.EncodeToString(plain), hex.EncodeToString(tc.Plain))
		}
	}
}

func testEncrypt(t *testing.T, testCases []testCase) {
	dst := make([]byte, BlockSize)

	for _, tc := range testCases {
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

	for _, tc := range testCases {
		c, err := NewCipher(tc.Key)
		if err != nil {
			t.Error(err)
		}

		c.Decrypt(dst, tc.Secure)
		if !bytes.Equal(dst, tc.Plain) {
			t.Errorf("decrypt failed.\nresult: %s\nanswer: %s", hex.EncodeToString(dst), hex.EncodeToString(tc.Secure))
		}
	}
}
