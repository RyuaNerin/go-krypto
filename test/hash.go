package test

import (
	"bytes"
	"encoding/hex"
	"hash"
	"io"
	"testing"
)

type HashTestFunc func(t *testing.T, path string, newHash func() hash.Hash)

func HashTest(t *testing.T, path string, newHash func() hash.Hash) {
	r, err := newHashTestCaseReader(path)
	if err != nil {
		t.Error(err)
		return
	}
	defer r.Close()

	for {
		_, _, msg, md, err := r.Next(true)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Error(err)
			t.FailNow()
			return
		}

		h := newHash()

		h.Write(msg)
		out := h.Sum(nil)

		if !bytes.Equal(out, md) {
			t.Errorf(
				`%s
MSG   : %s
MD    : %s
	
Eval  : %s`,
				path,
				hex.EncodeToString(msg),
				hex.EncodeToString(md),
				hex.EncodeToString(out),
			)
			t.FailNow()

			return
		}
	}
}
