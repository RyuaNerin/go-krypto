// Package hight implements HIGHT encryption, as defined in TTAS.KO-12.0040/R1
package hight

import (
	"crypto/cipher"
	"fmt"
)

const (
	// The HIGHT block size in bytes.
	BlockSize = 8
	// The HIGHT key size in bytes.
	KeySize = 16
)

const pdwRoundKeySize = 136

type hight struct {
	pdwRoundKey [pdwRoundKeySize]byte
}

// NewCipher creates and returns a new cipher.Block.
func NewCipher(key []byte) (cipher.Block, error) {
	l := len(key)
	if l != KeySize {
		return nil, KeySizeError(l)
	}

	block := new(hight)

	for i := 0; i < 4; i++ {
		block.pdwRoundKey[i+0] = key[i+12]
		block.pdwRoundKey[i+4] = key[i+0]
	}

	for i := 0; i < 8; i++ {
		for k := 0; k < 8; k++ {
			block.pdwRoundKey[8+16*i+k+0] = key[((k-i)&7)+0] + delta[16*i+k+0]
		}
		for k := 0; k < 8; k++ {
			block.pdwRoundKey[8+16*i+k+8] = key[((k-i)&7)+8] + delta[16*i+k+8]
		}
	}

	return block, nil
}

func (s *hight) BlockSize() int {
	return BlockSize
}

