package cavp

import (
	"bufio"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"hash"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/RyuaNerin/go-krypto/eckcdsa"

	"github.com/stretchr/testify/require"
)

type (
	// ECKCDSAKpg is an EC-KCDSA key-pair generation record.
	ECKCDSAKpg struct {
		Curve string `cavp:"Curve,header=eckcdsa-svt"`

		D  []byte `cavp:"d,hex"`
		Qx []byte `cavp:"Qx,hex"`
		Qy []byte `cavp:"Qy,hex"`
	}

	// ECKCDSAPkv is an EC-KCDSA public-key validation record.
	ECKCDSAPkv struct {
		Qx     []byte `cavp:"Qx,hex"`
		Qy     []byte `cavp:"Qy,hex"`
		Result string `cavp:"Result,optional"` // Output: P(pass) or F(fail).
	}

	// ECKCDSASgt is an EC-KCDSA signature-generation record.
	ECKCDSASgt struct {
		Curve  string `cavp:"Curve,header=eckcdsa-svt"`
		SHAAlg string `cavp:"SHAAlg,header=eckcdsa-svt"`

		M                      []byte `cavp:"M,hex"`
		Qx                     []byte `cavp:"Qx,hex,optional"` // Output.
		Qy                     []byte `cavp:"Qy,hex,optional"` // Output.
		R                      []byte `cavp:"R,hex,optional"`  // Output.
		S                      []byte `cavp:"S,hex,optional"`  // Output.
		SignVerificationResult string `cavp:"Sign Verification Result,string,optional"`
	}

	// ECKCDSASvt is an EC-KCDSA signature-verification record.
	ECKCDSASvt struct {
		Curve  string `cavp:"Curve,header=eckcdsa-svt"`
		SHAAlg string `cavp:"SHAAlg,header=eckcdsa-svt"`

		M      []byte `cavp:"M,hex"`
		Qx     []byte `cavp:"Qx,hex"`
		Qy     []byte `cavp:"Qy,hex"`
		R      []byte `cavp:"R,hex"`
		S      []byte `cavp:"S,hex"`
		Result string `cavp:"Result,optional"` // Output: P(pass) or F(fail).
	}
)

func ProcessECKCDSAKpg(input, output string, curve elliptic.Curve, iteration int) error {
	r, err := os.Open(input)
	if err != nil {
		return err
	}
	defer r.Close()
	if iteration < 0 {
		return fmt.Errorf("cavp: EC-KCDSA KPG: negative iteration count %d", iteration)
	}

	coordinateBytes := (curve.Params().BitSize + 7) / 8
	scalarBytes := len(curve.Params().N.Bytes())
	var lines []string
	var record []eckcdsaRequestLine
	lineNumber := 0
	requests := 0
	flush := func() error {
		if len(record) == 0 {
			return nil
		}
		var headers []string
		haveD := false
		for _, line := range record {
			trimmed := strings.TrimSpace(line.text)
			if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
				headers = append(headers, line.text)
				continue
			}
			key, _, found := strings.Cut(line.text, "=")
			if !found || !strings.EqualFold(strings.TrimSpace(key), "d") {
				continue
			}
			if haveD {
				return fmt.Errorf("cavp: EC-KCDSA KPG: line %d: duplicate d", line.number)
			}
			haveD = true
		}
		if !haveD {
			record = record[:0]
			return nil
		}

		lines = append(lines, headers...)
		requests++
		for range iteration {
			privateKey, err := eckcdsa.GenerateKey(curve, rand.Reader)
			if err != nil {
				return err
			}
			d, err := formatECKCDSAInteger(privateKey.D, scalarBytes)
			if err != nil {
				return err
			}
			qx, err := formatECKCDSAInteger(privateKey.X, coordinateBytes)
			if err != nil {
				return err
			}
			qy, err := formatECKCDSAInteger(privateKey.Y, coordinateBytes)
			if err != nil {
				return err
			}
			lines = append(lines, "d = "+d, "Qx = "+qx, "Qy = "+qy, "")
		}
		record = record[:0]
		return nil
	}

	scanner := newECKCDSAScanner(r)
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSuffix(scanner.Text(), "\r")
		if strings.TrimSpace(line) == "" {
			if err := flush(); err != nil {
				return err
			}
			continue
		}
		record = append(record, eckcdsaRequestLine{text: line, number: lineNumber})
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if err := flush(); err != nil {
		return err
	}
	if requests == 0 {
		return fmt.Errorf("cavp: EC-KCDSA KPG: no d request found")
	}
	return writeECKCDSALines(output, lines, true)
}

