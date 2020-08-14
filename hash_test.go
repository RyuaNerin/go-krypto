package kipher

import (
	"bytes"
	"encoding/csv"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"testing"

	"kipher/lsh256"
	"kipher/lsh512"
)

func test_hash(t *testing.T, path string, newHash func() hash.Hash) {
	var outLen, inLen int
	out := make([]byte, 64)
	in := make([]byte, 1024)

	fs, err := os.Open(path)
	if err != nil {
		t.Error(err)
		return
	}
	defer fs.Close()

	cr := csv.NewReader(fs)

	count := 0

	for {
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Error(path, err)
			return
		}

		outLen, err = hex.Decode(out, s2b(record[0]))
		if err != nil {
			t.Error(path, err)
			return
		}

		inLen, err = hex.Decode(in, s2b(record[1]))
		if err != nil {
			t.Error(path, err)
			return
		}

		count++

		h := newHash()
		h.Write(in[:inLen])

		buf := h.Sum(nil)

		if !bytes.Equal(buf, out[:outLen]) {
			t.Errorf(
				`%s
COUNT  : %d
INPUT
%s

Test   : %s
Want   : %s`,
				path,
				count-1,
				hex.Dump(in[:inLen]),
				hex.EncodeToString(buf),
				hex.EncodeToString(out[:outLen]),
			)

			return
		}
	}
}

func TestLSH256(t *testing.T) {
	test_hash(t, "test/LSH256.csv", func() hash.Hash { return lsh256.New() })
}
func TestLSH512(t *testing.T) {
	test_hash(t, "test/LSH512.csv", func() hash.Hash { return lsh512.New() })
}
