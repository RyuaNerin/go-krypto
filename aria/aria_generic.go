package aria

import (
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"math/bits"
)

func newCipherGo(key []byte) (cipher.Block, error) {
	ctx := new(ariaContext)
	ctx.rounds = (len(key)*8 + 256) / 32

	encKeySetup(ctx.ek[:], key)

	ctx.dk = ctx.ek
	decKeySetup(ctx.dk[:], ctx.rounds)
	return ctx, nil
}

func (s *ariaContext) BlockSize() int {
	return BlockSize
}

func (s *ariaContext) Encrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf("krypto/aria: invalid block size %d (src)", len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf("krypto/aria: invalid block size %d (dst)", len(dst)))
	}

	processGo(dst, src, s.ek[:], s.rounds)
}

func (s *ariaContext) Decrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf("krypto/aria: invalid block size %d (src)", len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf("krypto/aria: invalid block size %d (dst)", len(dst)))
	}

	processGo(dst, src, s.dk[:], s.rounds)
}

func encKeySetup(rk []byte, key []byte) {
	m := cimpl{
		rk:      rk,
		mk:      key,
		keyBits: len(key) * 8,
	}

	m.encKeySetup()
}

func decKeySetup(dk []byte, rounds int) {
	dl := func(i, o []byte) {
		var T byte

		T = i[3] ^ i[4] ^ i[9] ^ i[14]
		o[0] = i[6] ^ i[8] ^ i[13] ^ T
		o[5] = i[1] ^ i[10] ^ i[15] ^ T
		o[11] = i[2] ^ i[7] ^ i[12] ^ T
		o[14] = i[0] ^ i[5] ^ i[11] ^ T
		T = i[2] ^ i[5] ^ i[8] ^ i[15]
		o[1] = i[7] ^ i[9] ^ i[12] ^ T
		o[4] = i[0] ^ i[11] ^ i[14] ^ T
		o[10] = i[3] ^ i[6] ^ i[13] ^ T
		o[15] = i[1] ^ i[4] ^ i[10] ^ T
		T = i[1] ^ i[6] ^ i[11] ^ i[12]
		o[2] = i[4] ^ i[10] ^ i[15] ^ T
		o[7] = i[3] ^ i[8] ^ i[13] ^ T
		o[9] = i[0] ^ i[5] ^ i[14] ^ T
		o[12] = i[2] ^ i[7] ^ i[9] ^ T
		T = i[0] ^ i[7] ^ i[10] ^ i[13]
		o[3] = i[5] ^ i[11] ^ i[14] ^ T
		o[6] = i[2] ^ i[9] ^ i[12] ^ T
		o[8] = i[1] ^ i[4] ^ i[15] ^ T
		o[13] = i[3] ^ i[6] ^ i[8] ^ T
	}

	var t [16]byte

	for j := 0; j < 16; j++ {
		t[j] = dk[j]
		dk[j] = dk[16*rounds+j]
		dk[16*rounds+j] = t[j]
	}
	for i := 1; i <= rounds/2; i++ {
		dl(dk[i*16:], t[:])
		dl(dk[(rounds-i)*16:], dk[i*16:])
		for j := 0; j < 16; j++ {
			dk[(rounds-i)*16+j] = t[j]
		}
	}
}

