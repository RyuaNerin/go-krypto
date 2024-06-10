// Package kbkdf implements Key Based Key Derivation Functions, as defined in TTAK.KO-12.0272, TTAK.KO-12.0333, NIST SP 800-108.
package kbkdf

import (
	"github.com/RyuaNerin/go-krypto/internal"
)

// counterSize: 0 <= counterSize <= 8
func CounterMode(dst []byte, prf PRF, key, label, context []byte, counterSize, length int) []byte {
	if counterSize < 0 || 8 < counterSize {
		panic(msgInvalidCounterSize)
	}

	out, dst := internal.SliceForAppend(dst, length)

	//  1: n ← ⎡L/h⎤
	//  2: if (n > (2r - 1)) then
	//  3:     return ERROR_FLAG
	//  4: end if
	//  5: result(0) ← ∅
	//  6: for i = 1 to n do
	//  7:     K(i) ← HMAC(KI, [i]2 ‖ Label ‖ 0x00 ‖ Context ‖ [L]2)
	//  8:     result(i) ← result(i - 1) ‖ K(i)
	//  9: end do
	// 10: KO ← leftmost(result(n), L)
	// 11: return KO

	var Lr [8]byte
	L := fillL(Lr[:], uint64(length*8))
	I := make([]byte, counterSize)

	var K []byte

	for off := 0; off < length; {
		internal.IncCtr(I)

		K = prf.Sum(K[:0], key, I, label, []byte{0}, context, L)
		copy(dst[off:], K)

		off += len(K)
	}

	return out
}

// counterSize: 0 <= counterSize <= 8
func FeedbackMode(dst []byte, prf PRF, key, label, context, iv []byte, counterSize, length int) []byte {
	if counterSize < 0 || 8 < counterSize {
		panic(msgInvalidCounterSize)
	}

	out, dst := internal.SliceForAppend(dst, length)

	//  1: n ← ⎡L/h⎤
	//  2: if (n > (232 - 1)) then
	//  3:     return ERROR_FLAG
	//  4: end if
	//  5: result(0) ← ∅
	//  6: K(0) ← ∅
	//  7: for i = 1 to n do
	//  8:     K(i) ← HMAC(KI, K(i - 1) {‖ [i]2} ‖ Label ‖ 0x00 ‖ Context ‖ [L]2)
	//  9:     result(i) ← result(i - 1) ‖ K(i)
	// 10: end do
	// 11: KO ← leftmost(result(n), L)
	// 12: return KO

	var Lr [8]byte
	L := fillL(Lr[:], uint64(length*8))
	I := make([]byte, counterSize)

	K := internal.BytesClone(iv)

	for off := 0; off < length; {
		internal.IncCtr(I)

		K = prf.Sum(K[:0], key, K, I, label, []byte{0}, context, L)
		copy(dst[off:], K)
		off += len(K)
	}

	return out
}

// counterSize: 0 <= counterSize <= 8
func DoublePipeMode(dst []byte, prf PRF, key, label, context []byte, counterSize, length int) []byte {
	if counterSize < 0 || 8 < counterSize {
		panic(msgInvalidCounterSize)
	}

	out, dst := internal.SliceForAppend(dst, length)

	//  1: n ← ⎡L/h⎤
	//  2: if (n > (232 - 1)) then
	//  3:     return ERROR_FLAG
	//  4: end if
	//  5: result(0) ← ∅
	//  6: A(0) ← Label ‖ 0x00 ‖ Context ‖ [L]2
	//  7: for i = 1 to n do
	//  8:     A(i) ← HMAC(KI, A(i - 1))
	//  9:     K(i) ← HMAC(KI, A(i) {‖ [i]2} ‖ Label ‖ 0x00 ‖ Context ‖ [L]2)
	// 10:     result(i) ← result(i - 1) ‖ K(i)
	// 11: end do
	// 12: KO ← leftmost(result(n), L)
	// 13: return KO

	var Lr [8]byte
	L := fillL(Lr[:], uint64(length*8))
	I := make([]byte, counterSize)

	// off = 0
	internal.IncCtr(I)
	A := prf.Sum(nil, key, label, []byte{0}, context, L)
	K := prf.Sum(nil, key, A, I, label, []byte{0}, context, L)
	off := copy(dst, K)

	for off < length {
		internal.IncCtr(I)

		A = prf.Sum(A[:0], key, A)
		K = prf.Sum(K[:0], key, A, I, label, []byte{0}, context, L)
		copy(dst[off:], K)

		off += len(K)
	}

	return out
}

func fillL(dst []byte, v uint64) []byte {
	switch {
	case v < 1<<8:
		dst[0] = byte(v)
		return dst[:1]

	case v < 1<<16:
		dst[0] = byte(v >> 8)
		dst[1] = byte(v)
		return dst[:2]

	case v < 1<<24:
		dst[0] = byte(v >> 16)
		dst[1] = byte(v >> 8)
		dst[2] = byte(v)
		return dst[:3]

	case v < 1<<32:
		dst[0] = byte(v >> 24)
		dst[1] = byte(v >> 16)
		dst[2] = byte(v >> 8)
		dst[3] = byte(v)
		return dst[:4]

	case v < 1<<40:
		dst[0] = byte(v >> 32)
		dst[1] = byte(v >> 24)
		dst[2] = byte(v >> 16)
		dst[3] = byte(v >> 8)
		dst[4] = byte(v)
		return dst[:5]

	case v < 1<<48:
		dst[0] = byte(v >> 40)
		dst[1] = byte(v >> 32)
		dst[2] = byte(v >> 24)
		dst[3] = byte(v >> 16)
		dst[4] = byte(v >> 8)
		dst[5] = byte(v)
		return dst[:6]

	case v < 1<<56:
		dst[0] = byte(v >> 48)
		dst[1] = byte(v >> 40)
		dst[2] = byte(v >> 32)
		dst[3] = byte(v >> 24)
		dst[4] = byte(v >> 16)
		dst[5] = byte(v >> 8)
		dst[6] = byte(v)
		return dst[:7]

	default:
		dst[0] = byte(v >> 56)
		dst[1] = byte(v >> 48)
		dst[2] = byte(v >> 40)
		dst[3] = byte(v >> 32)
		dst[4] = byte(v >> 24)
		dst[5] = byte(v >> 16)
		dst[6] = byte(v >> 8)
		dst[7] = byte(v)
		return dst[:8]
	}
}
