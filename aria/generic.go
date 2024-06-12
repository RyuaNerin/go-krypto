//go:build amd64 || (!amd64 && !arm64) || purego
// +build amd64 !amd64,!arm64 purego

package aria

import (
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"math/bits"

	"github.com/RyuaNerin/go-krypto/internal/memory"
)

func newCipherGo(key []byte) (cipher.Block, error) {
	ctx := new(ariaContext)
	ctx.rounds = (len(key) + 32) / 4

	ctx.initRoundKey(key)
	return ctx, nil
}

func (ctx *ariaContext) BlockSize() int {
	return BlockSize
}

func (ctx *ariaContext) Encrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf(msgFormatInvalidBlockSizeSrc, len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf(msgFormatInvalidBlockSizeDst, len(dst)))
	}

	processGo(dst, src, ctx.ek[:], ctx.rounds)
}

func (ctx *ariaContext) Decrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf(msgFormatInvalidBlockSizeSrc, len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf(msgFormatInvalidBlockSizeDst, len(dst)))
	}

	processGo(dst, src, ctx.dk[:], ctx.rounds)
}

func (ctx *ariaContext) initRoundKey(key []byte) {
	encKeySetup(ctx.ek[:], key)

	ctx.dk = ctx.ek
	decKeySetup(ctx.dk[:], ctx.rounds)
}