func processGo(dst, src, rk []byte, round int) {
	m := cimpl{
		rk: rk,
	}

	m.vars.t0 = binary.BigEndian.Uint32(src[0*4:]) // WordLoad(WO(i,0), t0)
	m.vars.t1 = binary.BigEndian.Uint32(src[1*4:]) // WordLoad(WO(i,1), t1)
	m.vars.t2 = binary.BigEndian.Uint32(src[2*4:]) // WordLoad(WO(i,2), t2)
	m.vars.t3 = binary.BigEndian.Uint32(src[3*4:]) // WordLoad(WO(i,3), t3)

	if round > 12 {
		m.KXL()
		m.FO()
		m.KXL()
		m.FE()
	}
	if round > 14 {
		m.KXL()
		m.FO()
		m.KXL()
		m.FE()
	}

	m.KXL()
	m.FO()
	m.KXL()
	m.FE()
	m.KXL()
	m.FO()
	m.KXL()
	m.FE()
	m.KXL()
	m.FO()
	m.KXL()
	m.FE()

	m.KXL()
	m.FO()
	m.KXL()
	m.FE()
	m.KXL()
	m.FO()
	m.KXL()
	m.FE()
	m.KXL()
	m.FO()
	m.KXL()

	dst[0x0] = (byte)(x1[m.BRF(m.vars.t0, 0x18)]>>0) ^ rk[m.vars.rkIndex+0x3]
	dst[0x1] = (byte)(x2[m.BRF(m.vars.t0, 0x10)]>>8) ^ rk[m.vars.rkIndex+0x2]
	dst[0x2] = (byte)(s1[m.BRF(m.vars.t0, 0x08)]>>0) ^ rk[m.vars.rkIndex+0x1]
	dst[0x3] = (byte)(s2[m.BRF(m.vars.t0, 0x00)]>>0) ^ rk[m.vars.rkIndex+0x0]

	dst[0x4] = (byte)(x1[m.BRF(m.vars.t1, 0x18)]>>0) ^ rk[m.vars.rkIndex+0x7]
	dst[0x5] = (byte)(x2[m.BRF(m.vars.t1, 0x10)]>>8) ^ rk[m.vars.rkIndex+0x6]
	dst[0x6] = (byte)(s1[m.BRF(m.vars.t1, 0x08)]>>0) ^ rk[m.vars.rkIndex+0x5]
	dst[0x7] = (byte)(s2[m.BRF(m.vars.t1, 0x00)]>>0) ^ rk[m.vars.rkIndex+0x4]

	dst[0x8] = (byte)(x1[m.BRF(m.vars.t2, 0x18)]>>0) ^ rk[m.vars.rkIndex+0xB]
	dst[0x9] = (byte)(x2[m.BRF(m.vars.t2, 0x10)]>>8) ^ rk[m.vars.rkIndex+0xA]
	dst[0xA] = (byte)(s1[m.BRF(m.vars.t2, 0x08)]>>0) ^ rk[m.vars.rkIndex+0x9]
	dst[0xB] = (byte)(s2[m.BRF(m.vars.t2, 0x00)]>>0) ^ rk[m.vars.rkIndex+0x8]

	dst[0xC] = (byte)(x1[m.BRF(m.vars.t3, 0x18)]>>0) ^ rk[m.vars.rkIndex+0xF]
	dst[0xD] = (byte)(x2[m.BRF(m.vars.t3, 0x10)]>>8) ^ rk[m.vars.rkIndex+0xE]
	dst[0xE] = (byte)(s1[m.BRF(m.vars.t3, 0x08)]>>0) ^ rk[m.vars.rkIndex+0xD]
	dst[0xF] = (byte)(s2[m.BRF(m.vars.t3, 0x00)]>>0) ^ rk[m.vars.rkIndex+0xC]
}

type cimpl struct {
	rk      []byte
	mk      []byte // make key
	keyBits int

	vars struct {
		t0, t1, t2, t3 uint32
		w0, w1, w2, w3 [4]uint32
		q              int
		rkIndex        int
	}
}

