package internal

import "encoding/binary"

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

func ConsumeBigU16(b []byte) ([]byte, uint16) {
	return b[2:], binary.BigEndian.Uint16(b)
}

func ConsumeBigU32(b []byte) ([]byte, uint32) {
	return b[4:], binary.BigEndian.Uint32(b)
}

func ConsumeBigU64(b []byte) ([]byte, uint64) {
	return b[8:], binary.BigEndian.Uint64(b)
}
