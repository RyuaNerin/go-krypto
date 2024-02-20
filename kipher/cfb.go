package kipher

// Based on https://github.com/golang/go/blob/go1.21.6/src/crypto/cipher/cfb.go

import (
	"crypto/cipher"

	"github.com/RyuaNerin/go-krypto/internal/alias"
	"github.com/RyuaNerin/go-krypto/internal/subtle"
)

// NewCFBEncrypter returns a Stream which encrypts with cipher feedback mode,
// using the given Block. The iv must be the same length as the Block's block
// size.
func NewCFBEncrypter(block cipher.Block, iv []byte, cfbBlockByteSize int) cipher.Stream {
	if block.BlockSize() == cfbBlockByteSize {
		return cipher.NewCFBEncrypter(block, iv)
	}
	return newCFB(block, iv, false, cfbBlockByteSize)
}

// NewCFB8Decrypter returns a Stream which decrypts with cipher feedback mode,
// using the given Block. The iv must be the same length as the Block's block
// size.
func NewCFBDecrypter(block cipher.Block, iv []byte, cfbBlockByteSize int) cipher.Stream {
	if block.BlockSize() == cfbBlockByteSize {
		return cipher.NewCFBDecrypter(block, iv)
	}
	return newCFB(block, iv, true, cfbBlockByteSize)
}

func newCFB(block cipher.Block, iv []byte, decrypt bool, cfbBlockByteSize int) cipher.Stream {
	blockSize := block.BlockSize()
	if len(iv) != blockSize {
		// stack trace will indicate whether it was de or encryption
		panic("kipher.newCFB: IV length must equal block size")
	}
	x := &cfb{
		b:       block,
		out:     make([]byte, blockSize),
		next:    make([]byte, blockSize),
		outUsed: cfbBlockByteSize,
		bytes:   cfbBlockByteSize,
		decrypt: decrypt,
	}
	copy(x.next, iv)

	return x
}

type cfb struct {
	b       cipher.Block
	next    []byte
	out     []byte
	outUsed int
	decrypt bool

	bytes int
}

func (x *cfb) XORKeyStream(dst, src []byte) {
	if len(dst) < len(src) {
		panic("krypto/kipher: output smaller than input")
	}
	if alias.InexactOverlap(dst[:len(src)], src) {
		panic("krypto/kipher: invalid buffer overlap")
	}
	for len(src) > 0 {
		if x.outUsed == x.bytes {
			x.b.Encrypt(x.out, x.next)
			copy(x.next, x.next[x.bytes:])
			x.outUsed = 0
		}

		if x.decrypt {
			// We can precompute a larger segment of the
			// keystream on decryption. This will allow
			// larger batches for xor, and we should be
			// able to match CTR/OFB performance.
			copy(x.next[len(x.next)-x.bytes+x.outUsed:], src)
		}
		n := subtle.XORBytes(dst, src, x.out[x.outUsed:x.bytes])
		if !x.decrypt {
			copy(x.next[len(x.next)-x.bytes+x.outUsed:], dst)
		}
		dst = dst[n:]
		src = src[n:]
		x.outUsed += n
	}
}