func (m *cimpl) encKeySetup() {
	m.vars.w0[0] = binary.BigEndian.Uint32(m.mk[0*4:]) // WordLoad(WO(mk, 0), w0[0])
	m.vars.w0[1] = binary.BigEndian.Uint32(m.mk[1*4:]) // WordLoad(WO(mk, 1), w0[1])
	m.vars.w0[2] = binary.BigEndian.Uint32(m.mk[2*4:]) // WordLoad(WO(mk, 2), w0[2])
	m.vars.w0[3] = binary.BigEndian.Uint32(m.mk[3*4:]) // WordLoad(WO(mk, 3), w0[3])

	m.vars.q = (m.keyBits - 128) / 64
	m.vars.t0 = m.vars.w0[0] ^ krk[m.vars.q*4+0]
	m.vars.t1 = m.vars.w0[1] ^ krk[m.vars.q*4+1]
	m.vars.t2 = m.vars.w0[2] ^ krk[m.vars.q*4+2]
	m.vars.t3 = m.vars.w0[3] ^ krk[m.vars.q*4+3]
	m.FO()
	if m.keyBits > 128 {
		m.vars.w1[0] = binary.BigEndian.Uint32(m.mk[4*4:]) // WordLoad(WO(mk, 4), w1[0])
		m.vars.w1[1] = binary.BigEndian.Uint32(m.mk[5*4:]) // WordLoad(WO(mk, 5), w1[1])
		if m.keyBits > 192 {
			m.vars.w1[2] = binary.BigEndian.Uint32(m.mk[6*4:]) // WordLoad(WO(mk,6), w1[2]);
			m.vars.w1[3] = binary.BigEndian.Uint32(m.mk[7*4:]) // WordLoad(WO(mk,7), w1[3]);
		} else {
			m.vars.w1[2] = 0
			m.vars.w1[3] = 0
		}
	} else {
		m.vars.w1[0] = 0
		m.vars.w1[1] = 0
		m.vars.w1[2] = 0
		m.vars.w1[3] = 0
	}
	m.vars.w1[0] ^= m.vars.t0
	m.vars.w1[1] ^= m.vars.t1
	m.vars.w1[2] ^= m.vars.t2
	m.vars.w1[3] ^= m.vars.t3
	m.vars.t0 = m.vars.w1[0]
	m.vars.t1 = m.vars.w1[1]
	m.vars.t2 = m.vars.w1[2]
	m.vars.t3 = m.vars.w1[3]

	m.vars.q = (m.vars.q + 1) % 3
	m.vars.t0 ^= krk[m.vars.q*4+0]
	m.vars.t1 ^= krk[m.vars.q*4+1]
	m.vars.t2 ^= krk[m.vars.q*4+2]
	m.vars.t3 ^= krk[m.vars.q*4+3]
	m.FE()
	m.vars.t0 ^= m.vars.w0[0]
	m.vars.t1 ^= m.vars.w0[1]
	m.vars.t2 ^= m.vars.w0[2]
	m.vars.t3 ^= m.vars.w0[3]
	m.vars.w2[0] = m.vars.t0
	m.vars.w2[1] = m.vars.t1
	m.vars.w2[2] = m.vars.t2
	m.vars.w2[3] = m.vars.t3

	m.vars.q = (m.vars.q + 1) % 3
	m.vars.t0 ^= krk[m.vars.q*4+0]
	m.vars.t1 ^= krk[m.vars.q*4+1]
	m.vars.t2 ^= krk[m.vars.q*4+2]
	m.vars.t3 ^= krk[m.vars.q*4+3]
	m.FO()
	m.vars.w3[0] = m.vars.t0 ^ m.vars.w1[0]
	m.vars.w3[1] = m.vars.t1 ^ m.vars.w1[1]
	m.vars.w3[2] = m.vars.t2 ^ m.vars.w1[2]
	m.vars.w3[3] = m.vars.t3 ^ m.vars.w1[3]

	m.GSRK(m.vars.w0, m.vars.w1, 19)
	m.GSRK(m.vars.w1, m.vars.w2, 19)
	m.GSRK(m.vars.w2, m.vars.w3, 19)
	m.GSRK(m.vars.w3, m.vars.w0, 19)
	m.GSRK(m.vars.w0, m.vars.w1, 31)
	m.GSRK(m.vars.w1, m.vars.w2, 31)
	m.GSRK(m.vars.w2, m.vars.w3, 31)
	m.GSRK(m.vars.w3, m.vars.w0, 31)
	m.GSRK(m.vars.w0, m.vars.w1, 67)
	m.GSRK(m.vars.w1, m.vars.w2, 67)
	m.GSRK(m.vars.w2, m.vars.w3, 67)
	m.GSRK(m.vars.w3, m.vars.w0, 67)
	m.GSRK(m.vars.w0, m.vars.w1, 97)
	if m.keyBits > 128 {
		m.GSRK(m.vars.w1, m.vars.w2, 97)
		m.GSRK(m.vars.w2, m.vars.w3, 97)
	}
	if m.keyBits > 192 {
		m.GSRK(m.vars.w3, m.vars.w0, 97)
		m.GSRK(m.vars.w0, m.vars.w1, 109)
	}
}

func (cimpl) BRF(T uint32, R int) int {
	return int((T)>>R) & 0xFF
}

