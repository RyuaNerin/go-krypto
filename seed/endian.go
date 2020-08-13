package seed

import (
	"unsafe"
)

var littleEndian = false

func init() {
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		littleEndian = true
	case [2]byte{0xAB, 0xCD}:
		littleEndian = false
	default:
		panic("Could not determine native endianness.")
	}
}

func endianChange(dwS uint32) uint32 {
	return (((dwS << 8) | (dwS >> 24)) & uint32(0x00ff00ff)) |
		(((dwS << 24) | (dwS >> 8)) & uint32(0xff00ff00))
}
