package cavp

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"hash"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/RyuaNerin/go-krypto/pbkdf2"

	"github.com/stretchr/testify/require"
)

type PbkdfKAT struct {
	PRF       string `cavp:"PRF,string,header"`
	Iteration int    `cavp:"Iteration,int,header"`
	Count     int    `cavp:"COUNT,int"`
	Password  string `cavp:"Password,string"`
	Salt      []byte `cavp:"Salt,hex"`
	KeyLength int    `cavp:"KLen,int"`
	MasterKey []byte `cavp:"MK,hex,optional"`
}

func ProcessPbkdfKAT(input, output string, newHash func() hash.Hash) error {
	r, err := os.Open(input)
	if err != nil {
		return err
	}
	defer r.Close()

	var response, record []string
	iteration := 0
	flush := func() error {
		if len(record) == 0 {
			return nil
		}

		fields := make(map[string][]string)
		for _, raw := range record {
			key, value, ok := parsePbkdfRecordLine(raw)
			if ok {
				fields[strings.ToLower(key)] = append(fields[strings.ToLower(key)], value)
			}
		}

		var password string
		var salt []byte
		keyLength := -1
		if values := fields["password"]; len(values) != 0 {
			password = values[len(values)-1]
		}
		if values := fields["salt"]; len(values) != 0 {
			salt, err = hex.DecodeString(strings.Join(strings.Fields(values[len(values)-1]), ""))
			if err != nil {
				return fmt.Errorf("cavp: Salt: %w", err)
			}
		}
		if values := fields["klen"]; len(values) != 0 {
			keyLength, err = strconv.Atoi(strings.TrimSpace(values[len(values)-1]))
			if err != nil {
				return fmt.Errorf("cavp: KLen: %w", err)
			}
		}
		if iteration <= 0 || keyLength < 0 || keyLength%8 != 0 {
			return fmt.Errorf("cavp: invalid PBKDF parameters: Iteration=%d, KLen=%d", iteration, keyLength)
		}
		masterKey := pbkdf2.Generate(nil, []byte(password), salt, iteration, keyLength/8, newHash)
		for _, raw := range record {
			key, _, ok := parsePbkdfRecordLine(raw)
			if !ok || !strings.EqualFold(key, "MK") {
				response = append(response, raw)
			}
		}
		response = append(response, fmt.Sprintf("MK = %X", masterKey))
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
			if strings.EqualFold(entry.key, "Iteration") {
				iteration, err = strconv.Atoi(strings.TrimSpace(entry.val))
				if err != nil {
					return fmt.Errorf("cavp: Iteration: %w", err)
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

// parsePbkdfRecordLine removes the single format separator after '=' while
// leaving every other byte in the value intact. In particular, a password
// encoded as "Password =  value " has the value " value ".
func parsePbkdfRecordLine(line string) (key, value string, ok bool) {
	before, after, ok := strings.Cut(line, "=")
	if !ok {
		return strings.TrimSpace(line), "", false
	}
	key = strings.TrimSpace(before)
	value = after
	value = strings.TrimPrefix(value, " ")
	return key, value, true
}

func TestPbkdfKAT(t testing.TB, filename string, newHash func() hash.Hash) {
	t.Helper()

	var out []byte

	Test(t, filename, func(v PbkdfKAT) {
		out = pbkdf2.Generate(out[:0], []byte(v.Password), v.Salt, v.Iteration, v.KeyLength/8, newHash)

		require.Equalf(t, v.MasterKey, out, "MK mismatch for COUNT=%d, Iteration=%d, KLen=%d", v.Count, v.Iteration, v.KeyLength)
	})
}
