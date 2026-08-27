package cavp

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"
	"testing"

	kcdsa_internal "github.com/RyuaNerin/go-krypto/internal/kcdsa"
	"github.com/RyuaNerin/go-krypto/kcdsa"

	"github.com/stretchr/testify/require"
)

type (
	KCDSAKpg struct {
		PSize int `cavp:"|P|,header=plain"`
		QSize int `cavp:"|Q|,header=plain"`

		P     []byte `cavp:"P,optional"`         // Output
		Q     []byte `cavp:"Q,optional"`         // Output
		G     []byte `cavp:"G,optional"`         // Output
		X     []byte `cavp:"X,optional"`         // Output
		Y     []byte `cavp:"Y,optional"`         // Output
		Seed  []byte `cavp:"SEED,optional"`      // Output
		Count int    `cavp:"Count,hex,optional"` // Output
		XKey  []byte `cavp:"XKEY,optional"`      // Output
		Oupri []byte `cavp:"OUPRI,optional"`     // Output
		H     []byte `cavp:"h,optional"`         // Output
		J     []byte `cavp:"J,optional"`         // Output
	}

	KCDSASgt struct {
		PSize  int    `cavp:"|P|,header=plain"`
		QSize  int    `cavp:"|Q|,header=plain"`
		SHAAlg string `cavp:"SHAAlg,header=plain"`

		P []byte `cavp:"P,optional"` // Output
		Q []byte `cavp:"Q,optional"` // Output
		G []byte `cavp:"G,optional"` // Output
		M []byte `cavp:"M,hex"`
		X []byte `cavp:"X,optional"` // Output
		Y []byte `cavp:"Y,optional"` // Output
		R []byte `cavp:"R,optional"` // Output
		S []byte `cavp:"S,optional"` // Output
	}

	KCDSASvt struct {
		PSize  int    `cavp:"|P|,header=kcdsa-svt"`
		QSize  int    `cavp:"|Q|,header=kcdsa-svt"`
		SHAAlg string `cavp:"SHAAlg,header=kcdsa-svt"`

		P      []byte `cavp:"P,optional"`
		Q      []byte `cavp:"Q,optional"`
		G      []byte `cavp:"G,optional"`
		M      []byte `cavp:"M,hex"`
		X      []byte `cavp:"X,optional"`
		Y      []byte `cavp:"Y,optional"`
		R      []byte `cavp:"R,hex"`
		S      []byte `cavp:"S,hex"`
		Result string `cavp:"Result,optional"` // Output, P(pass) or F(fail)
	}
)

