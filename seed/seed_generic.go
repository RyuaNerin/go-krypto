package seed

import (
	"crypto/cipher"
	"encoding/binary"
	"fmt"

	"github.com/RyuaNerin/go-krypto/internal/golang.org/x/sys/cpu"
)

var byteOrder binary.ByteOrder

func init() {
	if cpu.IsBigEndian {
		byteOrder = binary.LittleEndian
	} else {
		byteOrder = binary.BigEndian
	}
}

type seed128 struct {
	pdwRoundKey [32]uint32
}

func new128(key []byte) cipher.Block {
	c := new(seed128)
	c.keySchedKey(key)
	return c
}

func (s *seed128) BlockSize() int {
	return BlockSize
}

func (s *seed128) keySchedKey(key []byte) {
	var T0, T1 uint32
	A := byteOrder.Uint32(key[0:])
	B := byteOrder.Uint32(key[4:])
	C := byteOrder.Uint32(key[8:])
	D := byteOrder.Uint32(key[12:])

	// RoundKeyUpdate0(K+0, A, B, C, D, KC0)
	{
		T0 = A + C - kc[0]
		T1 = B + kc[0] - D
		s.pdwRoundKey[0] = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		s.pdwRoundKey[1] = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = A
		A = (A >> 8) ^ (B << 24)
		B = (B >> 8) ^ (T0 << 24)
	}
	// RoundKeyUpdate1(K+ 2, A, B, C, D, KC1)
	{
		T0 = A + C - kc[1]
		T1 = B + kc[1] - D
		s.pdwRoundKey[2] = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		s.pdwRoundKey[3] = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = C
		C = (C << 8) ^ (D >> 24)
		D = (D << 8) ^ (T0 >> 24)
	}
	// RoundKeyUpdate0(K+4, A, B, C, D, KC2)
	{
		T0 = A + C - kc[2]
		T1 = B + kc[2] - D
		s.pdwRoundKey[4] = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		s.pdwRoundKey[5] = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = A
		A = (A >> 8) ^ (B << 24)
		B = (B >> 8) ^ (T0 << 24)
	}
	// RoundKeyUpdate1(K+ 6, A, B, C, D, KC3)
	{
		T0 = A + C - kc[3]
		T1 = B + kc[3] - D
		s.pdwRoundKey[6] = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		s.pdwRoundKey[7] = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = C
		C = (C << 8) ^ (D >> 24)
		D = (D << 8) ^ (T0 >> 24)
	}
	// RoundKeyUpdate0(K+8, A, B, C, D, KC4)
	{
		T0 = A + C - kc[4]
		T1 = B + kc[4] - D
		s.pdwRoundKey[8] = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		s.pdwRoundKey[9] = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = A
		A = (A >> 8) ^ (B << 24)
		B = (B >> 8) ^ (T0 << 24)
	}
	// RoundKeyUpdate1(K+ 10, A, B, C, D, KC5)
	{
		T0 = A + C - kc[5]
		T1 = B + kc[5] - D
		s.pdwRoundKey[10] = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		s.pdwRoundKey[11] = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = C
		C = (C << 8) ^ (D >> 24)
		D = (D << 8) ^ (T0 >> 24)
	}
	// RoundKeyUpdate0(K+12, A, B, C, D, KC6)
	{
		T0 = A + C - kc[6]
		T1 = B + kc[6] - D
		s.pdwRoundKey[12] = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		s.pdwRoundKey[13] = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = A
		A = (A >> 8) ^ (B << 24)
		B = (B >> 8) ^ (T0 << 24)
	}
	// RoundKeyUpdate1(K+ 14, A, B, C, D, KC7)
	{
		T0 = A + C - kc[7]
		T1 = B + kc[7] - D
		s.pdwRoundKey[14] = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		s.pdwRoundKey[15] = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = C
		C = (C << 8) ^ (D >> 24)
		D = (D << 8) ^ (T0 >> 24)
	}
	// RoundKeyUpdate0(K+16, A, B, C, D, KC8)
	{
		T0 = A + C - kc[8]
		T1 = B + kc[8] - D
		s.pdwRoundKey[16] = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		s.pdwRoundKey[17] = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = A
		A = (A >> 8) ^ (B << 24)
		B = (B >> 8) ^ (T0 << 24)
	}
	// RoundKeyUpdate1(K+ 18, A, B, C, D, KC9)
	{
		T0 = A + C - kc[9]
		T1 = B + kc[9] - D
		s.pdwRoundKey[18] = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		s.pdwRoundKey[19] = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = C
		C = (C << 8) ^ (D >> 24)
		D = (D << 8) ^ (T0 >> 24)
	}
	// RoundKeyUpdate0(K+20, A, B, C, D, KC10)
	{
		T0 = A + C - kc[10]
		T1 = B + kc[10] - D
		s.pdwRoundKey[20] = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		s.pdwRoundKey[21] = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = A
		A = (A >> 8) ^ (B << 24)
		B = (B >> 8) ^ (T0 << 24)
	}
	// RoundKeyUpdate1(K+ 22, A, B, C, D, KC11)
	{
		T0 = A + C - kc[11]
		T1 = B + kc[11] - D
		s.pdwRoundKey[22] = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		s.pdwRoundKey[23] = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = C
		C = (C << 8) ^ (D >> 24)
		D = (D << 8) ^ (T0 >> 24)
	}
	// RoundKeyUpdate0(K+24, A, B, C, D, KC12)
	{
		T0 = A + C - kc[12]
		T1 = B + kc[12] - D
		s.pdwRoundKey[24] = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		s.pdwRoundKey[25] = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = A
		A = (A >> 8) ^ (B << 24)
		B = (B >> 8) ^ (T0 << 24)
	}
	// RoundKeyUpdate1(K+ 26, A, B, C, D, KC13)
	{
		T0 = A + C - kc[13]
		T1 = B + kc[13] - D
		s.pdwRoundKey[26] = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		s.pdwRoundKey[27] = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = C
		C = (C << 8) ^ (D >> 24)
		D = (D << 8) ^ (T0 >> 24)
	}
	// RoundKeyUpdate0(K+28, A, B, C, D, KC14)
	{
		T0 = A + C - kc[14]
		T1 = B + kc[14] - D
		s.pdwRoundKey[28] = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		s.pdwRoundKey[29] = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = A
		A = (A >> 8) ^ (B << 24)
		B = (B >> 8) ^ (T0 << 24)
	}

	T0 = A + C - kc[15]
	T1 = B + kc[15] - D
	s.pdwRoundKey[30] = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
	s.pdwRoundKey[31] = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
}

