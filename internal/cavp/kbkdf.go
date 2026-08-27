package cavp

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/RyuaNerin/go-krypto/kbkdf"

	"github.com/stretchr/testify/require"
)

type KBKDFMode uint8

const (
	KBKDFModeCounter KBKDFMode = iota
	KBKDFModeDoublePipeline
	KBKDFModeFeedback
)

// Vector is a CAVP KBKDF record. RLen is zero for a [NO COUNTER] file.
type Vector struct {
	PRF     string `cavp:"PRF,string,header,optional"`
	RLen    int    `cavp:"RLEN,int,header,optional"`
	Count   int    `cavp:"COUNT,int"`
	L       int    `cavp:"L,int"`
	KI      []byte `cavp:"KI,hex"`
	Label   []byte `cavp:"Label,hex"`
	Context []byte `cavp:"Context,hex"`
	IV      []byte `cavp:"IV,hex,optional"`
	KO      []byte `cavp:"KO,hex,optional"`
}

func derive(v Vector, mode KBKDFMode, prf kbkdf.PRF) []byte {
	switch mode {
	case KBKDFModeCounter:
		return kbkdf.CounterMode(nil, prf, v.KI, v.Label, v.Context, v.RLen/8, v.L/8)
	case KBKDFModeDoublePipeline:
		return kbkdf.DoublePipeMode(nil, prf, v.KI, v.Label, v.Context, v.RLen/8, v.L/8)
	case KBKDFModeFeedback:
		return kbkdf.FeedbackMode(nil, prf, v.KI, v.Label, v.Context, v.IV, v.RLen/8, v.L/8)
	default:
		panic("invalid mode")
	}
}

func ProcessKbkdfKat(input, output string, mode KBKDFMode, prf kbkdf.PRF) error {
	r, err := os.Open(input)
	if err != nil {
		return err
	}
	defer r.Close()

	var response, record []string
	rlen := 0
	flush := func() error {
		if len(record) == 0 {
			return nil
		}

		fields := make(map[string][]string)
		for _, raw := range record {
			entry := parseLine(strings.TrimSpace(raw))
			key := strings.ToLower(entry.key)
			fields[key] = append(fields[key], entry.val)
		}

		v := Vector{RLen: rlen}
		if values := fields["count"]; len(values) != 0 {
			v.Count, err = strconv.Atoi(strings.TrimSpace(values[len(values)-1]))
			if err != nil {
				return fmt.Errorf("cavp: COUNT: %w", err)
			}
		}
		if values := fields["l"]; len(values) != 0 {
			v.L, err = strconv.Atoi(strings.TrimSpace(values[len(values)-1]))
			if err != nil {
				return fmt.Errorf("cavp: L: %w", err)
			}
		}
		for key, dst := range map[string]*[]byte{
			"ki":      &v.KI,
			"label":   &v.Label,
			"context": &v.Context,
			"iv":      &v.IV,
		} {
			values := fields[key]
			if len(values) == 0 {
				continue
			}
			*dst, err = hex.DecodeString(strings.Join(strings.Fields(values[len(values)-1]), ""))
			if err != nil {
				return fmt.Errorf("cavp: %s: %w", key, err)
			}
		}
		if v.L < 0 || v.L%8 != 0 || v.RLen < 0 || v.RLen%8 != 0 {
			return fmt.Errorf("cavp: invalid KBKDF lengths: L=%d, RLEN=%d", v.L, v.RLen)
		}
		for _, raw := range record {
			if !strings.EqualFold(parseLine(strings.TrimSpace(raw)).key, "KO") {
				response = append(response, raw)
			}
		}
		response = append(response, fmt.Sprintf("KO = %X", derive(v, mode, prf)))
		record = record[:0]
		return nil
	}

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	for scanner.Scan() {
		raw := strings.TrimSuffix(scanner.Text(), "\r")
		line := strings.TrimSpace(raw)
		if line == "" {
			if err := flush(); err != nil {
				return err
			}
			response = append(response, raw)
			continue
		}
		if len(record) == 0 && strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			entry := parseLine(strings.TrimSpace(line[1 : len(line)-1]))
			if strings.EqualFold(entry.key, "RLEN") {
				rlen, err = strconv.Atoi(entry.val)
				if err != nil {
					return fmt.Errorf("cavp: RLEN: %w", err)
				}
			}
			response = append(response, raw)
			continue
		}
		record = append(record, raw)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if err := flush(); err != nil {
		return err
	}
	for len(response) > 0 && strings.TrimSpace(response[len(response)-1]) == "" {
		response = response[:len(response)-1]
	}
	return os.WriteFile(output, []byte(strings.Join(response, "\r\n")), 0o644)
}

func TestKbkdfKat(t testing.TB, filename string, mode KBKDFMode, prf kbkdf.PRF) {
	t.Helper()

	Test(t, filename, func(v Vector) {
		got := derive(v, mode, prf)

		require.Equal(t, v.KO, got, "KO mismatch for COUNT=%d, L=%d, RLen=%d", v.Count, v.L, v.RLen)
	})
}