func ProcessECKCDSAPkv(input, output string, curve elliptic.Curve) error {
	r, err := os.Open(input)
	if err != nil {
		return err
	}
	defer r.Close()

	var lines []string
	var record []eckcdsaRequestLine
	records := 0
	lineNumber := 0
	flush := func() error {
		if len(record) == 0 {
			return nil
		}
		fields, err := parseECKCDSARecord(record, "Qx", "Qy")
		if err != nil {
			return err
		}
		qx, haveQx := fields["qx"]
		qy, haveQy := fields["qy"]
		if haveQx || haveQy {
			if !haveQx || !haveQy {
				return fmt.Errorf("cavp: EC-KCDSA PKV: incomplete record ending at line %d", record[len(record)-1].number)
			}
			result := "F"
			if curve.IsOnCurve(new(big.Int).SetBytes(qx), new(big.Int).SetBytes(qy)) {
				result = "P"
			}
			lines = appendECKCDSARawLines(lines, record)
			lines = append(lines, "Result = "+result)
			records++
		} else {
			lines = appendECKCDSARawLines(lines, record)
		}
		record = record[:0]
		return nil
	}

	scanner := newECKCDSAScanner(r)
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSuffix(scanner.Text(), "\r")
		if strings.TrimSpace(line) == "" {
			if err := flush(); err != nil {
				return err
			}
			lines = append(lines, line)
			continue
		}
		record = append(record, eckcdsaRequestLine{text: line, number: lineNumber})
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if err := flush(); err != nil {
		return err
	}
	if records == 0 {
		return fmt.Errorf("cavp: EC-KCDSA PKV: no records found")
	}
	return writeECKCDSALines(output, trimECKCDSATrailingBlankLines(lines), false)
}

func ProcessECKCDSASgt(input, output string, curve elliptic.Curve, h hash.Hash) error {
	r, err := os.Open(input)
	if err != nil {
		return err
	}
	defer r.Close()

	coordinateBytes := (curve.Params().BitSize + 7) / 8
	scalarBytes := len(curve.Params().N.Bytes())
	rBytes := min(h.Size(), scalarBytes)
	var lines []string
	var record []eckcdsaRequestLine
	records := 0
	lineNumber := 0
	flush := func() error {
		if len(record) == 0 {
			return nil
		}
		fields, err := parseECKCDSARecord(record, "M")
		if err != nil {
			return err
		}
		message, haveMessage := fields["m"]
		lines = appendECKCDSARawLines(lines, record)
		record = record[:0]
		if !haveMessage {
			return nil
		}

		privateKey, err := eckcdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			return err
		}
		rValue, sValue, err := eckcdsa.Sign(rand.Reader, privateKey, h, message)
		if err != nil {
			return err
		}
		if !eckcdsa.Verify(&privateKey.PublicKey, h, message, rValue, sValue) {
			return fmt.Errorf("cavp: generated EC-KCDSA signature did not verify")
		}
		qx, err := formatECKCDSAInteger(privateKey.X, coordinateBytes)
		if err != nil {
			return err
		}
		qy, err := formatECKCDSAInteger(privateKey.Y, coordinateBytes)
		if err != nil {
			return err
		}
		rHex, err := formatECKCDSAInteger(rValue, rBytes)
		if err != nil {
			return err
		}
		sHex, err := formatECKCDSAInteger(sValue, scalarBytes)
		if err != nil {
			return err
		}
		lines = append(lines, "Qx = "+qx, "Qy = "+qy, "R = "+rHex, "S = "+sHex)
		records++
		return nil
	}

	scanner := newECKCDSAScanner(r)
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSuffix(scanner.Text(), "\r")
		if strings.TrimSpace(line) == "" {
			if err := flush(); err != nil {
				return err
			}
			lines = append(lines, line)
			continue
		}
		record = append(record, eckcdsaRequestLine{text: line, number: lineNumber})
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if err := flush(); err != nil {
		return err
	}
	if records == 0 {
		return fmt.Errorf("cavp: EC-KCDSA SGT: no records found")
	}
	if len(lines) == 0 || lines[len(lines)-1] != "" {
		lines = append(lines, "")
	}
	return writeECKCDSALines(output, lines, true)
}

