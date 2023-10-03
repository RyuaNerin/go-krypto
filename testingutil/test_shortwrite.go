package testingutil

import (
	"hash"
	"testing"
)

func TestShortWrite(t *testing.T, h hash.Hash) {
	buf := make([]byte, shortWriteSize)
	for i := 1; i < shortWriteSize; i++ {
		rnd.Read(buf[:i])
		n, err := h.Write(buf[:i])
		if err != nil {
			t.Error(err)
		}
		if n != i {
			t.Fail()
		}

		rnd.Read(buf[:1])
		if buf[0]%5 == 0 {
			h.Reset()
		}
	}
}
