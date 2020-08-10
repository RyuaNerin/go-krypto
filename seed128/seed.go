package seed128

func endianChange(dwS uint32) uint32 {
	return (((dwS << 8) | (dwS >> 24)) & uint32(0x00ff00ff)) |
		(((dwS << 24) | (dwS >> 8)) & uint32(0xff00ff00))
}

/******************** Macros for Encryption and Decryption ********************/

func getB0(A uint32) uint32 { return 0xFF & (A >> 0) }
func getB1(A uint32) uint32 { return 0xFF & (A >> 8) }
func getB2(A uint32) uint32 { return 0xFF & (A >> 16) }
func getB3(A uint32) uint32 { return 0xFF & (A >> 24) }

// Round function F and adding output of F to L.
// L0, L1 : left input values at each round
// R0, R1 : right input values at each round
// K : round keys at each round
func seed_KeySched(L0, L1 *uint32, R0, R1 uint32, K []uint32, T0, T1 *uint32) {
	*T0 = R0 ^ (K)[0]
	*T1 = R1 ^ (K)[1]
	*T1 ^= *T0
	*T1 = ss0[getB0(*T1)] ^ ss1[getB1(*T1)] ^ ss2[getB2(*T1)] ^ ss3[getB3(*T1)]
	*T0 = (*T0 + *T1) & 0xffffffff
	*T0 = ss0[getB0(*T0)] ^ ss1[getB1(*T0)] ^ ss2[getB2(*T0)] ^ ss3[getB3(*T0)]
	*T1 = (*T1 + *T0) & 0xffffffff
	*T1 = ss0[getB0(*T1)] ^ ss1[getB1(*T1)] ^ ss2[getB2(*T1)] ^ ss3[getB3(*T1)]
	*T0 = (*T0 + *T1) & 0xffffffff
	*L0 ^= *T0
	*L1 ^= *T1
}

/********************************* Encryption *********************************/

