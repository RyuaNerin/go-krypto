package seed

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
)

// TTAS.KO-12.0004/R1

var (
	testCases = []struct {
		Key    []byte
		Plain  []byte
		Secure []byte
	}{
		// p. 21
		// Ⅰ.1. 참조구현값 1
		{
			Key:    internal.HB(`00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00`),
			Plain:  internal.HB(`00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F`),
			Secure: internal.HB(`5E BA C6 E0 05 4E 16 68 19 AF F1 CC 6D 34 6C DB`),
		},
		// p. 21
		// Ⅰ.2. 참조구현값 2
		{
			Key:    internal.HB(`00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F`),
			Plain:  internal.HB(`00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00`),
			Secure: internal.HB(`C1 1F 22 F2 01 40 50 50 84 48 35 97 E4 37 0F 43`),
		},
		// p. 22
		// Ⅰ.3. 참조구현값 3
		{
			Key:    internal.HB(`47 06 48 08 51 E6 1B E8 5D 74 BF B3 FD 95 61 85`),
			Plain:  internal.HB(`83 A2 F8 A2 88 64 1F B9 A4 E9 A5 CC 2F 13 1C 7D`),
			Secure: internal.HB(`EE 54 D1 3E BC AE 70 6D 22 6B C3 14 2C D4 0D 4A`),
		},
		// p. 24
		// Ⅰ.4. 참조구현값 4
		{
			Key:    internal.HB(`28 DB C3 BC 49 FF D8 7D CF A5 09 B1 1D 42 2B E7`),
			Plain:  internal.HB(`B4 1E 6B E2 EB A8 4A 14 8E 2E ED 84 59 3C 5E C7`),
			Secure: internal.HB(`9B 9B 7B FC D1 81 3C B9 5D 0B 36 18 F4 0F 51 22`),
		},
	}
)

func TestEncryptDecrypt(t *testing.T) {
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

func TestEncrypt(t *testing.T) {
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

func TestDecrypt(t *testing.T) {
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
