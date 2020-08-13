package seed

import (
	"crypto/cipher"
	"fmt"
)

type seed128 struct {
	round       int
	pdwRoundKey [32]uint32
}

func new128(key []byte) cipher.Block {
	block := new(seed128)

	var A uint32 = (uint32(key[3]) << 24) | (uint32(key[2]) << 16) | (uint32(key[1]) << 8) | (uint32(key[0]))
	var B uint32 = (uint32(key[7]) << 24) | (uint32(key[6]) << 16) | (uint32(key[5]) << 8) | (uint32(key[4]))
	var C uint32 = (uint32(key[11]) << 24) | (uint32(key[10]) << 16) | (uint32(key[9]) << 8) | (uint32(key[8]))
	var D uint32 = (uint32(key[15]) << 24) | (uint32(key[14]) << 16) | (uint32(key[13]) << 8) | (uint32(key[12]))

	if !littleEndian {
		A = endianChange(A)
		B = endianChange(B)
		C = endianChange(C)
		D = endianChange(D)
	}

	var T0, T1 uint32
	for i := 0; i < 16; i++ {
		if (i % 2) == 0 {
			T0 = A + C - kc[i]
			T1 = B - D + kc[i]
			block.pdwRoundKey[i*2+0] = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
			block.pdwRoundKey[i*2+1] = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
			T0 = A

			A = (A >> 8) ^ (B << 24)
			B = (B >> 8) ^ (T0 << 24)
		} else {
			T0 = A + C - kc[i]
			T1 = B - D + kc[i]
			block.pdwRoundKey[i*2+0] = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
			block.pdwRoundKey[i*2+1] = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
			T0 = C

			C = (C << 8) ^ (D >> 24)
			D = (D << 8) ^ (T0 >> 24)
		}
	}

	return block
}

func (s *seed128) BlockSize() int {
	return BlockSize
}

func (s *seed128) seed_KeySched(L0, L1 *uint32, R0, R1 uint32, ki int) {
	T0 := R0 ^ s.pdwRoundKey[ki+0]
	T1 := R1 ^ s.pdwRoundKey[ki+1]
	T1 ^= T0
	T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
	T0 = (T0 + T1) & 0xffffffff
	T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
	T1 = (T1 + T0) & 0xffffffff
	T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
	T0 = (T0 + T1) & 0xffffffff

	*L0 ^= T0
	*L1 ^= T1
}

func (s *seed128) Encrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf("kipher/seed128: invalid block size %d (src)", len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf("kipher/seed128: invalid block size %d (dst)", len(dst)))
	}

	var L0 uint32 = (uint32(src[3]) << 24) | (uint32(src[2]) << 16) | (uint32(src[1]) << 8) | (uint32(src[0]))
	var L1 uint32 = (uint32(src[7]) << 24) | (uint32(src[6]) << 16) | (uint32(src[5]) << 8) | (uint32(src[4]))
	var R0 uint32 = (uint32(src[11]) << 24) | (uint32(src[10]) << 16) | (uint32(src[9]) << 8) | (uint32(src[8]))
	var R1 uint32 = (uint32(src[15]) << 24) | (uint32(src[14]) << 16) | (uint32(src[13]) << 8) | (uint32(src[12]))

	// Reorder for big endian
	// Because SEED use little endian order in default
	if littleEndian {
		L0 = endianChange(L0)
		L1 = endianChange(L1)
		R0 = endianChange(R0)
		R1 = endianChange(R1)
	}

	s.seed_KeySched(&L0, &L1, R0, R1, 0)  // Round 1
	s.seed_KeySched(&R0, &R1, L0, L1, 2)  // Round 2
	s.seed_KeySched(&L0, &L1, R0, R1, 4)  // Round 3
	s.seed_KeySched(&R0, &R1, L0, L1, 6)  // Round 4
	s.seed_KeySched(&L0, &L1, R0, R1, 8)  // Round 5
	s.seed_KeySched(&R0, &R1, L0, L1, 10) // Round 6
	s.seed_KeySched(&L0, &L1, R0, R1, 12) // Round 7
	s.seed_KeySched(&R0, &R1, L0, L1, 14) // Round 8
	s.seed_KeySched(&L0, &L1, R0, R1, 16) // Round 9
	s.seed_KeySched(&R0, &R1, L0, L1, 18) // Round 10
	s.seed_KeySched(&L0, &L1, R0, R1, 20) // Round 11
	s.seed_KeySched(&R0, &R1, L0, L1, 22) // Round 12
	s.seed_KeySched(&L0, &L1, R0, R1, 24) // Round 13
	s.seed_KeySched(&R0, &R1, L0, L1, 26) // Round 14
	s.seed_KeySched(&L0, &L1, R0, R1, 28) // Round 15
	s.seed_KeySched(&R0, &R1, L0, L1, 30) // Round 16

	if littleEndian {
		L0 = endianChange(L0)
		L1 = endianChange(L1)
		R0 = endianChange(R0)
		R1 = endianChange(R1)
	}

	// Copying output values from last round to pbData
	dst[0] = byte((R0) & 0xFF)
	dst[1] = byte((R0 >> 8) & 0xFF)
	dst[2] = byte((R0 >> 16) & 0xFF)
	dst[3] = byte((R0 >> 24) & 0xFF)

	dst[4] = byte((R1) & 0xFF)
	dst[5] = byte((R1 >> 8) & 0xFF)
	dst[6] = byte((R1 >> 16) & 0xFF)
	dst[7] = byte((R1 >> 24) & 0xFF)

	dst[8] = byte((L0) & 0xFF)
	dst[9] = byte((L0 >> 8) & 0xFF)
	dst[10] = byte((L0 >> 16) & 0xFF)
	dst[11] = byte((L0 >> 24) & 0xFF)

	dst[12] = byte((L1) & 0xFF)
	dst[13] = byte((L1 >> 8) & 0xFF)
	dst[14] = byte((L1 >> 16) & 0xFF)
	dst[15] = byte((L1 >> 24) & 0xFF)
}