func encKeySetup(rk []byte, mk []byte) {
	keyBytes := len(mk)

	var w0, w1, w2, w3 [4]uint32

	w0[0] = binary.BigEndian.Uint32(mk[0*4:]) // WordLoad(WO(mk, 0), w0[0])
	w0[1] = binary.BigEndian.Uint32(mk[1*4:]) // WordLoad(WO(mk, 1), w0[1])
	w0[2] = binary.BigEndian.Uint32(mk[2*4:]) // WordLoad(WO(mk, 2), w0[2])
	w0[3] = binary.BigEndian.Uint32(mk[3*4:]) // WordLoad(WO(mk, 3), w0[3])

	q := (keyBytes - 16) / 8
	t0 := w0[0] ^ krk[q*4+0]
	t1 := w0[1] ^ krk[q*4+1]
	t2 := w0[2] ^ krk[q*4+2]
	t3 := w0[3] ^ krk[q*4+3]
	// FO(&t0, &t1, &t2, &t3)
	{
		// SBL1_M(T0, T1, T2, T3)
		{
			t0 = s1[c_brf(t0, 24)] ^ s2[c_brf(t0, 16)] ^ x1[c_brf(t0, 8)] ^ x2[c_brf(t0, 0)]
			t1 = s1[c_brf(t1, 24)] ^ s2[c_brf(t1, 16)] ^ x1[c_brf(t1, 8)] ^ x2[c_brf(t1, 0)]
			t2 = s1[c_brf(t2, 24)] ^ s2[c_brf(t2, 16)] ^ x1[c_brf(t2, 8)] ^ x2[c_brf(t2, 0)]
			t3 = s1[c_brf(t3, 24)] ^ s2[c_brf(t3, 16)] ^ x1[c_brf(t3, 8)] ^ x2[c_brf(t3, 0)]
		}
		c_mm(&t0, &t1, &t2, &t3) // inlining
		c_p(&t0, &t1, &t2, &t3)  // inlining
		c_mm(&t0, &t1, &t2, &t3) // inlining
	}
	if keyBytes > 16 {
		w1[0] = binary.BigEndian.Uint32(mk[4*4:]) // WordLoad(WO(mk, 4), w1[0])
		w1[1] = binary.BigEndian.Uint32(mk[5*4:]) // WordLoad(WO(mk, 5), w1[1])
		if keyBytes > 24 {
			w1[2] = binary.BigEndian.Uint32(mk[6*4:]) // WordLoad(WO(mk,6), w1[2]);
			w1[3] = binary.BigEndian.Uint32(mk[7*4:]) // WordLoad(WO(mk,7), w1[3]);
		} else {
			w1[2] = 0
			w1[3] = 0
		}
	} else {
		w1[0] = 0
		w1[1] = 0
		w1[2] = 0
		w1[3] = 0
	}
	w1[0] ^= t0
	w1[1] ^= t1
	w1[2] ^= t2
	w1[3] ^= t3
	t0 = w1[0]
	t1 = w1[1]
	t2 = w1[2]
	t3 = w1[3]

	q = (q + 1) % 3
	t0 ^= krk[q*4+0]
	t1 ^= krk[q*4+1]
	t2 ^= krk[q*4+2]
	t3 ^= krk[q*4+3]
	// FE(&t0, &t1, &t2, &t3)
	{
		// SBL2_M(T0, T1, T2, T3)
		{
			t0 = x1[c_brf(t0, 24)] ^ x2[c_brf(t0, 16)] ^ s1[c_brf(t0, 8)] ^ s2[c_brf(t0, 0)]
			t1 = x1[c_brf(t1, 24)] ^ x2[c_brf(t1, 16)] ^ s1[c_brf(t1, 8)] ^ s2[c_brf(t1, 0)]
			t2 = x1[c_brf(t2, 24)] ^ x2[c_brf(t2, 16)] ^ s1[c_brf(t2, 8)] ^ s2[c_brf(t2, 0)]
			t3 = x1[c_brf(t3, 24)] ^ x2[c_brf(t3, 16)] ^ s1[c_brf(t3, 8)] ^ s2[c_brf(t3, 0)]
		}
		c_mm(&t0, &t1, &t2, &t3) // inlining
		c_p(&t2, &t3, &t0, &t1)  // inlining
		c_mm(&t0, &t1, &t2, &t3) // inlining
	}
	t0 ^= w0[0]
	t1 ^= w0[1]
	t2 ^= w0[2]
	t3 ^= w0[3]
	w2[0] = t0
	w2[1] = t1
	w2[2] = t2
	w2[3] = t3

	q = (q + 1) % 3
	t0 ^= krk[q*4+0]
	t1 ^= krk[q*4+1]
	t2 ^= krk[q*4+2]
	t3 ^= krk[q*4+3]
	// FO(&t0, &t1, &t2, &t3)
	{
		// SBL1_M(T0, T1, T2, T3)
		{
			t0 = s1[c_brf(t0, 24)] ^ s2[c_brf(t0, 16)] ^ x1[c_brf(t0, 8)] ^ x2[c_brf(t0, 0)]
			t1 = s1[c_brf(t1, 24)] ^ s2[c_brf(t1, 16)] ^ x1[c_brf(t1, 8)] ^ x2[c_brf(t1, 0)]
			t2 = s1[c_brf(t2, 24)] ^ s2[c_brf(t2, 16)] ^ x1[c_brf(t2, 8)] ^ x2[c_brf(t2, 0)]
			t3 = s1[c_brf(t3, 24)] ^ s2[c_brf(t3, 16)] ^ x1[c_brf(t3, 8)] ^ x2[c_brf(t3, 0)]
		}
		c_mm(&t0, &t1, &t2, &t3) // inlining
		c_p(&t0, &t1, &t2, &t3)  // inlining
		c_mm(&t0, &t1, &t2, &t3) // inlining
	}
	w3[0] = t0 ^ w1[0]
	w3[1] = t1 ^ w1[1]
	w3[2] = t2 ^ w1[2]
	w3[3] = t3 ^ w1[3]

	// var rkIndex int
	// GSRK(rk, &rkIndex, w0, w1, 0, 1, 2, 3, 19, 13)
	{
		binary.LittleEndian.PutUint32(rk[0x000:], (w0[0])^((w1[0])>>19)^((w1[3])<<13))
		binary.LittleEndian.PutUint32(rk[0x004:], (w0[1])^((w1[1])>>19)^((w1[0])<<13))
		binary.LittleEndian.PutUint32(rk[0x008:], (w0[2])^((w1[2])>>19)^((w1[1])<<13))
		binary.LittleEndian.PutUint32(rk[0x00C:], (w0[3])^((w1[3])>>19)^((w1[2])<<13))
	}
	// GSRK(rk, &rkIndex, w1, w2, 0, 1, 2, 3, 19, 13)
	{
		binary.LittleEndian.PutUint32(rk[0x010:], (w1[0])^((w2[0])>>19)^((w2[3])<<13))
		binary.LittleEndian.PutUint32(rk[0x014:], (w1[1])^((w2[1])>>19)^((w2[0])<<13))
		binary.LittleEndian.PutUint32(rk[0x018:], (w1[2])^((w2[2])>>19)^((w2[1])<<13))
		binary.LittleEndian.PutUint32(rk[0x01C:], (w1[3])^((w2[3])>>19)^((w2[2])<<13))
	}
	// GSRK(rk, &rkIndex, w2, w3, 0, 1, 2, 3, 19, 13)
	{
		binary.LittleEndian.PutUint32(rk[0x020:], (w2[0])^((w3[0])>>19)^((w3[3])<<13))
		binary.LittleEndian.PutUint32(rk[0x024:], (w2[1])^((w3[1])>>19)^((w3[0])<<13))
		binary.LittleEndian.PutUint32(rk[0x028:], (w2[2])^((w3[2])>>19)^((w3[1])<<13))
		binary.LittleEndian.PutUint32(rk[0x02C:], (w2[3])^((w3[3])>>19)^((w3[2])<<13))
	}
	// GSRK(rk, &rkIndex, w3, w0, 0, 1, 2, 3, 19, 13)
	{
		binary.LittleEndian.PutUint32(rk[0x030:], (w3[0])^((w0[0])>>19)^((w0[3])<<13))
		binary.LittleEndian.PutUint32(rk[0x034:], (w3[1])^((w0[1])>>19)^((w0[0])<<13))
		binary.LittleEndian.PutUint32(rk[0x038:], (w3[2])^((w0[2])>>19)^((w0[1])<<13))
		binary.LittleEndian.PutUint32(rk[0x03C:], (w3[3])^((w0[3])>>19)^((w0[2])<<13))
	}

	// GSRK(rk, &rkIndex, w0, w1, 0, 1, 2, 3, 31, 1)
	{
		binary.LittleEndian.PutUint32(rk[0x040:], (w0[0])^((w1[0])>>31)^((w1[3])<<1))
		binary.LittleEndian.PutUint32(rk[0x044:], (w0[1])^((w1[1])>>31)^((w1[0])<<1))
		binary.LittleEndian.PutUint32(rk[0x048:], (w0[2])^((w1[2])>>31)^((w1[1])<<1))
		binary.LittleEndian.PutUint32(rk[0x04C:], (w0[3])^((w1[3])>>31)^((w1[2])<<1))
	}
	// GSRK(rk, &rkIndex, w1, w2, 0, 1, 2, 3, 31, 1)
	{
		binary.LittleEndian.PutUint32(rk[0x050:], (w1[0])^((w2[0])>>31)^((w2[3])<<1))
		binary.LittleEndian.PutUint32(rk[0x054:], (w1[1])^((w2[1])>>31)^((w2[0])<<1))
		binary.LittleEndian.PutUint32(rk[0x058:], (w1[2])^((w2[2])>>31)^((w2[1])<<1))
		binary.LittleEndian.PutUint32(rk[0x05C:], (w1[3])^((w2[3])>>31)^((w2[2])<<1))
	}
	// GSRK(rk, &rkIndex, w2, w3, 0, 1, 2, 3, 31, 1)
	{
		binary.LittleEndian.PutUint32(rk[0x060:], (w2[0])^((w3[0])>>31)^((w3[3])<<1))
		binary.LittleEndian.PutUint32(rk[0x064:], (w2[1])^((w3[1])>>31)^((w3[0])<<1))
		binary.LittleEndian.PutUint32(rk[0x068:], (w2[2])^((w3[2])>>31)^((w3[1])<<1))
		binary.LittleEndian.PutUint32(rk[0x06C:], (w2[3])^((w3[3])>>31)^((w3[2])<<1))
	}
	// GSRK(rk, &rkIndex, w3, w0, 0, 1, 2, 3, 31, 1)
	{
		binary.LittleEndian.PutUint32(rk[0x070:], (w3[0])^((w0[0])>>31)^((w0[3])<<1))
		binary.LittleEndian.PutUint32(rk[0x074:], (w3[1])^((w0[1])>>31)^((w0[0])<<1))
		binary.LittleEndian.PutUint32(rk[0x078:], (w3[2])^((w0[2])>>31)^((w0[1])<<1))
		binary.LittleEndian.PutUint32(rk[0x07C:], (w3[3])^((w0[3])>>31)^((w0[2])<<1))
	}

	// GSRK(rk, &rkIndex, w0, w1, 2, 3, 0, 1, 3, 29)
	{
		binary.LittleEndian.PutUint32(rk[0x080:], (w0[0])^((w1[2])>>3)^((w1[1])<<29))
		binary.LittleEndian.PutUint32(rk[0x084:], (w0[1])^((w1[3])>>3)^((w1[2])<<29))
		binary.LittleEndian.PutUint32(rk[0x088:], (w0[2])^((w1[0])>>3)^((w1[3])<<29))
		binary.LittleEndian.PutUint32(rk[0x08C:], (w0[3])^((w1[1])>>3)^((w1[0])<<29))
	}
	// GSRK(rk, &rkIndex, w1, w2, 2, 3, 0, 1, 3, 29)
	{
		binary.LittleEndian.PutUint32(rk[0x090:], (w1[0])^((w2[2])>>3)^((w2[1])<<29))
		binary.LittleEndian.PutUint32(rk[0x094:], (w1[1])^((w2[3])>>3)^((w2[2])<<29))
		binary.LittleEndian.PutUint32(rk[0x098:], (w1[2])^((w2[0])>>3)^((w2[3])<<29))
		binary.LittleEndian.PutUint32(rk[0x09C:], (w1[3])^((w2[1])>>3)^((w2[0])<<29))
	}
	// GSRK(rk, &rkIndex, w2, w3, 2, 3, 0, 1, 3, 29)
	{
		binary.LittleEndian.PutUint32(rk[0x0A0:], (w2[0])^((w3[2])>>3)^((w3[1])<<29))
		binary.LittleEndian.PutUint32(rk[0x0A4:], (w2[1])^((w3[3])>>3)^((w3[2])<<29))
		binary.LittleEndian.PutUint32(rk[0x0A8:], (w2[2])^((w3[0])>>3)^((w3[3])<<29))
		binary.LittleEndian.PutUint32(rk[0x0AC:], (w2[3])^((w3[1])>>3)^((w3[0])<<29))
	}
	// GSRK(rk, &rkIndex, w3, w0, 2, 3, 0, 1, 3, 29)
	{
		binary.LittleEndian.PutUint32(rk[0x0B0:], (w3[0])^((w0[2])>>3)^((w0[1])<<29))
		binary.LittleEndian.PutUint32(rk[0x0B4:], (w3[1])^((w0[3])>>3)^((w0[2])<<29))
		binary.LittleEndian.PutUint32(rk[0x0B8:], (w3[2])^((w0[0])>>3)^((w0[3])<<29))
		binary.LittleEndian.PutUint32(rk[0x0BC:], (w3[3])^((w0[1])>>3)^((w0[0])<<29))
	}

	// GSRK(rk, &rkIndex, w0, w1, 1, 2, 3, 0, 1, 31)
	{
		binary.LittleEndian.PutUint32(rk[0x0C0:], (w0[0])^((w1[1])>>1)^((w1[0])<<31))
		binary.LittleEndian.PutUint32(rk[0x0C4:], (w0[1])^((w1[2])>>1)^((w1[1])<<31))
		binary.LittleEndian.PutUint32(rk[0x0C8:], (w0[2])^((w1[3])>>1)^((w1[2])<<31))
		binary.LittleEndian.PutUint32(rk[0x0CC:], (w0[3])^((w1[0])>>1)^((w1[3])<<31))
	}
	if keyBytes > 16 {
		// GSRK(rk, &rkIndex, w1, w2, 1, 2, 3, 0, 1, 31)
		{
			binary.LittleEndian.PutUint32(rk[0x0D0:], (w1[0])^((w2[1])>>1)^((w2[0])<<31))
			binary.LittleEndian.PutUint32(rk[0x0D4:], (w1[1])^((w2[2])>>1)^((w2[1])<<31))
			binary.LittleEndian.PutUint32(rk[0x0D8:], (w1[2])^((w2[3])>>1)^((w2[2])<<31))
			binary.LittleEndian.PutUint32(rk[0x0DC:], (w1[3])^((w2[0])>>1)^((w2[3])<<31))
		}
		// GSRK(rk, &rkIndex, w2, w3, 1, 2, 3, 0, 1, 31)
		{
			binary.LittleEndian.PutUint32(rk[0x0E0:], (w2[0])^((w3[1])>>1)^((w3[0])<<31))
			binary.LittleEndian.PutUint32(rk[0x0E4:], (w2[1])^((w3[2])>>1)^((w3[1])<<31))
			binary.LittleEndian.PutUint32(rk[0x0E8:], (w2[2])^((w3[3])>>1)^((w3[2])<<31))
			binary.LittleEndian.PutUint32(rk[0x0EC:], (w2[3])^((w3[0])>>1)^((w3[3])<<31))
		}

		if keyBytes > 24 {
			// GSRK(rk, &rkIndex, w3, w0, 1, 2, 3, 0, 1, 31)
			{
				binary.LittleEndian.PutUint32(rk[0x0F0:], (w3[0])^((w0[1])>>1)^((w0[0])<<31))
				binary.LittleEndian.PutUint32(rk[0x0F4:], (w3[1])^((w0[2])>>1)^((w0[1])<<31))
				binary.LittleEndian.PutUint32(rk[0x0F8:], (w3[2])^((w0[3])>>1)^((w0[2])<<31))
				binary.LittleEndian.PutUint32(rk[0x0FC:], (w3[3])^((w0[0])>>1)^((w0[3])<<31))
			}

			// GSRK(rk, &rkIndex, w0, w1, 1, 2, 3, 0, 13, 19)
			{
				binary.LittleEndian.PutUint32(rk[0x100:], (w0[0])^((w1[1])>>13)^((w1[0])<<19))
				binary.LittleEndian.PutUint32(rk[0x104:], (w0[1])^((w1[2])>>13)^((w1[1])<<19))
				binary.LittleEndian.PutUint32(rk[0x108:], (w0[2])^((w1[3])>>13)^((w1[2])<<19))
				binary.LittleEndian.PutUint32(rk[0x10C:], (w0[3])^((w1[0])>>13)^((w1[3])<<19))
			}
		}
	}
}

