package internal

import (
	"encoding/binary"
	"io"
	"math/bits"
)

func SliceForAppend(in []byte, n int) (head, tail []byte) {
	if total := len(in) + n; cap(in) >= total {
		head = in[:total]
	} else {
		head = make([]byte, total)
		copy(head, in)
	}
	tail = head[len(in):]
	return
}

// without guarantee of data
func Grow(buf []byte, n int) []byte {
	if n < cap(buf) {
		if rem := len(buf) - n; rem > 0 {
			clear(buf[len(buf)-rem:])
		}
		return buf[:n]
	}
	clear(buf)

	return make([]byte, n)
}

// with guarantee of data
func IncreaseCapacity(src []byte, n int) []byte {
	if len(src)+n <= cap(src) {
		if rem := len(src) + n - len(src); rem > 0 {
			clear(src[len(src):])
		}
		return src
	}

	dst := make([]byte, len(src)+n)
	copy(dst, src)
	return dst[:len(src)]
}

func ResizeAndCopy(dst, src []byte) []byte {
	if len(src) <= cap(dst) {
		dst = dst[:len(src)]
		copy(dst, src)
		return dst
	}

	dst = make([]byte, len(src))
	copy(dst, src)
	return dst
}

// keep data
func Resize(arr []byte, n int) []byte {
	arrLen := len(arr)

	switch {
	case n == arrLen:
	case n <= cap(arr):
		if n > arrLen {
			clear(arr[arrLen:])
		}
		arr = arr[:n]

	default:
		arr2 := make([]byte, n)
		copy(arr2, arr)
		arr = arr2
	}
	return arr
}

// 0 0[0 0 0 0 0 0]
func RightMost(b []byte, bits int) []byte {
	n := BitsToBytes(bits)
	b = b[len(b)-n:]

	remain := bits % 8
	if remain > 0 {
		b[0] &= ((1 << remain) - 1)
	}

	return b
}

// [0 0 0 0 0 0]0 0
func LeftMost(b []byte, bits int) []byte {
	n := BitsToBytes(bits)
	b = b[:n]

	remain := bits % 8
	if remain > 0 {
		b[n-1] &= byte(0b_11111111 << (8 - remain))
	}

	return b
}

func Add(dst []byte, src ...[]byte) {
	if len(src) == 0 {
		clear(dst)
		return
	}

	v := src[0]
	if len(v) >= len(dst) {
		copy(dst, v[len(v)-len(dst):])
	} else {
		k := len(dst) - len(v)
		clear(dst[:k])
		copy(dst[k:], v)
	}

	for _, v := range src[1:] {
		addInto(dst, v)
	}
}

func addInto(dst, v []byte) {
	i, j := len(dst), len(v)
	if j > i {
		v = v[j-i:] // 초과분 선제 절단
		j = i
	}

	var carry uint64
	for j >= 8 {
		i, j = i-8, j-8
		a := dst[i : i+8 : i+8]
		b := v[j : j+8 : j+8]
		s, c := bits.Add64(
			binary.BigEndian.Uint64(a),
			binary.BigEndian.Uint64(b),
			carry,
		)
		binary.BigEndian.PutUint64(a, s)
		carry = c
	}

	for j > 0 {
		i, j = i-1, j-1
		s := uint64(dst[i]) + uint64(v[j]) + carry
		dst[i] = byte(s)
		carry = s >> 8
	}

	for i > 0 && carry != 0 {
		i--
		s := uint64(dst[i]) + carry
		dst[i] = byte(s)
		carry = s >> 8
	}
}

const is64bit = ^uint(0)>>63 != 0

func incCtr2(b []byte) bool {
	v := binary.BigEndian.Uint16(b) + 1
	binary.BigEndian.PutUint16(b, v)
	return v == 0
}

func incCtr4(b []byte) bool {
	v := binary.BigEndian.Uint32(b) + 1
	binary.BigEndian.PutUint32(b, v)
	return v == 0
}

func incCtr8(b []byte) bool {
	if is64bit {
		v := binary.BigEndian.Uint64(b) + 1
		binary.BigEndian.PutUint64(b, v)
		return v == 0
	}
	return incCtr4(b[4:8]) && incCtr4(b[0:4])
}

