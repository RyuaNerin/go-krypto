package kipher

import (
	"bytes"
	"crypto/cipher"
	"encoding/csv"
	"encoding/hex"
	"io"
	"kipher/aria"
	"kipher/hight"
	"kipher/lea"
	"kipher/seed"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"unsafe"
)

const (
	MaxBlockSize = 32
)

func decode(s string) []byte {
	b, _ := hex.DecodeString(strings.ReplaceAll(s, " ", ""))
	return b
}

var (
	regKey = regexp.MustCompilePOSIX("^KEY ?= ?([0-9A-Fa-f]+)$")
	regPT  = regexp.MustCompilePOSIX("^PT ?= ?([0-9A-Fa-f]+)$")
	regCT  = regexp.MustCompilePOSIX("^CT ?= ?([0-9A-Fa-f]+)$")
)

func s2b(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return b
}

func test(t *testing.T, do func(dst, src []byte), path string, count int, key, pt, ct, buf []byte, encryption bool) bool {
	input := pt
	answer := ct
	mode := "Encryption"
	if !encryption {
		input, answer = answer, input
		mode = "Decription"
	}

	do(buf, input)

	if !bytes.Equal(buf, answer) {
		t.Errorf(
			`%s
[%s]
COUNT : %d
KEY   : %s
PT    : %s
CT    : %s

Test  : %s
Want  : %s`,
			path,
			mode,
			count-1,
			hex.EncodeToString(key),
			hex.EncodeToString(pt),
			hex.EncodeToString(ct),
			hex.EncodeToString(buf),
			hex.EncodeToString(answer),
		)

		return false
	}

	return true
}

func testCipher(t *testing.T, path string, newCipher func(key []byte) (cipher.Block, error)) {
	var keyLen, ptLen, ctLen int
	key := make([]byte, 256/8)
	pt := make([]byte, 64)
	ct := make([]byte, 64)

	buf := make([]byte, 64)

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

		keyLen, err = hex.Decode(key, s2b(record[0]))
		if err != nil {
			t.Error(path, err)
			return
		}

		ptLen, err = hex.Decode(pt, s2b(record[1]))
		if err != nil {
			t.Error(path, err)
			return
		}

		ctLen, err = hex.Decode(ct, s2b(record[2]))
		if err != nil {
			t.Error(path, err)
			return
		}

		count++

		b, err := newCipher(key[:keyLen])
		if err != nil {
			t.Error(path, err)
			return
		}

		if !test(t, b.Encrypt, path, count, key[:keyLen], pt[:ptLen], ct[:ctLen], buf[:ptLen], true) {
			return
		}
		if !test(t, b.Decrypt, path, count, key[:keyLen], pt[:ptLen], ct[:ctLen], buf[:ptLen], false) {
			return
		}
	}
}

func TestARIA128(t *testing.T) {
	testCipher(t, "test/ARIA128.csv", func(key []byte) (cipher.Block, error) { return aria.NewCipher(key) })
}
func TestARIA192(t *testing.T) {
	testCipher(t, "test/ARIA192.csv", func(key []byte) (cipher.Block, error) { return aria.NewCipher(key) })
}
func TestARIA256(t *testing.T) {
	testCipher(t, "test/ARIA256.csv", func(key []byte) (cipher.Block, error) { return aria.NewCipher(key) })
}

func TestLEA128(t *testing.T) {
	testCipher(t, "test/LEA128.csv", func(key []byte) (cipher.Block, error) { return lea.NewCipher(key) })
}
func TestLEA192(t *testing.T) {
	testCipher(t, "test/LEA192.csv", func(key []byte) (cipher.Block, error) { return lea.NewCipher(key) })
}
func TestLEA256(t *testing.T) {
	testCipher(t, "test/LEA256.csv", func(key []byte) (cipher.Block, error) { return lea.NewCipher(key) })
}

func TestSEED128(t *testing.T) {
	testCipher(t, "test/SEED128.csv", func(key []byte) (cipher.Block, error) { return seed.NewCipher(key) })
}

func TestHIGHT(t *testing.T) {
	testCipher(t, "test/HIGHT.csv", func(key []byte) (cipher.Block, error) { return hight.NewCipher(key) })
}