func ProcessKCDSAKpg(
	input, output string,
	parameterSizes kcdsa.ParameterSizes,
	iteration int,
) error {
	r, err := os.Open(input)
	if err != nil {
		return err
	}
	defer r.Close()
	w, err := os.Create(output)
	if err != nil {
		return err
	}
	defer w.Close()

	domain, ok := kcdsa_internal.GetDomain(int(parameterSizes))
	if !ok {
		return fmt.Errorf("cavp: invalid KCDSA parameter sizes %d", parameterSizes)
	}

	// A KPG request contains two headers followed by human-readable directions.
	// Copy only the headers verbatim: the directions are not part of the response.
	scanner := newKCDSAScanner(r)
	var headers []string
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "|P|") || strings.HasPrefix(trimmed, "|Q|") {
			headers = append(headers, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if len(headers) != 2 {
		return fmt.Errorf("cavp: KCDSA KPG request has %d parameter headers, want 2", len(headers))
	}

	wr := bufio.NewWriter(w)
	for _, line := range headers {
		if err := writeKCDSALine(wr, line); err != nil {
			return err
		}
	}
	if err := writeKCDSALine(wr, ""); err != nil {
		return err
	}

	for range iteration {
		generated, err := kcdsa_internal.GenerateParametersTTAK(rand.Reader, domain)
		if err != nil {
			return err
		}
		privateKey := kcdsa.PrivateKey{
			PublicKey: kcdsa.PublicKey{
				Parameters: kcdsa.Parameters{
					P: generated.P,
					Q: generated.Q,
					G: generated.G,
					GenParameters: kcdsa.GenerationParameters{
						J:     generated.J,
						Seed:  generated.Seed,
						Count: generated.Count,
					},
				},
			},
		}

		xkey, oupri, err := kcdsa.GenerateKeyWithSeed(&privateKey, rand.Reader, nil, nil, parameterSizes)
		if err != nil {
			return err
		}

		values := []struct {
			name  string
			value string
		}{
			{"P", kcdsaFixedHex(privateKey.P.Bytes(), domain.A)},
			{"Q", kcdsaFixedHex(privateKey.Q.Bytes(), domain.B)},
			{"G", kcdsaFixedHex(privateKey.G.Bytes(), domain.A)},
			{"X", kcdsaFixedHex(privateKey.X.Bytes(), domain.B)},
			{"Y", kcdsaFixedHex(privateKey.Y.Bytes(), domain.A)},
			{"Seed", kcdsaFixedHex(privateKey.GenParameters.Seed, domain.B)},
			{"Count", fmt.Sprintf("%08X", uint32(privateKey.GenParameters.Count))},
			{"XKEY", kcdsaFixedHex(xkey, domain.B)},
			{"OUPRI", strings.ToUpper(hex.EncodeToString(oupri))},
			{"h", kcdsaFixedHex(generated.H.Bytes(), domain.A)},
			{"J", kcdsaFixedHex(generated.J.Bytes(), domain.A-domain.B)},
		}
		for _, value := range values {
			if err := writeKCDSALine(wr, value.name+" = "+value.value); err != nil {
				return err
			}
		}
		if err := writeKCDSALine(wr, ""); err != nil {
			return err
		}
	}

	return wr.Flush()
}

func TestKCDSAKpg(
	t testing.TB,
	filename string,
	parameterSizes kcdsa.ParameterSizes,
) {
	t.Helper()

	domain, ok := kcdsa_internal.GetDomain(int(parameterSizes))
	require.True(t, ok, "Invalid parameter sizes")

	h := domain.NewHash()

	var buf []byte

	Test(t, filename, func(v KCDSAKpg) {
		// J -> PQ -> (H -> G) -> X -> Y
		var gotJ big.Int

		buf, ok = kcdsa_internal.GenerateJ(&gotJ, buf, v.Seed, h, domain)
		require.Truef(t, ok, "Failed to generate J: seed=%s, J=%s", hex.EncodeToString(v.Seed), gotJ.Text(16))
		require.Equalf(t, new(big.Int).SetBytes(v.J).Text(16), gotJ.Text(16), "failed to generate J: seed=%s, J=%s", hex.EncodeToString(v.Seed), gotJ.Text(16))

		var gotP, gotQ big.Int
		gotP.SetBytes(v.P)
		gotQ.SetBytes(v.Q)

		ok := kcdsa_internal.RegeneratePQ(
			&gotP, &gotQ,
			domain,
			&gotJ,
			v.Seed,
			v.Count,
		)
		require.True(t, ok, "Failed to regenerate PQ")
		require.Equalf(t, new(big.Int).SetBytes(v.P).Text(16), gotP.Text(16), "failed to regenerate PQ: P mismatch\nJ=%s\nseed=%s", gotJ.Text(16), hex.EncodeToString(v.Seed))
		require.Equalf(t, new(big.Int).SetBytes(v.Q).Text(16), gotQ.Text(16), "failed to regenerate PQ: Q mismatch\nJ=%s\nseed=%s", gotJ.Text(16), hex.EncodeToString(v.Seed))

		var gotG big.Int
		gotG.SetBytes(v.G)

		var gotX big.Int
		gotX.SetBytes(v.X)

		var gotY big.Int
		kcdsa_internal.GenerateY(
			&gotY,
			&gotP, &gotQ, &gotG, &gotX,
		)
		require.Equalf(t, new(big.Int).SetBytes(v.Y).Text(16), gotY.Text(16), "failed to generate Y: Y mismatch\nJ=%s\nseed=%s", gotJ.Text(16), hex.EncodeToString(v.Seed))
	})
}

func ProcessKCDSASgt(
	input, output string,
	parameterSizes kcdsa.ParameterSizes,
	iteration int,
) error {
	r, err := os.Open(input)
	if err != nil {
		return err
	}
	defer r.Close()
	w, err := os.Create(output)
	if err != nil {
		return err
	}
	defer w.Close()

	domain, ok := kcdsa_internal.GetDomain(int(parameterSizes))
	if !ok {
		return fmt.Errorf("cavp: invalid KCDSA parameter sizes %d", parameterSizes)
	}

	scanner := newKCDSAScanner(r)
	wr := bufio.NewWriter(w)
	wroteMessage := false
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "M") || !strings.Contains(trimmed, "=") {
			// Before the first M these are the three headers and their blank
			// separator. Later blank lines are request record separators and are
			// regenerated below in the response's canonical shape.
			if !wroteMessage {
				if err := writeKCDSALine(wr, line); err != nil {
					return err
				}
			} else if trimmed != "" {
				return fmt.Errorf("cavp: unexpected KCDSA SGT line %q", line)
			}
			continue
		}

		message, err := parseKCDSAHexLine(line, "M")
		if err != nil {
			return err
		}
		wroteMessage = true

		for range iteration {
			var privateKey kcdsa.PrivateKey
			err := kcdsa.GenerateParameters(&privateKey.Parameters, rand.Reader, parameterSizes)
			if err != nil {
				return err
			}

			_, _, err = kcdsa.GenerateKeyWithSeed(&privateKey, rand.Reader, nil, nil, parameterSizes)
			if err != nil {
				return err
			}

			r, s, err := kcdsa.Sign(rand.Reader, &privateKey, parameterSizes, message)
			if err != nil {
				return err
			}

			values := []struct {
				name  string
				value string
			}{
				{"P", kcdsaFixedHex(privateKey.P.Bytes(), domain.A)},
				{"Q", kcdsaFixedHex(privateKey.Q.Bytes(), domain.B)},
				{"G", kcdsaFixedHex(privateKey.G.Bytes(), domain.A)},
				// M is deliberately emitted from the original request line below.
				{"X", kcdsaFixedHex(privateKey.X.Bytes(), domain.B)},
				{"Y", kcdsaFixedHex(privateKey.Y.Bytes(), domain.A)},
				{"R", kcdsaFixedHex(r.Bytes(), min(domain.LH, domain.B))},
				{"S", kcdsaFixedHex(s.Bytes(), domain.B)},
			}
			for i, value := range values {
				if i == 3 {
					if err := writeKCDSALine(wr, line); err != nil {
						return err
					}
				}
				if err := writeKCDSALine(wr, value.name+" = "+value.value); err != nil {
					return err
				}
			}
			if err := writeKCDSALine(wr, ""); err != nil {
				return err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if !wroteMessage {
		return fmt.Errorf("cavp: KCDSA SGT request contains no M record")
	}

	return wr.Flush()
}

func ProcessKCDSASvt(
	input, output string,
	parameterSizes kcdsa.ParameterSizes,
) error {
	r, err := os.Open(input)
	if err != nil {
		return err
	}
	defer r.Close()
	w, err := os.Create(output)
	if err != nil {
		return err
	}
	defer w.Close()

	scanner := newKCDSAScanner(r)
	wr := bufio.NewWriter(w)
	recordLines := make([]string, 0, 8)
	fields := make(map[string][]byte, 8)
	records := 0
	headerComplete := false

	flushRecord := func() error {
		if len(recordLines) == 0 {
			return nil
		}
		if err := writeKCDSASvtRecord(wr, recordLines, fields, parameterSizes, records > 0); err != nil {
			return err
		}
		records++
		recordLines = recordLines[:0]
		clear(fields)
		return nil
	}

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			if !headerComplete {
				if err := writeKCDSALine(wr, line); err != nil {
					return err
				}
				headerComplete = true
			} else if err := flushRecord(); err != nil {
				return err
			}
			// Requests contain two empty lines between records. Empty records are
			// ignored and the response separator is emitted by flushRecord.
			continue
		}

		if !headerComplete {
			if err := writeKCDSALine(wr, line); err != nil {
				return err
			}
			continue
		}

		recordLines = append(recordLines, line)
		name, value, known, err := parseKCDSASvtField(line)
		if err != nil {
			return err
		}
		if !known {
			continue
		}
		if _, duplicate := fields[name]; duplicate {
			return fmt.Errorf("cavp: duplicate KCDSA SVT field %s", name)
		}
		fields[name] = value
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	// The final request record may end at EOF rather than with a blank line.
	if err := flushRecord(); err != nil {
		return err
	}
	if records == 0 {
		return fmt.Errorf("cavp: KCDSA SVT request contains no records")
	}

	return wr.Flush()
}

var kcdsaSvtFieldNames = [...]string{"P", "Q", "G", "M", "X", "Y", "R", "S"}

func parseKCDSASvtField(line string) (name string, value []byte, known bool, err error) {
	namePart, valuePart, assignment := strings.Cut(line, "=")
	if !assignment {
		return "", nil, false, nil
	}
	name = strings.TrimSpace(namePart)
	for _, fieldName := range kcdsaSvtFieldNames {
		if name != fieldName {
			continue
		}
		value, err = hex.DecodeString(strings.TrimSpace(valuePart))
		if err != nil {
			return "", nil, false, fmt.Errorf("cavp: invalid KCDSA %s value: %w", name, err)
		}
		return name, value, true, nil
	}
	return "", nil, false, nil
}

func writeKCDSASvtRecord(
	w *bufio.Writer,
	lines []string,
	fields map[string][]byte,
	parameterSizes kcdsa.ParameterSizes,
	separator bool,
) error {
	for _, name := range kcdsaSvtFieldNames {
		if _, ok := fields[name]; !ok {
			return fmt.Errorf("cavp: incomplete KCDSA SVT record: missing %s", name)
		}
	}

	publicKey := kcdsa.PublicKey{
		Parameters: kcdsa.Parameters{
			P: new(big.Int).SetBytes(fields["P"]),
			Q: new(big.Int).SetBytes(fields["Q"]),
			G: new(big.Int).SetBytes(fields["G"]),
		},
		Y: new(big.Int).SetBytes(fields["Y"]),
	}
	result := "F"
	if kcdsa.Verify(
		&publicKey,
		parameterSizes,
		fields["M"],
		new(big.Int).SetBytes(fields["R"]),
		new(big.Int).SetBytes(fields["S"]),
	) {
		result = "P"
	}

	if separator {
		if _, err := w.WriteString("\r\n\r\n"); err != nil {
			return err
		}
	}
	for _, line := range lines {
		if err := writeKCDSALine(w, line); err != nil {
			return err
		}
	}
	_, err := w.WriteString("Result = " + result)
	return err
}

func newKCDSAScanner(r io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(r)
	// SVT messages are deliberately long. Leave ample headroom while retaining
	// Scanner's useful one-line-at-a-time processing model.
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)
	return scanner
}

