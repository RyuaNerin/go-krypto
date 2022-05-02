//go:build amd64

package lsh256

import (
	"hash"

	"github.com/RyuaNerin/go-krypto/internal/kryptoutil"
)

type simdSet struct {
	init   func(ctx *lsh256ContextAsmData)
	update func(ctx *lsh256ContextAsmData, data []byte, remain_msg_byte int) int
	final  func(ctx *lsh256ContextAsmData, hashval []byte)
}

var (
	simdSetDefault simdSet
)

func newContext(algType algType) hash.Hash {
	return newContextAsm(algType, simdSetDefault)
}

func newContextAsm(algType algType, simd simdSet) hash.Hash {
	ctx := new(lsh256ContextAsm)
	initContextAsm(ctx, algType, simd)
	return ctx
}

type lsh256ContextAsm struct {
	simd simdSet
	data lsh256ContextAsmData
}
type lsh256ContextAsmData struct {
	algtype            int
	remain_databytelen int
	cv_l               []byte
	cv_r               []byte
	last_block         []byte
}

func initContextAsm(ctx *lsh256ContextAsm, algtype algType, simd simdSet) {
	ctx.simd = simd
	ctx.data.algtype = int(algtype)
	ctx.data.cv_l = make([]byte, 32)
	ctx.data.cv_r = make([]byte, 32)
	ctx.data.last_block = make([]byte, 128) // LSH256_MSG_BLK_BYTE_LEN = 128
	ctx.Reset()
}

func (ctx *lsh256ContextAsm) Size() int {
	return ctx.data.algtype
}

func (ctx *lsh256ContextAsm) BlockSize() int {
	return BlockSize
}

func (ctx *lsh256ContextAsm) Reset() {
	ctx.data.remain_databytelen = 0
	ctx.simd.init(&ctx.data)
}

func (ctx *lsh256ContextAsm) Write(data []byte) (n int, err error) {
	databytelen := len(data)
	remain_msg_byte := ctx.data.remain_databytelen

	// LSH256_MSG_BLK_BYTE_LEN = 128
	if databytelen+remain_msg_byte < 128 {
		copy(ctx.data.last_block[remain_msg_byte:], data)
		ctx.data.remain_databytelen += databytelen
		return len(data), nil
	}

	// go 에서 사전 처리
	if remain_msg_byte > 0 {
		more_byte := 128 - remain_msg_byte
		copy(ctx.data.last_block[remain_msg_byte:], data[:more_byte])
	}

	databytelen = ctx.simd.update(&ctx.data, data, remain_msg_byte)

	if databytelen > 0 {
		copy(ctx.data.last_block, data[:databytelen])
		ctx.data.remain_databytelen = databytelen
	}

	return len(data), nil
}

func (ctx *lsh256ContextAsm) Sum(b []byte) []byte {
	remain_msg_byte := ctx.data.remain_databytelen
	ctx.data.last_block[remain_msg_byte] = 0x80
	kryptoutil.MemsetByte(ctx.data.last_block[remain_msg_byte+1:], 0)

	hash := make([]byte, Size)
	ctx.simd.final(&ctx.data, hash)

	return append(b, hash[:ctx.Size()]...)
}
