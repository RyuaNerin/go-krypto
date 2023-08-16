package avoutil

import (
	"encoding/hex"
	"strings"

	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
)

func AllocHex(name string, value string) Mem {
	mem := GLOBL(name, NOPTR|RODATA)

	value = strings.TrimPrefix(value, "0x")
	bytes, err := hex.DecodeString(value)
	if err != nil {
		panic(err)
	}
	for idx, v := range bytes {
		DATA(idx, U8(v))
	}

	return mem
}

func Alloc8(name string, values ...uint32) Mem {
	mem := GLOBL(name, NOPTR|RODATA)

	for idx, v := range values {
		DATA(idx, U8(v))
	}
	return mem
}

func Alloc32(name string, values ...uint32) Mem {
	mem := GLOBL(name, NOPTR|RODATA)

	for idx, v := range values {
		DATA(4*idx, U32(v))
	}
	return mem
}

func Alloc64(name string, values ...uint64) Mem {
	mem := GLOBL(name, NOPTR|RODATA)

	for idx, v := range values {
		DATA(8*idx, U64(v))
	}
	return mem
}
