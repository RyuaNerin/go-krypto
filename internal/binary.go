package internal

import "encoding/binary"

func AppendBigUint8[T ~byte](b []byte, v T) []byte {
	return append(
		b,
		byte(v),
	)
}

func AppendBigUint16[T ~uint16](b []byte, v T) []byte {
	return append(
		b,
		byte(v>>8),
		byte(v),
	)
}

func AppendBigUint32[T ~uint32](b []byte, v T) []byte {
	return append(
		b,
		byte(v>>24),
		byte(v>>16),
		byte(v>>8),
		byte(v),
	)
}

func AppendBigUint64[T ~uint64](b []byte, v T) []byte {
	return append(
		b,
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

func ConsumeBigU8(b []byte) ([]byte, uint8) {
	_ = b[0]
	return b[1:], b[0]
}

func ConsumeBigU16(b []byte) ([]byte, uint16) {
	_ = b[1]
	return b[2:], binary.BigEndian.Uint16(b)
}

func ConsumeBigU32(b []byte) ([]byte, uint32) {
	_ = b[3]
	return b[4:], binary.BigEndian.Uint32(b)
}

func ConsumeBigU64(b []byte) ([]byte, uint64) {
	_ = b[7]
	return b[8:], binary.BigEndian.Uint64(b)
}