func decKeySetup(rk []byte, rValue int) {
	a := memory.PUint32(rk)
	z := a

	aIdx := 0
	zIdx := rValue * 4

	var s0, s1, s2, s3 uint32

	WordM1 := func(X uint32) uint32 {
		return (X)<<8 ^ (X)>>8 ^ (X)<<16 ^ (X)>>16 ^ (X)<<24 ^ (X)>>24
	}

	t0 := a[aIdx+0]
	t1 := a[aIdx+1]
	t2 := a[aIdx+2]
	t3 := a[aIdx+3]
	a[aIdx+0] = z[zIdx+0]
	a[aIdx+1] = z[zIdx+1]
	a[aIdx+2] = z[zIdx+2]
	a[aIdx+3] = z[zIdx+3]
	z[zIdx+0] = t0
	z[zIdx+1] = t1
	z[zIdx+2] = t2
	z[zIdx+3] = t3

	aIdx += 4
	zIdx -= 4

	for aIdx < zIdx {
		t0 = WordM1(a[aIdx+0])
		t1 = WordM1(a[aIdx+1])
		t2 = WordM1(a[aIdx+2])
		t3 = WordM1(a[aIdx+3])

		c_mm(&t0, &t1, &t2, &t3)
		c_p(&t0, &t1, &t2, &t3)
		c_mm(&t0, &t1, &t2, &t3)

		s0 = t0
		s1 = t1
		s2 = t2
		s3 = t3

		t0 = WordM1(z[zIdx+0])
		t1 = WordM1(z[zIdx+1])
		t2 = WordM1(z[zIdx+2])
		t3 = WordM1(z[zIdx+3])

		c_mm(&t0, &t1, &t2, &t3)
		c_p(&t0, &t1, &t2, &t3)
		c_mm(&t0, &t1, &t2, &t3)

		a[aIdx+0] = t0
		a[aIdx+1] = t1
		a[aIdx+2] = t2
		a[aIdx+3] = t3
		z[zIdx+0] = s0
		z[zIdx+1] = s1
		z[zIdx+2] = s2
		z[zIdx+3] = s3

		aIdx += 4
		zIdx -= 4
	}
	t0 = WordM1(a[aIdx+0])
	t1 = WordM1(a[aIdx+1])
	t2 = WordM1(a[aIdx+2])
	t3 = WordM1(a[aIdx+3])
	c_mm(&t0, &t1, &t2, &t3)
	c_p(&t0, &t1, &t2, &t3)
	c_mm(&t0, &t1, &t2, &t3)
	z[zIdx+0] = t0
	z[zIdx+1] = t1
	z[zIdx+2] = t2
	z[zIdx+3] = t3
}

