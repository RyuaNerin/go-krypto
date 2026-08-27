package kipher

import (
	"crypto/cipher"
	"crypto/subtle"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/alias"
)

func NewCFBEncrypter(block krypto.Block, iv []byte, cfbBlockByteSize int) *CFBEncrypter {
	blockSize := block.BlockSize()

	on := make([]byte, blockSize*2)
	x := &CFBEncrypter{
		cipher:  block,
		out:     on[:blockSize],
		next:    on[blockSize:],
		outUsed: cfbBlockByteSize,
		bytes:   cfbBlockByteSize,
	}
	copy(x.next, iv)

	return x
}

func NewCFBDecrypter(block krypto.Block8, iv []byte, cfbBlockByteSize int) *CFBDecrypter {
	blockSize := block.BlockSize()

	on := make([]byte, blockSize*(8+9))
	x := &CFBDecrypter{
		cipher:    block,
		blockSize: blockSize,
		out:       on[:8*blockSize],
		next:      on[8*blockSize:],
		outUsed:   cfbBlockByteSize,
		bytes:     cfbBlockByteSize,
	}

	copy(x.next[8*blockSize:], iv)

	return x
}

type (
	CFBEncrypter struct {
		cipher cipher.Block

		next    []byte
		out     []byte
		outUsed int

		bytes int
	}
	CFBDecrypter struct {
		cipher    krypto.Block8
		blockSize int

		next    []byte // [N블록 입력 | 피드백 1블록] = 9블록
		out     []byte // 정확히 N블록
		outUsed int

		bytes int
	}
)

var (
	_ krypto.Stream = (*CFBDecrypter)(nil)
	_ krypto.Stream = (*CFBEncrypter)(nil)
)

func (x *CFBEncrypter) Zeroize() {
	x.cipher = nil
	clear(x.next)
	x.next = nil
	clear(x.out)
	x.out = nil
	x.outUsed = 0
	x.bytes = 0
}

func (x *CFBDecrypter) Zeroize() {
	x.cipher = nil
	x.blockSize = 0

	clear(x.next)
	x.next = nil
	clear(x.out)
	x.out = nil
	x.outUsed = 0
	x.bytes = 0
}

func (x *CFBEncrypter) XORKeyStream(dst, src []byte) {
	if len(dst) < len(src) {
		panic(msgSmallDst)
	}
	if alias.InexactOverlap(dst[:len(src)], src) {
		panic(msgBufferOverlap)
	}
	for len(src) > 0 {
		if x.outUsed == x.bytes {
			x.cipher.Encrypt(x.out, x.next)
			copy(x.next, x.next[x.bytes:])
			x.outUsed = 0
		}

		n := subtle.XORBytes(dst, src, x.out[x.outUsed:x.bytes])
		copy(x.next[len(x.next)-x.bytes+x.outUsed:], dst)
		dst = dst[n:]
		src = src[n:]
		x.outUsed += n
	}
}

func (x *CFBDecrypter) XORKeyStream(dst, src []byte) {
	if len(dst) < len(src) {
		panic(msgSmallDst)
	}
	if alias.InexactOverlap(dst[:len(src)], src) {
		panic(msgBufferOverlap)
	}

	for len(src) > 0 && x.outUsed != x.bytes {
		n := min(len(src), x.bytes-x.outUsed)
		feedback := x.next[8*x.blockSize:]
		copy(feedback[len(feedback)-x.bytes+x.outUsed:], src[:n])
		subtle.XORBytes(dst[:n], src[:n], x.out[x.outUsed:x.outUsed+n])

		dst = dst[n:]
		src = src[n:]
		x.outUsed += n
	}

	for len(src) >= 8*x.bytes {
		x.decryptBlocks8(dst, src)

		n := 8 * x.bytes
		dst = dst[n:]
		src = src[n:]
	}

	for len(src) >= 4*x.bytes {
		x.decryptBlocks4(dst, src)

		n := 4 * x.bytes
		dst = dst[n:]
		src = src[n:]
	}

	for len(src) > 0 {
		if x.outUsed == x.bytes {
			feedback := x.next[8*x.blockSize:]
			x.cipher.Encrypt(x.out[:x.blockSize], feedback)
			copy(feedback, feedback[x.bytes:])

			x.outUsed = 0
		}

		n := min(len(src), x.bytes-x.outUsed)
		feedback := x.next[8*x.blockSize:]
		copy(feedback[len(feedback)-x.bytes+x.outUsed:], src[:n])
		subtle.XORBytes(dst[:n], src[:n], x.out[x.outUsed:x.outUsed+n])

		dst = dst[n:]
		src = src[n:]
		x.outUsed += n
	}
}

func (x *CFBDecrypter) decryptBlocks8(dst, src []byte) {
	blockSize := x.blockSize
	feedback := x.next[8*blockSize:]

	copy(x.next[:blockSize], feedback)
	for i := 1; i < 8; i++ {
		previous := x.next[(i-1)*blockSize : i*blockSize]
		current := x.next[i*blockSize : (i+1)*blockSize]
		copy(current, previous[x.bytes:])
		copy(current[blockSize-x.bytes:], src[(i-1)*x.bytes:i*x.bytes])
	}

	last := x.next[7*blockSize : 8*blockSize]
	copy(feedback, last[x.bytes:])
	copy(feedback[blockSize-x.bytes:], src[7*x.bytes:8*x.bytes])

	x.cipher.Encrypt8(x.out, x.next[:8*blockSize])

	for i := range 8 {
		offset := i * x.bytes
		keystream := x.out[i*blockSize : i*blockSize+x.bytes]
		subtle.XORBytes(dst[offset:offset+x.bytes], src[offset:offset+x.bytes], keystream)
	}
}

func (x *CFBDecrypter) decryptBlocks4(dst, src []byte) {
	blockSize := x.blockSize
	feedback := x.next[8*blockSize:]

	copy(x.next[:blockSize], feedback)
	for i := 1; i < 4; i++ {
		previous := x.next[(i-1)*blockSize : i*blockSize]
		current := x.next[i*blockSize : (i+1)*blockSize]
		copy(current, previous[x.bytes:])
		copy(current[blockSize-x.bytes:], src[(i-1)*x.bytes:i*x.bytes])
	}

	last := x.next[3*blockSize : 4*blockSize]
	copy(feedback, last[x.bytes:])
	copy(feedback[blockSize-x.bytes:], src[3*x.bytes:4*x.bytes])

	x.cipher.Encrypt4(x.out[:4*blockSize], x.next[:4*blockSize])

	for i := range 4 {
		offset := i * x.bytes
		keystream := x.out[i*blockSize : i*blockSize+x.bytes]
		subtle.XORBytes(dst[offset:offset+x.bytes], src[offset:offset+x.bytes], keystream)
	}
}
