package main

import (
	"strconv"

	"github.com/RyuaNerin/go-krypto/internal/subtle"
)

func processBlock_MCT_ECB(cavp *cavpProcessor, newCipher funcNewBlockCipher, key, pt, ct, iv [][]byte, blockSize int) {
	res := cavpSection{
		{"COUNT", "", false},
		{"KEY", "", false},
		{"PT", "", false},
		{"CT", "", false},
	}
	var j int
	for i := 0; i < 100; i++ {
		res[0].Value = strconv.Itoa(i)

		res[1].Value = hexStr(key[i]) // Output Key[i];
		res[2].Value = hexStr(pt[0])  // Output PT[0];

		c, err := newCipher(key[i])
		if err != nil {
			panic(err)
		}
		for j = 0; j <= 999; j++ {
			// CT[j] = ARIA(Key[i], PT[j]);
			c.Encrypt(ct[j], pt[j])
			// PT[j+1] = CT[j];
			copy(pt[j+1], ct[j])
		}
		j = 999

		res[3].Value = hexStr(ct[j]) // Output CT[j];
		cavp.WriteValues(res)

		/**
		If ( keylen = 128 ) Key[i+1] = Key[i] xor CT[j];
		If ( keylen = 192 ) Key[i+1] = Key[i] xor (CT64[j-1]||CT[j]);
		If ( keylen = 256 ) Key[i+1] = Key[i] xor (CT[j-1]||CT[j]);
		*/
		{
			remain := len(key[0])
			ctIdx := j
			for remain > 0 {
				if remain < blockSize {
					subtle.XORBytes(key[i+1], key[i], ct[ctIdx][remain:])
				} else {
					subtle.XORBytes(key[i+1][remain-blockSize:], key[i][remain-blockSize:], ct[ctIdx])
				}
				remain -= blockSize
				ctIdx--
			}
		}

		// PT[0] = CT[j];
		copy(pt[0], ct[j])
	}
}

func processBlock_MCT_CBC(cavp *cavpProcessor, newCipher funcNewBlockCipher, key, pt, ct, iv [][]byte, blockSize int) {
	res := cavpSection{
		{"COUNT", "", false},
		{"KEY", "", false},
		{"IV", "", false},
		{"PT", "", false},
		{"CT", "", false},
	}

	tmp := make([]byte, blockSize)

	var j int
	for i := 0; i < 100; i++ {
		res[0].Value = strconv.Itoa(i)
		res[1].Value = hexStr(key[i]) // Output Key[i];
		res[2].Value = hexStr(iv[i])  // Output IV[i];
		res[3].Value = hexStr(pt[0])  // Output PT[0];

		c, err := newCipher(key[i])
		if err != nil {
			panic(err)
		}
		for j = 0; j <= 999; j++ {
			if j == 0 {
				// CT[j] = ARIA(Key[i], PT[j] xor IV[i]);
				subtle.XORBytes(tmp, pt[j], iv[i])
				c.Encrypt(ct[j], tmp)
				// PT[j+1] = IV[i];
				copy(pt[j+1], iv[i])
			} else {
				// CT[j] = ARIA(Key[i], PT[j] xor CT[j-1]);
				subtle.XORBytes(tmp, pt[j], ct[j-1])
				c.Encrypt(ct[j], tmp)
				// PT[j+1] = CT[j-1];
				copy(pt[j+1], ct[j-1])
			}
		}
		j = 999

		res[4].Value = hexStr(ct[j]) // Output CT[j];
		cavp.WriteValues(res)

		/**
		If ( keylen = 128 ) Key[i+1] = Key[i] xor CT[j];
		If ( keylen = 192 ) Key[i+1] = Key[i] xor (CT64[j-1]||CT[j]);
		If ( keylen = 256 ) Key[i+1] = Key[i] xor (CT[j-1]||CT[j]);
		*/
		{
			remain := len(key[i])
			ctIdx := j
			for remain > 0 {
				if remain < blockSize {
					subtle.XORBytes(key[i+1], key[i], ct[ctIdx][remain:])
				} else {
					subtle.XORBytes(key[i+1][remain-blockSize:], key[i][remain-blockSize:], ct[ctIdx])
				}
				remain -= blockSize
				ctIdx--
			}
		}

		// IV[i+1] = CT[j];
		copy(iv[i+1], ct[j])
		// PT[0] = CT[j-1];
		copy(pt[0], ct[j-1])
	}
}

