//go:build amd64 || arm64 || loong64 || mips || mips64 || mips64le || mipsle || ppc64 || ppc64le || riscv64

/**
| GOARCH   | Bits | GPR  |
|----------|------|------|
| arm64    | 64   | ~ 26 |
| ppc64    | 64   | ~ 26 |
| ppc64le  | 64   | ~ 26 |
| loong64  | 64   | ~ 26 |
| riscv64  | 64   | ~ 25 |
| mips64   | 64   | ~ 22 |
| mips64le | 64   | ~ 22 |
| mips     | 32   | ~ 22 |
| mipsle   | 32   | ~ 22 |
| amd64    | 64   | ~ 13 | benchmarked (~13 %)
| s390x    | 64   | ~ 12 |
| arm      | 32   | ~ 11 |
| 386      | 32   | ~ 7  |
| wasm     | -    |      |
*/

package aria

import (
	"encoding/binary"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/alias"
)

var _ krypto.Block4 = (*ariaContext)(nil)

func (ctx *ariaContext) Encrypt4(dst, src []byte) { //nolint:dupl
	ctx.z.ThrowIfZeroized()

	if len(src) < 4*BlockSize {
		panic(msgNotFullBlock)
	}
	if len(dst) < 4*BlockSize {
		panic(msgNotFullBlock)
	}
	if alias.InexactOverlap(dst[:4*BlockSize], src[:4*BlockSize]) {
		panic(msgBufferOverlap)
	}

	_ = dst[4*BlockSize-1]
	_ = src[4*BlockSize-1]

	process2Go((*block2)(dst[0*BlockSize:2*BlockSize]), (*block2)(src[0*BlockSize:2*BlockSize]), &ctx.ek, ctx.rounds)
	process2Go((*block2)(dst[2*BlockSize:4*BlockSize]), (*block2)(src[2*BlockSize:4*BlockSize]), &ctx.ek, ctx.rounds)
}

func (ctx *ariaContext) Decrypt4(dst, src []byte) { //nolint:dupl
	ctx.z.ThrowIfZeroized()

	if len(src) < 4*BlockSize {
		panic(msgNotFullBlock)
	}
	if len(dst) < 4*BlockSize {
		panic(msgNotFullBlock)
	}
	if alias.InexactOverlap(dst[:4*BlockSize], src[:4*BlockSize]) {
		panic(msgBufferOverlap)
	}

	_ = dst[4*BlockSize-1]
	_ = src[4*BlockSize-1]

	process2Go((*block2)(dst[0*BlockSize:2*BlockSize]), (*block2)(src[0*BlockSize:2*BlockSize]), &ctx.dk, ctx.rounds)
	process2Go((*block2)(dst[2*BlockSize:4*BlockSize]), (*block2)(src[2*BlockSize:4*BlockSize]), &ctx.dk, ctx.rounds)
}