func (s *seed128) Encrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf(msgInvalidBlockSizeSrcFormat, len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf(msgInvalidBlockSizeDstFormat, len(dst)))
	}

	var T0, T1 uint32
	L0 := byteOrder.Uint32(src[0:])
	L1 := byteOrder.Uint32(src[4:])
	R0 := byteOrder.Uint32(src[8:])
	R1 := byteOrder.Uint32(src[12:])

	// SEED_KeySched(L0, L1, R0, R1, K+0)
	{
		T0 = R0 ^ s.pdwRoundKey[0]
		T1 = R1 ^ s.pdwRoundKey[1]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		L0 ^= T0
		L1 ^= T1
	}
	// SEED_KeySched(R0, R1, L0, L1, K+2)
	{
		T0 = L0 ^ s.pdwRoundKey[2]
		T1 = L1 ^ s.pdwRoundKey[3]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		R0 ^= T0
		R1 ^= T1
	}
	// SEED_KeySched(L0, L1, R0, R1, K+4)
	{
		T0 = R0 ^ s.pdwRoundKey[4]
		T1 = R1 ^ s.pdwRoundKey[5]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		L0 ^= T0
		L1 ^= T1
	}
	// SEED_KeySched(R0, R1, L0, L1, K+6)
	{
		T0 = L0 ^ s.pdwRoundKey[6]
		T1 = L1 ^ s.pdwRoundKey[7]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		R0 ^= T0
		R1 ^= T1
	}
	// SEED_KeySched(L0, L1, R0, R1, K+8)
	{
		T0 = R0 ^ s.pdwRoundKey[8]
		T1 = R1 ^ s.pdwRoundKey[9]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		L0 ^= T0
		L1 ^= T1
	}
	// SEED_KeySched(R0, R1, L0, L1, K+10)
	{
		T0 = L0 ^ s.pdwRoundKey[10]
		T1 = L1 ^ s.pdwRoundKey[11]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		R0 ^= T0
		R1 ^= T1
	}
	// SEED_KeySched(L0, L1, R0, R1, K+12)
	{
		T0 = R0 ^ s.pdwRoundKey[12]
		T1 = R1 ^ s.pdwRoundKey[13]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		L0 ^= T0
		L1 ^= T1
	}
	// SEED_KeySched(R0, R1, L0, L1, K+14)
	{
		T0 = L0 ^ s.pdwRoundKey[14]
		T1 = L1 ^ s.pdwRoundKey[15]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		R0 ^= T0
		R1 ^= T1
	}
	// SEED_KeySched(L0, L1, R0, R1, K+16)
	{
		T0 = R0 ^ s.pdwRoundKey[16]
		T1 = R1 ^ s.pdwRoundKey[17]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		L0 ^= T0
		L1 ^= T1
	}
	// SEED_KeySched(R0, R1, L0, L1, K+18)
	{
		T0 = L0 ^ s.pdwRoundKey[18]
		T1 = L1 ^ s.pdwRoundKey[19]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		R0 ^= T0
		R1 ^= T1
	}
	// SEED_KeySched(L0, L1, R0, R1, K+20)
	{
		T0 = R0 ^ s.pdwRoundKey[20]
		T1 = R1 ^ s.pdwRoundKey[21]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		L0 ^= T0
		L1 ^= T1
	}
	// SEED_KeySched(R0, R1, L0, L1, K+22)
	{
		T0 = L0 ^ s.pdwRoundKey[22]
		T1 = L1 ^ s.pdwRoundKey[23]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		R0 ^= T0
		R1 ^= T1
	}
	// SEED_KeySched(L0, L1, R0, R1, K+24)
	{
		T0 = R0 ^ s.pdwRoundKey[24]
		T1 = R1 ^ s.pdwRoundKey[25]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		L0 ^= T0
		L1 ^= T1
	}
	// SEED_KeySched(R0, R1, L0, L1, K+26)
	{
		T0 = L0 ^ s.pdwRoundKey[26]
		T1 = L1 ^ s.pdwRoundKey[27]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		R0 ^= T0
		R1 ^= T1
	}
	// SEED_KeySched(L0, L1, R0, R1, K+28)
	{
		T0 = R0 ^ s.pdwRoundKey[28]
		T1 = R1 ^ s.pdwRoundKey[29]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		L0 ^= T0
		L1 ^= T1
	}
	// SEED_KeySched(R0, R1, L0, L1, K+30)
	{
		T0 = L0 ^ s.pdwRoundKey[30]
		T1 = L1 ^ s.pdwRoundKey[31]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		R0 ^= T0
		R1 ^= T1
	}

	byteOrder.PutUint32(dst[0:], R0)
	byteOrder.PutUint32(dst[4:], R1)
	byteOrder.PutUint32(dst[8:], L0)
	byteOrder.PutUint32(dst[12:], L1)
}

