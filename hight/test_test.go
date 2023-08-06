package hight

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func Test_Encrypt(t *testing.T) {
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

func Test_Decript(t *testing.T) {
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
