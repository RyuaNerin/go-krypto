package test

import (
	"encoding/hex"
	"io"
	"math/big"
	"testing"
)

type DSATestFunc func(
	t *testing.T,
	path string,
	sign func(p, q, g, x, y *big.Int, m []byte) (r, s *big.Int, err error),
	verify func(p, q, g, x, y *big.Int, m []byte, r, s *big.Int) (ok bool, err error),
)

func DSATest(
	t *testing.T,
	path string,
	sign func(p, q, g, x, y *big.Int, m []byte) (r, s *big.Int, err error),
	verify func(p, q, g, y *big.Int, m []byte, r, s *big.Int) (ok bool, err error),
) {
	reader, err := newDSATestCaseReader(path)
	if err != nil {
		t.Error(err)
		return
	}
	defer reader.Close()

	for {
		p, q, g, x, y, m, r, s, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Error(err)
			return
		}

		_, _, err = sign(p, q, g, x, y, m)
		if err != nil {
			t.Error(err)
			return
		}

		ok, err := verify(p, q, g, y, m, r, s)
		if err != nil {
			t.Error(err)
			return
		}
		if !ok {
			t.Errorf(`%s
P = %s
Q = %s
G = %s
X = %s
Y = %s
M = %s
R = %s
S = %s`,
				path,
				p.Text(16),
				q.Text(16),
				g.Text(16),
				x.Text(16),
				y.Text(16),
				hex.EncodeToString(m),
				r.Text(16),
				s.Text(16),
			)
			return
		}
	}
}
