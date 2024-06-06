package kipher

import (
	"bufio"
	"crypto/rand"

	. "github.com/RyuaNerin/testingutil"
)

var rnd = bufio.NewReaderSize(rand.Reader, 1<<15)

var as = []CipherSize{
	{Name: "2 blocks", Size: 2},
	{Name: "3 blocks", Size: 3},
	{Name: "4 Blocks", Size: 4},
	{Name: "7 Blocks", Size: 7},
	{Name: "8 Blocks", Size: 8},
	{Name: "13 Blocks", Size: 13},
	{Name: "16 Blocks", Size: 16},
}