func ProcessECKCDSASvt(input, output string, curve elliptic.Curve, h hash.Hash) error {
	r, err := os.Open(input)
	if err != nil {
		return err
	}
	defer r.Close()

	var lines []string
	var record []eckcdsaRequestLine
	records := 0
	lineNumber := 0
	flush := func() error {
		if len(record) == 0 {
			return nil
		}
		fields, err := parseECKCDSARecord(record, "M", "Qx", "Qy", "R", "S")
		if err != nil {
			return err
		}
		_, haveM := fields["m"]
		_, haveQx := fields["qx"]
		_, haveQy := fields["qy"]
		_, haveR := fields["r"]
		_, haveS := fields["s"]
		if !haveM && !haveQx && !haveQy && !haveR && !haveS {
			lines = appendECKCDSARawLines(lines, record)
			record = record[:0]
			return nil
		}
		if !haveM || !haveQx || !haveQy || !haveR || !haveS {
			return fmt.Errorf("cavp: EC-KCDSA SVT: incomplete record ending at line %d", record[len(record)-1].number)
		}

		pub := eckcdsa.PublicKey{
			Curve: curve,
			X:     new(big.Int).SetBytes(fields["qx"]),
			Y:     new(big.Int).SetBytes(fields["qy"]),
		}
		verified := eckcdsa.Verify(
			&pub,
			h,
			fields["m"],
			new(big.Int).SetBytes(fields["r"]),
			new(big.Int).SetBytes(fields["s"]),
		)
		result := "F"
		if verified {
			result = "P"
		}
		lines = appendECKCDSARawLines(lines, record)
		lines = append(lines, "Result = "+result)
		record = record[:0]
		records++
		return nil
	}

	scanner := newECKCDSAScanner(r)
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSuffix(scanner.Text(), "\r")
		if strings.TrimSpace(line) == "" {
			if err := flush(); err != nil {
				return err
			}
			lines = append(lines, line)
			continue
		}
		record = append(record, eckcdsaRequestLine{text: line, number: lineNumber})
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if err := flush(); err != nil {
		return err
	}
	if records == 0 {
		return fmt.Errorf("cavp: EC-KCDSA SVT: no records found")
	}
	return writeECKCDSALines(output, trimECKCDSATrailingBlankLines(lines), false)
}

type eckcdsaRequestLine struct {
	text   string
	number int
}

func parseECKCDSARecord(record []eckcdsaRequestLine, wantedKeys ...string) (map[string][]byte, error) {
	wanted := make(map[string]string, len(wantedKeys))
	for _, key := range wantedKeys {
		wanted[strings.ToLower(key)] = key
	}
	fields := make(map[string][]byte, len(wantedKeys))
	seen := make(map[string]bool, len(wantedKeys))
	for _, line := range record {
		key, value, found := strings.Cut(line.text, "=")
		if !found {
			continue
		}
		key = strings.ToLower(strings.TrimSpace(key))
		canonical, ok := wanted[key]
		if !ok {
			continue
		}
		if seen[key] {
			return nil, fmt.Errorf("cavp: EC-KCDSA: line %d: duplicate %s", line.number, canonical)
		}
		decoded, err := hex.DecodeString(strings.TrimSpace(value))
		if err != nil {
			return nil, fmt.Errorf("cavp: EC-KCDSA: line %d: invalid %s: %w", line.number, canonical, err)
		}
		fields[key] = decoded
		seen[key] = true
	}
	return fields, nil
}

func appendECKCDSARawLines(dst []string, record []eckcdsaRequestLine) []string {
	for _, line := range record {
		dst = append(dst, line.text)
	}
	return dst
}

func newECKCDSAScanner(r *os.File) *bufio.Scanner {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)
	return scanner
}

func formatECKCDSAInteger(value *big.Int, size int) (string, error) {
	if value == nil || value.Sign() < 0 || len(value.Bytes()) > size {
		return "", fmt.Errorf("cavp: EC-KCDSA: integer does not fit in %d bytes", size)
	}
	buf := make([]byte, size)
	value.FillBytes(buf)
	return strings.ToUpper(hex.EncodeToString(buf)), nil
}

func trimECKCDSATrailingBlankLines(lines []string) []string {
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func writeECKCDSALines(output string, lines []string, terminate bool) (err error) {
	w, err := os.Create(output)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := w.Close(); err == nil {
			err = closeErr
		}
	}()

	bw := bufio.NewWriter(w)
	for i, line := range lines {
		if i > 0 {
			if _, err = bw.WriteString("\r\n"); err != nil {
				return err
			}
		}
		if _, err = bw.WriteString(line); err != nil {
			return err
		}
	}
	if terminate {
		if _, err = bw.WriteString("\r\n"); err != nil {
			return err
		}
	}
	return bw.Flush()
}

// TestECKCDSASvt checks the expected result of every SVT response.
func TestECKCDSASvt(t testing.TB, input string, curve elliptic.Curve, h hash.Hash) {
	t.Helper()

	Test(t, input, func(v ECKCDSASvt) {
		pub := eckcdsa.PublicKey{
			Curve: curve,
			X:     new(big.Int).SetBytes(v.Qx),
			Y:     new(big.Int).SetBytes(v.Qy),
		}
		ok := eckcdsa.Verify(
			&pub,
			h,
			v.M,
			new(big.Int).SetBytes(v.R),
			new(big.Int).SetBytes(v.S),
		)
		switch strings.TrimSpace(v.Result) {
		case "P":
			require.True(t, ok, "expected pass but got fail")
		case "F":
			require.False(t, ok, "expected fail but got pass")
		default:
			require.Failf(t, "invalid result", "expected P or F, got %q", v.Result)
		}
	})
}
