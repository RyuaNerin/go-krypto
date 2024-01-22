package main

import (
	"bufio"
	"io"
	"os"
	"text/template"
)

func main() {
	w := bufio.NewWriter(os.Stdout)

	w.WriteString("---------- key ----------")
	keySchedKey(w)
	w.WriteString("\n")

	w.WriteString("---------- encrypt ----------")
	encrypt(w)
	w.WriteString("\n")

	w.WriteString("---------- decrypt ----------")
	decrypt(w)
	w.WriteString("\n")

	w.Flush()
}

var (
	m = template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}

	SEED_KeySched_tmpl = template.Must(template.New("").Funcs(m).Parse(`
// SEED_KeySched({{ .L0 }}, {{ .L1 }}, {{ .R0 }}, {{ .R1 }}, K+{{ .K }})
{
	T0 = {{.R0}} ^ s.pdwRoundKey[{{add .K 0}}]
	T1 = {{.R1}} ^ s.pdwRoundKey[{{add .K 1}}]
	T1 ^= T0
	T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
	T0 = (T0 + T1) & 0xffffffff
	T0 = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
	T1 = (T1 + T0) & 0xffffffff
	T1 = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
	T0 = (T0 + T1) & 0xffffffff
	{{.L0}} ^= T0
	{{.L1}} ^= T1
}`))
	RoundKeyUpdate0_tmpl = template.Must(template.New("").Funcs(m).Parse(`
// RoundKeyUpdate0(K+{{ .K }}, A, B, C, D, KC{{ .KC }})
{
    T0 = A + C - kc[{{.KC}}]
    T1 = B + kc[{{.KC}}] - D
    s.pdwRoundKey[{{add .K 0}}] = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
    s.pdwRoundKey[{{add .K 1}}] = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
    T0 = A
    A = (A>>8) ^ (B<<24)
    B = (B>>8) ^ (T0<<24)
}`))
	RoundKeyUpdate1_tmpl = template.Must(template.New("").Funcs(m).Parse(`
// RoundKeyUpdate1(K+ {{ .K }}, A, B, C, D, KC{{ .KC }})
{
    T0 = A + C - kc[{{.KC}}]
    T1 = B + kc[{{.KC}}] - D
    s.pdwRoundKey[{{add .K 0}}] = ss0[getB0(T0)] ^ ss1[getB1(T0)] ^ ss2[getB2(T0)] ^ ss3[getB3(T0)]
    s.pdwRoundKey[{{add .K 1}}] = ss0[getB0(T1)] ^ ss1[getB1(T1)] ^ ss2[getB2(T1)] ^ ss3[getB3(T1)]
    T0 = C
    C = (C<<8) ^ (D>>24)
    D = (D<<8) ^ (T0>>24)
}`))

	L0, L1, R0, R1 = "L0", "L1", "R0", "R1"
	A, B, C, D     = "A", "B", "C", "D"
)

func SEED_KeySched(w io.Writer, L0, L1, R0, R1 string, K int) {
	SEED_KeySched_tmpl.Execute(w, map[string]any{
		"L0": L0,
		"L1": L1,
		"R0": R0,
		"R1": R1,
		"K":  K,
	})
}
func RoundKeyUpdate0(w io.Writer, K int, A, B, C, D string, KC int) {
	RoundKeyUpdate0_tmpl.Execute(w, map[string]any{
		"K":  K,
		"A":  A,
		"B":  B,
		"C":  C,
		"D":  D,
		"KC": KC,
	})
}
func RoundKeyUpdate1(w io.Writer, K int, A, B, C, D string, KC int) {
	RoundKeyUpdate1_tmpl.Execute(w, map[string]any{
		"K":  K,
		"A":  A,
		"B":  B,
		"C":  C,
		"D":  D,
		"KC": KC,
	})
}

