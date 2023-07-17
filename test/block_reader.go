package test

import (
	"bufio"
	"io"
	"os"
	"strconv"
)

type blockTestCaseReader struct {
	fs *os.File
	br *bufio.Reader

	count int
	key   byteBuf
	iv    byteBuf
	pt    byteBuf
	ct    byteBuf
}

var (
	blockTestCaseReaderSelector = selector{
		"COUNT = ": func(line string, r interface{}) (err error) {
			r.(*blockTestCaseReader).count, err = strconv.Atoi(line)
			return err
		},
		`KEY = `: func(line string, r interface{}) (err error) {
			return r.(*blockTestCaseReader).key.parseHex(line)
		},
		`IV = `: func(line string, r interface{}) (err error) {
			return r.(*blockTestCaseReader).iv.parseHex(line)
		},
		`PT = `: func(line string, r interface{}) (err error) {
			return r.(*blockTestCaseReader).pt.parseHex(line)
		},
		`CT = `: func(line string, r interface{}) (err error) {
			return r.(*blockTestCaseReader).ct.parseHex(line)
		},
	}
)

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

func (r *blockTestCaseReader) Next() (int, []byte, []byte, []byte, []byte, error) {
	r.key.reset()
	r.iv.reset()
	r.pt.reset()
	r.ct.reset()

	var err error
	var lineRaw []byte
	for {
		lineRaw, _, err = r.br.ReadLine()
		if err == io.EOF {
			return 0, nil, nil, nil, nil, err
		}
		if err != nil {
			return 0, nil, nil, nil, nil, err
		}

		err = blockTestCaseReaderSelector.Select(b2s(lineRaw), r)
		if err != nil {
			return 0, nil, nil, nil, nil, err
		}

		if len(lineRaw) == 0 && r.key.data != nil && r.iv.data != nil && r.pt.data != nil && r.ct.data != nil {
			return r.count, r.key.data, r.iv.data, r.pt.data, r.ct.data, nil
		}
	}
}
