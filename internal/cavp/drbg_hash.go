package cavp

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type DrbgHashKAT struct {
	PredictionResistance     bool `cavp:"PredictionResistance,header"`
	EntropyInputLen          int  `cavp:"EntropyInputLen,header"`
	NonceLen                 int  `cavp:"NonceLen,header"`
	PersonalizationStringLen int  `cavp:"PersonalizationStringLen,header"`
	AdditionalInputLen       int  `cavp:"AdditionalInputLen,header"`
	ReturnedBitsLen          int  `cavp:"ReturnedBitsLen,header"`

	Count                 int       `cavp:"COUNT"`
	EntropyInput          []byte    `cavp:"EntropyInput,hex"`
	Nonce                 []byte    `cavp:"Nonce,hex"`
	PersonalizationString []byte    `cavp:"PersonalizationString,hex"`
	EntropyInputReseed    []byte    `cavp:"EntropyInputReseed,hex,optional"`
	AdditionalInputReseed []byte    `cavp:"AdditionalInputReseed,hex,optional"`
	EntropyInputPR        [2][]byte `cavp:"EntropyInputPR,hex,optional"`
	AdditionalInput       [2][]byte `cavp:"AdditionalInput,hex"`
	ReturnedBits          []byte    `cavp:"ReturnedBits,hex,optional"`
}

type drbgConfig struct {
	predictionResistance     bool
	entropyInputLen          int
	nonceLen                 int
	personalizationStringLen int
	additionalInputLen       int
	returnedBitsLen          int
}

type drbgRecord struct {
	value                DrbgHashKAT
	lines                []string
	fields               map[string][]drbgLine
	additionalInputCount int
	entropyInputPRCount  int
	seen                 map[string]bool
}

type drbgLine struct {
	key string
	val string
}

func ProcessDrbgHashKAT(input, output string, newHash func() hash.Hash) error {
	return processDrbgKATFile(input, output, func(v DrbgHashKAT) ([]byte, error) {
		return drbgHash(v.vector(), newHash)
	})
}

func TestDrbgHashKAT(
	t testing.TB,
	input string,
	newHash func() hash.Hash,
) {
	t.Helper()
	Test(t, input, func(v DrbgHashKAT) {
		got, err := drbgHash(v.vector(), newHash)
		require.NoError(t, err)
		require.Equal(t, v.ReturnedBits, got, "COUNT = %d", v.Count)
	})
}

func (v DrbgHashKAT) vector() drbgVector {
	return drbgVector{
		EntropyInput:          v.EntropyInput,
		Nonce:                 v.Nonce,
		PersonalizationString: v.PersonalizationString,
		EntropyInputReseed:    v.EntropyInputReseed,
		AdditionalInputReseed: v.AdditionalInputReseed,
		EntropyInputPR:        v.EntropyInputPR,
		AdditionalInput:       v.AdditionalInput,
		ReturnedBitsLen:       v.ReturnedBitsLen,
		PredictionResistance:  v.PredictionResistance,
	}
}

func processDrbgKATFile(
	input, output string,
	generate func(DrbgHashKAT) ([]byte, error),
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

	bw := bufio.NewWriter(w)
	response := drbgResponseLines{w: bw}
	err = scanDrbgKAT(
		r,
		func(line string) error {
			return response.write(line)
		},
		func(record drbgRecord) error {
			if err := record.validate(false); err != nil {
				return err
			}
			returnedBits, err := generate(record.value)
			if err != nil {
				return err
			}
			for _, line := range record.lines {
				entry, _ := parseDrbgLine(line)
				if strings.EqualFold(entry.key, "ReturnedBits") {
					continue
				}
				if err = response.write(line); err != nil {
					return err
				}
			}
			return response.write(fmt.Sprintf("ReturnedBits = %X", returnedBits))
		},
	)
	if err != nil {
		return err
	}
	return bw.Flush()
}

type drbgResponseLines struct {
	w       io.Writer
	started bool
	pending []string
}

