package hight

import (
	"bytes"
	"encoding/hex"
	"strings"
	"testing"
)

var (
	testCase = []struct {
		key       []byte
		plain     []byte
		encrypted []byte
	}{
		{
			decode("88E34F8F 081779F1 E9F39437 0AD40589"),
			decode("D76D0D18 327EC562"),
			decode("E4BC2E31 2277E4DD"),
		},
	}
)

func decode(s string) []byte {
	b, _ := hex.DecodeString(strings.ReplaceAll(s, " ", ""))
	return b
}

func TestEncrypt(t *testing.T) {
	for i, tc := range testCase {
		c, err := NewCipher(tc.key)
		if err != nil {
			t.Errorf("[%d] %+v", i, err)
			continue
		}

		dst := make([]byte, c.BlockSize())
		c.Encrypt(dst, tc.plain)

		if !bytes.Equal(dst, tc.encrypted) {
			t.Errorf("[%d] Did not match\nTest : %sWant : %s", i, hex.Dump(dst), hex.Dump(tc.encrypted))
		}
	}
}

func TestDecrypt(t *testing.T) {
	for i, tc := range testCase {
		c, err := NewCipher(tc.key)
		if err != nil {
			t.Errorf("[%d] %+v", i, err)
			continue
		}

		dst := make([]byte, c.BlockSize())
		c.Decrypt(dst, tc.encrypted)

		if !bytes.Equal(dst, tc.plain) {
			t.Errorf("[%d] Did not match\nTest : %sWant : %s", i, hex.Dump(dst), hex.Dump(tc.plain))
		}
	}
}
