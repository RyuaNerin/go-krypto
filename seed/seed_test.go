package seed

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
			decode("00000000 00000000 00000000 00000000"),
			decode("00010203 04050607 08090A0B 0C0D0E0F"),
			decode("5EBAC6E0 054E1668 19AFF1CC 6D346CDB"),
		},
		{
			decode("00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000"),
			decode("00010203 04050607 08090A0B 0C0D0E0F"),
			decode("C609214B E64E38CB EC8E8F0A FEBA74DF"),
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