func (s *seed128) Decrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf("kipher/seed128: invalid block size %d (src)", len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf("kipher/seed128: invalid block size %d (dst)", len(dst)))
	}

	var L0 uint32 = (uint32(src[3]) << 24) | (uint32(src[2]) << 16) | (uint32(src[1]) << 8) | (uint32(src[0]))
	var L1 uint32 = (uint32(src[7]) << 24) | (uint32(src[6]) << 16) | (uint32(src[5]) << 8) | (uint32(src[4]))
	var R0 uint32 = (uint32(src[11]) << 24) | (uint32(src[10]) << 16) | (uint32(src[9]) << 8) | (uint32(src[8]))
	var R1 uint32 = (uint32(src[15]) << 24) | (uint32(src[14]) << 16) | (uint32(src[13]) << 8) | (uint32(src[12]))

	// Reorder for big endian
	if littleEndian {
		L0 = endianChange(L0)
		L1 = endianChange(L1)
		R0 = endianChange(R0)
		R1 = endianChange(R1)
	}

	s.seed_KeySched(&L0, &L1, R0, R1, 30) // Round 1
	s.seed_KeySched(&R0, &R1, L0, L1, 28) // Round 2
	s.seed_KeySched(&L0, &L1, R0, R1, 26) // Round 3
	s.seed_KeySched(&R0, &R1, L0, L1, 24) // Round 4
	s.seed_KeySched(&L0, &L1, R0, R1, 22) // Round 5
	s.seed_KeySched(&R0, &R1, L0, L1, 20) // Round 6
	s.seed_KeySched(&L0, &L1, R0, R1, 18) // Round 7
	s.seed_KeySched(&R0, &R1, L0, L1, 16) // Round 8
	s.seed_KeySched(&L0, &L1, R0, R1, 14) // Round 9
	s.seed_KeySched(&R0, &R1, L0, L1, 12) // Round 10
	s.seed_KeySched(&L0, &L1, R0, R1, 10) // Round 11
	s.seed_KeySched(&R0, &R1, L0, L1, 8)  // Round 12
	s.seed_KeySched(&L0, &L1, R0, R1, 6)  // Round 13
	s.seed_KeySched(&R0, &R1, L0, L1, 4)  // Round 14
	s.seed_KeySched(&L0, &L1, R0, R1, 2)  // Round 15
	s.seed_KeySched(&R0, &R1, L0, L1, 0)  // Round 16

	if littleEndian {
		L0 = endianChange(L0)
		L1 = endianChange(L1)
		R0 = endianChange(R0)
		R1 = endianChange(R1)
	}

	dst[0] = byte((R0) & 0xFF)
	dst[1] = byte((R0 >> 8) & 0xFF)
	dst[2] = byte((R0 >> 16) & 0xFF)
	dst[3] = byte((R0 >> 24) & 0xFF)

	dst[4] = byte((R1) & 0xFF)
	dst[5] = byte((R1 >> 8) & 0xFF)
	dst[6] = byte((R1 >> 16) & 0xFF)
	dst[7] = byte((R1 >> 24) & 0xFF)

	dst[8] = byte((L0) & 0xFF)
	dst[9] = byte((L0 >> 8) & 0xFF)
	dst[10] = byte((L0 >> 16) & 0xFF)
	dst[11] = byte((L0 >> 24) & 0xFF)

	dst[12] = byte((L1) & 0xFF)
	dst[13] = byte((L1 >> 8) & 0xFF)
	dst[14] = byte((L1 >> 16) & 0xFF)
	dst[15] = byte((L1 >> 24) & 0xFF)
}
