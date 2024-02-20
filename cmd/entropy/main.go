package main

import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/randutil"
)

func main() {
	generate(32, "crypto.rand", rand.Reader)
	generate(64, "crypto.rand", rand.Reader)

	csprng := newCSPRNG(rand.Reader)
	generate(32, "csprng", csprng)
	generate(64, "csprng", csprng)
}

func generate(bits int, prefix string, rand io.Reader) {
	fs, err := os.Create(fmt.Sprintf("entropy_%s_%s_%dbits.txt", prefix, runtime.GOOS, bits))
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

	var buf []byte
	for i := 0; i < 25_0000; i++ {
		buf, err = internal.ReadBits(buf, rand, bits)
		if err != nil {
			panic(err)
		}

		hw.Write(buf)
	}
}

func newCSPRNG(r io.Reader) io.Reader {
	v := make([]byte, 1)
	r.Read(v)

	d := make([]byte, v[0])
	if _, err := io.ReadFull(r, d); err != nil {
		panic(err)
	}

	hashInput := make([]byte, v[0])
	if _, err := io.ReadFull(r, hashInput); err != nil {
		panic(err)
	}

	hash := sha256.Sum256(hashInput)
	csprng, err := randutil.MixedCSPRNG(rand.Reader, d, hash[:])
	if err != nil {
		panic(err)
	}

	return csprng

}