func (w *drbgResponseLines) write(line string) error {
	if strings.TrimSpace(line) == "" {
		w.pending = append(w.pending, line)
		return nil
	}

	if w.started {
		if _, err := io.WriteString(w.w, "\r\n"); err != nil {
			return err
		}
	}
	for _, pending := range w.pending {
		if _, err := io.WriteString(w.w, pending); err != nil {
			return err
		}
		if _, err := io.WriteString(w.w, "\r\n"); err != nil {
			return err
		}
	}
	w.pending = w.pending[:0]

	if _, err := io.WriteString(w.w, line); err != nil {
		return err
	}
	w.started = true
	return nil
}

func scanDrbgKAT(
	r io.Reader,
	lineFn func(string) error,
	recordFn func(drbgRecord) error,
) error {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)

	var config drbgConfig
	var record *drbgRecord
	flush := func() error {
		if record == nil {
			return nil
		}
		if err := record.decode(); err != nil {
			return err
		}
		err := recordFn(*record)
		record = nil
		return err
	}

	for scanner.Scan() {
		raw := strings.TrimRight(scanner.Text(), "\r")
		line := strings.TrimSpace(raw)
		entry, header := parseDrbgLine(line)

		if line == "" {
			if err := flush(); err != nil {
				return err
			}
			if lineFn != nil {
				if err := lineFn(raw); err != nil {
					return err
				}
			}
			continue
		}

		if record != nil {
			record.add(raw, entry)
			continue
		}

		if header {
			if err := config.set(entry); err != nil {
				return err
			}
			if lineFn != nil {
				if err := lineFn(raw); err != nil {
					return err
				}
			}
			continue
		}

		record = newDrbgKATRecord(config)
		record.add(raw, entry)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return flush()
}

func newDrbgKATRecord(config drbgConfig) *drbgRecord {
	return &drbgRecord{
		value: DrbgHashKAT{
			PredictionResistance:     config.predictionResistance,
			EntropyInputLen:          config.entropyInputLen,
			NonceLen:                 config.nonceLen,
			PersonalizationStringLen: config.personalizationStringLen,
			AdditionalInputLen:       config.additionalInputLen,
			ReturnedBitsLen:          config.returnedBitsLen,
		},
		fields: make(map[string][]drbgLine),
		seen:   make(map[string]bool),
	}
}

func (r *drbgRecord) add(raw string, entry drbgLine) {
	r.lines = append(r.lines, raw)
	key := strings.ToLower(entry.key)
	r.fields[key] = append(r.fields[key], entry)
}

// decode fills the computational vector only after the complete record has
// been collected. Values for repeated fields remain in their source order.
func (r *drbgRecord) decode() error {
	keys := []string{
		"count",
		"entropyinput",
		"nonce",
		"personalizationstring",
		"entropyinputreseed",
		"additionalinputreseed",
		"entropyinputpr",
		"additionalinput",
		"returnedbits",
	}
	for _, key := range keys {
		for _, entry := range r.fields[key] {
			if err := r.set(entry); err != nil {
				return err
			}
		}
	}
	return nil
}

func parseDrbgLine(line string) (drbgLine, bool) {
	line = strings.TrimSpace(line)
	header := strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]")
	if header {
		line = strings.TrimSpace(line[1 : len(line)-1])
	}

	index := strings.IndexByte(line, '=')
	if index < 0 {
		return drbgLine{key: strings.TrimSpace(line)}, header
	}
	return drbgLine{
		key: strings.TrimSpace(line[:index]),
		val: strings.TrimSpace(line[index+1:]),
	}, header
}

func (c *drbgConfig) set(entry drbgLine) error {
	key := strings.ToLower(entry.key)
	switch key {
	case "predictionresistance":
		value, err := strconv.ParseBool(entry.val)
		if err != nil {
			return fmt.Errorf("cavp: PredictionResistance: %w", err)
		}
		c.predictionResistance = value
	case "entropyinputlen":
		return parseDrbgInt(entry.val, &c.entropyInputLen, entry.key)
	case "noncelen":
		return parseDrbgInt(entry.val, &c.nonceLen, entry.key)
	case "personalizationstringlen":
		return parseDrbgInt(entry.val, &c.personalizationStringLen, entry.key)
	case "additionalinputlen":
		return parseDrbgInt(entry.val, &c.additionalInputLen, entry.key)
	case "returnedbitslen":
		return parseDrbgInt(entry.val, &c.returnedBitsLen, entry.key)
	}
	return nil
}

