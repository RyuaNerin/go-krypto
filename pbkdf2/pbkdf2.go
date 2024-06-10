// Package pbkdf2 implements HMAC-based Key Derivation Functions, as defined in TTAK.KO-12.0334-Part1
package pbkdf2

import (
	"crypto/hmac"
	"encoding/binary"
	"hash"
	"runtime"
	"sync"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/subtle"
)

// Generate a key from the password, salt and iteration count,
// then returns a []byte of length keylen.
func Generate(dst, password, salt []byte, iteration, keyLen int, h func() hash.Hash) []byte {
	hh := hmac.New(h, password)
	hLen := hh.Size()

	U := make([]byte, hLen, hLen*2)
	T := U[hLen : hLen*2]

	out, dst := internal.SliceForAppend(dst, keyLen)

	var i [4]byte
	for off := 0; off < keyLen; off += hLen {
		internal.IncCtr(i[:])

		hh.Reset()
		hh.Write(salt)
		hh.Write(i[:])
		U = hh.Sum(U[:0])

		copy(T, U)

		for iter := 1; iter < iteration; iter++ {
			hh.Reset()
			hh.Write(U)
			U = hh.Sum(U[:0])

			subtle.XORBytes(T, T, U)
		}

		copy(dst[off:], T)
	}

	return out
}

// Generate a key from the password, salt and iteration count,
// then returns a []byte of length keylen.
//
// nThreads: the number of goroutines to use. If nThreads <= 0, the number of CPUs is used.
func GenerateParallel(dst, password, salt []byte, iteration, keyLen int, h func() hash.Hash, nThreads int) []byte {
	hLen := h().Size()

	if nThreads <= 0 {
		nThreads = runtime.NumCPU()
	}

	iters := (keyLen + hLen - 1) / hLen
	if iters < nThreads {
		nThreads = iters
	}

	out, dst := internal.SliceForAppend(dst, iters*hLen)

	type iter struct {
		off int
		ctr uint32
	}
	ch := make(chan iter, nThreads)

	var w sync.WaitGroup
	w.Add(nThreads)
	for i := 0; i < nThreads; i++ {
		go func() {
			defer w.Done()

			hh := hmac.New(h, password)

			U := make([]byte, hLen)

			var ctr [4]byte
			for it := range ch {
				binary.BigEndian.PutUint32(ctr[:], it.ctr)

				hh.Reset()
				hh.Write(salt)
				hh.Write(ctr[:])
				U = hh.Sum(U[:0])

				copy(dst[it.off:], U)

				for iter := 1; iter < iteration; iter++ {
					hh.Reset()
					hh.Write(U)
					U = hh.Sum(U[:0])

					subtle.XORBytes(dst[it.off:], dst[it.off:], U)
				}
			}
		}()
	}

	var ctr uint32
	for off := 0; off < keyLen; off += hLen {
		ctr++
		ch <- iter{off, ctr}
	}

	close(ch)

	w.Wait()
	return out[len(out)-len(dst)+keyLen:]
}
