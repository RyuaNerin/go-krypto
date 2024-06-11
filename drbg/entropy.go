package drbg

import (
	"io"
)

type Entropy interface {
	Get() ([]byte, error)
}

func NewEntropy(rand io.Reader, length int) Entropy {
	return newEntropy(rand, length)
}

type entropy struct {
	rand io.Reader
	buf  []byte
}

func newEntropy(rand io.Reader, length int) *entropy {
	e := &entropy{
		rand: rand,
		buf:  make([]byte, length),
	}

	return e
}

func (e *entropy) Get() (p []byte, err error) {
	_, err = io.ReadFull(e.rand, e.buf)
	if err != nil {
		return nil, err
	}

	return e.buf, nil
}
