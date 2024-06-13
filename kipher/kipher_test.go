package kipher_test

import (
	"bufio"
	"crypto/rand"
)

var rnd = bufio.NewReaderSize(rand.Reader, 1<<15)

const (
	keySize = 16
	iter    = 64 * 1024
	blocks  = (8 + 4 + 1)
)