func processBlock_MCT_CFB(sz int) func(cavp *cavpProcessor, newCipher funcNewBlockCipher, key, pt, ct, iv [][]byte, blockSize int) {
	sz /= 8

	return func(cavp *cavpProcessor, newCipher funcNewBlockCipher, key, pt, ct, iv [][]byte, blockSize int) {
		res := cavpSection{
			{"COUNT", "", false},
			{"KEY", "", false},
			{"IV", "", false},
			{"PT", "", false},
			{"CT", "", false},
		}
		cf := make([][]byte, len(iv))
		for idx := range cf {
			cf[idx] = make([]byte, blockSize)
		}

		// log.Println("KEY", hex.EncodeToString(key[0]))
		// log.Println("IV", hex.EncodeToString(iv[0]))
		// log.Println("PT", hex.EncodeToString(pt[0]))

		tmp := make([]byte, blockSize)

		var jj int
		for i := 0; i < 100; i++ {
			// log.Println("==================================================")
			// log.Println("COUNT", i)
			res[0].Value = strconv.Itoa(i)
			res[1].Value = hexStr(key[i]) // Output Key[i];
			res[2].Value = hexStr(iv[i])  // Output IV[i];
			res[3].Value = hexStr(pt[0])  // Output PT[0];

			c, err := newCipher(key[i])
			if err != nil {
				panic(err)
			}
			for jj = 0; jj <= 999; jj++ {
				if jj == 0 {
					// CT[j] = PT[j] xor ARIA(Key[i], IV[i]);
					c.Encrypt(tmp, iv[i])
					subtle.XORBytes(ct[jj], pt[jj], tmp)

					// PT[j+1] = ByteJ(IV[i]);
					copy(pt[jj+1], iv[i])

					// CF[j+1] = LSB120(IV[i]) || CT[j];
					copy(cf[jj+1], iv[i][sz:])
					copy(cf[jj+1][blockSize-sz:], ct[jj])
				} else {
					// CT[j] = PT[j] xor ARIA(Key[i], CF[j]);
					c.Encrypt(tmp, cf[jj])
					subtle.XORBytes(ct[jj], pt[jj], tmp)

					if jj < blockSize/sz {
						// PT[j+1] = ByteJ(IV[i]);
						copy(pt[jj+1], iv[i][sz*jj:])
					} else { // else
						// PT[j+1] = CT[j-16];
						copy(pt[jj+1], ct[jj-blockSize/sz])
					}
					// CF[j+1] = LSB120(CF[j]) || CT[j];
					copy(cf[jj+1], cf[jj][sz:])
					copy(cf[jj+1][blockSize-sz:], ct[jj])
				}

				// log.Println("jj", jj)
				// log.Println("TMP", hex.EncodeToString(tmp))
				// log.Println("CF", hex.EncodeToString(cf[jj]))
				// log.Println("CT", hex.EncodeToString(ct[jj]))
				// log.Println("PT", hex.EncodeToString(pt[jj]))
			}
			jj = 999

			res[4].Value = hexStr(ct[jj]) // Output CT[j];
			cavp.WriteValues(res)

			{
				remain := len(key[0])
				ctIdx := jj
				for remain > 0 {
					if remain < sz {
						remain -= subtle.XORBytes(key[i+1], key[i], ct[ctIdx][remain:])
					} else {
						remain -= subtle.XORBytes(key[i+1][remain-sz:], key[i][remain-sz:], ct[ctIdx])
					}
					ctIdx--
				}
			}

			// IV[i+1] = (CT[j-15] || CT[j-14] || ... || CT[j]);
			for k := 0; k < blockSize/sz; k++ {
				copy(iv[i+1][blockSize-sz*(k+1):], ct[jj-k])
			}

			// PT[0] = CT[j-16];
			copy(pt[0], ct[jj-blockSize/sz])
		}
	}
}