func (r *drbgRecord) set(entry drbgLine) error {
	key := strings.ToLower(entry.key)
	if key == "" || strings.HasPrefix(key, "#") {
		return nil
	}
	r.seen[key] = true

	switch key {
	case "count":
		return parseDrbgInt(entry.val, &r.value.Count, entry.key)
	case "entropyinput":
		return parseDrbgHex(entry.val, &r.value.EntropyInput, entry.key)
	case "nonce":
		return parseDrbgHex(entry.val, &r.value.Nonce, entry.key)
	case "personalizationstring":
		return parseDrbgHex(entry.val, &r.value.PersonalizationString, entry.key)
	case "entropyinputreseed":
		return parseDrbgHex(entry.val, &r.value.EntropyInputReseed, entry.key)
	case "additionalinputreseed":
		return parseDrbgHex(entry.val, &r.value.AdditionalInputReseed, entry.key)
	case "entropyinputpr":
		if r.entropyInputPRCount >= len(r.value.EntropyInputPR) {
			return fmt.Errorf("cavp: too many EntropyInputPR values")
		}
		err := parseDrbgHex(entry.val, &r.value.EntropyInputPR[r.entropyInputPRCount], entry.key)
		r.entropyInputPRCount++
		return err
	case "additionalinput":
		if r.additionalInputCount >= len(r.value.AdditionalInput) {
			return fmt.Errorf("cavp: too many AdditionalInput values")
		}
		err := parseDrbgHex(entry.val, &r.value.AdditionalInput[r.additionalInputCount], entry.key)
		r.additionalInputCount++
		return err
	case "returnedbits":
		return parseDrbgHex(entry.val, &r.value.ReturnedBits, entry.key)
	}
	return nil
}

func parseDrbgInt(value string, dst *int, key string) error {
	n, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return fmt.Errorf("cavp: %s: %w", key, err)
	}
	*dst = n
	return nil
}

func parseDrbgHex(value string, dst *[]byte, key string) error {
	value = strings.Join(strings.Fields(value), "")
	if len(value)%2 != 0 {
		value = "0" + value
	}
	b, err := hex.DecodeString(value)
	if err != nil {
		return fmt.Errorf("cavp: %s: %w", key, err)
	}
	*dst = b
	return nil
}

func (r drbgRecord) validate(returnedBits bool) error {
	if r.value.ReturnedBitsLen < 0 || r.value.ReturnedBitsLen%8 != 0 {
		return fmt.Errorf("cavp: ReturnedBitsLen must be a non-negative multiple of 8")
	}
	if !r.seen["entropyinput"] {
		return errors.New("cavp: missing EntropyInput")
	}
	if r.additionalInputCount != len(r.value.AdditionalInput) {
		return fmt.Errorf("cavp: got %d AdditionalInput values, want %d", r.additionalInputCount, len(r.value.AdditionalInput))
	}
	if r.value.PredictionResistance {
		if r.entropyInputPRCount != len(r.value.EntropyInputPR) {
			return fmt.Errorf("cavp: got %d EntropyInputPR values, want %d", r.entropyInputPRCount, len(r.value.EntropyInputPR))
		}
	} else if !r.seen["entropyinputreseed"] || !r.seen["additionalinputreseed"] {
		return errors.New("cavp: missing reseed input")
	}
	if returnedBits && !r.seen["returnedbits"] {
		return errors.New("cavp: missing ReturnedBits")
	}
	if returnedBits && len(r.value.ReturnedBits) != r.value.ReturnedBitsLen/8 {
		return fmt.Errorf("cavp: ReturnedBits is %d bytes, want %d", len(r.value.ReturnedBits), r.value.ReturnedBitsLen/8)
	}
	return nil
}
