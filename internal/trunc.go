package internal

// 0 0[0 0 0 0 0 0]
func TruncateLeft(b []byte, bits int) []byte {
	bytes := Bytes(bits)
	b = b[len(b)-bytes:]

	remain := bits % 8
	if remain > 0 {
		b[0] = b[0] & ((1 << remain) - 1)
	}

	return b
}

// [0 0 0 0 0 0]0 0
func TruncateRight(b []byte, bits int) []byte {
	bytes := Bytes(bits)
	b = b[:bytes]

	remain := bits % 8
	if remain > 0 {
		b[0] = b[0] & byte(0b_11111111<<(8-remain))
	}

	return b
}
