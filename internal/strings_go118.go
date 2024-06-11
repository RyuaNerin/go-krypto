//go:build go1.18
// +build go1.18

package internal

import (
	"strings"
)

func StringClone(s string) string {
	return strings.Clone(s)
}
