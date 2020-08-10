package hight

/*************** Function *************************************************/

func hight_KeySched(UserKey []byte, UserKeyLen int, RoundKey []byte) {
	var i, j int32

	for i = 0; i < 4; i++ {
		RoundKey[i] = UserKey[i+12]
		RoundKey[i+4] = UserKey[i]
	}

	for i = 0; i < 8; i++ {
		for j = 0; j < 8; j++ {
			RoundKey[8+16*i+j] = byte(UserKey[(j-i)&7] + delta[16*i+j])
		}
		// Use "&7"  instead of the "%8" for Performance

		for j = 0; j < 8; j++ {
			RoundKey[8+16*i+j+8] = byte(UserKey[((j-i)&7)+8] + delta[16*i+j+8])
		}
	}
}

/*************** Encryption*************************************************/

func hight_ENC(RoundKey []byte, XX []uint32, k, i0, i1, i2, i3, i4, i5, i6, i7 int) {
	XX[i0] = (XX[i0] ^ uint32(hight_F0[XX[i1]]+RoundKey[4*k+3])) & 0xFF
	XX[i2] = (XX[i2] + uint32(hight_F1[XX[i3]]^RoundKey[4*k+2])) & 0xFF
	XX[i4] = (XX[i4] ^ uint32(hight_F0[XX[i5]]+RoundKey[4*k+1])) & 0xFF
	XX[i6] = (XX[i6] + uint32(hight_F1[XX[i7]]^RoundKey[4*k+0])) & 0xFF
}

func hight_Encrypt(RoundKey, Data []byte) {
	XX := make([]uint32, 8)

	// First Round
	XX[1] = uint32(Data[1])
	XX[3] = uint32(Data[3])
	XX[5] = uint32(Data[5])
	XX[7] = uint32(Data[7])

	XX[0] = uint32((Data[0] + RoundKey[0]) & 0xFF)
	XX[2] = uint32((Data[2] ^ RoundKey[1]))
	XX[4] = uint32((Data[4] + RoundKey[2]) & 0xFF)
	XX[6] = uint32((Data[6] ^ RoundKey[3]))

	hight_ENC(RoundKey, XX, 2, 7, 6, 5, 4, 3, 2, 1, 0)
	hight_ENC(RoundKey, XX, 3, 6, 5, 4, 3, 2, 1, 0, 7)
	hight_ENC(RoundKey, XX, 4, 5, 4, 3, 2, 1, 0, 7, 6)
	hight_ENC(RoundKey, XX, 5, 4, 3, 2, 1, 0, 7, 6, 5)
	hight_ENC(RoundKey, XX, 6, 3, 2, 1, 0, 7, 6, 5, 4)
	hight_ENC(RoundKey, XX, 7, 2, 1, 0, 7, 6, 5, 4, 3)
	hight_ENC(RoundKey, XX, 8, 1, 0, 7, 6, 5, 4, 3, 2)
	hight_ENC(RoundKey, XX, 9, 0, 7, 6, 5, 4, 3, 2, 1)
	hight_ENC(RoundKey, XX, 10, 7, 6, 5, 4, 3, 2, 1, 0)
	hight_ENC(RoundKey, XX, 11, 6, 5, 4, 3, 2, 1, 0, 7)
	hight_ENC(RoundKey, XX, 12, 5, 4, 3, 2, 1, 0, 7, 6)
	hight_ENC(RoundKey, XX, 13, 4, 3, 2, 1, 0, 7, 6, 5)
	hight_ENC(RoundKey, XX, 14, 3, 2, 1, 0, 7, 6, 5, 4)
	hight_ENC(RoundKey, XX, 15, 2, 1, 0, 7, 6, 5, 4, 3)
	hight_ENC(RoundKey, XX, 16, 1, 0, 7, 6, 5, 4, 3, 2)
	hight_ENC(RoundKey, XX, 17, 0, 7, 6, 5, 4, 3, 2, 1)
	hight_ENC(RoundKey, XX, 18, 7, 6, 5, 4, 3, 2, 1, 0)
	hight_ENC(RoundKey, XX, 19, 6, 5, 4, 3, 2, 1, 0, 7)
	hight_ENC(RoundKey, XX, 20, 5, 4, 3, 2, 1, 0, 7, 6)
	hight_ENC(RoundKey, XX, 21, 4, 3, 2, 1, 0, 7, 6, 5)
	hight_ENC(RoundKey, XX, 22, 3, 2, 1, 0, 7, 6, 5, 4)
	hight_ENC(RoundKey, XX, 23, 2, 1, 0, 7, 6, 5, 4, 3)
	hight_ENC(RoundKey, XX, 24, 1, 0, 7, 6, 5, 4, 3, 2)
	hight_ENC(RoundKey, XX, 25, 0, 7, 6, 5, 4, 3, 2, 1)
	hight_ENC(RoundKey, XX, 26, 7, 6, 5, 4, 3, 2, 1, 0)
	hight_ENC(RoundKey, XX, 27, 6, 5, 4, 3, 2, 1, 0, 7)
	hight_ENC(RoundKey, XX, 28, 5, 4, 3, 2, 1, 0, 7, 6)
	hight_ENC(RoundKey, XX, 29, 4, 3, 2, 1, 0, 7, 6, 5)
	hight_ENC(RoundKey, XX, 30, 3, 2, 1, 0, 7, 6, 5, 4)
	hight_ENC(RoundKey, XX, 31, 2, 1, 0, 7, 6, 5, 4, 3)
	hight_ENC(RoundKey, XX, 32, 1, 0, 7, 6, 5, 4, 3, 2)
	hight_ENC(RoundKey, XX, 33, 0, 7, 6, 5, 4, 3, 2, 1)

	// Final Round
	Data[1] = byte(XX[2])
	Data[3] = byte(XX[4])
	Data[5] = byte(XX[6])
	Data[7] = byte(XX[0])

	Data[0] = byte((XX[1] + uint32(RoundKey[4])))
	Data[2] = byte((XX[3] ^ uint32(RoundKey[5])))
	Data[4] = byte((XX[5] + uint32(RoundKey[6])))
	Data[6] = byte((XX[7] ^ uint32(RoundKey[7])))
}

