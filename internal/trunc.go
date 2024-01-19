package internal

func TruncateLeft(b []byte, bits int) []byte {
	bytes := Bytes(bits)
	b = b[len(b)-bytes:]

	remain := bits % 8
	if remain > 0 {
		b[0] = b[0] & ((1 << remain) - 1)
	}

	return b
}
