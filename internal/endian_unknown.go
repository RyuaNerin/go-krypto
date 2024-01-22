//go:build !armbe && !arm64be && !m68k && !mips && !mips64 && !mips64p32 && !ppc && !ppc64 && !s390 && !s390x && !shbe && !sparc && !sparc64 && !386 && !amd64 && !amd64p32 && !alpha && !arm && !arm64 && !loong64 && !mipsle && !mips64le && !mips64p32le && !nios2 && !ppc64le && !riscv && !riscv64 && !sh && !wasm
// +build !armbe,!arm64be,!m68k,!mips,!mips64,!mips64p32,!ppc,!ppc64,!s390,!s390x,!shbe,!sparc,!sparc64,!386,!amd64,!amd64p32,!alpha,!arm,!arm64,!loong64,!mipsle,!mips64le,!mips64p32le,!nios2,!ppc64le,!riscv,!riscv64,!sh,!wasm

package internal

import (
	"unsafe"
)

var (
	IsLittleEndian bool
	IsBigEndian    bool
)

func init() {
	// https://stackoverflow.com/a/53286786

	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		IsLittleEndian = true
		IsBigEndian = false
	case [2]byte{0xAB, 0xCD}:
		IsLittleEndian = false
		IsBigEndian = true
	default:
		panic("Could not determine native endianness.")
	}
}
