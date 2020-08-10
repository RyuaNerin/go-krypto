package aria

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
			decode("00112233 44556677 8899aabb ccddeeff 00112233 44556677"),
			decode("11111111 aaaaaaaa 11111111 bbbbbbbb"),
			decode("8d147062 5f59ebac b0e55b53 4b3e462b"),
		},
		{
			decode("00010203 04050607 08090a0b 0c0d0e0f"),
			decode("00112233 44556677 8899aabb ccddeeff"),
			decode("d718fbd6 ab644c73 9da95f3b e6451778"),
		},
		{
			decode("00112233 44556677 8899aabb ccddeeff 00112233 44556677 8899aabb ccddeeff"),
			decode("11111111 aaaaaaaa 11111111 bbbbbbbb"),
			decode("58a875e6 044ad7ff fa4f5842 0f7f442d"),
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