func seed_Encrypt(pbData []byte, pdwRoundKey []uint32) {
	var L0, L1, R0, R1 uint32 // Iuput/output values at each rounds
	var T0, T1 uint32         // Temporary variables for round function F
	K := pdwRoundKey          // Pointer of round keys

	// Set up input values for first round
	L0 = (uint32(pbData[3]) << 24) | (uint32(pbData[2]) << 16) | (uint32(pbData[1]) << 8) | (uint32(pbData[0]))
	L1 = (uint32(pbData[7]) << 24) | (uint32(pbData[6]) << 16) | (uint32(pbData[5]) << 8) | (uint32(pbData[4]))
	R0 = (uint32(pbData[11]) << 24) | (uint32(pbData[10]) << 16) | (uint32(pbData[9]) << 8) | (uint32(pbData[8]))
	R1 = (uint32(pbData[15]) << 24) | (uint32(pbData[14]) << 16) | (uint32(pbData[13]) << 8) | (uint32(pbData[12]))

	// Reorder for big endian
	// Because SEED use little endian order in default
	if littleEndian {
		L0 = endianChange(L0)
		L1 = endianChange(L1)
		R0 = endianChange(R0)
		R1 = endianChange(R1)
	}

	seed_KeySched(&L0, &L1, R0, R1, K[0:0+2], &T0, &T1)   // Round 1
	seed_KeySched(&R0, &R1, L0, L1, K[2:2+2], &T0, &T1)   // Round 2
	seed_KeySched(&L0, &L1, R0, R1, K[4:4+2], &T0, &T1)   // Round 3
	seed_KeySched(&R0, &R1, L0, L1, K[6:6+2], &T0, &T1)   // Round 4
	seed_KeySched(&L0, &L1, R0, R1, K[8:8+2], &T0, &T1)   // Round 5
	seed_KeySched(&R0, &R1, L0, L1, K[10:10+2], &T0, &T1) // Round 6
	seed_KeySched(&L0, &L1, R0, R1, K[12:12+2], &T0, &T1) // Round 7
	seed_KeySched(&R0, &R1, L0, L1, K[14:14+2], &T0, &T1) // Round 8
	seed_KeySched(&L0, &L1, R0, R1, K[16:16+2], &T0, &T1) // Round 9
	seed_KeySched(&R0, &R1, L0, L1, K[18:18+2], &T0, &T1) // Round 10
	seed_KeySched(&L0, &L1, R0, R1, K[20:20+2], &T0, &T1) // Round 11
	seed_KeySched(&R0, &R1, L0, L1, K[22:22+2], &T0, &T1) // Round 12
	seed_KeySched(&L0, &L1, R0, R1, K[24:24+2], &T0, &T1) // Round 13
	seed_KeySched(&R0, &R1, L0, L1, K[26:26+2], &T0, &T1) // Round 14
	seed_KeySched(&L0, &L1, R0, R1, K[28:28+2], &T0, &T1) // Round 15
	seed_KeySched(&R0, &R1, L0, L1, K[30:30+2], &T0, &T1) // Round 16

	if littleEndian {
		L0 = endianChange(L0)
		L1 = endianChange(L1)
		R0 = endianChange(R0)
		R1 = endianChange(R1)
	}

	// Copying output values from last round to pbData
	pbData[0] = byte((R0) & 0xFF)
	pbData[1] = byte((R0 >> 8) & 0xFF)
	pbData[2] = byte((R0 >> 16) & 0xFF)
	pbData[3] = byte((R0 >> 24) & 0xFF)

	pbData[4] = byte((R1) & 0xFF)
	pbData[5] = byte((R1 >> 8) & 0xFF)
	pbData[6] = byte((R1 >> 16) & 0xFF)
	pbData[7] = byte((R1 >> 24) & 0xFF)

	pbData[8] = byte((L0) & 0xFF)
	pbData[9] = byte((L0 >> 8) & 0xFF)
	pbData[10] = byte((L0 >> 16) & 0xFF)
	pbData[11] = byte((L0 >> 24) & 0xFF)

	pbData[12] = byte((L1) & 0xFF)
	pbData[13] = byte((L1 >> 8) & 0xFF)
	pbData[14] = byte((L1 >> 16) & 0xFF)
	pbData[15] = byte((L1 >> 24) & 0xFF)
}

/********************************* Decryption *********************************/

