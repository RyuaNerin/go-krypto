package test

import (
	"bytes"
	"crypto/cipher"
	"encoding/hex"
	"io"
	"testing"
)

type CipherMode int

const (
	CipherModeECB CipherMode = iota
	CipherModeCBC
	CipherModeOFB
	CipherModeCTR
	CipherModeCFB
)

type BlockTestFunc func(t *testing.T, path string, cipherMode CipherMode, newBlock func(key []byte) (cipher.Block, error))

func BlockTest(t *testing.T, path string, cipherMode CipherMode, newBlock func(key []byte) (cipher.Block, error)) {
	r, err := newBlockTestCaseReader(path)
	if err != nil {
		t.Error(err)
		return
	}
	defer r.Close()

	var dst []byte
	for {
		count, key, iv, pt, ct, err := r.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Error(err)
			t.FailNow()
			return
		}

		block, err := newBlock(key)
		if err != nil {
			t.Error(err)
			t.FailNow()
			return
		}

		if len(dst) < len(pt) {
			dst = make([]byte, len(pt))
		}

		switch cipherMode {
		case CipherModeECB:
			for i := 0; i < len(pt); i += block.BlockSize() {
				block.Encrypt(dst[i:], pt[i:])
			}

		case CipherModeCBC:
			bm := cipher.NewCBCEncrypter(block, iv)
			bm.CryptBlocks(dst, pt)

		case CipherModeCFB:
			fallthrough
		case CipherModeCTR:
			fallthrough
		case CipherModeOFB:
			var stream cipher.Stream
			switch cipherMode {
			case CipherModeCFB:
				stream = cipher.NewCFBEncrypter(block, iv)
			case CipherModeCTR:
				stream = cipher.NewCTR(block, iv)
			case CipherModeOFB:
				stream = cipher.NewOFB(block, iv)
			}
			stream.XORKeyStream(dst, pt)
		}

		if !bytes.Equal(dst[:len(ct)], ct) {
			t.Errorf(
				`%s
COUNT : %d
KEY   : %s
IV    : %s
PT    : %s
CT    : %s
	
Eval  : %s`,
				path,
				count,
				hex.EncodeToString(key),
				hex.EncodeToString(iv),
				hex.EncodeToString(pt),
				hex.EncodeToString(ct),
				hex.EncodeToString(dst[:len(ct)]),
			)
			t.FailNow()
			return
		}
	}
}