func processBlock_MCT_OFB(cavp *cavpProcessor, newCipher funcNewBlockCipher, key, pt, ct, iv [][]byte, blockSize int) {
	res := cavpSection{
		{"COUNT", "", false},
		{"KEY", "", false},
		{"IV", "", false},
		{"PT", "", false},
		{"CT", "", false},
	}
	ot := make([][]byte, 1000)
	for idx := range ot {
		ot[idx] = make([]byte, blockSize)
	}

	var jj int
	for i := 0; i < 100; i++ {
		res[0].Value = strconv.Itoa(i)
		res[1].Value = hexStr(key[i]) // Output Key[i];
		res[2].Value = hexStr(iv[i])  // Output IV[i];
		res[3].Value = hexStr(pt[0])  // Output PT[0];

		c, err := newCipher(key[i])
		if err != nil {
			panic(err)
		}
		for jj = 0; jj <= 999; jj++ {
			if jj == 0 {
				// OT[j] = SEED(Key[i], IV[i]);
				c.Encrypt(ot[jj], iv[i])
				// CT[j] = PT[j] xor OT[j];
				subtle.XORBytes(ct[jj], pt[jj], ot[jj])
				// PT[j+1] = IV[i];
				copy(pt[jj+1], iv[i])
			} else {
				// OT[j] = SEED(Key[i], OT[j-1]);
				c.Encrypt(ot[jj], ot[jj-1])
				// CT[j] = PT[j] xor OT[j];
				subtle.XORBytes(ct[jj], pt[jj], ot[jj])
				// PT[j+1] = CT[j-1];
				copy(pt[jj+1], ct[jj-1])
			}
		}
		jj = 999

		res[4].Value = hexStr(ct[jj]) // Output CT[j];
		cavp.WriteValues(res)

		/**
		If ( keylen = 128 ) Key[i+1] = Key[i] xor CT[j];
		If ( keylen = 192 ) Key[i+1] = Key[i] xor (CT64[j-1]||CT[j]);
		If ( keylen = 256 ) Key[i+1] = Key[i] xor (CT[j-1]||CT[j]);
		*/
		{
			remain := len(key[0])
			ctIdx := jj
			for remain > 0 {
				if remain < blockSize {
					subtle.XORBytes(key[i+1], key[i], ct[ctIdx][remain:])
				} else {
					subtle.XORBytes(key[i+1][remain-blockSize:], key[i][remain-blockSize:], ct[ctIdx])
				}
				remain -= blockSize
				ctIdx--
			}
		}

		// IV[i+1] = CT[j];
		copy(iv[i+1], ct[jj])

		// PT[0] = CT[j-1];
		copy(pt[0], ct[jj-1])
	}
}

func processBlock_MCT_CTR(cavp *cavpProcessor, newCipher funcNewBlockCipher, key, pt, ct, iv [][]byte, blockSize int) {
	res := cavpSection{
		{"COUNT", "", false},
		{"KEY", "", false},
		{"CTR", "", false},
		{"PT", "", false},
		{"CT", "", false},
	}
	tmp := make([]byte, blockSize)

	ctr := iv

	var j int
	for i := 0; i < 100; i++ {
		res[0].Value = strconv.Itoa(i)
		res[1].Value = hexStr(key[i]) // Output Key[i];
		res[2].Value = hexStr(iv[0])  // Output CTR[0];
		res[3].Value = hexStr(pt[0])  // Output PT[0];

		c, err := newCipher(key[i])
		if err != nil {
			panic(err)
		}
		for j = 0; j <= 999; j++ {
			// CT[j] = PT[j] xor SEED(Key[i], CTR[0]);
			c.Encrypt(tmp, ctr[0])
			subtle.XORBytes(ct[j], pt[j], tmp)

			// CTR[0] = (CTR[0] + 1) mod 2^128
			for i := len(ctr[0]) - 1; i >= 0; i-- {
				ctr[0][i]++
				if ctr[0][i] != 0 {
					break
				}
			}

			// PT[j+1] = CT[j];
			copy(pt[j+1], ct[j])
		}
		j = 999

		res[4].Value = hexStr(ct[j]) // Output CT[j];
		cavp.WriteValues(res)

		// Key[i+1] = Key[i] xor CT[j];
		{
			remain := len(key[0])
			ctIdx := j
			for remain > 0 {
				if remain < blockSize {
					subtle.XORBytes(key[i+1], key[i], ct[ctIdx][remain:])
				} else {
					subtle.XORBytes(key[i+1][remain-blockSize:], key[i][remain-blockSize:], ct[ctIdx])
				}
				remain -= blockSize
				ctIdx--
			}
		}
		// PT[0] = CT[j];
		copy(pt[0], ct[j])
	}
}