func (s *hight) Encrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf(msgInvalidBlockSizeSrcFormat, len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf(msgInvalidBlockSizeDstFormat, len(dst)))
	}

	XX := []byte{
		src[0] + s.pdwRoundKey[0],
		src[1],
		src[2] ^ s.pdwRoundKey[1],
		src[3],
		src[4] + s.pdwRoundKey[2],
		src[5],
		src[6] ^ s.pdwRoundKey[3],
		src[7],
	}

	// HIGHT_ENC(2,  7,6,5,4,3,2,1,0)
	{
		XX[7] = (XX[7] ^ (hight_F0[XX[6]] + s.pdwRoundKey[4*2+3]))
		XX[5] = (XX[5] + (hight_F1[XX[4]] ^ s.pdwRoundKey[4*2+2]))
		XX[3] = (XX[3] ^ (hight_F0[XX[2]] + s.pdwRoundKey[4*2+1]))
		XX[1] = (XX[1] + (hight_F1[XX[0]] ^ s.pdwRoundKey[4*2+0]))
	}
	// HIGHT_ENC(3,  6,5,4,3,2,1,0,7)
	{
		XX[6] = (XX[6] ^ (hight_F0[XX[5]] + s.pdwRoundKey[4*3+3]))
		XX[4] = (XX[4] + (hight_F1[XX[3]] ^ s.pdwRoundKey[4*3+2]))
		XX[2] = (XX[2] ^ (hight_F0[XX[1]] + s.pdwRoundKey[4*3+1]))
		XX[0] = (XX[0] + (hight_F1[XX[7]] ^ s.pdwRoundKey[4*3+0]))
	}
	// HIGHT_ENC(4,  5,4,3,2,1,0,7,6)
	{
		XX[5] = (XX[5] ^ (hight_F0[XX[4]] + s.pdwRoundKey[4*4+3]))
		XX[3] = (XX[3] + (hight_F1[XX[2]] ^ s.pdwRoundKey[4*4+2]))
		XX[1] = (XX[1] ^ (hight_F0[XX[0]] + s.pdwRoundKey[4*4+1]))
		XX[7] = (XX[7] + (hight_F1[XX[6]] ^ s.pdwRoundKey[4*4+0]))
	}
	// HIGHT_ENC(5,  4,3,2,1,0,7,6,5)
	{
		XX[4] = (XX[4] ^ (hight_F0[XX[3]] + s.pdwRoundKey[4*5+3]))
		XX[2] = (XX[2] + (hight_F1[XX[1]] ^ s.pdwRoundKey[4*5+2]))
		XX[0] = (XX[0] ^ (hight_F0[XX[7]] + s.pdwRoundKey[4*5+1]))
		XX[6] = (XX[6] + (hight_F1[XX[5]] ^ s.pdwRoundKey[4*5+0]))
	}
	// HIGHT_ENC(6,  3,2,1,0,7,6,5,4)
	{
		XX[3] = (XX[3] ^ (hight_F0[XX[2]] + s.pdwRoundKey[4*6+3]))
		XX[1] = (XX[1] + (hight_F1[XX[0]] ^ s.pdwRoundKey[4*6+2]))
		XX[7] = (XX[7] ^ (hight_F0[XX[6]] + s.pdwRoundKey[4*6+1]))
		XX[5] = (XX[5] + (hight_F1[XX[4]] ^ s.pdwRoundKey[4*6+0]))
	}
	// HIGHT_ENC(7,  2,1,0,7,6,5,4,3)
	{
		XX[2] = (XX[2] ^ (hight_F0[XX[1]] + s.pdwRoundKey[4*7+3]))
		XX[0] = (XX[0] + (hight_F1[XX[7]] ^ s.pdwRoundKey[4*7+2]))
		XX[6] = (XX[6] ^ (hight_F0[XX[5]] + s.pdwRoundKey[4*7+1]))
		XX[4] = (XX[4] + (hight_F1[XX[3]] ^ s.pdwRoundKey[4*7+0]))
	}
	// HIGHT_ENC(8,  1,0,7,6,5,4,3,2)
	{
		XX[1] = (XX[1] ^ (hight_F0[XX[0]] + s.pdwRoundKey[4*8+3]))
		XX[7] = (XX[7] + (hight_F1[XX[6]] ^ s.pdwRoundKey[4*8+2]))
		XX[5] = (XX[5] ^ (hight_F0[XX[4]] + s.pdwRoundKey[4*8+1]))
		XX[3] = (XX[3] + (hight_F1[XX[2]] ^ s.pdwRoundKey[4*8+0]))
	}
	// HIGHT_ENC(9,  0,7,6,5,4,3,2,1)
	{
		XX[0] = (XX[0] ^ (hight_F0[XX[7]] + s.pdwRoundKey[4*9+3]))
		XX[6] = (XX[6] + (hight_F1[XX[5]] ^ s.pdwRoundKey[4*9+2]))
		XX[4] = (XX[4] ^ (hight_F0[XX[3]] + s.pdwRoundKey[4*9+1]))
		XX[2] = (XX[2] + (hight_F1[XX[1]] ^ s.pdwRoundKey[4*9+0]))
	}
	// HIGHT_ENC(10,  7,6,5,4,3,2,1,0)
	{
		XX[7] = (XX[7] ^ (hight_F0[XX[6]] + s.pdwRoundKey[4*10+3]))
		XX[5] = (XX[5] + (hight_F1[XX[4]] ^ s.pdwRoundKey[4*10+2]))
		XX[3] = (XX[3] ^ (hight_F0[XX[2]] + s.pdwRoundKey[4*10+1]))
		XX[1] = (XX[1] + (hight_F1[XX[0]] ^ s.pdwRoundKey[4*10+0]))
	}
	// HIGHT_ENC(11,  6,5,4,3,2,1,0,7)
	{
		XX[6] = (XX[6] ^ (hight_F0[XX[5]] + s.pdwRoundKey[4*11+3]))
		XX[4] = (XX[4] + (hight_F1[XX[3]] ^ s.pdwRoundKey[4*11+2]))
		XX[2] = (XX[2] ^ (hight_F0[XX[1]] + s.pdwRoundKey[4*11+1]))
		XX[0] = (XX[0] + (hight_F1[XX[7]] ^ s.pdwRoundKey[4*11+0]))
	}
	// HIGHT_ENC(12,  5,4,3,2,1,0,7,6)
	{
		XX[5] = (XX[5] ^ (hight_F0[XX[4]] + s.pdwRoundKey[4*12+3]))
		XX[3] = (XX[3] + (hight_F1[XX[2]] ^ s.pdwRoundKey[4*12+2]))
		XX[1] = (XX[1] ^ (hight_F0[XX[0]] + s.pdwRoundKey[4*12+1]))
		XX[7] = (XX[7] + (hight_F1[XX[6]] ^ s.pdwRoundKey[4*12+0]))
	}
	// HIGHT_ENC(13,  4,3,2,1,0,7,6,5)
	{
		XX[4] = (XX[4] ^ (hight_F0[XX[3]] + s.pdwRoundKey[4*13+3]))
		XX[2] = (XX[2] + (hight_F1[XX[1]] ^ s.pdwRoundKey[4*13+2]))
		XX[0] = (XX[0] ^ (hight_F0[XX[7]] + s.pdwRoundKey[4*13+1]))
		XX[6] = (XX[6] + (hight_F1[XX[5]] ^ s.pdwRoundKey[4*13+0]))
	}
	// HIGHT_ENC(14,  3,2,1,0,7,6,5,4)
	{
		XX[3] = (XX[3] ^ (hight_F0[XX[2]] + s.pdwRoundKey[4*14+3]))
		XX[1] = (XX[1] + (hight_F1[XX[0]] ^ s.pdwRoundKey[4*14+2]))
		XX[7] = (XX[7] ^ (hight_F0[XX[6]] + s.pdwRoundKey[4*14+1]))
		XX[5] = (XX[5] + (hight_F1[XX[4]] ^ s.pdwRoundKey[4*14+0]))
	}
	// HIGHT_ENC(15,  2,1,0,7,6,5,4,3)
	{
		XX[2] = (XX[2] ^ (hight_F0[XX[1]] + s.pdwRoundKey[4*15+3]))
		XX[0] = (XX[0] + (hight_F1[XX[7]] ^ s.pdwRoundKey[4*15+2]))
		XX[6] = (XX[6] ^ (hight_F0[XX[5]] + s.pdwRoundKey[4*15+1]))
		XX[4] = (XX[4] + (hight_F1[XX[3]] ^ s.pdwRoundKey[4*15+0]))
	}
	// HIGHT_ENC(16,  1,0,7,6,5,4,3,2)
	{
		XX[1] = (XX[1] ^ (hight_F0[XX[0]] + s.pdwRoundKey[4*16+3]))
		XX[7] = (XX[7] + (hight_F1[XX[6]] ^ s.pdwRoundKey[4*16+2]))
		XX[5] = (XX[5] ^ (hight_F0[XX[4]] + s.pdwRoundKey[4*16+1]))
		XX[3] = (XX[3] + (hight_F1[XX[2]] ^ s.pdwRoundKey[4*16+0]))
	}
	// HIGHT_ENC(17,  0,7,6,5,4,3,2,1)
	{
		XX[0] = (XX[0] ^ (hight_F0[XX[7]] + s.pdwRoundKey[4*17+3]))
		XX[6] = (XX[6] + (hight_F1[XX[5]] ^ s.pdwRoundKey[4*17+2]))
		XX[4] = (XX[4] ^ (hight_F0[XX[3]] + s.pdwRoundKey[4*17+1]))
		XX[2] = (XX[2] + (hight_F1[XX[1]] ^ s.pdwRoundKey[4*17+0]))
	}
	// HIGHT_ENC(18,  7,6,5,4,3,2,1,0)
	{
		XX[7] = (XX[7] ^ (hight_F0[XX[6]] + s.pdwRoundKey[4*18+3]))
		XX[5] = (XX[5] + (hight_F1[XX[4]] ^ s.pdwRoundKey[4*18+2]))
		XX[3] = (XX[3] ^ (hight_F0[XX[2]] + s.pdwRoundKey[4*18+1]))
		XX[1] = (XX[1] + (hight_F1[XX[0]] ^ s.pdwRoundKey[4*18+0]))
	}
	// HIGHT_ENC(19,  6,5,4,3,2,1,0,7)
	{
		XX[6] = (XX[6] ^ (hight_F0[XX[5]] + s.pdwRoundKey[4*19+3]))
		XX[4] = (XX[4] + (hight_F1[XX[3]] ^ s.pdwRoundKey[4*19+2]))
		XX[2] = (XX[2] ^ (hight_F0[XX[1]] + s.pdwRoundKey[4*19+1]))
		XX[0] = (XX[0] + (hight_F1[XX[7]] ^ s.pdwRoundKey[4*19+0]))
	}
	// HIGHT_ENC(20,  5,4,3,2,1,0,7,6)
	{
		XX[5] = (XX[5] ^ (hight_F0[XX[4]] + s.pdwRoundKey[4*20+3]))
		XX[3] = (XX[3] + (hight_F1[XX[2]] ^ s.pdwRoundKey[4*20+2]))
		XX[1] = (XX[1] ^ (hight_F0[XX[0]] + s.pdwRoundKey[4*20+1]))
		XX[7] = (XX[7] + (hight_F1[XX[6]] ^ s.pdwRoundKey[4*20+0]))
	}
	// HIGHT_ENC(21,  4,3,2,1,0,7,6,5)
	{
		XX[4] = (XX[4] ^ (hight_F0[XX[3]] + s.pdwRoundKey[4*21+3]))
		XX[2] = (XX[2] + (hight_F1[XX[1]] ^ s.pdwRoundKey[4*21+2]))
		XX[0] = (XX[0] ^ (hight_F0[XX[7]] + s.pdwRoundKey[4*21+1]))
		XX[6] = (XX[6] + (hight_F1[XX[5]] ^ s.pdwRoundKey[4*21+0]))
	}
	// HIGHT_ENC(22,  3,2,1,0,7,6,5,4)
	{
		XX[3] = (XX[3] ^ (hight_F0[XX[2]] + s.pdwRoundKey[4*22+3]))
		XX[1] = (XX[1] + (hight_F1[XX[0]] ^ s.pdwRoundKey[4*22+2]))
		XX[7] = (XX[7] ^ (hight_F0[XX[6]] + s.pdwRoundKey[4*22+1]))
		XX[5] = (XX[5] + (hight_F1[XX[4]] ^ s.pdwRoundKey[4*22+0]))
	}
	// HIGHT_ENC(23,  2,1,0,7,6,5,4,3)
	{
		XX[2] = (XX[2] ^ (hight_F0[XX[1]] + s.pdwRoundKey[4*23+3]))
		XX[0] = (XX[0] + (hight_F1[XX[7]] ^ s.pdwRoundKey[4*23+2]))
		XX[6] = (XX[6] ^ (hight_F0[XX[5]] + s.pdwRoundKey[4*23+1]))
		XX[4] = (XX[4] + (hight_F1[XX[3]] ^ s.pdwRoundKey[4*23+0]))
	}
	// HIGHT_ENC(24,  1,0,7,6,5,4,3,2)
	{
		XX[1] = (XX[1] ^ (hight_F0[XX[0]] + s.pdwRoundKey[4*24+3]))
		XX[7] = (XX[7] + (hight_F1[XX[6]] ^ s.pdwRoundKey[4*24+2]))
		XX[5] = (XX[5] ^ (hight_F0[XX[4]] + s.pdwRoundKey[4*24+1]))
		XX[3] = (XX[3] + (hight_F1[XX[2]] ^ s.pdwRoundKey[4*24+0]))
	}
	// HIGHT_ENC(25,  0,7,6,5,4,3,2,1)
	{
		XX[0] = (XX[0] ^ (hight_F0[XX[7]] + s.pdwRoundKey[4*25+3]))
		XX[6] = (XX[6] + (hight_F1[XX[5]] ^ s.pdwRoundKey[4*25+2]))
		XX[4] = (XX[4] ^ (hight_F0[XX[3]] + s.pdwRoundKey[4*25+1]))
		XX[2] = (XX[2] + (hight_F1[XX[1]] ^ s.pdwRoundKey[4*25+0]))
	}
	// HIGHT_ENC(26,  7,6,5,4,3,2,1,0)
	{
		XX[7] = (XX[7] ^ (hight_F0[XX[6]] + s.pdwRoundKey[4*26+3]))
		XX[5] = (XX[5] + (hight_F1[XX[4]] ^ s.pdwRoundKey[4*26+2]))
		XX[3] = (XX[3] ^ (hight_F0[XX[2]] + s.pdwRoundKey[4*26+1]))
		XX[1] = (XX[1] + (hight_F1[XX[0]] ^ s.pdwRoundKey[4*26+0]))
	}
	// HIGHT_ENC(27,  6,5,4,3,2,1,0,7)
	{
		XX[6] = (XX[6] ^ (hight_F0[XX[5]] + s.pdwRoundKey[4*27+3]))
		XX[4] = (XX[4] + (hight_F1[XX[3]] ^ s.pdwRoundKey[4*27+2]))
		XX[2] = (XX[2] ^ (hight_F0[XX[1]] + s.pdwRoundKey[4*27+1]))
		XX[0] = (XX[0] + (hight_F1[XX[7]] ^ s.pdwRoundKey[4*27+0]))
	}
	// HIGHT_ENC(28,  5,4,3,2,1,0,7,6)
	{
		XX[5] = (XX[5] ^ (hight_F0[XX[4]] + s.pdwRoundKey[4*28+3]))
		XX[3] = (XX[3] + (hight_F1[XX[2]] ^ s.pdwRoundKey[4*28+2]))
		XX[1] = (XX[1] ^ (hight_F0[XX[0]] + s.pdwRoundKey[4*28+1]))
		XX[7] = (XX[7] + (hight_F1[XX[6]] ^ s.pdwRoundKey[4*28+0]))
	}
	// HIGHT_ENC(29,  4,3,2,1,0,7,6,5)
	{
		XX[4] = (XX[4] ^ (hight_F0[XX[3]] + s.pdwRoundKey[4*29+3]))
		XX[2] = (XX[2] + (hight_F1[XX[1]] ^ s.pdwRoundKey[4*29+2]))
		XX[0] = (XX[0] ^ (hight_F0[XX[7]] + s.pdwRoundKey[4*29+1]))
		XX[6] = (XX[6] + (hight_F1[XX[5]] ^ s.pdwRoundKey[4*29+0]))
	}
	// HIGHT_ENC(30,  3,2,1,0,7,6,5,4)
	{
		XX[3] = (XX[3] ^ (hight_F0[XX[2]] + s.pdwRoundKey[4*30+3]))
		XX[1] = (XX[1] + (hight_F1[XX[0]] ^ s.pdwRoundKey[4*30+2]))
		XX[7] = (XX[7] ^ (hight_F0[XX[6]] + s.pdwRoundKey[4*30+1]))
		XX[5] = (XX[5] + (hight_F1[XX[4]] ^ s.pdwRoundKey[4*30+0]))
	}
	// HIGHT_ENC(31,  2,1,0,7,6,5,4,3)
	{
		XX[2] = (XX[2] ^ (hight_F0[XX[1]] + s.pdwRoundKey[4*31+3]))
		XX[0] = (XX[0] + (hight_F1[XX[7]] ^ s.pdwRoundKey[4*31+2]))
		XX[6] = (XX[6] ^ (hight_F0[XX[5]] + s.pdwRoundKey[4*31+1]))
		XX[4] = (XX[4] + (hight_F1[XX[3]] ^ s.pdwRoundKey[4*31+0]))
	}
	// HIGHT_ENC(32,  1,0,7,6,5,4,3,2)
	{
		XX[1] = (XX[1] ^ (hight_F0[XX[0]] + s.pdwRoundKey[4*32+3]))
		XX[7] = (XX[7] + (hight_F1[XX[6]] ^ s.pdwRoundKey[4*32+2]))
		XX[5] = (XX[5] ^ (hight_F0[XX[4]] + s.pdwRoundKey[4*32+1]))
		XX[3] = (XX[3] + (hight_F1[XX[2]] ^ s.pdwRoundKey[4*32+0]))
	}
	// HIGHT_ENC(33,  0,7,6,5,4,3,2,1)
	{
		XX[0] = (XX[0] ^ (hight_F0[XX[7]] + s.pdwRoundKey[4*33+3]))
		XX[6] = (XX[6] + (hight_F1[XX[5]] ^ s.pdwRoundKey[4*33+2]))
		XX[4] = (XX[4] ^ (hight_F0[XX[3]] + s.pdwRoundKey[4*33+1]))
		XX[2] = (XX[2] + (hight_F1[XX[1]] ^ s.pdwRoundKey[4*33+0]))
	}

	// Final Round
	dst[0] = XX[1] + s.pdwRoundKey[4]
	dst[1] = XX[2]
	dst[2] = XX[3] ^ s.pdwRoundKey[5]
	dst[3] = XX[4]
	dst[4] = XX[5] + s.pdwRoundKey[6]
	dst[5] = XX[6]
	dst[6] = XX[7] ^ s.pdwRoundKey[7]
	dst[7] = XX[0]
}