func writeKCDSALine(w *bufio.Writer, line string) error {
	_, err := w.WriteString(line + "\r\n")
	return err
}

func parseKCDSAHexLine(line, wantName string) ([]byte, error) {
	name, value, ok := strings.Cut(line, "=")
	if !ok || strings.TrimSpace(name) != wantName {
		return nil, fmt.Errorf("cavp: expected KCDSA %s line, got %q", wantName, line)
	}
	decoded, err := hex.DecodeString(strings.TrimSpace(value))
	if err != nil {
		return nil, fmt.Errorf("cavp: invalid KCDSA %s value: %w", wantName, err)
	}
	return decoded, nil
}

func kcdsaFixedHex(value []byte, bits int) string {
	width := (bits + 7) / 8
	// Generated KCDSA values are constrained to bits, but never discard the
	// most-significant bytes if a caller violates that invariant.
	if len(value) > width {
		width = len(value)
	}
	encoded := make([]byte, width*2)
	for i := range encoded {
		encoded[i] = '0'
	}
	hex.Encode(encoded[len(encoded)-len(value)*2:], value)
	return strings.ToUpper(string(encoded))
}

func TestKCDSASvt(
	t testing.TB,
	filename string,
	parameterSizes kcdsa.ParameterSizes,
) {
	t.Helper()

	Test(t, filename, func(v KCDSASvt) {
		privateKey := kcdsa.PrivateKey{
			PublicKey: kcdsa.PublicKey{
				Parameters: kcdsa.Parameters{
					P: new(big.Int).SetBytes(v.P),
					Q: new(big.Int).SetBytes(v.Q),
					G: new(big.Int).SetBytes(v.G),
				},
				Y: new(big.Int).SetBytes(v.Y),
			},
			X: new(big.Int).SetBytes(v.X),
		}

		ok := kcdsa.Verify(
			&privateKey.PublicKey,
			parameterSizes,
			v.M,
			new(big.Int).SetBytes(v.R),
			new(big.Int).SetBytes(v.S),
		)
		if ok {
			require.Equal(t, "P", v.Result, "Expected Pass, but failed")
		} else {
			require.Equal(t, "F", v.Result, "Expected Fail, but passed")
		}
	})
}
