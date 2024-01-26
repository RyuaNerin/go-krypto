//go:build exclude
// +build exclude

package main

import (
	"bufio"
	"io"
	"os"
	"text/template"
)

var w io.Writer

func main() {
	bw := bufio.NewWriter(os.Stdout)
	w = bw

	bw.WriteString("---------- encrypt ----------")
	encrypt()
	bw.WriteString("\n")

	bw.WriteString("---------- decrypt ----------")
	decrypt()
	bw.WriteString("\n")

	bw.Flush()
}

var (
	m = template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}

	HIGHT_ENC_tmpl = template.Must(template.New("").Funcs(m).Parse(`
// HIGHT_ENC({{.k}},  {{.i0}},{{.i1}},{{.i2}},{{.i3}},{{.i4}},{{.i5}},{{.i6}},{{.i7}})
{
	XX[{{.i0}}] = (XX[{{.i0}}] ^ (hight_F0[XX[{{.i1}}]] + s.pdwRoundKey[4*{{.k}}+3]))
	XX[{{.i2}}] = (XX[{{.i2}}] + (hight_F1[XX[{{.i3}}]] ^ s.pdwRoundKey[4*{{.k}}+2]))
	XX[{{.i4}}] = (XX[{{.i4}}] ^ (hight_F0[XX[{{.i5}}]] + s.pdwRoundKey[4*{{.k}}+1]))
	XX[{{.i6}}] = (XX[{{.i6}}] + (hight_F1[XX[{{.i7}}]] ^ s.pdwRoundKey[4*{{.k}}+0]))
}`))
	HIGHT_DEC_tmpl = template.Must(template.New("").Funcs(m).Parse(`
// HIGHT_DEC({{.k}},  {{.i0}},{{.i1}},{{.i2}},{{.i3}},{{.i4}},{{.i5}},{{.i6}},{{.i7}})
{
	XX[{{.i1}}] = (XX[{{.i1}}] - (hight_F1[XX[{{.i2}}]] ^ s.pdwRoundKey[4*{{.k}}+2]))
	XX[{{.i3}}] = (XX[{{.i3}}] ^ (hight_F0[XX[{{.i4}}]] + s.pdwRoundKey[4*{{.k}}+1]))
	XX[{{.i5}}] = (XX[{{.i5}}] - (hight_F1[XX[{{.i6}}]] ^ s.pdwRoundKey[4*{{.k}}+0]))
	XX[{{.i7}}] = (XX[{{.i7}}] ^ (hight_F0[XX[{{.i0}}]] + s.pdwRoundKey[4*{{.k}}+3]))
}`))

	L0, L1, R0, R1 = "L0", "L1", "R0", "R1"
	A, B, C, D     = "A", "B", "C", "D"
)

func HIGHT_ENC(k, i0, i1, i2, i3, i4, i5, i6, i7 int) {
	HIGHT_ENC_tmpl.Execute(w, map[string]any{
		"k":  k,
		"i0": i0,
		"i1": i1,
		"i2": i2,
		"i3": i3,
		"i4": i4,
		"i5": i5,
		"i6": i6,
		"i7": i7,
	})
}
func HIGHT_DEC(k, i0, i1, i2, i3, i4, i5, i6, i7 int) {
	HIGHT_DEC_tmpl.Execute(w, map[string]any{
		"k":  k,
		"i0": i0,
		"i1": i1,
		"i2": i2,
		"i3": i3,
		"i4": i4,
		"i5": i5,
		"i6": i6,
		"i7": i7,
	})
}