func (s *hight) Decrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf(msgInvalidBlockSizeSrcFormat, len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf(msgInvalidBlockSizeDstFormat, len(dst)))
	}

	XX := []byte{
		src[7],
		src[0] - s.pdwRoundKey[4],
		src[1],
		src[2] ^ s.pdwRoundKey[5],
		src[3],
		src[4] - s.pdwRoundKey[6],
		src[5],
		src[6] ^ s.pdwRoundKey[7],
	}

	// HIGHT_DEC(33,  7,6,5,4,3,2,1,0)
	{
		XX[6] = (XX[6] - (hight_F1[XX[5]] ^ s.pdwRoundKey[4*33+2]))
		XX[4] = (XX[4] ^ (hight_F0[XX[3]] + s.pdwRoundKey[4*33+1]))
		XX[2] = (XX[2] - (hight_F1[XX[1]] ^ s.pdwRoundKey[4*33+0]))
		XX[0] = (XX[0] ^ (hight_F0[XX[7]] + s.pdwRoundKey[4*33+3]))
	}
	// HIGHT_DEC(32,  0,7,6,5,4,3,2,1)
	{
		XX[7] = (XX[7] - (hight_F1[XX[6]] ^ s.pdwRoundKey[4*32+2]))
		XX[5] = (XX[5] ^ (hight_F0[XX[4]] + s.pdwRoundKey[4*32+1]))
		XX[3] = (XX[3] - (hight_F1[XX[2]] ^ s.pdwRoundKey[4*32+0]))
		XX[1] = (XX[1] ^ (hight_F0[XX[0]] + s.pdwRoundKey[4*32+3]))
	}
	// HIGHT_DEC(31,  1,0,7,6,5,4,3,2)
	{
		XX[0] = (XX[0] - (hight_F1[XX[7]] ^ s.pdwRoundKey[4*31+2]))
		XX[6] = (XX[6] ^ (hight_F0[XX[5]] + s.pdwRoundKey[4*31+1]))
		XX[4] = (XX[4] - (hight_F1[XX[3]] ^ s.pdwRoundKey[4*31+0]))
		XX[2] = (XX[2] ^ (hight_F0[XX[1]] + s.pdwRoundKey[4*31+3]))
	}
	// HIGHT_DEC(30,  2,1,0,7,6,5,4,3)
	{
		XX[1] = (XX[1] - (hight_F1[XX[0]] ^ s.pdwRoundKey[4*30+2]))
		XX[7] = (XX[7] ^ (hight_F0[XX[6]] + s.pdwRoundKey[4*30+1]))
		XX[5] = (XX[5] - (hight_F1[XX[4]] ^ s.pdwRoundKey[4*30+0]))
		XX[3] = (XX[3] ^ (hight_F0[XX[2]] + s.pdwRoundKey[4*30+3]))
	}
	// HIGHT_DEC(29,  3,2,1,0,7,6,5,4)
	{
		XX[2] = (XX[2] - (hight_F1[XX[1]] ^ s.pdwRoundKey[4*29+2]))
		XX[0] = (XX[0] ^ (hight_F0[XX[7]] + s.pdwRoundKey[4*29+1]))
		XX[6] = (XX[6] - (hight_F1[XX[5]] ^ s.pdwRoundKey[4*29+0]))
		XX[4] = (XX[4] ^ (hight_F0[XX[3]] + s.pdwRoundKey[4*29+3]))
	}
	// HIGHT_DEC(28,  4,3,2,1,0,7,6,5)
	{
		XX[3] = (XX[3] - (hight_F1[XX[2]] ^ s.pdwRoundKey[4*28+2]))
		XX[1] = (XX[1] ^ (hight_F0[XX[0]] + s.pdwRoundKey[4*28+1]))
		XX[7] = (XX[7] - (hight_F1[XX[6]] ^ s.pdwRoundKey[4*28+0]))
		XX[5] = (XX[5] ^ (hight_F0[XX[4]] + s.pdwRoundKey[4*28+3]))
	}
	// HIGHT_DEC(27,  5,4,3,2,1,0,7,6)
	{
		XX[4] = (XX[4] - (hight_F1[XX[3]] ^ s.pdwRoundKey[4*27+2]))
		XX[2] = (XX[2] ^ (hight_F0[XX[1]] + s.pdwRoundKey[4*27+1]))
		XX[0] = (XX[0] - (hight_F1[XX[7]] ^ s.pdwRoundKey[4*27+0]))
		XX[6] = (XX[6] ^ (hight_F0[XX[5]] + s.pdwRoundKey[4*27+3]))
	}
	// HIGHT_DEC(26,  6,5,4,3,2,1,0,7)
	{
		XX[5] = (XX[5] - (hight_F1[XX[4]] ^ s.pdwRoundKey[4*26+2]))
		XX[3] = (XX[3] ^ (hight_F0[XX[2]] + s.pdwRoundKey[4*26+1]))
		XX[1] = (XX[1] - (hight_F1[XX[0]] ^ s.pdwRoundKey[4*26+0]))
		XX[7] = (XX[7] ^ (hight_F0[XX[6]] + s.pdwRoundKey[4*26+3]))
	}
	// HIGHT_DEC(25,  7,6,5,4,3,2,1,0)
	{
		XX[6] = (XX[6] - (hight_F1[XX[5]] ^ s.pdwRoundKey[4*25+2]))
		XX[4] = (XX[4] ^ (hight_F0[XX[3]] + s.pdwRoundKey[4*25+1]))
		XX[2] = (XX[2] - (hight_F1[XX[1]] ^ s.pdwRoundKey[4*25+0]))
		XX[0] = (XX[0] ^ (hight_F0[XX[7]] + s.pdwRoundKey[4*25+3]))
	}
	// HIGHT_DEC(24,  0,7,6,5,4,3,2,1)
	{
		XX[7] = (XX[7] - (hight_F1[XX[6]] ^ s.pdwRoundKey[4*24+2]))
		XX[5] = (XX[5] ^ (hight_F0[XX[4]] + s.pdwRoundKey[4*24+1]))
		XX[3] = (XX[3] - (hight_F1[XX[2]] ^ s.pdwRoundKey[4*24+0]))
		XX[1] = (XX[1] ^ (hight_F0[XX[0]] + s.pdwRoundKey[4*24+3]))
	}
	// HIGHT_DEC(23,  1,0,7,6,5,4,3,2)
	{
		XX[0] = (XX[0] - (hight_F1[XX[7]] ^ s.pdwRoundKey[4*23+2]))
		XX[6] = (XX[6] ^ (hight_F0[XX[5]] + s.pdwRoundKey[4*23+1]))
		XX[4] = (XX[4] - (hight_F1[XX[3]] ^ s.pdwRoundKey[4*23+0]))
		XX[2] = (XX[2] ^ (hight_F0[XX[1]] + s.pdwRoundKey[4*23+3]))
	}
	// HIGHT_DEC(22,  2,1,0,7,6,5,4,3)
	{
		XX[1] = (XX[1] - (hight_F1[XX[0]] ^ s.pdwRoundKey[4*22+2]))
		XX[7] = (XX[7] ^ (hight_F0[XX[6]] + s.pdwRoundKey[4*22+1]))
		XX[5] = (XX[5] - (hight_F1[XX[4]] ^ s.pdwRoundKey[4*22+0]))
		XX[3] = (XX[3] ^ (hight_F0[XX[2]] + s.pdwRoundKey[4*22+3]))
	}
	// HIGHT_DEC(21,  3,2,1,0,7,6,5,4)
	{
		XX[2] = (XX[2] - (hight_F1[XX[1]] ^ s.pdwRoundKey[4*21+2]))
		XX[0] = (XX[0] ^ (hight_F0[XX[7]] + s.pdwRoundKey[4*21+1]))
		XX[6] = (XX[6] - (hight_F1[XX[5]] ^ s.pdwRoundKey[4*21+0]))
		XX[4] = (XX[4] ^ (hight_F0[XX[3]] + s.pdwRoundKey[4*21+3]))
	}
	// HIGHT_DEC(20,  4,3,2,1,0,7,6,5)
	{
		XX[3] = (XX[3] - (hight_F1[XX[2]] ^ s.pdwRoundKey[4*20+2]))
		XX[1] = (XX[1] ^ (hight_F0[XX[0]] + s.pdwRoundKey[4*20+1]))
		XX[7] = (XX[7] - (hight_F1[XX[6]] ^ s.pdwRoundKey[4*20+0]))
		XX[5] = (XX[5] ^ (hight_F0[XX[4]] + s.pdwRoundKey[4*20+3]))
	}
	// HIGHT_DEC(19,  5,4,3,2,1,0,7,6)
	{
		XX[4] = (XX[4] - (hight_F1[XX[3]] ^ s.pdwRoundKey[4*19+2]))
		XX[2] = (XX[2] ^ (hight_F0[XX[1]] + s.pdwRoundKey[4*19+1]))
		XX[0] = (XX[0] - (hight_F1[XX[7]] ^ s.pdwRoundKey[4*19+0]))
		XX[6] = (XX[6] ^ (hight_F0[XX[5]] + s.pdwRoundKey[4*19+3]))
	}
	// HIGHT_DEC(18,  6,5,4,3,2,1,0,7)
	{
		XX[5] = (XX[5] - (hight_F1[XX[4]] ^ s.pdwRoundKey[4*18+2]))
		XX[3] = (XX[3] ^ (hight_F0[XX[2]] + s.pdwRoundKey[4*18+1]))
		XX[1] = (XX[1] - (hight_F1[XX[0]] ^ s.pdwRoundKey[4*18+0]))
		XX[7] = (XX[7] ^ (hight_F0[XX[6]] + s.pdwRoundKey[4*18+3]))
	}
	// HIGHT_DEC(17,  7,6,5,4,3,2,1,0)
	{
		XX[6] = (XX[6] - (hight_F1[XX[5]] ^ s.pdwRoundKey[4*17+2]))
		XX[4] = (XX[4] ^ (hight_F0[XX[3]] + s.pdwRoundKey[4*17+1]))
		XX[2] = (XX[2] - (hight_F1[XX[1]] ^ s.pdwRoundKey[4*17+0]))
		XX[0] = (XX[0] ^ (hight_F0[XX[7]] + s.pdwRoundKey[4*17+3]))
	}
	// HIGHT_DEC(16,  0,7,6,5,4,3,2,1)
	{
		XX[7] = (XX[7] - (hight_F1[XX[6]] ^ s.pdwRoundKey[4*16+2]))
		XX[5] = (XX[5] ^ (hight_F0[XX[4]] + s.pdwRoundKey[4*16+1]))
		XX[3] = (XX[3] - (hight_F1[XX[2]] ^ s.pdwRoundKey[4*16+0]))
		XX[1] = (XX[1] ^ (hight_F0[XX[0]] + s.pdwRoundKey[4*16+3]))
	}
	// HIGHT_DEC(15,  1,0,7,6,5,4,3,2)
	{
		XX[0] = (XX[0] - (hight_F1[XX[7]] ^ s.pdwRoundKey[4*15+2]))
		XX[6] = (XX[6] ^ (hight_F0[XX[5]] + s.pdwRoundKey[4*15+1]))
		XX[4] = (XX[4] - (hight_F1[XX[3]] ^ s.pdwRoundKey[4*15+0]))
		XX[2] = (XX[2] ^ (hight_F0[XX[1]] + s.pdwRoundKey[4*15+3]))
	}
	// HIGHT_DEC(14,  2,1,0,7,6,5,4,3)
	{
		XX[1] = (XX[1] - (hight_F1[XX[0]] ^ s.pdwRoundKey[4*14+2]))
		XX[7] = (XX[7] ^ (hight_F0[XX[6]] + s.pdwRoundKey[4*14+1]))
		XX[5] = (XX[5] - (hight_F1[XX[4]] ^ s.pdwRoundKey[4*14+0]))
		XX[3] = (XX[3] ^ (hight_F0[XX[2]] + s.pdwRoundKey[4*14+3]))
	}
	// HIGHT_DEC(13,  3,2,1,0,7,6,5,4)
	{
		XX[2] = (XX[2] - (hight_F1[XX[1]] ^ s.pdwRoundKey[4*13+2]))
		XX[0] = (XX[0] ^ (hight_F0[XX[7]] + s.pdwRoundKey[4*13+1]))
		XX[6] = (XX[6] - (hight_F1[XX[5]] ^ s.pdwRoundKey[4*13+0]))
		XX[4] = (XX[4] ^ (hight_F0[XX[3]] + s.pdwRoundKey[4*13+3]))
	}
	// HIGHT_DEC(12,  4,3,2,1,0,7,6,5)
	{
		XX[3] = (XX[3] - (hight_F1[XX[2]] ^ s.pdwRoundKey[4*12+2]))
		XX[1] = (XX[1] ^ (hight_F0[XX[0]] + s.pdwRoundKey[4*12+1]))
		XX[7] = (XX[7] - (hight_F1[XX[6]] ^ s.pdwRoundKey[4*12+0]))
		XX[5] = (XX[5] ^ (hight_F0[XX[4]] + s.pdwRoundKey[4*12+3]))
	}
	// HIGHT_DEC(11,  5,4,3,2,1,0,7,6)
	{
		XX[4] = (XX[4] - (hight_F1[XX[3]] ^ s.pdwRoundKey[4*11+2]))
		XX[2] = (XX[2] ^ (hight_F0[XX[1]] + s.pdwRoundKey[4*11+1]))
		XX[0] = (XX[0] - (hight_F1[XX[7]] ^ s.pdwRoundKey[4*11+0]))
		XX[6] = (XX[6] ^ (hight_F0[XX[5]] + s.pdwRoundKey[4*11+3]))
	}
	// HIGHT_DEC(10,  6,5,4,3,2,1,0,7)
	{
		XX[5] = (XX[5] - (hight_F1[XX[4]] ^ s.pdwRoundKey[4*10+2]))
		XX[3] = (XX[3] ^ (hight_F0[XX[2]] + s.pdwRoundKey[4*10+1]))
		XX[1] = (XX[1] - (hight_F1[XX[0]] ^ s.pdwRoundKey[4*10+0]))
		XX[7] = (XX[7] ^ (hight_F0[XX[6]] + s.pdwRoundKey[4*10+3]))
	}
	// HIGHT_DEC(9,  7,6,5,4,3,2,1,0)
	{
		XX[6] = (XX[6] - (hight_F1[XX[5]] ^ s.pdwRoundKey[4*9+2]))
		XX[4] = (XX[4] ^ (hight_F0[XX[3]] + s.pdwRoundKey[4*9+1]))
		XX[2] = (XX[2] - (hight_F1[XX[1]] ^ s.pdwRoundKey[4*9+0]))
		XX[0] = (XX[0] ^ (hight_F0[XX[7]] + s.pdwRoundKey[4*9+3]))
	}
	// HIGHT_DEC(8,  0,7,6,5,4,3,2,1)
	{
		XX[7] = (XX[7] - (hight_F1[XX[6]] ^ s.pdwRoundKey[4*8+2]))
		XX[5] = (XX[5] ^ (hight_F0[XX[4]] + s.pdwRoundKey[4*8+1]))
		XX[3] = (XX[3] - (hight_F1[XX[2]] ^ s.pdwRoundKey[4*8+0]))
		XX[1] = (XX[1] ^ (hight_F0[XX[0]] + s.pdwRoundKey[4*8+3]))
	}
	// HIGHT_DEC(7,  1,0,7,6,5,4,3,2)
	{
		XX[0] = (XX[0] - (hight_F1[XX[7]] ^ s.pdwRoundKey[4*7+2]))
		XX[6] = (XX[6] ^ (hight_F0[XX[5]] + s.pdwRoundKey[4*7+1]))
		XX[4] = (XX[4] - (hight_F1[XX[3]] ^ s.pdwRoundKey[4*7+0]))
		XX[2] = (XX[2] ^ (hight_F0[XX[1]] + s.pdwRoundKey[4*7+3]))
	}
	// HIGHT_DEC(6,  2,1,0,7,6,5,4,3)
	{
		XX[1] = (XX[1] - (hight_F1[XX[0]] ^ s.pdwRoundKey[4*6+2]))
		XX[7] = (XX[7] ^ (hight_F0[XX[6]] + s.pdwRoundKey[4*6+1]))
		XX[5] = (XX[5] - (hight_F1[XX[4]] ^ s.pdwRoundKey[4*6+0]))
		XX[3] = (XX[3] ^ (hight_F0[XX[2]] + s.pdwRoundKey[4*6+3]))
	}
	// HIGHT_DEC(5,  3,2,1,0,7,6,5,4)
	{
		XX[2] = (XX[2] - (hight_F1[XX[1]] ^ s.pdwRoundKey[4*5+2]))
		XX[0] = (XX[0] ^ (hight_F0[XX[7]] + s.pdwRoundKey[4*5+1]))
		XX[6] = (XX[6] - (hight_F1[XX[5]] ^ s.pdwRoundKey[4*5+0]))
		XX[4] = (XX[4] ^ (hight_F0[XX[3]] + s.pdwRoundKey[4*5+3]))
	}
	// HIGHT_DEC(4,  4,3,2,1,0,7,6,5)
	{
		XX[3] = (XX[3] - (hight_F1[XX[2]] ^ s.pdwRoundKey[4*4+2]))
		XX[1] = (XX[1] ^ (hight_F0[XX[0]] + s.pdwRoundKey[4*4+1]))
		XX[7] = (XX[7] - (hight_F1[XX[6]] ^ s.pdwRoundKey[4*4+0]))
		XX[5] = (XX[5] ^ (hight_F0[XX[4]] + s.pdwRoundKey[4*4+3]))
	}
	// HIGHT_DEC(3,  5,4,3,2,1,0,7,6)
	{
		XX[4] = (XX[4] - (hight_F1[XX[3]] ^ s.pdwRoundKey[4*3+2]))
		XX[2] = (XX[2] ^ (hight_F0[XX[1]] + s.pdwRoundKey[4*3+1]))
		XX[0] = (XX[0] - (hight_F1[XX[7]] ^ s.pdwRoundKey[4*3+0]))
		XX[6] = (XX[6] ^ (hight_F0[XX[5]] + s.pdwRoundKey[4*3+3]))
	}
	// HIGHT_DEC(2,  6,5,4,3,2,1,0,7)
	{
		XX[5] = (XX[5] - (hight_F1[XX[4]] ^ s.pdwRoundKey[4*2+2]))
		XX[3] = (XX[3] ^ (hight_F0[XX[2]] + s.pdwRoundKey[4*2+1]))
		XX[1] = (XX[1] - (hight_F1[XX[0]] ^ s.pdwRoundKey[4*2+0]))
		XX[7] = (XX[7] ^ (hight_F0[XX[6]] + s.pdwRoundKey[4*2+3]))
	}

	// Final Round
	dst[0] = XX[0] - s.pdwRoundKey[0]
	dst[1] = XX[1]
	dst[2] = XX[2] ^ s.pdwRoundKey[1]
	dst[3] = XX[3]
	dst[4] = XX[4] - s.pdwRoundKey[2]
	dst[5] = XX[5]
	dst[6] = XX[6] ^ s.pdwRoundKey[3]
	dst[7] = XX[7]
}