func processGo(dst, src, rk []byte, round int) {
	t0 := binary.BigEndian.Uint32(src[0*4:]) // WordLoad(WO(i,0), t0)
	t1 := binary.BigEndian.Uint32(src[1*4:]) // WordLoad(WO(i,1), t1)
	t2 := binary.BigEndian.Uint32(src[2*4:]) // WordLoad(WO(i,2), t2)
	t3 := binary.BigEndian.Uint32(src[3*4:]) // WordLoad(WO(i,3), t3)

	if round > 12 {
		// KXL(rk, &rkIndex, &t0, &t1, &t2, &t3)
		{
			t0 ^= binary.LittleEndian.Uint32(rk[0x00:])
			t1 ^= binary.LittleEndian.Uint32(rk[0x04:])
			t2 ^= binary.LittleEndian.Uint32(rk[0x08:])
			t3 ^= binary.LittleEndian.Uint32(rk[0x0C:])
		}
		// FO(&t0, &t1, &t2, &t3)
		{
			// SBL1_M(T0, T1, T2, T3)
			{
				t0 = s1[c_brf(t0, 24)] ^ s2[c_brf(t0, 16)] ^ x1[c_brf(t0, 8)] ^ x2[c_brf(t0, 0)]
				t1 = s1[c_brf(t1, 24)] ^ s2[c_brf(t1, 16)] ^ x1[c_brf(t1, 8)] ^ x2[c_brf(t1, 0)]
				t2 = s1[c_brf(t2, 24)] ^ s2[c_brf(t2, 16)] ^ x1[c_brf(t2, 8)] ^ x2[c_brf(t2, 0)]
				t3 = s1[c_brf(t3, 24)] ^ s2[c_brf(t3, 16)] ^ x1[c_brf(t3, 8)] ^ x2[c_brf(t3, 0)]
			}
			c_mm(&t0, &t1, &t2, &t3) // inlining
			c_p(&t0, &t1, &t2, &t3)  // inlining
			c_mm(&t0, &t1, &t2, &t3) // inlining
		}
		// KXL(rk, &rkIndex, &t0, &t1, &t2, &t3)
		{
			t0 ^= binary.LittleEndian.Uint32(rk[0x10:])
			t1 ^= binary.LittleEndian.Uint32(rk[0x14:])
			t2 ^= binary.LittleEndian.Uint32(rk[0x18:])
			t3 ^= binary.LittleEndian.Uint32(rk[0x1C:])
		}
		// FE(&t0, &t1, &t2, &t3)
		{
			// SBL2_M(T0, T1, T2, T3)
			{
				t0 = x1[c_brf(t0, 24)] ^ x2[c_brf(t0, 16)] ^ s1[c_brf(t0, 8)] ^ s2[c_brf(t0, 0)]
				t1 = x1[c_brf(t1, 24)] ^ x2[c_brf(t1, 16)] ^ s1[c_brf(t1, 8)] ^ s2[c_brf(t1, 0)]
				t2 = x1[c_brf(t2, 24)] ^ x2[c_brf(t2, 16)] ^ s1[c_brf(t2, 8)] ^ s2[c_brf(t2, 0)]
				t3 = x1[c_brf(t3, 24)] ^ x2[c_brf(t3, 16)] ^ s1[c_brf(t3, 8)] ^ s2[c_brf(t3, 0)]
			}
			c_mm(&t0, &t1, &t2, &t3) // inlining
			c_p(&t2, &t3, &t0, &t1)  // inlining
			c_mm(&t0, &t1, &t2, &t3) // inlining
		}

		rk = rk[32:]

		if round > 14 {
			// KXL(rk, &rkIndex, &t0, &t1, &t2, &t3)
			{
				t0 ^= binary.LittleEndian.Uint32(rk[0x00:])
				t1 ^= binary.LittleEndian.Uint32(rk[0x04:])
				t2 ^= binary.LittleEndian.Uint32(rk[0x08:])
				t3 ^= binary.LittleEndian.Uint32(rk[0x0C:])
			}
			// FO(&t0, &t1, &t2, &t3)
			{
				// SBL1_M(T0, T1, T2, T3)
				{
					t0 = s1[c_brf(t0, 24)] ^ s2[c_brf(t0, 16)] ^ x1[c_brf(t0, 8)] ^ x2[c_brf(t0, 0)]
					t1 = s1[c_brf(t1, 24)] ^ s2[c_brf(t1, 16)] ^ x1[c_brf(t1, 8)] ^ x2[c_brf(t1, 0)]
					t2 = s1[c_brf(t2, 24)] ^ s2[c_brf(t2, 16)] ^ x1[c_brf(t2, 8)] ^ x2[c_brf(t2, 0)]
					t3 = s1[c_brf(t3, 24)] ^ s2[c_brf(t3, 16)] ^ x1[c_brf(t3, 8)] ^ x2[c_brf(t3, 0)]
				}
				c_mm(&t0, &t1, &t2, &t3) // inlining
				c_p(&t0, &t1, &t2, &t3)  // inlining
				c_mm(&t0, &t1, &t2, &t3) // inlining
			}
			// KXL(rk, &rkIndex, &t0, &t1, &t2, &t3)
			{
				t0 ^= binary.LittleEndian.Uint32(rk[0x10:])
				t1 ^= binary.LittleEndian.Uint32(rk[0x14:])
				t2 ^= binary.LittleEndian.Uint32(rk[0x18:])
				t3 ^= binary.LittleEndian.Uint32(rk[0x1C:])
			}
			// FE(&t0, &t1, &t2, &t3)
			{
				// SBL2_M(T0, T1, T2, T3)
				{
					t0 = x1[c_brf(t0, 24)] ^ x2[c_brf(t0, 16)] ^ s1[c_brf(t0, 8)] ^ s2[c_brf(t0, 0)]
					t1 = x1[c_brf(t1, 24)] ^ x2[c_brf(t1, 16)] ^ s1[c_brf(t1, 8)] ^ s2[c_brf(t1, 0)]
					t2 = x1[c_brf(t2, 24)] ^ x2[c_brf(t2, 16)] ^ s1[c_brf(t2, 8)] ^ s2[c_brf(t2, 0)]
					t3 = x1[c_brf(t3, 24)] ^ x2[c_brf(t3, 16)] ^ s1[c_brf(t3, 8)] ^ s2[c_brf(t3, 0)]
				}
				c_mm(&t0, &t1, &t2, &t3) // inlining
				c_p(&t2, &t3, &t0, &t1)  // inlining
				c_mm(&t0, &t1, &t2, &t3) // inlining
			}

			rk = rk[32:]
		}
	}

	// KXL(rk, &rkIndex, &t0, &t1, &t2, &t3)
	{
		t0 ^= binary.LittleEndian.Uint32(rk[0x00:])
		t1 ^= binary.LittleEndian.Uint32(rk[0x04:])
		t2 ^= binary.LittleEndian.Uint32(rk[0x08:])
		t3 ^= binary.LittleEndian.Uint32(rk[0x0C:])
	}
	// FO(&t0, &t1, &t2, &t3)
	{
		// SBL1_M(T0, T1, T2, T3)
		{
			t0 = s1[c_brf(t0, 24)] ^ s2[c_brf(t0, 16)] ^ x1[c_brf(t0, 8)] ^ x2[c_brf(t0, 0)]
			t1 = s1[c_brf(t1, 24)] ^ s2[c_brf(t1, 16)] ^ x1[c_brf(t1, 8)] ^ x2[c_brf(t1, 0)]
			t2 = s1[c_brf(t2, 24)] ^ s2[c_brf(t2, 16)] ^ x1[c_brf(t2, 8)] ^ x2[c_brf(t2, 0)]
			t3 = s1[c_brf(t3, 24)] ^ s2[c_brf(t3, 16)] ^ x1[c_brf(t3, 8)] ^ x2[c_brf(t3, 0)]
		}
		c_mm(&t0, &t1, &t2, &t3) // inlining
		c_p(&t0, &t1, &t2, &t3)  // inlining
		c_mm(&t0, &t1, &t2, &t3) // inlining
	}
	// KXL(rk, &rkIndex, &t0, &t1, &t2, &t3)
	{
		t0 ^= binary.LittleEndian.Uint32(rk[0x10:])
		t1 ^= binary.LittleEndian.Uint32(rk[0x14:])
		t2 ^= binary.LittleEndian.Uint32(rk[0x18:])
		t3 ^= binary.LittleEndian.Uint32(rk[0x1C:])
	}
	// FE(&t0, &t1, &t2, &t3)
	{
		// SBL2_M(T0, T1, T2, T3)
		{
			t0 = x1[c_brf(t0, 24)] ^ x2[c_brf(t0, 16)] ^ s1[c_brf(t0, 8)] ^ s2[c_brf(t0, 0)]
			t1 = x1[c_brf(t1, 24)] ^ x2[c_brf(t1, 16)] ^ s1[c_brf(t1, 8)] ^ s2[c_brf(t1, 0)]
			t2 = x1[c_brf(t2, 24)] ^ x2[c_brf(t2, 16)] ^ s1[c_brf(t2, 8)] ^ s2[c_brf(t2, 0)]
			t3 = x1[c_brf(t3, 24)] ^ x2[c_brf(t3, 16)] ^ s1[c_brf(t3, 8)] ^ s2[c_brf(t3, 0)]
		}
		c_mm(&t0, &t1, &t2, &t3) // inlining
		c_p(&t2, &t3, &t0, &t1)  // inlining
		c_mm(&t0, &t1, &t2, &t3) // inlining
	}
	// KXL(rk, &rkIndex, &t0, &t1, &t2, &t3)
	{
		t0 ^= binary.LittleEndian.Uint32(rk[0x20:])
		t1 ^= binary.LittleEndian.Uint32(rk[0x24:])
		t2 ^= binary.LittleEndian.Uint32(rk[0x28:])
		t3 ^= binary.LittleEndian.Uint32(rk[0x2C:])
	}
	// FO(&t0, &t1, &t2, &t3)
	{
		// SBL1_M(T0, T1, T2, T3)
		{
			t0 = s1[c_brf(t0, 24)] ^ s2[c_brf(t0, 16)] ^ x1[c_brf(t0, 8)] ^ x2[c_brf(t0, 0)]
			t1 = s1[c_brf(t1, 24)] ^ s2[c_brf(t1, 16)] ^ x1[c_brf(t1, 8)] ^ x2[c_brf(t1, 0)]
			t2 = s1[c_brf(t2, 24)] ^ s2[c_brf(t2, 16)] ^ x1[c_brf(t2, 8)] ^ x2[c_brf(t2, 0)]
			t3 = s1[c_brf(t3, 24)] ^ s2[c_brf(t3, 16)] ^ x1[c_brf(t3, 8)] ^ x2[c_brf(t3, 0)]
		}
		c_mm(&t0, &t1, &t2, &t3) // inlining
		c_p(&t0, &t1, &t2, &t3)  // inlining
		c_mm(&t0, &t1, &t2, &t3) // inlining
	}
	// KXL(rk, &rkIndex, &t0, &t1, &t2, &t3)
	{
		t0 ^= binary.LittleEndian.Uint32(rk[0x30:])
		t1 ^= binary.LittleEndian.Uint32(rk[0x34:])
		t2 ^= binary.LittleEndian.Uint32(rk[0x38:])
		t3 ^= binary.LittleEndian.Uint32(rk[0x3C:])
	}
	// FE(&t0, &t1, &t2, &t3)
	{
		// SBL2_M(T0, T1, T2, T3)
		{
			t0 = x1[c_brf(t0, 24)] ^ x2[c_brf(t0, 16)] ^ s1[c_brf(t0, 8)] ^ s2[c_brf(t0, 0)]
			t1 = x1[c_brf(t1, 24)] ^ x2[c_brf(t1, 16)] ^ s1[c_brf(t1, 8)] ^ s2[c_brf(t1, 0)]
			t2 = x1[c_brf(t2, 24)] ^ x2[c_brf(t2, 16)] ^ s1[c_brf(t2, 8)] ^ s2[c_brf(t2, 0)]
			t3 = x1[c_brf(t3, 24)] ^ x2[c_brf(t3, 16)] ^ s1[c_brf(t3, 8)] ^ s2[c_brf(t3, 0)]
		}
		c_mm(&t0, &t1, &t2, &t3) // inlining
		c_p(&t2, &t3, &t0, &t1)  // inlining
		c_mm(&t0, &t1, &t2, &t3) // inlining
	}
	// KXL(rk, &rkIndex, &t0, &t1, &t2, &t3)
	{
		t0 ^= binary.LittleEndian.Uint32(rk[0x40:])
		t1 ^= binary.LittleEndian.Uint32(rk[0x44:])
		t2 ^= binary.LittleEndian.Uint32(rk[0x48:])
		t3 ^= binary.LittleEndian.Uint32(rk[0x4C:])
	}
	// FO(&t0, &t1, &t2, &t3)
	{
		// SBL1_M(T0, T1, T2, T3)
		{
			t0 = s1[c_brf(t0, 24)] ^ s2[c_brf(t0, 16)] ^ x1[c_brf(t0, 8)] ^ x2[c_brf(t0, 0)]
			t1 = s1[c_brf(t1, 24)] ^ s2[c_brf(t1, 16)] ^ x1[c_brf(t1, 8)] ^ x2[c_brf(t1, 0)]
			t2 = s1[c_brf(t2, 24)] ^ s2[c_brf(t2, 16)] ^ x1[c_brf(t2, 8)] ^ x2[c_brf(t2, 0)]
			t3 = s1[c_brf(t3, 24)] ^ s2[c_brf(t3, 16)] ^ x1[c_brf(t3, 8)] ^ x2[c_brf(t3, 0)]
		}
		c_mm(&t0, &t1, &t2, &t3) // inlining
		c_p(&t0, &t1, &t2, &t3)  // inlining
		c_mm(&t0, &t1, &t2, &t3) // inlining
	}
	// KXL(rk, &rkIndex, &t0, &t1, &t2, &t3)
	{
		t0 ^= binary.LittleEndian.Uint32(rk[0x50:])
		t1 ^= binary.LittleEndian.Uint32(rk[0x54:])
		t2 ^= binary.LittleEndian.Uint32(rk[0x58:])
		t3 ^= binary.LittleEndian.Uint32(rk[0x5C:])
	}
	// FE(&t0, &t1, &t2, &t3)
	{
		// SBL2_M(T0, T1, T2, T3)
		{
			t0 = x1[c_brf(t0, 24)] ^ x2[c_brf(t0, 16)] ^ s1[c_brf(t0, 8)] ^ s2[c_brf(t0, 0)]
			t1 = x1[c_brf(t1, 24)] ^ x2[c_brf(t1, 16)] ^ s1[c_brf(t1, 8)] ^ s2[c_brf(t1, 0)]
			t2 = x1[c_brf(t2, 24)] ^ x2[c_brf(t2, 16)] ^ s1[c_brf(t2, 8)] ^ s2[c_brf(t2, 0)]
			t3 = x1[c_brf(t3, 24)] ^ x2[c_brf(t3, 16)] ^ s1[c_brf(t3, 8)] ^ s2[c_brf(t3, 0)]
		}
		c_mm(&t0, &t1, &t2, &t3) // inlining
		c_p(&t2, &t3, &t0, &t1)  // inlining
		c_mm(&t0, &t1, &t2, &t3) // inlining
	}

	// KXL(rk, &rkIndex, &t0, &t1, &t2, &t3)
	{
		t0 ^= binary.LittleEndian.Uint32(rk[0x60:])
		t1 ^= binary.LittleEndian.Uint32(rk[0x64:])
		t2 ^= binary.LittleEndian.Uint32(rk[0x68:])
		t3 ^= binary.LittleEndian.Uint32(rk[0x6C:])
	}
	// FO(&t0, &t1, &t2, &t3)
	{
		// SBL1_M(T0, T1, T2, T3)
		{
			t0 = s1[c_brf(t0, 24)] ^ s2[c_brf(t0, 16)] ^ x1[c_brf(t0, 8)] ^ x2[c_brf(t0, 0)]
			t1 = s1[c_brf(t1, 24)] ^ s2[c_brf(t1, 16)] ^ x1[c_brf(t1, 8)] ^ x2[c_brf(t1, 0)]
			t2 = s1[c_brf(t2, 24)] ^ s2[c_brf(t2, 16)] ^ x1[c_brf(t2, 8)] ^ x2[c_brf(t2, 0)]
			t3 = s1[c_brf(t3, 24)] ^ s2[c_brf(t3, 16)] ^ x1[c_brf(t3, 8)] ^ x2[c_brf(t3, 0)]
		}
		c_mm(&t0, &t1, &t2, &t3) // inlining
		c_p(&t0, &t1, &t2, &t3)  // inlining
		c_mm(&t0, &t1, &t2, &t3) // inlining
	}
	// KXL(rk, &rkIndex, &t0, &t1, &t2, &t3)
	{
		t0 ^= binary.LittleEndian.Uint32(rk[0x70:])
		t1 ^= binary.LittleEndian.Uint32(rk[0x74:])
		t2 ^= binary.LittleEndian.Uint32(rk[0x78:])
		t3 ^= binary.LittleEndian.Uint32(rk[0x7C:])
	}
	// FE(&t0, &t1, &t2, &t3)
	{
		// SBL2_M(T0, T1, T2, T3)
		{
			t0 = x1[c_brf(t0, 24)] ^ x2[c_brf(t0, 16)] ^ s1[c_brf(t0, 8)] ^ s2[c_brf(t0, 0)]
			t1 = x1[c_brf(t1, 24)] ^ x2[c_brf(t1, 16)] ^ s1[c_brf(t1, 8)] ^ s2[c_brf(t1, 0)]
			t2 = x1[c_brf(t2, 24)] ^ x2[c_brf(t2, 16)] ^ s1[c_brf(t2, 8)] ^ s2[c_brf(t2, 0)]
			t3 = x1[c_brf(t3, 24)] ^ x2[c_brf(t3, 16)] ^ s1[c_brf(t3, 8)] ^ s2[c_brf(t3, 0)]
		}
		c_mm(&t0, &t1, &t2, &t3) // inlining
		c_p(&t2, &t3, &t0, &t1)  // inlining
		c_mm(&t0, &t1, &t2, &t3) // inlining
	}
	// KXL(rk, &rkIndex, &t0, &t1, &t2, &t3)
	{
		t0 ^= binary.LittleEndian.Uint32(rk[0x80:])
		t1 ^= binary.LittleEndian.Uint32(rk[0x84:])
		t2 ^= binary.LittleEndian.Uint32(rk[0x88:])
		t3 ^= binary.LittleEndian.Uint32(rk[0x8C:])
	}
	// FO(&t0, &t1, &t2, &t3)
	{
		// SBL1_M(T0, T1, T2, T3)
		{
			t0 = s1[c_brf(t0, 24)] ^ s2[c_brf(t0, 16)] ^ x1[c_brf(t0, 8)] ^ x2[c_brf(t0, 0)]
			t1 = s1[c_brf(t1, 24)] ^ s2[c_brf(t1, 16)] ^ x1[c_brf(t1, 8)] ^ x2[c_brf(t1, 0)]
			t2 = s1[c_brf(t2, 24)] ^ s2[c_brf(t2, 16)] ^ x1[c_brf(t2, 8)] ^ x2[c_brf(t2, 0)]
			t3 = s1[c_brf(t3, 24)] ^ s2[c_brf(t3, 16)] ^ x1[c_brf(t3, 8)] ^ x2[c_brf(t3, 0)]
		}
		c_mm(&t0, &t1, &t2, &t3) // inlining
		c_p(&t0, &t1, &t2, &t3)  // inlining
		c_mm(&t0, &t1, &t2, &t3) // inlining
	}
	// KXL(rk, &rkIndex, &t0, &t1, &t2, &t3)
	{
		t0 ^= binary.LittleEndian.Uint32(rk[0x90:])
		t1 ^= binary.LittleEndian.Uint32(rk[0x94:])
		t2 ^= binary.LittleEndian.Uint32(rk[0x98:])
		t3 ^= binary.LittleEndian.Uint32(rk[0x9C:])
	}
	// FE(&t0, &t1, &t2, &t3)
	{
		// SBL2_M(T0, T1, T2, T3)
		{
			t0 = x1[c_brf(t0, 24)] ^ x2[c_brf(t0, 16)] ^ s1[c_brf(t0, 8)] ^ s2[c_brf(t0, 0)]
			t1 = x1[c_brf(t1, 24)] ^ x2[c_brf(t1, 16)] ^ s1[c_brf(t1, 8)] ^ s2[c_brf(t1, 0)]
			t2 = x1[c_brf(t2, 24)] ^ x2[c_brf(t2, 16)] ^ s1[c_brf(t2, 8)] ^ s2[c_brf(t2, 0)]
			t3 = x1[c_brf(t3, 24)] ^ x2[c_brf(t3, 16)] ^ s1[c_brf(t3, 8)] ^ s2[c_brf(t3, 0)]
		}
		c_mm(&t0, &t1, &t2, &t3) // inlining
		c_p(&t2, &t3, &t0, &t1)  // inlining
		c_mm(&t0, &t1, &t2, &t3) // inlining
	}
	// KXL(rk, &rkIndex, &t0, &t1, &t2, &t3)
	{
		t0 ^= binary.LittleEndian.Uint32(rk[0xA0:])
		t1 ^= binary.LittleEndian.Uint32(rk[0xA4:])
		t2 ^= binary.LittleEndian.Uint32(rk[0xA8:])
		t3 ^= binary.LittleEndian.Uint32(rk[0xAC:])
	}
	// FO(&t0, &t1, &t2, &t3)
	{
		// SBL1_M(T0, T1, T2, T3)
		{
			t0 = s1[c_brf(t0, 24)] ^ s2[c_brf(t0, 16)] ^ x1[c_brf(t0, 8)] ^ x2[c_brf(t0, 0)]
			t1 = s1[c_brf(t1, 24)] ^ s2[c_brf(t1, 16)] ^ x1[c_brf(t1, 8)] ^ x2[c_brf(t1, 0)]
			t2 = s1[c_brf(t2, 24)] ^ s2[c_brf(t2, 16)] ^ x1[c_brf(t2, 8)] ^ x2[c_brf(t2, 0)]
			t3 = s1[c_brf(t3, 24)] ^ s2[c_brf(t3, 16)] ^ x1[c_brf(t3, 8)] ^ x2[c_brf(t3, 0)]
		}
		c_mm(&t0, &t1, &t2, &t3) // inlining
		c_p(&t0, &t1, &t2, &t3)  // inlining
		c_mm(&t0, &t1, &t2, &t3) // inlining
	}
	// KXL(rk, &rkIndex, &t0, &t1, &t2, &t3)
	{
		t0 ^= binary.LittleEndian.Uint32(rk[0xB0:])
		t1 ^= binary.LittleEndian.Uint32(rk[0xB4:])
		t2 ^= binary.LittleEndian.Uint32(rk[0xB8:])
		t3 ^= binary.LittleEndian.Uint32(rk[0xBC:])
	}

	dst[0x0] = (byte)(x1[c_brf(t0, 0x18)]>>0) ^ rk[0xC3]
	dst[0x1] = (byte)(x2[c_brf(t0, 0x10)]>>8) ^ rk[0xC2]
	dst[0x2] = (byte)(s1[c_brf(t0, 0x08)]>>0) ^ rk[0xC1]
	dst[0x3] = (byte)(s2[c_brf(t0, 0x00)]>>0) ^ rk[0xC0]

	dst[0x4] = (byte)(x1[c_brf(t1, 0x18)]>>0) ^ rk[0xC7]
	dst[0x5] = (byte)(x2[c_brf(t1, 0x10)]>>8) ^ rk[0xC6]
	dst[0x6] = (byte)(s1[c_brf(t1, 0x08)]>>0) ^ rk[0xC5]
	dst[0x7] = (byte)(s2[c_brf(t1, 0x00)]>>0) ^ rk[0xC4]

	dst[0x8] = (byte)(x1[c_brf(t2, 0x18)]>>0) ^ rk[0xCB]
	dst[0x9] = (byte)(x2[c_brf(t2, 0x10)]>>8) ^ rk[0xCA]
	dst[0xA] = (byte)(s1[c_brf(t2, 0x08)]>>0) ^ rk[0xC9]
	dst[0xB] = (byte)(s2[c_brf(t2, 0x00)]>>0) ^ rk[0xC8]

	dst[0xC] = (byte)(x1[c_brf(t3, 0x18)]>>0) ^ rk[0xCF]
	dst[0xD] = (byte)(x2[c_brf(t3, 0x10)]>>8) ^ rk[0xCE]
	dst[0xE] = (byte)(s1[c_brf(t3, 0x08)]>>0) ^ rk[0xCD]
	dst[0xF] = (byte)(s2[c_brf(t3, 0x00)]>>0) ^ rk[0xCC]
}

func c_brf(T uint32, R int) int { // inlinable
	return int((T)>>R) & 0xFF
}

func c_mm(T0, T1, T2, T3 *uint32) { // inlinable
	(*T1) ^= (*T2)
	(*T2) ^= (*T3)
	(*T0) ^= (*T1)
	(*T3) ^= (*T1)
	(*T2) ^= (*T0)
	(*T1) ^= (*T2)
}

func c_p(_, T1, T2, T3 *uint32) { // inlinable
	*T1 = (((*T1) << 8) & uint32(0xff00ff00)) ^ (((*T1) >> 8) & uint32(0x00ff00ff))
	*T2 = (((*T2) << 16) & uint32(0xffff0000)) ^ (((*T2) >> 16) & uint32(0x0000ffff))
	*T3 = bits.ReverseBytes32(*T3)
}
