package aria

import (
	"crypto/cipher"
	"encoding/binary"
	"math/bits"

	"github.com/RyuaNerin/go-krypto/internal/alias"
)

func newCipher(key []byte) (cipher.Block, error) {
	ctx := new(ariaContext)
	ctx.z.SetValid()
	initRoundKey(ctx, key)
	return ctx, nil
}

func (ctx *ariaContext) Encrypt(dst, src []byte) {
	ctx.z.ThrowIfZeroized()

	if len(src) < BlockSize {
		panic(msgNotFullBlock)
	}
	if len(dst) < BlockSize {
		panic(msgNotFullBlock)
	}
	if alias.InexactOverlap(dst[:BlockSize], src[:BlockSize]) {
		panic(msgBufferOverlap)
	}

	_ = dst[BlockSize-1]
	_ = src[BlockSize-1]

	processGo((*block1)(dst), (*block1)(src), &ctx.ek, ctx.rounds)
}

func (ctx *ariaContext) Decrypt(dst, src []byte) {
	ctx.z.ThrowIfZeroized()

	if len(src) < BlockSize {
		panic(msgNotFullBlock)
	}
	if len(dst) < BlockSize {
		panic(msgNotFullBlock)
	}
	if alias.InexactOverlap(dst[:BlockSize], src[:BlockSize]) {
		panic(msgBufferOverlap)
	}

	_ = dst[BlockSize-1]
	_ = src[BlockSize-1]

	processGo((*block1)(dst), (*block1)(src), &ctx.dk, ctx.rounds)
}

func initRoundKey(ctx *ariaContext, key []byte) {
	ctx.rounds = (len(key) + 32) / 4

	encKeySetup(&ctx.ek, key)

	ctx.dk = ctx.ek
	decKeySetup(&ctx.dk, int(ctx.rounds))
}

