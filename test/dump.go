package test

import (
	"fmt"
	"strings"
)

func DumpByteArray(name string, a []byte, b []byte) string {
	var sb strings.Builder
	sb.WriteString(name)
	sb.WriteByte('\n')

	for i := 0; i < len(a); i++ {
		fmt.Fprintf(&sb, "[%3d] = %02x / %02x", i, a[i], b[i])
		if a[i] != b[i] {
			sb.WriteString("  <<<")
		}
		sb.WriteByte('\n')
	}

	return sb.String()
}
func DumpUint32Array(name string, a []uint32, b []uint32) string {
	var sb strings.Builder
	sb.WriteString(name)
	sb.WriteByte('\n')

	for i := 0; i < len(a); i++ {
		fmt.Fprintf(&sb, "[%3d] = %08x / %08x", i, a[i], b[i])
		if a[i] != b[i] {
			sb.WriteString("  <<<")
		}
		sb.WriteByte('\n')
	}

	return sb.String()
}
