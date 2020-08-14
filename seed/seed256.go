package seed

import (
	"crypto/cipher"
	"fmt"
)

type seed256 struct {
	round       int
	pdwRoundKey [48]uint32
}

var (
	seed256rot = [...]byte{12, 9, 9, 11, 11, 12}
)

func new256(key []byte) cipher.Block {
	block := new(seed256)

	var A uint32 = (uint32(key[3]) << 24) | (uint32(key[2]) << 16) | (uint32(key[1]) << 8) | (uint32(key[0]))
	var B uint32 = (uint32(key[7]) << 24) | (uint32(key[6]) << 16) | (uint32(key[5]) << 8) | (uint32(key[4]))
	var C uint32 = (uint32(key[11]) << 24) | (uint32(key[10]) << 16) | (uint32(key[9]) << 8) | (uint32(key[8]))
	var D uint32 = (uint32(key[15]) << 24) | (uint32(key[14]) << 16) | (uint32(key[13]) << 8) | (uint32(key[12]))
	var E uint32 = (uint32(key[19]) << 24) | (uint32(key[18]) << 16) | (uint32(key[17]) << 8) | (uint32(key[16]))
	var F uint32 = (uint32(key[23]) << 24) | (uint32(key[22]) << 16) | (uint32(key[21]) << 8) | (uint32(key[20]))
	var G uint32 = (uint32(key[27]) << 24) | (uint32(key[26]) << 16) | (uint32(key[25]) << 8) | (uint32(key[24]))
	var H uint32 = (uint32(key[31]) << 24) | (uint32(key[30]) << 16) | (uint32(key[29]) << 8) | (uint32(key[28]))

	if !littleEndian {
		A = endianChange(A)
		B = endianChange(B)
		C = endianChange(C)
		D = endianChange(D)
		E = endianChange(E)
		F = endianChange(F)
		G = endianChange(G)
		H = endianChange(H)
	}

	var T0, T1 uint32
	var rot byte

	T0 = (((A + C) ^ E) - F) ^ kc[0]
	T1 = (((B - D) ^ G) + H) ^ kc[0]
	block.pdwRoundKey[0] = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
	block.pdwRoundKey[1] = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]

	for i := 1; i < 24; i++ {
		rot = seed256rot[i%6]

		if ((i + 1) % 2) == 0 {
			T0 = D
			D = (D >> rot) ^ (C << (32 - rot))
			C = (C >> rot) ^ (B << (32 - rot))
			B = (B >> rot) ^ (A << (32 - rot))
			A = (A >> rot) ^ (T0 << (32 - rot))
		} else {
			T0 = E
			E = (E << rot) ^ (F >> (32 - rot))
			F = (F << rot) ^ (G >> (32 - rot))
			G = (G << rot) ^ (H >> (32 - rot))
			H = (H << rot) ^ (T0 >> (32 - rot))
		}

		T0 = (((A + C) ^ E) - F) ^ kc[i]
		T1 = (((B - D) ^ G) + H) ^ kc[i]

		block.pdwRoundKey[i*2+0] = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
		block.pdwRoundKey[i*2+1] = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
	}

	return block
}

func (s *seed256) BlockSize() int {
	return BlockSize
}

func (s *seed256) seed_KeySched(L0, L1 *uint32, R0, R1 uint32, ki int) {
	T0 := R0 ^ s.pdwRoundKey[ki+0]
	T1 := R1 ^ s.pdwRoundKey[ki+1]
	T1 ^= T0
	T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
	T0 += T1
	T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
	T1 += T0
	T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
	T0 += T1
	*L0 ^= T0
	*L1 ^= T1
}

func (s *seed256) Encrypt(dst, src []byte) {
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

	if !littleEndian {
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
	s.seed_KeySched(&L0, &L1, R0, R1, 32) // Round 17
	s.seed_KeySched(&R0, &R1, L0, L1, 34) // Round 18
	s.seed_KeySched(&L0, &L1, R0, R1, 36) // Round 19
	s.seed_KeySched(&R0, &R1, L0, L1, 38) // Round 20
	s.seed_KeySched(&L0, &L1, R0, R1, 40) // Round 21
	s.seed_KeySched(&R0, &R1, L0, L1, 42) // Round 22
	s.seed_KeySched(&L0, &L1, R0, R1, 44) // Round 23
	s.seed_KeySched(&R0, &R1, L0, L1, 46) // Round 24

	if !littleEndian {
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

func (s *seed256) Decrypt(dst, src []byte) {
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
	if !littleEndian {
		L0 = endianChange(L0)
		L1 = endianChange(L1)
		R0 = endianChange(R0)
		R1 = endianChange(R1)
	}

	s.seed_KeySched(&L0, &L1, R0, R1, 46) // Round 1
	s.seed_KeySched(&R0, &R1, L0, L1, 44) // Round 2
	s.seed_KeySched(&L0, &L1, R0, R1, 42) // Round 3
	s.seed_KeySched(&R0, &R1, L0, L1, 40) // Round 4
	s.seed_KeySched(&L0, &L1, R0, R1, 38) // Round 5
	s.seed_KeySched(&R0, &R1, L0, L1, 36) // Round 6
	s.seed_KeySched(&L0, &L1, R0, R1, 34) // Round 7
	s.seed_KeySched(&R0, &R1, L0, L1, 32) // Round 8
	s.seed_KeySched(&L0, &L1, R0, R1, 30) // Round 9
	s.seed_KeySched(&R0, &R1, L0, L1, 28) // Round 10
	s.seed_KeySched(&L0, &L1, R0, R1, 26) // Round 11
	s.seed_KeySched(&R0, &R1, L0, L1, 24) // Round 12
	s.seed_KeySched(&L0, &L1, R0, R1, 22) // Round 13
	s.seed_KeySched(&R0, &R1, L0, L1, 20) // Round 14
	s.seed_KeySched(&L0, &L1, R0, R1, 18) // Round 15
	s.seed_KeySched(&R0, &R1, L0, L1, 16) // Round 16
	s.seed_KeySched(&L0, &L1, R0, R1, 14) // Round 17
	s.seed_KeySched(&R0, &R1, L0, L1, 12) // Round 18
	s.seed_KeySched(&L0, &L1, R0, R1, 10) // Round 19
	s.seed_KeySched(&R0, &R1, L0, L1, 8)  // Round 20
	s.seed_KeySched(&L0, &L1, R0, R1, 6)  // Round 21
	s.seed_KeySched(&R0, &R1, L0, L1, 4)  // Round 22
	s.seed_KeySched(&L0, &L1, R0, R1, 2)  // Round 23
	s.seed_KeySched(&R0, &R1, L0, L1, 0)  // Round 24

	if !littleEndian {
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
