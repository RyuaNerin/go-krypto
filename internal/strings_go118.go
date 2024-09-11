//go:build gc && go1.18
// +build gc,go1.18

package internal

import (
	"strings"
)

func StringClone(s string) string {
	return strings.Clone(s)
}
