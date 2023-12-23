package testingutil

import (
	"hash"
	"testing"
)

const shortWriteSize = 1024

// Hash Test ShortWrite
func HTSW(
	t *testing.T,
	h hash.Hash,
	skip bool,
) {
	if skip {
		t.Skip()
		return
	}

	buf := make([]byte, shortWriteSize)
	for i := 1; i < shortWriteSize; i++ {
		rnd.Read(buf[:i])
		n, err := h.Write(buf[:i])
		if err != nil {
			t.Error(err)
			return
		}
		if n != i {
			t.Fail()
			return
		}

		rnd.Read(buf[:1])
		if buf[0]%5 == 0 {
			h.Reset()
		}
	}
}

// Hash Test ShortWrite All
func HTSWA(
	t *testing.T,
	sizes []CipherSize,
	newHash func(size int) hash.Hash,
	skip bool,
) {
	if skip {
		t.Skip()
		return
	}

	TA(t, sizes, func(t *testing.T, size int) {
		h := newHash(size)
		HTSW(t, h, false)
	}, false)
}
