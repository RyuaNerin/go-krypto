package kipher

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"io"
	"testing"
)

func TestCCM(t *testing.T) {
	input := make([]byte, 1024)
	nonce := make([]byte, 12)
	additionalData := make([]byte, 1024)

	io.ReadFull(rand.Reader, input)
	io.ReadFull(rand.Reader, nonce)
	io.ReadFull(rand.Reader, additionalData)

	key := make([]byte, 256/8)
	b, err := aes.NewCipher(key)
	if err != nil {
		t.Error(err)
		return
	}

	aead, err := NewCCM(b, len(nonce), 16)
	if err != nil {
		t.Error(err)
		return
	}

	seal := aead.Seal(nil, nonce, input, additionalData)

	output, err := aead.Open(nil, nonce, seal, additionalData)
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(input, output) {
		t.Error("output not equals to input")
		return
	}

	_, err = aead.Open(nil, nonce, seal, additionalData[:len(additionalData)-1])
	if err == nil {
		t.Error("unexpected success")
		return
	}
}
