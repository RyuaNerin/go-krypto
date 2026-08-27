// Package pbkdf2 implements HMAC-based Key Derivation Functions, as defined in TTAK.KO-12.0334, NIST SP 800-132, RFC 2898(PKCS #5)
package pbkdf2

import (
	"crypto/subtle"
	"encoding/binary"
	"hash"
	"runtime"
	"sync"

	"github.com/RyuaNerin/go-krypto/hmac"
	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

// Key derives a PBKDF2 key from the given password, salt, and iteration count,
// using h as the hash function, and returns a key of length keyLen.
func Key(password, salt []byte, iter, keyLen int, h func() hash.Hash) []byte {
	return Generate(nil, password, salt, iter, keyLen, h)
}

// Generate derives a PBKDF2 key from the given password, salt, and iteration count,
// using h as the hash function. The result has length keyLen and is appended to
// dst; if dst is nil, a new slice is allocated.
func Generate(
	dst []byte,
	password, salt []byte,
	iteration, keyLen int,
	h func() hash.Hash,
) []byte {
	if iteration <= 0 {
		panic(msgInvalidIterationCount)
	}
	if keyLen <= 0 {
		panic(msgInvalidKeyLength)
	}

	hh := hmac.New(h, password)
	defer zeroize.Zeroize(hh)

	hLen := hh.Size()

	checkKeyLength(keyLen, hLen)

	U := make([]byte, hLen, hLen*2)
	T := U[hLen : hLen*2]
	defer clear(U[:cap(U)])

	out, dst := internal.SliceForAppend(dst, keyLen)

	var ctr [4]byte
	defer clear(ctr[:])
	for len(dst) > 0 {
		internal.IncCtr(ctr[:])

		hh.Reset()
		hh.Write(salt)
		hh.Write(ctr[:])
		U = hh.Sum(U[:0])

		copy(T, U)

		for iter := 1; iter < iteration; iter++ {
			hh.Reset()
			hh.Write(U)
			U = hh.Sum(U[:0])

			subtle.XORBytes(T, T, U)
		}

		dst = dst[copy(dst, T):]
	}

	return out
}

// GenerateParallel derives a PBKDF2 key concurrently across goroutines. It behaves
// like Generate, deriving a key of length keyLen from the given password, salt,
// and iteration count using h.
//
// nThreads is the number of goroutines to use. If nThreads <= 0, the number of
// CPUs is used.
func GenerateParallel(
	dst []byte,
	password, salt []byte,
	iteration, keyLen int,
	h func() hash.Hash,
	nThreads int,
) []byte {
	if iteration <= 0 {
		panic(msgInvalidIterationCount)
	}
	if keyLen <= 0 {
		panic(msgInvalidKeyLength)
	}

	hLen := h().Size()
	checkKeyLength(keyLen, hLen)

	if nThreads <= 0 {
		nThreads = runtime.NumCPU()
	}
	nThreads = min(nThreads, internal.CeilDiv(keyLen, hLen))

	if nThreads <= 2 {
		return Generate(dst, password, salt, iteration, keyLen, h)
	}
	return generateParallel(
		dst,
		password, salt,
		h,
		iteration, keyLen, nThreads,
		hLen,
	)
}

func generateParallel(
	dst []byte,
	password, salt []byte,
	h func() hash.Hash,
	iteration, keyLen, nThreads int,
	hLen int,
) []byte {
	out, dst := internal.SliceForAppend(dst, keyLen)

	type iter struct {
		off int
		ctr uint32
	}
	ch := make(chan iter, nThreads)

	var w sync.WaitGroup
	for range nThreads {
		w.Go(func() {
			hh := hmac.New(h, password)
			defer zeroize.Zeroize(hh)

			U := make([]byte, hLen, hLen*2)
			T := U[hLen : hLen*2]
			defer clear(U[:cap(U)])

			var ctr [4]byte
			defer clear(ctr[:])
			for it := range ch {
				binary.BigEndian.PutUint32(ctr[:], it.ctr)

				hh.Reset()
				hh.Write(salt)
				hh.Write(ctr[:])
				U = hh.Sum(U[:0])

				copy(T, U)

				for iter := 1; iter < iteration; iter++ {
					hh.Reset()
					hh.Write(U)
					U = hh.Sum(U[:0])

					subtle.XORBytes(T, T, U)
				}

				copy(dst[it.off:], T)
			}
		})
	}

	var ctr uint32
	for off := 0; off < keyLen; off += hLen {
		ctr++
		ch <- iter{off, ctr}
	}

	close(ch)

	w.Wait()
	return out
}

// Check before allocating output or choosing a parallel path. Computing the
// block count as a quotient avoids overflowing keyLen+hLen-1.
func checkKeyLength(keyLen, hLen int) {
	const maxBlocks = uint64(1<<32 - 1)
	if hLen <= 0 || uint64(keyLen-1)/uint64(hLen) >= maxBlocks {
		panic(msgKeyLengthTooLong)
	}
}
