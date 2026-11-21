package lsh512

import (
	"encoding/binary"
	"hash"
)

func newContextGo(size int) hash.Hash {
	ctx := &lsh512Context{
		lsh512ContextData: lsh512ContextData{
			algType: uint64(size),
		},
	}
	ctx.z.SetValid()
	ctx.Reset()

	return ctx
}

func sumGo(size int, data []byte) [Size]byte {
	ctx := lsh512Context{
		lsh512ContextData: lsh512ContextData{
			algType: uint64(size),
		},
	}
	defer ctx.Zeroize()

	ctx.z.SetValid()
	ctx.Reset()
	ctx.Write(data)

	return ctx.checkSum()
}

func (ctx *lsh512Context) Size() int {
	return int(ctx.algType)
}

func (ctx *lsh512Context) BlockSize() int {
	return BlockSize
}

func (ctx *lsh512Context) Write(p []byte) (n int, err error) {
	if err = ctx.z.ReturnIfZeroized(); err != nil {
		return 0, err
	}

	if len(p) == 0 {
		return 0, nil
	}
	n = len(p)

	if ctx.remainDataByteLen != 0 {
		remain := int(ctx.remainDataByteLen)
		if remain < 0 || uint64(remain) != ctx.remainDataByteLen || remain >= BlockSize {
			panic(msgInvalidInternalState)
		}
		gap := BlockSize - remain
		if len(p) < gap {
			nth := copy(ctx.lastBlock[remain:], p)
			ctx.remainDataByteLen += uint64(nth)
			return n, nil
		}
		copy(ctx.lastBlock[remain:], p[:gap])
		ctx.compress(&ctx.lastBlock)
		ctx.remainDataByteLen = 0
		p = p[gap:]
	}

	for len(p) >= BlockSize {
		ctx.compress((*block)(p[:BlockSize]))
		p = p[BlockSize:]
	}

	if len(p) != 0 {
		ctx.remainDataByteLen = uint64(copy(ctx.lastBlock[:], p))
	}
	return n, nil
}

func (ctx *lsh512Context) Sum(p []byte) []byte {
	ctx.z.ThrowIfZeroized()

	b0 := ctx.lsh512ContextData
	defer b0.Zeroize()
	hash := b0.checkSum()
	size := int(ctx.algType)
	if size < 0 || uint64(size) != ctx.algType || size > len(hash) {
		panic(msgInvalidInternalState)
	}
	return append(p, hash[:size]...)
}

func (ctx *lsh512ContextData) checkSum() (out [Size]byte) {
	remain := int(ctx.remainDataByteLen)
	if remain < 0 || uint64(remain) != ctx.remainDataByteLen || remain >= BlockSize {
		panic(msgInvalidInternalState)
	}
	ctx.lastBlock[remain] = 0x80
	clear(ctx.lastBlock[remain+1:])
	ctx.compress(&ctx.lastBlock)

	binary.LittleEndian.PutUint64(out[0:], ctx.cvL[0]^ctx.cvR[0])
	binary.LittleEndian.PutUint64(out[8:], ctx.cvL[1]^ctx.cvR[1])
	binary.LittleEndian.PutUint64(out[16:], ctx.cvL[2]^ctx.cvR[2])
	binary.LittleEndian.PutUint64(out[24:], ctx.cvL[3]^ctx.cvR[3])
	binary.LittleEndian.PutUint64(out[32:], ctx.cvL[4]^ctx.cvR[4])
	binary.LittleEndian.PutUint64(out[40:], ctx.cvL[5]^ctx.cvR[5])
	binary.LittleEndian.PutUint64(out[48:], ctx.cvL[6]^ctx.cvR[6])
	binary.LittleEndian.PutUint64(out[56:], ctx.cvL[7]^ctx.cvR[7])

	return out
}