func encKeySetup(rk *roundKey, mk []byte) {
	keyBytes := len(mk)

	var key [32]byte
	copy(key[:], mk)
	defer clear(key[:])

	var t0, t1, t2, t3 uint32
	var w0, w1, w2, w3 [4]uint32

	w0[0] = binary.BigEndian.Uint32(key[0:])
	w0[1] = binary.BigEndian.Uint32(key[4:])
	w0[2] = binary.BigEndian.Uint32(key[8:])
	w0[3] = binary.BigEndian.Uint32(key[12:])

	q := uint(keyBytes-16) / 8 % 3
	KRK := (*[4]uint32)(krk[q*4:])
	t0, t1, t2, t3 = w0[0]^KRK[0], w0[1]^KRK[1], w0[2]^KRK[2], w0[3]^KRK[3]
	{ // FO
		t0, t1, t2, t3 = sbl1_m(t0), sbl1_m(t1), sbl1_m(t2), sbl1_m(t3) // SBL1_M
		t0, t1, t2, t3 = mm(t0, t1, t2, t3)                             // MM
		t1, t2, t3 = p(t1, t2, t3)                                      // P
		t0, t1, t2, t3 = mm(t0, t1, t2, t3)                             // MM
	}

	if keyBytes > 16 {
		w1[0] = binary.BigEndian.Uint32(key[16:]) // WordLoad(WO(mk, 4), w1[0])
		w1[1] = binary.BigEndian.Uint32(key[20:]) // WordLoad(WO(mk, 5), w1[1])
		if keyBytes > 24 {
			w1[2] = binary.BigEndian.Uint32(key[24:]) // WordLoad(WO(mk,6), w1[2]);
			w1[3] = binary.BigEndian.Uint32(key[28:]) // WordLoad(WO(mk,7), w1[3]);
		}
	}
	w1[0] ^= t0
	w1[1] ^= t1
	w1[2] ^= t2
	w1[3] ^= t3

	q = (q + 1) % 3
	KRK = (*[4]uint32)(krk[q*4:])
	t0, t1, t2, t3 = w1[0]^KRK[0], w1[1]^KRK[1], w1[2]^KRK[2], w1[3]^KRK[3]
	{ // FE
		t0, t1, t2, t3 = sbl2_m(t0), sbl2_m(t1), sbl2_m(t2), sbl2_m(t3) // SBL2_M
		t0, t1, t2, t3 = mm(t0, t1, t2, t3)                             // MM
		t3, t0, t1 = p(t3, t0, t1)                                      // P
		t0, t1, t2, t3 = mm(t0, t1, t2, t3)                             // MM
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
	KRK = (*[4]uint32)(krk[q*4:])
	t0, t1, t2, t3 = w2[0]^KRK[0], w2[1]^KRK[1], w2[2]^KRK[2], w2[3]^KRK[3]
	{ // FO
		t0, t1, t2, t3 = sbl1_m(t0), sbl1_m(t1), sbl1_m(t2), sbl1_m(t3) // SBL1_M
		t0, t1, t2, t3 = mm(t0, t1, t2, t3)                             // MM
		t1, t2, t3 = p(t1, t2, t3)                                      // P
		t0, t1, t2, t3 = mm(t0, t1, t2, t3)                             // MM
	}
	w3[0] = t0 ^ w1[0]
	w3[1] = t1 ^ w1[1]
	w3[2] = t2 ^ w1[2]
	w3[3] = t3 ^ w1[3]

	// q = 4 - n/32
	// r = n % 32
	// (q+0)%4, (q+1)%4, (q+2)%4, (q+3)%4, r

	// n=19, q=4, q0=0, q1=1, q2=2, q3=3, r=19
	gsrk((*[4]uint32)(rk[0:4]), &w0, w1[0], w1[1], w1[2], w1[3], 19)
	gsrk((*[4]uint32)(rk[4:8]), &w1, w2[0], w2[1], w2[2], w2[3], 19)
	gsrk((*[4]uint32)(rk[8:12]), &w2, w3[0], w3[1], w3[2], w3[3], 19)
	gsrk((*[4]uint32)(rk[12:16]), &w3, w0[0], w0[1], w0[2], w0[3], 19)

	// n=31, q=4, q0=0, q1=1, q2=2, q3=3, r=31
	gsrk((*[4]uint32)(rk[16:20]), &w0, w1[0], w1[1], w1[2], w1[3], 31)
	gsrk((*[4]uint32)(rk[20:24]), &w1, w2[0], w2[1], w2[2], w2[3], 31)
	gsrk((*[4]uint32)(rk[24:28]), &w2, w3[0], w3[1], w3[2], w3[3], 31)
	gsrk((*[4]uint32)(rk[28:32]), &w3, w0[0], w0[1], w0[2], w0[3], 31)

	// n=67, q=3, q0=2, q1=3, q2=0, q3=1, r=3
	gsrk((*[4]uint32)(rk[32:36]), &w0, w1[2], w1[3], w1[0], w1[1], 3)
	gsrk((*[4]uint32)(rk[36:40]), &w1, w2[2], w2[3], w2[0], w2[1], 3)
	gsrk((*[4]uint32)(rk[40:44]), &w2, w3[2], w3[3], w3[0], w3[1], 3)
	gsrk((*[4]uint32)(rk[44:48]), &w3, w0[2], w0[3], w0[0], w0[1], 3)

	// n=97, q=1, q0=1, q1=2, q2=3, q3=0, r=1
	gsrk((*[4]uint32)(rk[48:52]), &w0, w1[1], w1[2], w1[3], w1[0], 1)
	if keyBytes > 16 {
		gsrk((*[4]uint32)(rk[52:56]), &w1, w2[1], w2[2], w2[3], w2[0], 1)
		gsrk((*[4]uint32)(rk[56:60]), &w2, w3[1], w3[2], w3[3], w3[0], 1)
		if keyBytes > 24 {
			gsrk((*[4]uint32)(rk[60:64]), &w3, w0[1], w0[2], w0[3], w0[0], 1)
			// n=109, q=1, q0=1, q1=2, q2=3, q3=0, r=13
			gsrk((*[4]uint32)(rk[64:68]), &w0, w1[1], w1[2], w1[3], w1[0], 13)
		}
	}
}

func decKeySetup(rk *roundKey, rounds int) {
	a, z := 0, rounds*4
	ra := (*[4]uint32)(rk[a : a+4])
	rz := (*[4]uint32)(rk[z : z+4])

	ra[0], rz[0] = rz[0], ra[0]
	ra[1], rz[1] = rz[1], ra[1]
	ra[2], rz[2] = rz[2], ra[2]
	ra[3], rz[3] = rz[3], ra[3]

	a += 4
	z -= 4

	for a < z {
		ra = (*[4]uint32)(rk[a : a+4])
		rz = (*[4]uint32)(rk[z : z+4])

		s0, s1, s2, s3 := wordM1(ra[0]), wordM1(ra[1]), wordM1(ra[2]), wordM1(ra[3])
		s0, s1, s2, s3 = mm(s0, s1, s2, s3) // MM
		s1, s2, s3 = p(s1, s2, s3)          // P
		s0, s1, s2, s3 = mm(s0, s1, s2, s3) // MM

		t0, t1, t2, t3 := wordM1(rz[0]), wordM1(rz[1]), wordM1(rz[2]), wordM1(rz[3])
		t0, t1, t2, t3 = mm(t0, t1, t2, t3) // MM
		t1, t2, t3 = p(t1, t2, t3)          // P
		t0, t1, t2, t3 = mm(t0, t1, t2, t3) // MM

		ra[0], ra[1], ra[2], ra[3] = t0, t1, t2, t3
		rz[0], rz[1], rz[2], rz[3] = s0, s1, s2, s3

		a += 4
		z -= 4
	}

	ra = (*[4]uint32)(rk[a : a+4])
	rz = (*[4]uint32)(rk[z : z+4])

	t0, t1, t2, t3 := wordM1(ra[0]), wordM1(ra[1]), wordM1(ra[2]), wordM1(ra[3])
	t0, t1, t2, t3 = mm(t0, t1, t2, t3) // MM
	t1, t2, t3 = p(t1, t2, t3)          // P
	t0, t1, t2, t3 = mm(t0, t1, t2, t3) // MM

	rz[0], rz[1], rz[2], rz[3] = t0, t1, t2, t3
}

func processGo(dst, src *block1, rk *roundKey, rounds int) {
	t0 := binary.BigEndian.Uint32(src[0:4])
	t1 := binary.BigEndian.Uint32(src[4:8])
	t2 := binary.BigEndian.Uint32(src[8:12])
	t3 := binary.BigEndian.Uint32(src[12:16])

	var k *[len(rk) - 16]uint32
	if rounds <= 12 {
		k = (*[len(rk) - 16]uint32)(rk[0:])
	} else {
		{ // KXL
			t0, t1, t2, t3 = t0^rk[0], t1^rk[1], t2^rk[2], t3^rk[3]
		}
		{ // FO
			t0, t1, t2, t3 = sbl1_m(t0), sbl1_m(t1), sbl1_m(t2), sbl1_m(t3) // SBL1_M
			t0, t1, t2, t3 = mm(t0, t1, t2, t3)                             // MM
			t1, t2, t3 = p(t1, t2, t3)                                      // P
			t0, t1, t2, t3 = mm(t0, t1, t2, t3)                             // MM
		}

		{ // KXL
			t0, t1, t2, t3 = t0^rk[4], t1^rk[5], t2^rk[6], t3^rk[7]
		}
		{ // FE
			t0, t1, t2, t3 = sbl2_m(t0), sbl2_m(t1), sbl2_m(t2), sbl2_m(t3) // SBL2_M
			t0, t1, t2, t3 = mm(t0, t1, t2, t3)                             // MM
			t3, t0, t1 = p(t3, t0, t1)                                      // P
			t0, t1, t2, t3 = mm(t0, t1, t2, t3)                             // MM
		}

		if rounds <= 14 {
			k = (*[len(rk) - 16]uint32)(rk[8:])
		} else {
			k = (*[len(rk) - 16]uint32)(rk[16:])
			{ // KXL
				t0, t1, t2, t3 = t0^rk[8], t1^rk[9], t2^rk[10], t3^rk[11]
			}
			{ // FO
				t0, t1, t2, t3 = sbl1_m(t0), sbl1_m(t1), sbl1_m(t2), sbl1_m(t3) // SBL1_M
				t0, t1, t2, t3 = mm(t0, t1, t2, t3)                             // MM
				t1, t2, t3 = p(t1, t2, t3)                                      // P
				t0, t1, t2, t3 = mm(t0, t1, t2, t3)                             // MM
			}

			{ // KXL
				t0, t1, t2, t3 = t0^rk[12], t1^rk[13], t2^rk[14], t3^rk[15]
			}
			{ // FE
				t0, t1, t2, t3 = sbl2_m(t0), sbl2_m(t1), sbl2_m(t2), sbl2_m(t3) // SBL2_M
				t0, t1, t2, t3 = mm(t0, t1, t2, t3)                             // MM
				t3, t0, t1 = p(t3, t0, t1)                                      // P
				t0, t1, t2, t3 = mm(t0, t1, t2, t3)                             // MM
			}
		}
	}

	for idx := 0; idx < 44; idx += 8 {
		k := (*[8]uint32)(k[idx : idx+8])

		{ // KXL
			t0, t1, t2, t3 = t0^k[0], t1^k[1], t2^k[2], t3^k[3]
		}
		{ // FO-a
			t0, t1, t2, t3 = sbl1_m(t0), sbl1_m(t1), sbl1_m(t2), sbl1_m(t3) // SBL1_M
			t0, t1, t2, t3 = mm(t0, t1, t2, t3)                             // MM
			t1, t2, t3 = p(t1, t2, t3)                                      // P
			t0, t1, t2, t3 = mm(t0, t1, t2, t3)                             // MM
		}

		if idx == 40 {
			break
		}

		{ // KXL
			t0, t1, t2, t3 = t0^k[4], t1^k[5], t2^k[6], t3^k[7]
		}
		{ // FE-a
			t0, t1, t2, t3 = sbl2_m(t0), sbl2_m(t1), sbl2_m(t2), sbl2_m(t3) // SBL2_M
			t0, t1, t2, t3 = mm(t0, t1, t2, t3)                             // MM
			t3, t0, t1 = p(t3, t0, t1)                                      // P
			t0, t1, t2, t3 = mm(t0, t1, t2, t3)                             // MM
		}
	}

	t0 = finalRound(t0^k[44]) ^ k[48]
	t1 = finalRound(t1^k[45]) ^ k[49]
	t2 = finalRound(t2^k[46]) ^ k[50]
	t3 = finalRound(t3^k[47]) ^ k[51]

	binary.BigEndian.PutUint32(dst[0:], t0)
	binary.BigEndian.PutUint32(dst[4:], t1)
	binary.BigEndian.PutUint32(dst[8:], t2)
	binary.BigEndian.PutUint32(dst[12:], t3)
}

func sbl1_m(t uint32) uint32 {
	return s1[t>>24] ^ s2[(t>>16)&0xff] ^ x1[(t>>8)&0xff] ^ x2[t&0xff]
}

func sbl2_m(t uint32) uint32 {
	return x1[t>>24] ^ x2[(t>>16)&0xff] ^ s1[(t>>8)&0xff] ^ s2[t&0xff]
}

func finalRound(t uint32) uint32 {
	return (x1[t>>24] & 0xff000000) |
		(x2[(t>>16)&0xff] & 0x00ff0000) |
		(s1[(t>>8)&0xff] & 0x0000ff00) |
		(s2[t&0xff] & 0x000000ff)
}

func mm(t0, t1, t2, t3 uint32) (uint32, uint32, uint32, uint32) {
	t1 ^= t2
	t2 ^= t3
	t0 ^= t1
	t3 ^= t1
	t2 ^= t0
	t1 ^= t2
	return t0, t1, t2, t3
}

func p(t1, t2, t3 uint32) (uint32, uint32, uint32) {
	t1 = bits.ReverseBytes32(bits.RotateLeft32(t1, 16)) // (T1) = (((T1)<< 8)&0xff00ff00) ^ (((T1)>> 8)&0x00ff00ff);
	t2 = bits.RotateLeft32(t2, 16)                      // (T2) = (((T2)<<16)&0xffff0000) ^ (((T2)>>16)&0x0000ffff)
	t3 = bits.ReverseBytes32(t3)                        // ReverseWord((T3));
	return t1, t2, t3
}

func wordM1(x uint32) (y uint32) {
	//y = x<<8 ^ x>>8 ^ x<<16 ^ x>>16 ^ x<<24 ^ x>>24
	y = bits.RotateLeft32(x, 8) ^ bits.RotateLeft32(x, 16) ^ bits.RotateLeft32(x, 24)
	return y
}

func gsrk(rk *[4]uint32, x *[4]uint32, w0, w1, w2, w3 uint32, r uint) {
	rk[0] = x[0] ^ uint32((uint64(w3)<<32|uint64(w0))>>r)
	rk[1] = x[1] ^ uint32((uint64(w0)<<32|uint64(w1))>>r)
	rk[2] = x[2] ^ uint32((uint64(w1)<<32|uint64(w2))>>r)
	rk[3] = x[3] ^ uint32((uint64(w2)<<32|uint64(w3))>>r)
}