// Same as encrypt, except that round keys are applied in reverse order
func seed_Decrypt(pbData []byte, pdwRoundKey []uint32) {
	var L0, L1, R0, R1 uint32 // Iuput/output values at each rounds
	var T0, T1 uint32         // Temporary variables for round function F
	K := pdwRoundKey          // Pointer of round keys

	// Set up input values for first round
	L0 = (uint32(pbData[3]) << 24) | (uint32(pbData[2]) << 16) | (uint32(pbData[1]) << 8) | (uint32(pbData[0]))
	L1 = (uint32(pbData[7]) << 24) | (uint32(pbData[6]) << 16) | (uint32(pbData[5]) << 8) | (uint32(pbData[4]))
	R0 = (uint32(pbData[11]) << 24) | (uint32(pbData[10]) << 16) | (uint32(pbData[9]) << 8) | (uint32(pbData[8]))
	R1 = (uint32(pbData[15]) << 24) | (uint32(pbData[14]) << 16) | (uint32(pbData[13]) << 8) | (uint32(pbData[12]))

	// Reorder for big endian
	if littleEndian {
		L0 = endianChange(L0)
		L1 = endianChange(L1)
		R0 = endianChange(R0)
		R1 = endianChange(R1)
	}
	//printf("%08X %08X %08X %08X\n",L0,L1,R0,R1);

	seed_KeySched(&L0, &L1, R0, R1, K[30:30+2], &T0, &T1) // Round 1
	seed_KeySched(&R0, &R1, L0, L1, K[28:28+2], &T0, &T1) // Round 2
	seed_KeySched(&L0, &L1, R0, R1, K[26:26+2], &T0, &T1) // Round 3
	seed_KeySched(&R0, &R1, L0, L1, K[24:24+2], &T0, &T1) // Round 4
	seed_KeySched(&L0, &L1, R0, R1, K[22:22+2], &T0, &T1) // Round 5
	seed_KeySched(&R0, &R1, L0, L1, K[20:20+2], &T0, &T1) // Round 6
	seed_KeySched(&L0, &L1, R0, R1, K[18:18+2], &T0, &T1) // Round 7
	seed_KeySched(&R0, &R1, L0, L1, K[16:16+2], &T0, &T1) // Round 8
	seed_KeySched(&L0, &L1, R0, R1, K[14:14+2], &T0, &T1) // Round 9
	seed_KeySched(&R0, &R1, L0, L1, K[12:12+2], &T0, &T1) // Round 10
	seed_KeySched(&L0, &L1, R0, R1, K[10:10+2], &T0, &T1) // Round 11
	seed_KeySched(&R0, &R1, L0, L1, K[8:8+2], &T0, &T1)   // Round 12
	seed_KeySched(&L0, &L1, R0, R1, K[6:6+2], &T0, &T1)   // Round 13
	seed_KeySched(&R0, &R1, L0, L1, K[4:4+2], &T0, &T1)   // Round 14
	seed_KeySched(&L0, &L1, R0, R1, K[2:2+2], &T0, &T1)   // Round 15
	seed_KeySched(&R0, &R1, L0, L1, K[0:0+2], &T0, &T1)   // Round 16

	if littleEndian {
		L0 = endianChange(L0)
		L1 = endianChange(L1)
		R0 = endianChange(R0)
		R1 = endianChange(R1)
	}

	// Copy output values from last round to pbData
	pbData[0] = byte((R0) & 0xFF)
	pbData[1] = byte((R0 >> 8) & 0xFF)
	pbData[2] = byte((R0 >> 16) & 0xFF)
	pbData[3] = byte((R0 >> 24) & 0xFF)

	pbData[4] = byte((R1) & 0xFF)
	pbData[5] = byte((R1 >> 8) & 0xFF)
	pbData[6] = byte((R1 >> 16) & 0xFF)
	pbData[7] = byte((R1 >> 24) & 0xFF)

	pbData[8] = byte((L0) & 0xFF)
	pbData[9] = byte((L0 >> 8) & 0xFF)
	pbData[10] = byte((L0 >> 16) & 0xFF)
	pbData[11] = byte((L0 >> 24) & 0xFF)

	pbData[12] = byte((L1) & 0xFF)
	pbData[13] = byte((L1 >> 8) & 0xFF)
	pbData[14] = byte((L1 >> 16) & 0xFF)
	pbData[15] = byte((L1 >> 24) & 0xFF)
}

/************************ Constants for Key schedule **************************/

/************************** Macros for Key schedule ***************************/

func roundKeyUpdate0(K []uint32, A *uint32, B *uint32, C *uint32, D *uint32, KC uint32, T0 *uint32, T1 *uint32) {
	*T0 = *A + *C - KC
	*T1 = *B + KC - *D
	(K)[0] = ss0[getB0(*T0)] ^ ss1[getB1(*T0)] ^ ss2[getB2(*T0)] ^ ss3[getB3(*T0)]
	(K)[1] = ss0[getB0(*T1)] ^ ss1[getB1(*T1)] ^ ss2[getB2(*T1)] ^ ss3[getB3(*T1)]
	*T0 = *A
	*A = (*A >> 8) ^ (*B << 24)
	*B = (*B >> 8) ^ (*T0 << 24)
}

func roundKeyUpdate1(K []uint32, A *uint32, B *uint32, C *uint32, D *uint32, KC uint32, T0 *uint32, T1 *uint32) {
	*T0 = *A + *C - KC
	*T1 = *B + KC - *D
	(K)[0] = ss0[getB0(*T0)] ^ ss1[getB1(*T0)] ^ ss2[getB2(*T0)] ^ ss3[getB3(*T0)]
	(K)[1] = ss0[getB0(*T1)] ^ ss1[getB1(*T1)] ^ ss2[getB2(*T1)] ^ ss3[getB3(*T1)]
	*T0 = *C
	*C = (*C << 8) ^ (*D >> 24)
	*D = (*D << 8) ^ (*T0 >> 24)
}

