package hight

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
)

// TTAK.KO-12.0040_R1

var (
	testCases = []struct {
		Key    []byte
		Plain  []byte
		Secure []byte
	}{
		// p. 21
		// Ⅰ.1. 참조구현값 1
		{
			Key:    internal.Reverse(internal.HB(`00 11 22 33 44 55 66 77 88 99 aa bb cc dd ee ff`)),
			Plain:  internal.Reverse(internal.HB(`00 00 00 00 00 00 00 00`)),
			Secure: internal.Reverse(internal.HB(`00 f4 18 ae d9 4f 03 f2`)),
		},
		// p. 22
		// Ⅰ.2. 참조구현값 2
		{
			Key:    internal.Reverse(internal.HB(`ff ee dd cc bb aa 99 88 77 66 55 44 33 22 11 00`)),
			Plain:  internal.Reverse(internal.HB(`00 11 22 33 44 55 66 77`)),
			Secure: internal.Reverse(internal.HB(`23 ce 9f 72 e5 43 e6 d8`)),
		},
		// p. 23
		// Ⅰ.3. 참조구현값 3
		{
			Key:    internal.Reverse(internal.HB(`00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f`)),
			Plain:  internal.Reverse(internal.HB(`01 23 45 67 89 ab cd ef`)),
			Secure: internal.Reverse(internal.HB(`7a 6f b2 a2 8d 23 f4 66`)),
		},
		// p. 24
		// Ⅰ.4. 참조구현값 4
		{
			Key:    internal.Reverse(internal.HB(`28 db c3 bc 49 ff d8 7d cf a5 09 b1 1d 42 2b e7`)),
			Plain:  internal.Reverse(internal.HB(`b4 1e 6b e2 eb a8 4a 14`)),
			Secure: internal.Reverse(internal.HB(`cc 04 7a 75 20 9c 1f c6`)),
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