func (m *cimpl) SBL1_M() {
	m.vars.t0 = s1[m.BRF(m.vars.t0, 24)] ^ s2[m.BRF(m.vars.t0, 16)] ^ x1[m.BRF(m.vars.t0, 8)] ^ x2[m.BRF(m.vars.t0, 0)]
	m.vars.t1 = s1[m.BRF(m.vars.t1, 24)] ^ s2[m.BRF(m.vars.t1, 16)] ^ x1[m.BRF(m.vars.t1, 8)] ^ x2[m.BRF(m.vars.t1, 0)]
	m.vars.t2 = s1[m.BRF(m.vars.t2, 24)] ^ s2[m.BRF(m.vars.t2, 16)] ^ x1[m.BRF(m.vars.t2, 8)] ^ x2[m.BRF(m.vars.t2, 0)]
	m.vars.t3 = s1[m.BRF(m.vars.t3, 24)] ^ s2[m.BRF(m.vars.t3, 16)] ^ x1[m.BRF(m.vars.t3, 8)] ^ x2[m.BRF(m.vars.t3, 0)]
}
func (m *cimpl) SBL2_M() {
	m.vars.t0 = x1[m.BRF(m.vars.t0, 24)] ^ x2[m.BRF(m.vars.t0, 16)] ^ s1[m.BRF(m.vars.t0, 8)] ^ s2[m.BRF(m.vars.t0, 0)]
	m.vars.t1 = x1[m.BRF(m.vars.t1, 24)] ^ x2[m.BRF(m.vars.t1, 16)] ^ s1[m.BRF(m.vars.t1, 8)] ^ s2[m.BRF(m.vars.t1, 0)]
	m.vars.t2 = x1[m.BRF(m.vars.t2, 24)] ^ x2[m.BRF(m.vars.t2, 16)] ^ s1[m.BRF(m.vars.t2, 8)] ^ s2[m.BRF(m.vars.t2, 0)]
	m.vars.t3 = x1[m.BRF(m.vars.t3, 24)] ^ x2[m.BRF(m.vars.t3, 16)] ^ s1[m.BRF(m.vars.t3, 8)] ^ s2[m.BRF(m.vars.t3, 0)]
}
func (m *cimpl) MM() {
	(m.vars.t1) ^= (m.vars.t2)
	(m.vars.t2) ^= (m.vars.t3)
	(m.vars.t0) ^= (m.vars.t1)
	(m.vars.t3) ^= (m.vars.t1)
	(m.vars.t2) ^= (m.vars.t0)
	(m.vars.t1) ^= (m.vars.t2)
}
func (m *cimpl) P(T0, T1, T2, T3 *uint32) {
	*T1 = (((*T1) << 8) & uint32(0xff00ff00)) ^ (((*T1) >> 8) & uint32(0x00ff00ff))
	*T2 = (((*T2) << 16) & uint32(0xffff0000)) ^ (((*T2) >> 16) & uint32(0x0000ffff))
	*T3 = bits.ReverseBytes32(*T3)
}

func (m *cimpl) KXL() {
	m.vars.t0 ^= binary.LittleEndian.Uint32(m.rk[m.vars.rkIndex+0*4:])
	m.vars.t1 ^= binary.LittleEndian.Uint32(m.rk[m.vars.rkIndex+1*4:])
	m.vars.t2 ^= binary.LittleEndian.Uint32(m.rk[m.vars.rkIndex+2*4:])
	m.vars.t3 ^= binary.LittleEndian.Uint32(m.rk[m.vars.rkIndex+3*4:])
	m.vars.rkIndex += 16
}
func (m *cimpl) FO() {
	m.SBL1_M()
	m.MM()
	m.P(&m.vars.t0, &m.vars.t1, &m.vars.t2, &m.vars.t3)
	m.MM()
}
func (m *cimpl) FE() {
	m.SBL2_M()
	m.MM()
	m.P(&m.vars.t2, &m.vars.t3, &m.vars.t0, &m.vars.t1)
	m.MM()
}
func (m *cimpl) GSRK(X [4]uint32, Y [4]uint32, n int) {
	m.vars.q = 4 - ((n) / 32)
	r := (n) % 32
	binary.LittleEndian.PutUint32(m.rk[m.vars.rkIndex+0*4:], (X[0])^((Y[(m.vars.q+0)%4])>>r)^((Y[(m.vars.q+3)%4])<<(32-r))) // WO(rk, 0) =
	binary.LittleEndian.PutUint32(m.rk[m.vars.rkIndex+1*4:], (X[1])^((Y[(m.vars.q+1)%4])>>r)^((Y[(m.vars.q+0)%4])<<(32-r))) // WO(rk, 1) =
	binary.LittleEndian.PutUint32(m.rk[m.vars.rkIndex+2*4:], (X[2])^((Y[(m.vars.q+2)%4])>>r)^((Y[(m.vars.q+1)%4])<<(32-r))) // WO(rk, 2) =
	binary.LittleEndian.PutUint32(m.rk[m.vars.rkIndex+3*4:], (X[3])^((Y[(m.vars.q+3)%4])>>r)^((Y[(m.vars.q+2)%4])<<(32-r))) // WO(rk, 3) =
	m.vars.rkIndex += 16
}
