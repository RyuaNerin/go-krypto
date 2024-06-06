package kipher

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"encoding/hex"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
	igcm "github.com/RyuaNerin/go-krypto/internal/gcm"
)

func Test_GCM(t *testing.T) {
	const iter = 256
	const maxLen = 16 * igcm.GCMBlockSize * igcm.GCMBlockSize

	key := make([]byte, 32)

	nonce := make([]byte, maxLen)
	input := make([]byte, maxLen)
	additional := make([]byte, maxLen)

	dstCipher := make([]byte, maxLen+igcm.GCMBlockSize)
	dstKipher := make([]byte, maxLen+igcm.GCMBlockSize)

	openedKipher := make([]byte, maxLen)

	nextLen := func(max int) int {
		var lRaw [4]byte
		rnd.Read(lRaw[:])
		lRaw[0] &= 0x7F
		return int(binary.BigEndian.Uint32(lRaw[:])) % max
	}

	for i := 0; i < iter; i++ {
		keySize := 16 + nextLen(2)*8
		nonceSize := 1 + nextLen(maxLen-1)
		inputSize := 1 + nextLen(maxLen-1)
		additionalSize := nextLen(maxLen)

		rnd.Read(key[:keySize])
		rnd.Read(nonce[:nonceSize])
		rnd.Read(input[:inputSize])
		rnd.Read(additional[:additionalSize])

		b, _ := aes.NewCipher(key[:keySize])
		kb := internal.WrapBlock(b) // ignore gcmAble

		gcmCipher, err := cipher.NewGCMWithNonceSize(kb, nonceSize)
		if err != nil {
			t.Error(err)
			return
		}
		gcmKipher, err := newGCMWithNonceAndTagSize(kb, nonceSize, igcm.GCMTagSize)
		if err != nil {
			t.Error(err)
			return
		}

		dstCipher = gcmCipher.Seal(dstCipher[:0], nonce[:nonceSize], input[:inputSize], additional[:additionalSize])
		dstKipher = gcmKipher.Seal(dstKipher[:0], nonce[:nonceSize], input[:inputSize], additional[:additionalSize])

		if !bytes.Equal(dstKipher, dstCipher) {
			t.Errorf("failed to Seal\nexpect: %s\nactual: %s", hex.EncodeToString(dstCipher), hex.EncodeToString(dstKipher))
			return
		}

		openedKipher, err = gcmKipher.Open(openedKipher[:0], nonce[:nonceSize], dstKipher, additional[:additionalSize])
		if err != nil {
			t.Error(err)
			return
		}
		if !bytes.Equal(openedKipher, input[:inputSize]) {
			t.Errorf("failed to Open\nexpect: %v\nactual: %v", hex.EncodeToString(input[:inputSize]), hex.EncodeToString(openedKipher))
			return
		}
	}
}
