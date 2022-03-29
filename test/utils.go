package test

import (
	"encoding/hex"
	"reflect"
	"strings"
	"unsafe"
)

func b2s(b []byte) (s string) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sh.Data = bh.Data
	sh.Len = bh.Len
	return s
}

func s2b(s string) (b []byte) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return b
}

type selector map[string]func(line string, r interface{}) (err error)

func (s selector) Select(line string, r interface{}) error {
	for prefix, f := range s {
		if strings.HasPrefix(line, prefix) {
			return f(strings.TrimPrefix(line, prefix), r)
		}
	}
	return nil
}

type byteBuf struct {
	data []byte
	buf  []byte
}

func (buf *byteBuf) reset() {
	buf.data = nil
}

func (buf *byteBuf) parseHex(line string) error {
	l := len(line) / 2
	if len(buf.buf) < l {
		buf.buf = make([]byte, l)
	}
	_, err := hex.Decode(buf.buf, s2b(line))
	buf.data = buf.buf[:l]
	return err
}
