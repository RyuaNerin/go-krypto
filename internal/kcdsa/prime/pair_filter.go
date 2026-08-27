package prime

import (
	"math/big"
	"math/bits"
	"slices"

	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

type PairFilter struct {
	bad [len(smallPrimes)]uint16
}

func NewPairFilter(J *big.Int) *PairFilter {
	f := new(PairFilter)

	twoJ := new(big.Int).Lsh(J, 1)
	defer zeroize.Zeroize(twoJ)
	for idx, prime := range smallPrimes[:] {
		a := modUint(twoJ, prime)
		if a == 0 {
			f.bad[idx] = prime
			continue
		}
		f.bad[idx] = (prime - modInv(a, prime)) % prime
	}
	return f
}

func (f *PairFilter) MaybePrimePair(q *big.Int) bool {
	for idx, prime := range smallPrimes[:] {
		qr := modUint(q, prime)
		if qr == 0 || qr == f.bad[idx] {
			return false
		}
	}
	return true
}

func modUint(x *big.Int, m uint16) uint16 {
	d := uint64(m)
	var r uint64
	w := x.Bits()
	for _, v := range slices.Backward(w) {
		v := uint64(v)
		if bits.UintSize == 64 {
			r = (r<<32 | v>>32) % d
		}
		r = (r<<32 | v&0xffffffff) % d
	}
	return uint16(r)
}

func modInv(a, m uint16) uint16 {
	t, newT := int64(0), int64(1)
	r, newR := int64(m), int64(a%m)
	for newR != 0 {
		q := r / newR
		t, newT = newT, t-q*newT
		r, newR = newR, r-q*newR
	}
	if t < 0 {
		t += int64(m)
	}
	return uint16(t)
}
