package gmac

import (
	"bytes"
	"crypto/aes"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"

	. "github.com/RyuaNerin/testingutil"
)

func Test_GMAC_ShortWrite(t *testing.T) {
	c, err := aes.NewCipher(make([]byte, 16))
	if err != nil {
		t.Fatal(err)
	}

	h, err := NewGMAC(c, nil)
	if err != nil {
		t.Fatal(err)
	}

	HTSW(t, h, false)
}

func Test_GMAC(t *testing.T) {
	var (
		/**
		https://cryptopp.com/wiki/GMAC

		$ ./test.exe
		Key: 00000000000000000000000000000000
		 IV: 00000000000000000000000000000000
		Message: Yoda said, do or do not. There is no try.
		GMAC: E7EE2C63B4DC328EED4A86B3FB3490AF
		*/
		key    = internal.HB(`00000000000000000000000000000000`)
		iv     = internal.HB(`00000000000000000000000000000000`)
		msg    = []byte(`Yoda said, do or do not. There is no try.`)
		expect = internal.HB(`E7EE2C63B4DC328EED4A86B3FB3490AF`)
	)

	c, err := aes.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}

	h, err := NewGMAC(c, iv)
	if err != nil {
		t.Fatal(err)
	}

	h.Write(msg)
	actual := h.Sum(nil)
	if !bytes.Equal(actual, expect) {
		t.Errorf("expect: %x, actual: %x", expect, actual)
		return
	}

	h.Reset()
	h.Write(msg)
	actual = h.Sum(actual[:0])
	if !bytes.Equal(actual, expect) {
		t.Errorf("expect: %x, actual: %x", expect, actual)
		return
	}
}
