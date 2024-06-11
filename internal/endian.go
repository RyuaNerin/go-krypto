package internal

func AppendBigUint8(b []byte, v byte) []byte {
	return append(b,
		v,
	)
}

func AppendBigUint16(b []byte, v uint16) []byte {
	return append(b,
		byte(v>>8),
		byte(v),
	)
}

func AppendBigUint32(b []byte, v uint32) []byte {
	return append(b,
		byte(v>>24),
		byte(v>>16),
		byte(v>>8),
		byte(v),
	)
}

func AppendBigUint64(b []byte, v uint64) []byte {
	return append(b,
		byte(v>>56),
		byte(v>>48),
		byte(v>>40),
		byte(v>>32),
		byte(v>>24),
		byte(v>>16),
		byte(v>>8),
		byte(v),
	)
}
