//go:build armbe || arm64be || m68k || mips || mips64 || mips64p32 || ppc || ppc64 || s390 || s390x || shbe || sparc || sparc64
// +build armbe arm64be m68k mips mips64 mips64p32 ppc ppc64 s390 s390x shbe sparc sparc64

package internal

const (
	IsLittleEndian = false
	IsBigEndian    = true
)