func (s *seed128) Decrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf("krypto/seed: invalid block size %d (src)", len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf("krypto/seed: invalid block size %d (dst)", len(dst)))
	}

	var T0, T1 uint32
	L0 := byteOrder.Uint32(src[0:])
	L1 := byteOrder.Uint32(src[4:])
	R0 := byteOrder.Uint32(src[8:])
	R1 := byteOrder.Uint32(src[12:])

	// SEED_KeySched(L0, L1, R0, R1, K+30)
	{
		T0 = R0 ^ s.pdwRoundKey[30]
		T1 = R1 ^ s.pdwRoundKey[31]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		L0 ^= T0
		L1 ^= T1
	}
	// SEED_KeySched(R0, R1, L0, L1, K+28)
	{
		T0 = L0 ^ s.pdwRoundKey[28]
		T1 = L1 ^ s.pdwRoundKey[29]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		R0 ^= T0
		R1 ^= T1
	}
	// SEED_KeySched(L0, L1, R0, R1, K+26)
	{
		T0 = R0 ^ s.pdwRoundKey[26]
		T1 = R1 ^ s.pdwRoundKey[27]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		L0 ^= T0
		L1 ^= T1
	}
	// SEED_KeySched(R0, R1, L0, L1, K+24)
	{
		T0 = L0 ^ s.pdwRoundKey[24]
		T1 = L1 ^ s.pdwRoundKey[25]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		R0 ^= T0
		R1 ^= T1
	}
	// SEED_KeySched(L0, L1, R0, R1, K+22)
	{
		T0 = R0 ^ s.pdwRoundKey[22]
		T1 = R1 ^ s.pdwRoundKey[23]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		L0 ^= T0
		L1 ^= T1
	}
	// SEED_KeySched(R0, R1, L0, L1, K+20)
	{
		T0 = L0 ^ s.pdwRoundKey[20]
		T1 = L1 ^ s.pdwRoundKey[21]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		R0 ^= T0
		R1 ^= T1
	}
	// SEED_KeySched(L0, L1, R0, R1, K+18)
	{
		T0 = R0 ^ s.pdwRoundKey[18]
		T1 = R1 ^ s.pdwRoundKey[19]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		L0 ^= T0
		L1 ^= T1
	}
	// SEED_KeySched(R0, R1, L0, L1, K+16)
	{
		T0 = L0 ^ s.pdwRoundKey[16]
		T1 = L1 ^ s.pdwRoundKey[17]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		R0 ^= T0
		R1 ^= T1
	}
	// SEED_KeySched(L0, L1, R0, R1, K+14)
	{
		T0 = R0 ^ s.pdwRoundKey[14]
		T1 = R1 ^ s.pdwRoundKey[15]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		L0 ^= T0
		L1 ^= T1
	}
	// SEED_KeySched(R0, R1, L0, L1, K+12)
	{
		T0 = L0 ^ s.pdwRoundKey[12]
		T1 = L1 ^ s.pdwRoundKey[13]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		R0 ^= T0
		R1 ^= T1
	}
	// SEED_KeySched(L0, L1, R0, R1, K+10)
	{
		T0 = R0 ^ s.pdwRoundKey[10]
		T1 = R1 ^ s.pdwRoundKey[11]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		L0 ^= T0
		L1 ^= T1
	}
	// SEED_KeySched(R0, R1, L0, L1, K+8)
	{
		T0 = L0 ^ s.pdwRoundKey[8]
		T1 = L1 ^ s.pdwRoundKey[9]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		R0 ^= T0
		R1 ^= T1
	}
	// SEED_KeySched(L0, L1, R0, R1, K+6)
	{
		T0 = R0 ^ s.pdwRoundKey[6]
		T1 = R1 ^ s.pdwRoundKey[7]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		L0 ^= T0
		L1 ^= T1
	}
	// SEED_KeySched(R0, R1, L0, L1, K+4)
	{
		T0 = L0 ^ s.pdwRoundKey[4]
		T1 = L1 ^ s.pdwRoundKey[5]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		R0 ^= T0
		R1 ^= T1
	}
	// SEED_KeySched(L0, L1, R0, R1, K+2)
	{
		T0 = R0 ^ s.pdwRoundKey[2]
		T1 = R1 ^ s.pdwRoundKey[3]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		L0 ^= T0
		L1 ^= T1
	}
	// SEED_KeySched(R0, R1, L0, L1, K+0)
	{
		T0 = L0 ^ s.pdwRoundKey[0]
		T1 = L1 ^ s.pdwRoundKey[1]
		T1 ^= T0
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		T1 = (T1 + T0) & 0xffffffff
		T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
		T0 = (T0 + T1) & 0xffffffff
		R0 ^= T0
		R1 ^= T1
	}

	byteOrder.PutUint32(dst[0:], R0)
	byteOrder.PutUint32(dst[4:], R1)
	byteOrder.PutUint32(dst[8:], L0)
	byteOrder.PutUint32(dst[12:], L1)
}

func getB0(aa uint32) int { return int(byte(aa >> 0)) }
func getB1(aa uint32) int { return int(byte(aa >> 8)) }
func getB2(aa uint32) int { return int(byte(aa >> 16)) }
func getB3(aa uint32) int { return int(byte(aa >> 24)) }