/***************Decryption *************************************************/

// Same as encrypt, except that round keys are applied in reverse order

func hight_DEC(RoundKey []byte, XX []uint32, k, i0, i1, i2, i3, i4, i5, i6, i7 int) {
	XX[i1] = (XX[i1] - uint32(hight_F1[XX[i2]]^RoundKey[4*k+2])) & 0xFF
	XX[i3] = (XX[i3] ^ uint32(hight_F0[XX[i4]]+RoundKey[4*k+1])) & 0xFF
	XX[i5] = (XX[i5] - uint32(hight_F1[XX[i6]]^RoundKey[4*k+0])) & 0xFF
	XX[i7] = (XX[i7] ^ uint32(hight_F0[XX[i0]]+RoundKey[4*k+3])) & 0xFF
}

func hight_Decrypt(RoundKey, Data []byte) {
	XX := make([]uint32, 8)

	XX[2] = uint32(Data[1])
	XX[4] = uint32(Data[3])
	XX[6] = uint32(Data[5])
	XX[0] = uint32(Data[7])

	XX[1] = uint32((Data[0] - RoundKey[4]))
	XX[3] = uint32((Data[2] ^ RoundKey[5]))
	XX[5] = uint32((Data[4] - RoundKey[6]))
	XX[7] = uint32((Data[6] ^ RoundKey[7]))

	hight_DEC(RoundKey, XX, 33, 7, 6, 5, 4, 3, 2, 1, 0)
	hight_DEC(RoundKey, XX, 32, 0, 7, 6, 5, 4, 3, 2, 1)
	hight_DEC(RoundKey, XX, 31, 1, 0, 7, 6, 5, 4, 3, 2)
	hight_DEC(RoundKey, XX, 30, 2, 1, 0, 7, 6, 5, 4, 3)
	hight_DEC(RoundKey, XX, 29, 3, 2, 1, 0, 7, 6, 5, 4)
	hight_DEC(RoundKey, XX, 28, 4, 3, 2, 1, 0, 7, 6, 5)
	hight_DEC(RoundKey, XX, 27, 5, 4, 3, 2, 1, 0, 7, 6)
	hight_DEC(RoundKey, XX, 26, 6, 5, 4, 3, 2, 1, 0, 7)
	hight_DEC(RoundKey, XX, 25, 7, 6, 5, 4, 3, 2, 1, 0)
	hight_DEC(RoundKey, XX, 24, 0, 7, 6, 5, 4, 3, 2, 1)
	hight_DEC(RoundKey, XX, 23, 1, 0, 7, 6, 5, 4, 3, 2)
	hight_DEC(RoundKey, XX, 22, 2, 1, 0, 7, 6, 5, 4, 3)
	hight_DEC(RoundKey, XX, 21, 3, 2, 1, 0, 7, 6, 5, 4)
	hight_DEC(RoundKey, XX, 20, 4, 3, 2, 1, 0, 7, 6, 5)
	hight_DEC(RoundKey, XX, 19, 5, 4, 3, 2, 1, 0, 7, 6)
	hight_DEC(RoundKey, XX, 18, 6, 5, 4, 3, 2, 1, 0, 7)
	hight_DEC(RoundKey, XX, 17, 7, 6, 5, 4, 3, 2, 1, 0)
	hight_DEC(RoundKey, XX, 16, 0, 7, 6, 5, 4, 3, 2, 1)
	hight_DEC(RoundKey, XX, 15, 1, 0, 7, 6, 5, 4, 3, 2)
	hight_DEC(RoundKey, XX, 14, 2, 1, 0, 7, 6, 5, 4, 3)
	hight_DEC(RoundKey, XX, 13, 3, 2, 1, 0, 7, 6, 5, 4)
	hight_DEC(RoundKey, XX, 12, 4, 3, 2, 1, 0, 7, 6, 5)
	hight_DEC(RoundKey, XX, 11, 5, 4, 3, 2, 1, 0, 7, 6)
	hight_DEC(RoundKey, XX, 10, 6, 5, 4, 3, 2, 1, 0, 7)
	hight_DEC(RoundKey, XX, 9, 7, 6, 5, 4, 3, 2, 1, 0)
	hight_DEC(RoundKey, XX, 8, 0, 7, 6, 5, 4, 3, 2, 1)
	hight_DEC(RoundKey, XX, 7, 1, 0, 7, 6, 5, 4, 3, 2)
	hight_DEC(RoundKey, XX, 6, 2, 1, 0, 7, 6, 5, 4, 3)
	hight_DEC(RoundKey, XX, 5, 3, 2, 1, 0, 7, 6, 5, 4)
	hight_DEC(RoundKey, XX, 4, 4, 3, 2, 1, 0, 7, 6, 5)
	hight_DEC(RoundKey, XX, 3, 5, 4, 3, 2, 1, 0, 7, 6)
	hight_DEC(RoundKey, XX, 2, 6, 5, 4, 3, 2, 1, 0, 7)

	Data[1] = byte((XX[1]))
	Data[3] = byte((XX[3]))
	Data[5] = byte((XX[5]))
	Data[7] = byte((XX[7]))

	Data[0] = byte((XX[0] - uint32(RoundKey[0])))
	Data[2] = byte((XX[2] ^ uint32(RoundKey[1])))
	Data[4] = byte((XX[4] - uint32(RoundKey[2])))
	Data[6] = byte((XX[6] ^ uint32(RoundKey[3])))
}

/*************** END OF FILE **********************************************/