func encrypt(w io.Writer) {
	SEED_KeySched(w, L0, L1, R0, R1, 0)  // K)    // Round 1
	SEED_KeySched(w, R0, R1, L0, L1, 2)  // K+2)  // Round 2
	SEED_KeySched(w, L0, L1, R0, R1, 4)  // K+4)  // Round 3
	SEED_KeySched(w, R0, R1, L0, L1, 6)  // K+6)  // Round 4
	SEED_KeySched(w, L0, L1, R0, R1, 8)  // K+8)  // Round 5
	SEED_KeySched(w, R0, R1, L0, L1, 10) // K+10) // Round 6
	SEED_KeySched(w, L0, L1, R0, R1, 12) // K+12) // Round 7
	SEED_KeySched(w, R0, R1, L0, L1, 14) // K+14) // Round 8
	SEED_KeySched(w, L0, L1, R0, R1, 16) // K+16) // Round 9
	SEED_KeySched(w, R0, R1, L0, L1, 18) // K+18) // Round 10
	SEED_KeySched(w, L0, L1, R0, R1, 20) // K+20) // Round 11
	SEED_KeySched(w, R0, R1, L0, L1, 22) // K+22) // Round 12
	SEED_KeySched(w, L0, L1, R0, R1, 24) // K+24) // Round 13
	SEED_KeySched(w, R0, R1, L0, L1, 26) // K+26) // Round 14
	SEED_KeySched(w, L0, L1, R0, R1, 28) // K+28) // Round 15
	SEED_KeySched(w, R0, R1, L0, L1, 30) // K+30) // Round 16
}

func decrypt(w io.Writer) {
	SEED_KeySched(w, L0, L1, R0, R1, 30) //K+30) // Round 1
	SEED_KeySched(w, R0, R1, L0, L1, 28) //K+28) // Round 2
	SEED_KeySched(w, L0, L1, R0, R1, 26) //K+26) // Round 3
	SEED_KeySched(w, R0, R1, L0, L1, 24) //K+24) // Round 4
	SEED_KeySched(w, L0, L1, R0, R1, 22) //K+22) // Round 5
	SEED_KeySched(w, R0, R1, L0, L1, 20) //K+20) // Round 6
	SEED_KeySched(w, L0, L1, R0, R1, 18) //K+18) // Round 7
	SEED_KeySched(w, R0, R1, L0, L1, 16) //K+16) // Round 8
	SEED_KeySched(w, L0, L1, R0, R1, 14) //K+14) // Round 9
	SEED_KeySched(w, R0, R1, L0, L1, 12) //K+12) // Round 10
	SEED_KeySched(w, L0, L1, R0, R1, 10) //K+10) // Round 11
	SEED_KeySched(w, R0, R1, L0, L1, 8)  //K+8)  // Round 12
	SEED_KeySched(w, L0, L1, R0, R1, 6)  //K+6)  // Round 13
	SEED_KeySched(w, R0, R1, L0, L1, 4)  //K+4)  // Round 14
	SEED_KeySched(w, L0, L1, R0, R1, 2)  //K+2)  // Round 15
	SEED_KeySched(w, R0, R1, L0, L1, 0)  //K+0)  // Round 16
}

func keySchedKey(w io.Writer) {
	RoundKeyUpdate0(w, 0, A, B, C, D, 0)   //KC0)     // K_1,0 and K_1,1
	RoundKeyUpdate1(w, 2, A, B, C, D, 1)   //KC1)   // K_2,0 and K_2,1
	RoundKeyUpdate0(w, 4, A, B, C, D, 2)   //KC2)   // K_3,0 and K_3,1
	RoundKeyUpdate1(w, 6, A, B, C, D, 3)   //KC3)   // K_4,0 and K_4,1
	RoundKeyUpdate0(w, 8, A, B, C, D, 4)   //KC4)   // K_5,0 and K_5,1
	RoundKeyUpdate1(w, 10, A, B, C, D, 5)  //KC5)  // K_6,0 and K_6,1
	RoundKeyUpdate0(w, 12, A, B, C, D, 6)  //KC6)  // K_7,0 and K_7,1
	RoundKeyUpdate1(w, 14, A, B, C, D, 7)  //KC7)  // K_8,0 and K_8,1
	RoundKeyUpdate0(w, 16, A, B, C, D, 8)  //KC8)  // K_9,0 and K_9,1
	RoundKeyUpdate1(w, 18, A, B, C, D, 9)  //KC9)  // K_10,0 and K_10,1
	RoundKeyUpdate0(w, 20, A, B, C, D, 10) //KC10) // K_11,0 and K_11,1
	RoundKeyUpdate1(w, 22, A, B, C, D, 11) //KC11) // K_12,0 and K_12,1
	RoundKeyUpdate0(w, 24, A, B, C, D, 12) //KC12) // K_13,0 and K_13,1
	RoundKeyUpdate1(w, 26, A, B, C, D, 13) //KC13) // K_14,0 and K_14,1
	RoundKeyUpdate0(w, 28, A, B, C, D, 14) //KC14) // K_15,0 and K_15,1

	RoundKeyUpdate1(w, 30, A, B, C, D, 15) // K_16,0 and K_16,1
}
