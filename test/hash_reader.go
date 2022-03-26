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

	seedBuf    []byte
	msgReadBuf []byte
	msgBuf     []byte
	mdBuf      []byte
}

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

func (r *hashTestCaseReader) Next(needMsg bool) (seed []byte, count int, msg []byte, md []byte, err error) {
	seed = r.seedBuf
	count = -1

	var msgLen int = -1
	var line []byte
	var n int
	for {
		line, err = r.br.Peek(8)
		if err == io.EOF {
			return
		}
		if err != nil {
			return
		}

		if bytes.HasPrefix(line, prefixMsg) {
			r.br.Discard(len(prefixMsg))

			l := msgLen / 8
			if l == 0 {
				msg = make([]byte, 0)
				line, _, err = r.br.ReadLine()
				if err == io.EOF {
					return
				}
				if err != nil {
					return
				}
			} else if l > 0 {
				if len(r.msgBuf) < l {
					r.msgBuf = make([]byte, l)
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

					if err == io.EOF {
						return
					}
					if err != nil {
						return
					}
				}

				line, _, err = r.br.ReadLine()
				if err == io.EOF {
					return
				}
				if err != nil {
					return
				}

				_, err = hex.Decode(r.msgBuf, r.msgReadBuf[:l*2])
				msg = r.msgBuf[:l]
			} else {
				line, _, err = r.br.ReadLine()
				if err == io.EOF {
					return
				}
				if err != nil {
					return
				}

				line = bytes.TrimPrefix(line, prefixMsg)
				l := len(line) / 2
				if len(r.msgBuf) < l {
					r.msgBuf = make([]byte, l)
				}
				_, err = hex.Decode(r.msgBuf, line)
				msg = r.msgBuf[:l]
			}
		} else if len(line) > 0 {
			line, _, err = r.br.ReadLine()
			if err == io.EOF {
				return
			}
			if err != nil {
				return
			}

			switch {
			case bytes.HasPrefix(line, prefixSeed):
				line = bytes.TrimPrefix(line, prefixSeed)
				l := len(line) / 2
				if len(r.seedBuf) < l {
					r.seedBuf = make([]byte, l)
				}
				_, err = hex.Decode(r.seedBuf, line)
				seed = r.seedBuf[:l]

			case bytes.HasPrefix(line, prefixLen):
				line = bytes.TrimPrefix(line, prefixLen)
				msgLen, err = strconv.Atoi(string(line))

			case bytes.HasPrefix(line, prefixCount):
				line = bytes.TrimPrefix(line, prefixCount)
				count, err = strconv.Atoi(string(line))

			case bytes.HasPrefix(line, prefixMD):
				line = bytes.TrimPrefix(line, prefixMD)
				l := len(line) / 2
				if len(r.mdBuf) < l {
					r.mdBuf = make([]byte, l)
				}
				_, err = hex.Decode(r.mdBuf, line)
				md = r.mdBuf[:l]
			}
		} else {
			_, _, err = r.br.ReadLine()
		}
		if err == io.EOF {
			return
		}
		if err != nil {
			return
		}

		if len(line) == 0 && ((needMsg && msg != nil) || (!needMsg && count >= 0)) && md != nil {
			return
		}
	}
}
