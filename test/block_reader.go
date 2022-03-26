package test

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"io"
	"os"
	"strconv"
)

type blockTestCaseReader struct {
	fs *os.File
	br *bufio.Reader

	keyBuf []byte
	ivBuf  []byte
	ptBuf  []byte
	ctBuf  []byte
}

func newBlockTestCaseReader(path string) (r *blockTestCaseReader, err error) {
	fs, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	br := bufio.NewReader(fs)

	return &blockTestCaseReader{
		fs: fs,
		br: br,
	}, nil
}

func (tc *blockTestCaseReader) Close() {
	tc.fs.Close()
}

func (r *blockTestCaseReader) Next() (count int, key []byte, iv []byte, pt []byte, ct []byte, err error) {
	var line []byte
	for {
		line, _, err = r.br.ReadLine()
		if err == io.EOF {
			return
		}
		if err != nil {
			return
		}

		switch {
		case bytes.HasPrefix(line, prefixCount):
			line = bytes.TrimPrefix(line, prefixCount)
			count, err = strconv.Atoi(string(line))

		case bytes.HasPrefix(line, prefixKey):
			line = bytes.TrimPrefix(line, prefixKey)
			l := len(line) / 2
			if len(r.keyBuf) < l {
				r.keyBuf = make([]byte, l)
			}
			_, err = hex.Decode(r.keyBuf, line)
			key = r.keyBuf[:l]

		case bytes.HasPrefix(line, prefixIV):
			line = bytes.TrimPrefix(line, prefixIV)
			l := len(line) / 2
			if len(r.ivBuf) < l {
				r.ivBuf = make([]byte, l)
			}
			_, err = hex.Decode(r.ivBuf, line)
			iv = r.ivBuf[:l]

		case bytes.HasPrefix(line, prefixPT):
			line = bytes.TrimPrefix(line, prefixPT)
			l := len(line) / 2
			if len(r.ptBuf) < l {
				r.ptBuf = make([]byte, l)
			}
			_, err = hex.Decode(r.ptBuf, line)
			pt = r.ptBuf[:l]

		case bytes.HasPrefix(line, prefixCT):
			line = bytes.TrimPrefix(line, prefixCT)
			l := len(line) / 2
			if len(r.ctBuf) < l {
				r.ctBuf = make([]byte, l)
			}
			_, err = hex.Decode(r.ctBuf, line)
			ct = r.ctBuf[:l]
		}
		if err != nil {
			return
		}

		if len(line) == 0 && key != nil && iv != nil && pt != nil && ct != nil {
			return
		}
	}
}
