package has160

import "testing"

func Test_ShortWrite(t *testing.T) {
	h := New()

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
}
