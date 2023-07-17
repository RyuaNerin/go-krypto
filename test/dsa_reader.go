package test

import (
	"bufio"
	"errors"
	"math/big"
	"os"
)

type dsaTestCaseReader struct {
	fs *os.File
	br *bufio.Reader

	p *big.Int
	q *big.Int
	g *big.Int
	m byteBuf
	x *big.Int
	y *big.Int
	r *big.Int
	s *big.Int
}

var (
	dsaTestCaseReaderSelector = selector{
		`P = `: func(line string, r interface{}) (err error) {
			_, ok := r.(*dsaTestCaseReader).p.SetString(line, 16)
			if !ok {
				return errors.New("invalid line")
			}
			return nil
		},
		`Q = `: func(line string, r interface{}) (err error) {
			_, ok := r.(*dsaTestCaseReader).q.SetString(line, 16)
			if !ok {
				return errors.New("invalid line")
			}
			return nil
		},
		`G = `: func(line string, r interface{}) (err error) {
			_, ok := r.(*dsaTestCaseReader).g.SetString(line, 16)
			if !ok {
				return errors.New("invalid line")
			}
			return nil
		},
		`M = `: func(line string, r interface{}) (err error) {
			return r.(*dsaTestCaseReader).m.parseHex(line)
		},
		`X = `: func(line string, r interface{}) (err error) {
			_, ok := r.(*dsaTestCaseReader).x.SetString(line, 16)
			if !ok {
				return errors.New("invalid line")
			}
			return nil
		},
		`Y = `: func(line string, r interface{}) (err error) {
			_, ok := r.(*dsaTestCaseReader).y.SetString(line, 16)
			if !ok {
				return errors.New("invalid line")
			}
			return nil
		},
		`R = `: func(line string, r interface{}) (err error) {
			_, ok := r.(*dsaTestCaseReader).r.SetString(line, 16)
			if !ok {
				return errors.New("invalid line")
			}
			return nil
		},
		`S = `: func(line string, r interface{}) (err error) {
			_, ok := r.(*dsaTestCaseReader).s.SetString(line, 16)
			if !ok {
				return errors.New("invalid line")
			}
			return nil
		},
	}
)

func newDSATestCaseReader(path string) (r *dsaTestCaseReader, err error) {
	fs, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	br := bufio.NewReader(fs)

	return &dsaTestCaseReader{
		fs: fs,
		br: br,
		p:  new(big.Int),
		q:  new(big.Int),
		g:  new(big.Int),
		x:  new(big.Int),
		y:  new(big.Int),
		r:  new(big.Int),
		s:  new(big.Int),
	}, nil
}

func (tc *dsaTestCaseReader) Close() {
	tc.fs.Close()
}

func (r *dsaTestCaseReader) Next() (
	*big.Int, // P
	*big.Int, // Q
	*big.Int, // R
	*big.Int, // X
	*big.Int, // Y
	[]byte, // M
	*big.Int, // R
	*big.Int, // S
	error,
) {
	for {
		line, _, err := r.br.ReadLine()
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, err
		}

		err = dsaTestCaseReaderSelector.Select(b2s(line), r)
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, err
		}

		if len(line) == 0 && len(r.m.data) > 0 {
			return r.p, r.q, r.g, r.x, r.y, r.m.data, r.r, r.s, nil
		}
	}
}