func IncCtr(b []byte) {
	switch len(b) {
	case 1:
		b[0]++
	case 2:
		incCtr2(b)
	case 3:
		if incCtr2(b[1:]) {
			b[0]++
		}
	case 4:
		incCtr4(b)
	case 5:
		if incCtr4(b[1:]) {
			b[0]++
		}
	case 6:
		if incCtr4(b[2:]) {
			incCtr2(b)
		}
	case 7:
		if incCtr4(b[3:]) && incCtr2(b[1:]) {
			b[0]++
		}
	case 8:
		incCtr8(b)
	case 9:
		if incCtr8(b[1:]) {
			b[0]++
		}
	case 10:
		if incCtr8(b[2:]) {
			incCtr2(b)
		}
	case 11:
		if incCtr8(b[3:]) && incCtr2(b[1:]) {
			b[0]++
		}
	case 12:
		if incCtr8(b[4:]) {
			incCtr4(b)
		}
	default:
		i := len(b)
		if is64bit {
			for i >= 8 {
				i -= 8
				if !incCtr8(b[i:]) {
					return
				}
			}
		} else {
			for i >= 4 {
				i -= 4
				if !incCtr4(b[i:]) {
					return
				}
			}
		}
		for i--; i >= 0; i-- {
			b[i]++
			if b[i] != 0 {
				return
			}
		}
	}
}

func FillCtr(dst, ctr []byte) {
	switch len(ctr) {
	case 4:
		v := binary.BigEndian.Uint32(ctr)
		for len(dst) >= 4 {
			v++
			binary.BigEndian.PutUint32(dst, v)
			dst = dst[4:]
		}
		if len(dst) != 0 {
			panic(msgCounterOutputNotMultipleOfSize)
		}
		binary.BigEndian.PutUint32(ctr, v)

	case 8:
		v := binary.BigEndian.Uint64(ctr)
		for len(dst) >= 8 {
			v++
			binary.BigEndian.PutUint64(dst, v)
			dst = dst[8:]
		}
		if len(dst) != 0 {
			panic(msgCounterOutputNotMultipleOfSize)
		}
		binary.BigEndian.PutUint64(ctr, v)

	case 12:
		hi := binary.BigEndian.Uint64(ctr)
		lo := binary.BigEndian.Uint32(ctr[8:])
		for len(dst) >= 12 {
			var c uint32
			lo, c = bits.Add32(lo, 1, 0)
			hi += uint64(c)
			binary.BigEndian.PutUint64(dst[0:], hi)
			binary.BigEndian.PutUint32(dst[8:], lo)
			dst = dst[12:]
		}
		if len(dst) != 0 {
			panic(msgCounterOutputNotMultipleOfSize)
		}
		binary.BigEndian.PutUint64(ctr, hi)
		binary.BigEndian.PutUint32(ctr[8:], lo)

	case 16:
		hi := binary.BigEndian.Uint64(ctr)
		lo := binary.BigEndian.Uint64(ctr[8:])
		for len(dst) >= 16 {
			var c uint64
			lo, c = bits.Add64(lo, 1, 0)
			hi += c
			binary.BigEndian.PutUint64(dst[0:], hi)
			binary.BigEndian.PutUint64(dst[8:], lo)
			dst = dst[16:]
		}
		if len(dst) != 0 {
			panic(msgCounterOutputNotMultipleOfSize)
		}
		binary.BigEndian.PutUint64(ctr, hi)
		binary.BigEndian.PutUint64(ctr[8:], lo)

	default:
		n := len(ctr)
		if n == 0 {
			panic(msgEmptyCounter)
		}
		for len(dst) >= n {
			IncCtr(ctr)
			copy(dst, ctr)
			dst = dst[n:]
		}
		if len(dst) != 0 {
			panic(msgCounterOutputNotMultipleOfSize)
		}
	}
}

// resize dst, ReadFull, cut from right
func ReadBits(dst []byte, rand io.Reader, bits int) ([]byte, error) {
	n := BitsToBytes(bits)

	dst = Grow(dst, n)

	if _, err := io.ReadFull(rand, dst); err != nil {
		return dst, err
	}

	if rem := bits & 0x07; rem != 0 {
		dst[0] &= byte(1<<rem) - 1
	}

	return dst, nil
}

// resize dst, ReadFull, cut from right
func ReadBytes(dst []byte, rand io.Reader, n int) ([]byte, error) {
	dst = Grow(dst, n)

	if _, err := io.ReadFull(rand, dst); err != nil {
		return dst, err
	}

	return dst, nil
}
