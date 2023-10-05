package lea

import "crypto/cipher"

func newCipherSimple(key []byte) (cipher.Block, error) {
	block, err := NewCipher(key)
	if err != nil {
		return nil, err
	}
	return &blockWrap{block}, nil
}

type blockWrap struct {
	b cipher.Block
}

func (bw *blockWrap) BlockSize() int {
	return bw.b.BlockSize()
}

func (bw *blockWrap) Encrypt(dst, src []byte) {
	bw.b.Encrypt(dst, src)
}

func (bw *blockWrap) Decrypt(dst, src []byte) {
	bw.b.Decrypt(dst, src)
}