/******************************** Key Schedule ********************************/

func seed_KeySchedKey(pdwRoundKey []uint32, pbUserKey []byte) {
	var A, B, C, D uint32
	var T0, T1 uint32
	K := pdwRoundKey

	A = (uint32(pbUserKey[3]) << 24) | (uint32(pbUserKey[2]) << 16) | (uint32(pbUserKey[1]) << 8) | (uint32(pbUserKey[0]))
	B = (uint32(pbUserKey[7]) << 24) | (uint32(pbUserKey[6]) << 16) | (uint32(pbUserKey[5]) << 8) | (uint32(pbUserKey[4]))
	C = (uint32(pbUserKey[11]) << 24) | (uint32(pbUserKey[10]) << 16) | (uint32(pbUserKey[9]) << 8) | (uint32(pbUserKey[8]))
	D = (uint32(pbUserKey[15]) << 24) | (uint32(pbUserKey[14]) << 16) | (uint32(pbUserKey[13]) << 8) | (uint32(pbUserKey[12]))

	// Reorder for big endian
	if !littleEndian {
		A = endianChange(A)
		B = endianChange(B)
		C = endianChange(C)
		D = endianChange(D)
	}

	// i-th round keys( K_i,0 and K_i,1 ) are denoted as K[2*(i-1)] and K[2*i-1], respectively
	roundKeyUpdate0(K[0:2], &A, &B, &C, &D, kc0, &T0, &T1)      // K_1,0 and K_1,1
	roundKeyUpdate1(K[2:2+2], &A, &B, &C, &D, kc1, &T0, &T1)    // K_2,0 and K_2,1
	roundKeyUpdate0(K[4:4+2], &A, &B, &C, &D, kc2, &T0, &T1)    // K_3,0 and K_3,1
	roundKeyUpdate1(K[6:6+2], &A, &B, &C, &D, kc3, &T0, &T1)    // K_4,0 and K_4,1
	roundKeyUpdate0(K[8:8+2], &A, &B, &C, &D, kc4, &T0, &T1)    // K_5,0 and K_5,1
	roundKeyUpdate1(K[10:10+2], &A, &B, &C, &D, kc5, &T0, &T1)  // K_6,0 and K_6,1
	roundKeyUpdate0(K[12:12+2], &A, &B, &C, &D, kc6, &T0, &T1)  // K_7,0 and K_7,1
	roundKeyUpdate1(K[14:14+2], &A, &B, &C, &D, kc7, &T0, &T1)  // K_8,0 and K_8,1
	roundKeyUpdate0(K[16:16+2], &A, &B, &C, &D, kc8, &T0, &T1)  // K_9,0 and K_9,1
	roundKeyUpdate1(K[18:18+2], &A, &B, &C, &D, kc9, &T0, &T1)  // K_10,0 and K_10,1
	roundKeyUpdate0(K[20:20+2], &A, &B, &C, &D, kc10, &T0, &T1) // K_11,0 and K_11,1
	roundKeyUpdate1(K[22:22+2], &A, &B, &C, &D, kc11, &T0, &T1) // K_12,0 and K_12,1
	roundKeyUpdate0(K[24:24+2], &A, &B, &C, &D, kc12, &T0, &T1) // K_13,0 and K_13,1
	roundKeyUpdate1(K[26:26+2], &A, &B, &C, &D, kc13, &T0, &T1) // K_14,0 and K_14,1
	roundKeyUpdate0(K[28:28+2], &A, &B, &C, &D, kc14, &T0, &T1) // K_15,0 and K_15,1

	T0 = A + C - kc15
	T1 = B - D + kc15
	K[30] = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ // K_16,0
		ss2[getB2(T0)] ^ ss3[getB3(T0)]
	K[31] = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ // K_16,1
		ss2[getB2(T1)] ^ ss3[getB3(T1)]
}

/*********************************** END **************************************/