func encrypt() {
	HIGHT_ENC(2, 7, 6, 5, 4, 3, 2, 1, 0)
	HIGHT_ENC(3, 6, 5, 4, 3, 2, 1, 0, 7)
	HIGHT_ENC(4, 5, 4, 3, 2, 1, 0, 7, 6)
	HIGHT_ENC(5, 4, 3, 2, 1, 0, 7, 6, 5)
	HIGHT_ENC(6, 3, 2, 1, 0, 7, 6, 5, 4)
	HIGHT_ENC(7, 2, 1, 0, 7, 6, 5, 4, 3)
	HIGHT_ENC(8, 1, 0, 7, 6, 5, 4, 3, 2)
	HIGHT_ENC(9, 0, 7, 6, 5, 4, 3, 2, 1)
	HIGHT_ENC(10, 7, 6, 5, 4, 3, 2, 1, 0)
	HIGHT_ENC(11, 6, 5, 4, 3, 2, 1, 0, 7)
	HIGHT_ENC(12, 5, 4, 3, 2, 1, 0, 7, 6)
	HIGHT_ENC(13, 4, 3, 2, 1, 0, 7, 6, 5)
	HIGHT_ENC(14, 3, 2, 1, 0, 7, 6, 5, 4)
	HIGHT_ENC(15, 2, 1, 0, 7, 6, 5, 4, 3)
	HIGHT_ENC(16, 1, 0, 7, 6, 5, 4, 3, 2)
	HIGHT_ENC(17, 0, 7, 6, 5, 4, 3, 2, 1)
	HIGHT_ENC(18, 7, 6, 5, 4, 3, 2, 1, 0)
	HIGHT_ENC(19, 6, 5, 4, 3, 2, 1, 0, 7)
	HIGHT_ENC(20, 5, 4, 3, 2, 1, 0, 7, 6)
	HIGHT_ENC(21, 4, 3, 2, 1, 0, 7, 6, 5)
	HIGHT_ENC(22, 3, 2, 1, 0, 7, 6, 5, 4)
	HIGHT_ENC(23, 2, 1, 0, 7, 6, 5, 4, 3)
	HIGHT_ENC(24, 1, 0, 7, 6, 5, 4, 3, 2)
	HIGHT_ENC(25, 0, 7, 6, 5, 4, 3, 2, 1)
	HIGHT_ENC(26, 7, 6, 5, 4, 3, 2, 1, 0)
	HIGHT_ENC(27, 6, 5, 4, 3, 2, 1, 0, 7)
	HIGHT_ENC(28, 5, 4, 3, 2, 1, 0, 7, 6)
	HIGHT_ENC(29, 4, 3, 2, 1, 0, 7, 6, 5)
	HIGHT_ENC(30, 3, 2, 1, 0, 7, 6, 5, 4)
	HIGHT_ENC(31, 2, 1, 0, 7, 6, 5, 4, 3)
	HIGHT_ENC(32, 1, 0, 7, 6, 5, 4, 3, 2)
	HIGHT_ENC(33, 0, 7, 6, 5, 4, 3, 2, 1)
}

func decrypt() {
	HIGHT_DEC(33, 7, 6, 5, 4, 3, 2, 1, 0)
	HIGHT_DEC(32, 0, 7, 6, 5, 4, 3, 2, 1)
	HIGHT_DEC(31, 1, 0, 7, 6, 5, 4, 3, 2)
	HIGHT_DEC(30, 2, 1, 0, 7, 6, 5, 4, 3)
	HIGHT_DEC(29, 3, 2, 1, 0, 7, 6, 5, 4)
	HIGHT_DEC(28, 4, 3, 2, 1, 0, 7, 6, 5)
	HIGHT_DEC(27, 5, 4, 3, 2, 1, 0, 7, 6)
	HIGHT_DEC(26, 6, 5, 4, 3, 2, 1, 0, 7)
	HIGHT_DEC(25, 7, 6, 5, 4, 3, 2, 1, 0)
	HIGHT_DEC(24, 0, 7, 6, 5, 4, 3, 2, 1)
	HIGHT_DEC(23, 1, 0, 7, 6, 5, 4, 3, 2)
	HIGHT_DEC(22, 2, 1, 0, 7, 6, 5, 4, 3)
	HIGHT_DEC(21, 3, 2, 1, 0, 7, 6, 5, 4)
	HIGHT_DEC(20, 4, 3, 2, 1, 0, 7, 6, 5)
	HIGHT_DEC(19, 5, 4, 3, 2, 1, 0, 7, 6)
	HIGHT_DEC(18, 6, 5, 4, 3, 2, 1, 0, 7)
	HIGHT_DEC(17, 7, 6, 5, 4, 3, 2, 1, 0)
	HIGHT_DEC(16, 0, 7, 6, 5, 4, 3, 2, 1)
	HIGHT_DEC(15, 1, 0, 7, 6, 5, 4, 3, 2)
	HIGHT_DEC(14, 2, 1, 0, 7, 6, 5, 4, 3)
	HIGHT_DEC(13, 3, 2, 1, 0, 7, 6, 5, 4)
	HIGHT_DEC(12, 4, 3, 2, 1, 0, 7, 6, 5)
	HIGHT_DEC(11, 5, 4, 3, 2, 1, 0, 7, 6)
	HIGHT_DEC(10, 6, 5, 4, 3, 2, 1, 0, 7)
	HIGHT_DEC(9, 7, 6, 5, 4, 3, 2, 1, 0)
	HIGHT_DEC(8, 0, 7, 6, 5, 4, 3, 2, 1)
	HIGHT_DEC(7, 1, 0, 7, 6, 5, 4, 3, 2)
	HIGHT_DEC(6, 2, 1, 0, 7, 6, 5, 4, 3)
	HIGHT_DEC(5, 3, 2, 1, 0, 7, 6, 5, 4)
	HIGHT_DEC(4, 4, 3, 2, 1, 0, 7, 6, 5)
	HIGHT_DEC(3, 5, 4, 3, 2, 1, 0, 7, 6)
	HIGHT_DEC(2, 6, 5, 4, 3, 2, 1, 0, 7)
}
