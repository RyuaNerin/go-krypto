package lea

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
			decode("07AB6305 B025D83F 79ADDAA6 3AC8AD00"),
			decode("F28AE325 6AAD23B4 15E02806 3B610C60"),
			decode("64D908FC B7EBFEF9 0FD67010 6DE7C7C5"),
		},
		{
			decode("1437AF53 3069BD75 25C1560C 78BAD2A1 E534671C 007EF27C"),
			decode("1CB4F4CB 6C4BDB51 68EA8409 727BFD51"),
			decode("69725C6D F912F8B7 0EB511E6 663C5870"),
		},
		{
			decode("4F6779E2 BD1E9319 C63015AC FFEFD7A7 91F0ED59 DF1B7007 69FE82E2 F0668C35"),
			decode("DC31CAE3 DA5E0A11 C966B020 D7CFFEDE"),
			decode("EDA20420 98F667E8 57A02DB8 CAA7DFF2"),
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
