package main

import (
	"bufio"
	"encoding/hex"
	"io"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/RyuaNerin/go-krypto/internal"
)

type cavpRow struct {
	Key     string // optional
	Value   string
	Section bool
}

func (cv cavpRow) Int() int {
	v, err := strconv.Atoi(cv.Value)
	if err != nil {
		panic(err)
	}
	return v
}

func (cv cavpRow) Bool() bool {
	return strings.EqualFold(cv.Value, "true") || cv.Value == "1"
}

func (cv cavpRow) Hex() []byte {
	if len(cv.Value)%2 == 1 {
		cv.Value = "0" + cv.Value
	}
	v, err := hex.DecodeString(cv.Value)
	if err != nil {
		panic(err)
	}
	return v
}

func (cv cavpRow) BigInt() *big.Int {
	return new(big.Int).SetBytes(cv.Hex())
}

type cavpSection []cavpRow

func (cv cavpSection) ContainsKey(key string) bool {
	for _, v := range cv {
		if strings.EqualFold(v.Key, key) {
			return true
		}
	}
	return false
}

func (cv cavpSection) ContainsValue(value string) bool {
	for _, v := range cv {
		if strings.EqualFold(v.Value, value) {
			return true
		}
	}
	return false
}

func (cv cavpSection) Int(name string) int {
	for _, v := range cv {
		if strings.EqualFold(v.Key, name) {
			return v.Int()
		}
	}
	panic("not found")
}

func (cv cavpSection) Bool(name string) bool {
	for _, v := range cv {
		if strings.EqualFold(v.Key, name) {
			return v.Bool()
		}
	}
	panic("not found")
}

func (cv cavpSection) Value(name string) string {
	for _, v := range cv {
		if strings.EqualFold(v.Key, name) {
			return v.Value
		}
	}
	panic("not found")
}

func (cv cavpSection) Hex(name string) []byte {
	for _, v := range cv {
		if strings.EqualFold(v.Key, name) {
			return v.Hex()
		}
	}
	panic("not found")
}

func (cv cavpSection) HexList(name string) (lst [][]byte) {
	for _, v := range cv {
		if strings.EqualFold(v.Key, name) {
			lst = append(lst, v.Hex())
		}
	}
	if len(lst) == 0 {
		panic("not found")
	}
	return lst
}

func (cv cavpSection) BigInt(name string) *big.Int {
	for _, v := range cv {
		if v.Key == name {
			return v.BigInt()
		}
	}
	panic("not found")
}

type cavpProcessor struct {
	r *os.File
	w *os.File

	filename string

	br      *bufio.Reader
	bw      *bufio.Writer
	lineBuf strings.Builder

	lineMax int
	lineCur int
	eof     bool

	cs cavpSection

	emptyLine bool
}

func NewCavp(path, filename string) *cavpProcessor {
	var err error
	p := &cavpProcessor{
		filename: filename,
	}

	////////////////////////////////////////////////////////////
	// Input

	if p.r, err = os.Open(path); err != nil {
		panic(err)
	}

	////////////////////////////////////////////////////////////
	// Output

	if p.w, err = os.Create(filepath.Join(filepath.Dir(path), strings.ReplaceAll(filepath.Base(path), ".req", ".rsp"))); err != nil {
		panic(err)
	}
	if err = p.w.Truncate(0); err != nil {
		panic(err)
	}
	if _, err = p.w.Seek(0, io.SeekStart); err != nil {
		panic(err)
	}

	////////////////////////////////////////////////////////////

	p.br = bufio.NewReader(p.r)
	p.bw = bufio.NewWriter(p.w)
	p.lineBuf.Grow(4096)

	////////////////////////////////////////////////////////////

	eof := false
	for !eof {
		_, prefix, err := p.br.ReadLine()
		if err == io.EOF {
			prefix = false
			eof = true
			err = nil
		}
		if err != nil {
			panic(err)
		}

		if prefix {
			continue
		}

		p.lineMax++
	}

	if _, err = p.r.Seek(0, io.SeekStart); err != nil {
		panic(err)
	}
	p.br.Reset(p.r)

	return p
}

func (p *cavpProcessor) Close() {
	p.bw.Flush()
	p.w.Close()
	p.r.Close()
}

func (p *cavpProcessor) ReadLine() string {
	for !p.eof {
		lineRaw, prefix, err := p.br.ReadLine()
		if err == io.EOF {
			p.eof = true
			prefix = false
			err = nil
		}
		if err != nil {
			panic(err)
		}
		p.lineBuf.Write(lineRaw)

		if prefix {
			continue
		}
		break
	}

	s := internal.StringClone(p.lineBuf.String())
	p.lineBuf.Reset()

	p.lineCur++

	p.verbose()
	return s
}

func (p *cavpProcessor) WriteLine(s string) {
	if s == "" {
		if p.emptyLine {
			p.emptyLine = false
			return
		}
		p.emptyLine = true
	} else {
		p.emptyLine = false
	}
	p.bw.WriteString(s)
	p.bw.WriteByte('\n')
}

func (p *cavpProcessor) Next() bool {
	p.cs = p.cs[:0]

	for {
		for !p.eof {
			lineRaw, prefix, err := p.br.ReadLine()
			if err == io.EOF {
				p.eof = true
				prefix = false
				err = nil
			}
			if err != nil {
				panic(err)
			}
			p.lineBuf.Write(lineRaw)

			if prefix {
				continue
			}

			line := strings.TrimSpace(p.lineBuf.String())
			p.lineBuf.Reset()

			p.lineCur++
			if len(line) == 0 {
				break
			}

			if strings.HasPrefix(line, "#") {
				p.cs = append(p.cs, cavpRow{"", internal.StringClone(line), false})
			} else {
				idx := strings.Index(line, "=")
				if idx == -1 {
					p.cs = append(p.cs, cavpRow{"", internal.StringClone(line), false})
				} else {
					bo := strings.HasPrefix(line, "[")
					bc := strings.HasSuffix(line, "]")

					var section bool
					switch {
					case bo && bc:
						section = true
						line = strings.TrimSuffix(strings.TrimPrefix(line, "["), "]")
						idx = strings.Index(line, "=")
						fallthrough
					case !bo && !bc:
						p.cs = append(
							p.cs,
							cavpRow{
								Key:     internal.StringClone(strings.TrimSpace(line[:idx])),
								Value:   internal.StringClone(strings.TrimSpace(line[idx+1:])),
								Section: section,
							},
						)

					case !bo && bc:
						fallthrough
					case bo && !bc:
						p.cs = append(p.cs, cavpRow{"", internal.StringClone(line), false})
					}
				}
			}
		}

		p.verbose()

		if p.eof || len(p.cs) > 0 {
			break
		}
	}

	return len(p.cs) > 0
}

func (p *cavpProcessor) ReadValues() cavpSection {
	return p.cs
}

func (p *cavpProcessor) WriteValues(lst cavpSection) {
	if len(lst) == 0 {
		return
	}

	for _, v := range lst {
		if v.Section {
			p.bw.WriteByte('[')
		}
		if v.Key != "" {
			p.bw.WriteString(v.Key)
			p.bw.WriteString(" = ")
			p.bw.WriteString(v.Value)
		} else if v.Value != "" {
			p.bw.WriteString(v.Value)
		}
		if v.Section {
			p.bw.WriteByte(']')
		}
		p.bw.WriteByte('\n')
	}
	p.bw.WriteByte('\n')
	p.bw.Flush()
}

func (p *cavpProcessor) verbose() {
	log.Printf("        [%d / %d] %s\n", p.lineCur, p.lineMax, p.filename)
}