// process2Go processes.
// it maybe be fater than twise of processGo.
func process2Go(dst, src *block2, rk *roundKey, rounds int) {
	t0a := binary.BigEndian.Uint32(src[0:4])
	t1a := binary.BigEndian.Uint32(src[4:8])
	t2a := binary.BigEndian.Uint32(src[8:12])
	t3a := binary.BigEndian.Uint32(src[12:16])
	t0b := binary.BigEndian.Uint32(src[16:20])
	t1b := binary.BigEndian.Uint32(src[20:24])
	t2b := binary.BigEndian.Uint32(src[24:28])
	t3b := binary.BigEndian.Uint32(src[28:32])

	var k *[len(rk) - 16]uint32
	if rounds <= 12 {
		k = (*[len(rk) - 16]uint32)(rk[0:])
	} else {
		{ // KXL
			t0a, t1a, t2a, t3a = t0a^rk[0], t1a^rk[1], t2a^rk[2], t3a^rk[3]
			t0b, t1b, t2b, t3b = t0b^rk[0], t1b^rk[1], t2b^rk[2], t3b^rk[3]
		}
		{ // FO
			t0a, t1a, t2a, t3a = sbl1_m(t0a), sbl1_m(t1a), sbl1_m(t2a), sbl1_m(t3a) // SBL1_M   a
			t0b, t1b, t2b, t3b = sbl1_m(t0b), sbl1_m(t1b), sbl1_m(t2b), sbl1_m(t3b) // SBL1_M     b
			t0a, t1a, t2a, t3a = mm(t0a, t1a, t2a, t3a)                             // MM       a
			t0b, t1b, t2b, t3b = mm(t0b, t1b, t2b, t3b)                             // MM         b
			t1a, t2a, t3a = p(t1a, t2a, t3a)                                        // P        a
			t1b, t2b, t3b = p(t1b, t2b, t3b)                                        // P          b
			t0a, t1a, t2a, t3a = mm(t0a, t1a, t2a, t3a)                             // MM       a
			t0b, t1b, t2b, t3b = mm(t0b, t1b, t2b, t3b)                             // MM         b
		}

		{ // KXL
			t0a, t1a, t2a, t3a = t0a^rk[4], t1a^rk[5], t2a^rk[6], t3a^rk[7]
			t0b, t1b, t2b, t3b = t0b^rk[4], t1b^rk[5], t2b^rk[6], t3b^rk[7]
		}
		{ // FE
			t0a, t1a, t2a, t3a = sbl2_m(t0a), sbl2_m(t1a), sbl2_m(t2a), sbl2_m(t3a) // SBL2_M   a
			t0b, t1b, t2b, t3b = sbl2_m(t0b), sbl2_m(t1b), sbl2_m(t2b), sbl2_m(t3b) // SBL2_M     b
			t0a, t1a, t2a, t3a = mm(t0a, t1a, t2a, t3a)                             // MM       a
			t0b, t1b, t2b, t3b = mm(t0b, t1b, t2b, t3b)                             // MM         b
			t3a, t0a, t1a = p(t3a, t0a, t1a)                                        // P        a
			t3b, t0b, t1b = p(t3b, t0b, t1b)                                        // P          b
			t0a, t1a, t2a, t3a = mm(t0a, t1a, t2a, t3a)                             // MM       a
			t0b, t1b, t2b, t3b = mm(t0b, t1b, t2b, t3b)                             // MM         b
		}

		if rounds <= 14 {
			k = (*[len(rk) - 16]uint32)(rk[8:])
		} else {
			k = (*[len(rk) - 16]uint32)(rk[16:])
			{ // KXL
				t0a, t1a, t2a, t3a = t0a^rk[8], t1a^rk[9], t2a^rk[10], t3a^rk[11]
				t0b, t1b, t2b, t3b = t0b^rk[8], t1b^rk[9], t2b^rk[10], t3b^rk[11]
			}
			{ // FO
				t0a, t1a, t2a, t3a = sbl1_m(t0a), sbl1_m(t1a), sbl1_m(t2a), sbl1_m(t3a) // SBL1_M   a
				t0b, t1b, t2b, t3b = sbl1_m(t0b), sbl1_m(t1b), sbl1_m(t2b), sbl1_m(t3b) // SBL1_M     b
				t0a, t1a, t2a, t3a = mm(t0a, t1a, t2a, t3a)                             // MM       a
				t0b, t1b, t2b, t3b = mm(t0b, t1b, t2b, t3b)                             // MM         b
				t1a, t2a, t3a = p(t1a, t2a, t3a)                                        // P        a
				t1b, t2b, t3b = p(t1b, t2b, t3b)                                        // P          b
				t0a, t1a, t2a, t3a = mm(t0a, t1a, t2a, t3a)                             // MM       a
				t0b, t1b, t2b, t3b = mm(t0b, t1b, t2b, t3b)                             // MM         b
			}

			{ // KXL
				t0a, t1a, t2a, t3a = t0a^rk[12], t1a^rk[13], t2a^rk[14], t3a^rk[15]
				t0b, t1b, t2b, t3b = t0b^rk[12], t1b^rk[13], t2b^rk[14], t3b^rk[15]
			}
			{ // FE
				t0a, t1a, t2a, t3a = sbl2_m(t0a), sbl2_m(t1a), sbl2_m(t2a), sbl2_m(t3a) // SBL2_M   a
				t0b, t1b, t2b, t3b = sbl2_m(t0b), sbl2_m(t1b), sbl2_m(t2b), sbl2_m(t3b) // SBL2_M     b
				t0a, t1a, t2a, t3a = mm(t0a, t1a, t2a, t3a)                             // MM       a
				t0b, t1b, t2b, t3b = mm(t0b, t1b, t2b, t3b)                             // MM         b
				t3a, t0a, t1a = p(t3a, t0a, t1a)                                        // P        a
				t3b, t0b, t1b = p(t3b, t0b, t1b)                                        // P          b
				t0a, t1a, t2a, t3a = mm(t0a, t1a, t2a, t3a)                             // MM       a
				t0b, t1b, t2b, t3b = mm(t0b, t1b, t2b, t3b)                             // MM         b
			}
		}
	}

	for idx := 0; idx < 40; idx += 8 {
		k := (*[8]uint32)(k[idx : idx+8])

		{ // KXL
			t0a, t1a, t2a, t3a = t0a^k[0], t1a^k[1], t2a^k[2], t3a^k[3]
			t0b, t1b, t2b, t3b = t0b^k[0], t1b^k[1], t2b^k[2], t3b^k[3]
		}
		{ // FO
			t0a, t1a, t2a, t3a = sbl1_m(t0a), sbl1_m(t1a), sbl1_m(t2a), sbl1_m(t3a) // SBL1_M   a
			t0b, t1b, t2b, t3b = sbl1_m(t0b), sbl1_m(t1b), sbl1_m(t2b), sbl1_m(t3b) // SBL1_M     b
			t0a, t1a, t2a, t3a = mm(t0a, t1a, t2a, t3a)                             // MM       a
			t0b, t1b, t2b, t3b = mm(t0b, t1b, t2b, t3b)                             // MM         b
			t1a, t2a, t3a = p(t1a, t2a, t3a)                                        // P        a
			t1b, t2b, t3b = p(t1b, t2b, t3b)                                        // P          b
			t0a, t1a, t2a, t3a = mm(t0a, t1a, t2a, t3a)                             // MM       a
			t0b, t1b, t2b, t3b = mm(t0b, t1b, t2b, t3b)                             // MM         b
		}

		if idx == 40 {
			break
		}

		{ // KXL
			t0a, t1a, t2a, t3a = t0a^k[4], t1a^k[5], t2a^k[6], t3a^k[7]
			t0b, t1b, t2b, t3b = t0b^k[4], t1b^k[5], t2b^k[6], t3b^k[7]
		}
		{ // FE
			t0a, t1a, t2a, t3a = sbl2_m(t0a), sbl2_m(t1a), sbl2_m(t2a), sbl2_m(t3a) // SBL2_M  a
			t0b, t1b, t2b, t3b = sbl2_m(t0b), sbl2_m(t1b), sbl2_m(t2b), sbl2_m(t3b) // SBL2_M    b
			t0a, t1a, t2a, t3a = mm(t0a, t1a, t2a, t3a)                             // MM      a
			t0b, t1b, t2b, t3b = mm(t0b, t1b, t2b, t3b)                             // MM        b
			t3a, t0a, t1a = p(t3a, t0a, t1a)                                        // P       a
			t3b, t0b, t1b = p(t3b, t0b, t1b)                                        // P         b
			t0a, t1a, t2a, t3a = mm(t0a, t1a, t2a, t3a)                             // MM      a
			t0b, t1b, t2b, t3b = mm(t0b, t1b, t2b, t3b)                             // MM        b
		}
	}

	{ // KXL
		t0a, t1a, t2a, t3a = t0a^k[40], t1a^k[41], t2a^k[42], t3a^k[43]
		t0b, t1b, t2b, t3b = t0b^k[40], t1b^k[41], t2b^k[42], t3b^k[43]
	}
	{ // FO
		t0a, t1a, t2a, t3a = sbl1_m(t0a), sbl1_m(t1a), sbl1_m(t2a), sbl1_m(t3a) // SBL1_M   a
		t0b, t1b, t2b, t3b = sbl1_m(t0b), sbl1_m(t1b), sbl1_m(t2b), sbl1_m(t3b) // SBL1_M     b
		t0a, t1a, t2a, t3a = mm(t0a, t1a, t2a, t3a)                             // MM       a
		t0b, t1b, t2b, t3b = mm(t0b, t1b, t2b, t3b)                             // MM         b
		t1a, t2a, t3a = p(t1a, t2a, t3a)                                        // P        a
		t1b, t2b, t3b = p(t1b, t2b, t3b)                                        // P          b
		t0a, t1a, t2a, t3a = mm(t0a, t1a, t2a, t3a)                             // MM       a
		t0b, t1b, t2b, t3b = mm(t0b, t1b, t2b, t3b)                             // MM         b
	}

	binary.BigEndian.PutUint32(dst[0:4], finalRound(t0a^k[44])^k[48])
	binary.BigEndian.PutUint32(dst[4:8], finalRound(t1a^k[45])^k[49])
	binary.BigEndian.PutUint32(dst[8:12], finalRound(t2a^k[46])^k[50])
	binary.BigEndian.PutUint32(dst[12:16], finalRound(t3a^k[47])^k[51])
	binary.BigEndian.PutUint32(dst[16:20], finalRound(t0b^k[44])^k[48])
	binary.BigEndian.PutUint32(dst[20:24], finalRound(t1b^k[45])^k[49])
	binary.BigEndian.PutUint32(dst[24:28], finalRound(t2b^k[46])^k[50])
	binary.BigEndian.PutUint32(dst[28:32], finalRound(t3b^k[47])^k[51])
}
