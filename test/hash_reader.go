package test

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"io"
	"os"
	"strconv"
)

type hashTestCaseReader struct {
	fs *os.File
	br *bufio.Reader

	count      int
	len_       int
	seed       byteBuf
	md         byteBuf
	msg        byteBuf
	msgReadBuf []byte
}

var (
	hashTestCaseReaderSelector = selector{
		`Seed = `: func(line string, r interface{}) (err error) {
			return r.(*hashTestCaseReader).seed.parseHex(line)
		},
		`Len = `: func(line string, r interface{}) (err error) {
			r.(*hashTestCaseReader).len_, err = strconv.Atoi(line)
			return err
		},
		`COUNT = `: func(line string, r interface{}) (err error) {
			r.(*hashTestCaseReader).count, err = strconv.Atoi(line)
			return err
		},
		`MD = `: func(line string, r interface{}) (err error) {
			return r.(*hashTestCaseReader).md.parseHex(line)
		},
	}
)

func newHashTestCaseReader(path string) (r *hashTestCaseReader, err error) {
	fs, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	br := bufio.NewReader(fs)

	return &hashTestCaseReader{
		fs: fs,
		br: br,
	}, nil
}

func (tc *hashTestCaseReader) Close() {
	tc.fs.Close()
}

func (r *hashTestCaseReader) Next(needMsg bool) ([]byte, int, []byte, []byte, error) {
	var msgLen int = -1
	var line []byte
	var n int
	var err error
	for {
		line, err = r.br.Peek(8)
		if err == io.EOF {
			return nil, 0, nil, nil, err
		}
		if err != nil {
			return nil, 0, nil, nil, err
		}

		if bytes.HasPrefix(line, prefixMsg) {
			r.br.Discard(len(prefixMsg))

			l := msgLen / 8
			switch {
			case msgLen == -1:
				line, _, err = r.br.ReadLine()
				if err != nil {
					return nil, 0, nil, nil, err
				}
				err = r.msg.parseHex(b2s(line))

			case l == 0:
				if len(r.msg.buf) < 1 {
					r.msg.buf = make([]byte, 1)
				}
				r.msg.data = r.msg.buf[:0]
				line, _, err = r.br.ReadLine()
				if err != nil {
					return nil, 0, nil, nil, err
				}

			case l > 0:
				if len(r.msg.buf) < l {
					r.msg.buf = make([]byte, l)
				}
				if len(r.msgReadBuf) < l*2 {
					r.msgReadBuf = make([]byte, l*2)
				}

				offset := 0
				remain := l * 2
				for remain > 0 {
					n, err = r.br.Read(r.msgReadBuf[offset : offset+remain])

					offset += n
					remain -= n

					if err != nil {
						return nil, 0, nil, nil, err
					}
				}

				line, _, err = r.br.ReadLine()
				if err != nil {
					return nil, 0, nil, nil, err
				}

				_, err = hex.Decode(r.msg.buf, r.msgReadBuf[:l*2])
				r.msg.data = r.msg.buf[:l]
			}
		} else if len(line) > 0 {
			line, _, err = r.br.ReadLine()
			if err != nil {
				return nil, 0, nil, nil, err
			}

			err = blockTestCaseReaderSelector.Select(b2s(line), r)
			if err != nil {
				return nil, 0, nil, nil, err
			}
		} else {
			_, _, err = r.br.ReadLine()
			if err != nil {
				return nil, 0, nil, nil, err
			}
		}

		if len(line) == 0 && ((needMsg && len(r.msg.data) > 0) || (!needMsg && r.count >= 0)) && len(r.md.data) > 0 {
			return r.seed.data, r.count, r.msg.data, r.md.data, nil
		}
	}
}
