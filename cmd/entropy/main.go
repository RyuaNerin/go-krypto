package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/RyuaNerin/go-krypto/internal"
)

func main() {
	generate(32, rand.Reader)
	generate(64, rand.Reader)
}

func generate(bits int, rand io.Reader) {
	fs, err := os.Create(fmt.Sprintf("entropy_%s_%dbits.txt", runtime.GOOS, bits))
	if err != nil {
		panic(err)
	}
	defer fs.Close()

	fs.Truncate(0)
	fs.Seek(0, io.SeekStart)

	w := bufio.NewWriter(fs)
	defer w.Flush()

	fmt.Fprintf(w, "Len = %d, Num = 250000\n\n", bits)

	hw := hex.NewEncoder(w)

	var b []byte
	for i := 0; i < 25_0000; i++ {
		b, err := internal.ReadBits(b[:0], rand, bits)
		if err != nil {
			panic(err)
		}

		hw.Write(b)
	}
}
