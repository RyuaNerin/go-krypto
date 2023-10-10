package testingutil

import (
	"os"
	"strings"
)

func init() {
	if len(os.Args) < 2 || !strings.Contains(os.Args[1], "test") {
		print("warning: testingutil package is only for go test.")
	}
}
