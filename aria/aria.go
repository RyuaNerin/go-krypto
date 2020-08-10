package aria

func dl(i, o []byte) {
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

// Right-rotate 128 bit source string s by n bits and XOR it to target string t
func rotXOR(s []byte, n int, t []byte, ti int) {
	q := n / 8
	n %= 8
	for i := 0; i < 16; i++ {
		t[ti+(q+i)%16] ^= (s[i] >> n)
		if n != 0 {
			t[ti+(q+i+1)%16] ^= (s[i] << (8 - n))
		}
	}
}

// Encryption round key generation rountine
// w0 : master key, e : encryption round keys
func encKeySetup(w0 []byte, e []byte, keyBits int) int {
	var i int
	var q int
	var R int = (keyBits + 256) / 32
	var t, w1, w2, w3 [16]byte

	q = (keyBits - 128) / 64
	for i = 0; i < 16; i++ {
		t[i] = s[i%4][krk[q][i]^w0[i]]
	}

	dl(t[:], w1[:])
	if R == 14 {
		for i = 0; i < 8; i++ {
			w1[i] ^= w0[16+i]
		}
	} else if R == 16 {
		for i = 0; i < 16; i++ {
			w1[i] ^= w0[16+i]
		}
	}

	//q = (q==2)? 0 : (q+1);
	if q == 2 {
		q = 0
	} else {
		q = q + 1
	}

	for i = 0; i < 16; i++ {

		t[i] = s[(2+i)%4][krk[q][i]^w1[i]]
	}
	dl(t[:], w2[:])
	for i = 0; i < 16; i++ {
		w2[i] ^= w0[i]
	}

	//q = (q==2)? 0 : (q+1);
	if q == 2 {
		q = 0
	} else {
		q = (q + 1)
	}
	for i = 0; i < 16; i++ {
		t[i] = s[i%4][krk[q][i]^w2[i]]
	}
	dl(t[:], w3[:])
	for i = 0; i < 16; i++ {
		w3[i] ^= w1[i]
	}

	for i = 0; i < 16*(R+1); i++ {
		e[i] = 0
	}

	rotXOR(w0, 0, e, 0)
	rotXOR(w1[:], 19, e, 0)
	rotXOR(w1[:], 0, e, 16)
	rotXOR(w2[:], 19, e, 16)
	rotXOR(w2[:], 0, e, 32)
	rotXOR(w3[:], 19, e, 32)
	rotXOR(w3[:], 0, e, 48)
	rotXOR(w0, 19, e, 48)
	rotXOR(w0, 0, e, 64)
	rotXOR(w1[:], 31, e, 64)
	rotXOR(w1[:], 0, e, 80)
	rotXOR(w2[:], 31, e, 80)
	rotXOR(w2[:], 0, e, 96)
	rotXOR(w3[:], 31, e, 96)
	rotXOR(w3[:], 0, e, 112)
	rotXOR(w0, 31, e, 112)
	rotXOR(w0, 0, e, 128)
	rotXOR(w1[:], 67, e, 128)
	rotXOR(w1[:], 0, e, 144)
	rotXOR(w2[:], 67, e, 144)
	rotXOR(w2[:], 0, e, 160)
	rotXOR(w3[:], 67, e, 160)
	rotXOR(w3[:], 0, e, 176)
	rotXOR(w0, 67, e, 176)
	rotXOR(w0, 0, e, 192)
	rotXOR(w1[:], 97, e, 192)
	if R > 12 {
		rotXOR(w1[:], 0, e, 208)
		rotXOR(w2[:], 97, e, 208)
		rotXOR(w2[:], 0, e, 224)
		rotXOR(w3[:], 97, e, 224)
	}
	if R > 14 {
		rotXOR(w3[:], 0, e, 240)
		rotXOR(w0, 97, e, 240)
		rotXOR(w0, 0, e, 256)
		rotXOR(w1[:], 109, e, 256)
	}
	return R
}

// Decryption round key generation rountine
// w0 : maskter key, d : decryption round keys
func decKeySetup(w0 []byte, d []byte, keyBits int) int {
	var i, j, R int
	var t [16]byte

	R = encKeySetup(w0, d, keyBits)
	for j = 0; j < 16; j++ {
		t[j] = d[j]
		d[j] = d[16*R+j]
		d[16*R+j] = t[j]
	}
	for i = 1; i <= R/2; i++ {
		dl(d[i*16:], t[:])
		dl(d[(R-i)*16:], d[i*16:])
		for j = 0; j < 16; j++ {
			d[(R-i)*16+j] = t[j]
		}
	}
	return R
}

// Encryption and decryption rountine
// p: plain text, e: round keys, c: ciphertext
func crypt(p []byte, R int, e []byte, c []byte) {
	var i, j int
	var t [16]byte

	copy(c, p[:16])

	ei := 0
	for i = 0; i < R/2; i++ {
		for j = 0; j < 16; j++ {
			t[j] = s[j%4][e[ei+j]^c[j]]
		}
		dl(t[:], c)
		ei += 16
		for j = 0; j < 16; j++ {
			t[j] = s[(2+j)%4][e[ei+j]^c[j]]
		}
		dl(t[:], c)
		ei += 16
	}
	dl(c, t[:])
	for j = 0; j < 16; j++ {
		c[j] = e[ei+j] ^ t[j]
	}
}
