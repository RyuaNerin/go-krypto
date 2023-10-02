package lsh512

import (
	"hash"
	"testing"
)

func Test_ShortWrite(t *testing.T) {
	testSize(
		t,
		func(t *testing.T, size int) {
			var h hash.Hash
			switch size {
			case Size:
				h = New()
			case Size384:
				h = New384()
			case Size256:
				h = New256()
			case Size224:
				h = New224()
			}

			buf := make([]byte, 8*1024)
			for i := 1; i < 8*1024; i++ {
				rnd.Read(buf[:i])
				n, err := h.Write(buf[:i])
				if err != nil {
					t.Error(err)
				}
				if n != i {
					t.Fail()
				}

				rnd.Read(buf[:1])
				switch buf[0] % 5 {
				case 0:
					h.Reset()
				case 1:
					h.Sum(buf[:0])
				}
			}
		},
	)
}
