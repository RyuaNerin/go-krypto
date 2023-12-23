package testingutil

import (
	"bufio"
	"crypto/rand"
)

var rnd = bufio.NewReaderSize(rand.Reader, 1<<15)

const (
	continusBlockTestIter = 64 * 1024
	continusHashTestIter  = 8 * 1024
)
