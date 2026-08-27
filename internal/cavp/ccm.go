package cavp

import (
	"bufio"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/RyuaNerin/go-krypto/kipher"

	"github.com/stretchr/testify/require"
)

type (
	// CCM Generation Encryption
	CCMGe struct {
		Count      int    `cavp:"COUNT"`
		Key        []byte `cavp:"K,hex"`
		Nonce      []byte `cavp:"N,hex"`
		Adata      []byte `cavp:"A,hex"`
		PlainText  []byte `cavp:"P,hex"`
		TagSize    int    `cavp:"Tlen"`
		CipherText []byte `cavp:"C,hex,optional"` // Output
	}

	// CCM Decryption Verification
	CCMDv struct {
		Count      int    `cavp:"COUNT"`
		Key        []byte `cavp:"K,hex"`
		Nonce      []byte `cavp:"N,hex"`
		Adata      []byte `cavp:"A,hex"`
		CipherText []byte `cavp:"C,hex"`
		TagSize    int    `cavp:"Tlen"`
		PlainText  []byte `cavp:"P,hex,optional"`   // Output
		Invalid    bool   `cavp:"INVALID,optional"` // Output
	}
)

func ProcessCCMGe(
	input, output string,
	newCipher func([]byte) (cipher.Block, error),
) error {
	r, err := os.Open(input)
	if err != nil {
		return err
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)
	response := make([]string, 0, 1024)
	record := make([]string, 0, 8)
	flush := func() error {
		if len(record) == 0 {
			return nil
		}
		v, isRecord, parseErr := parseCCMRecord(record, false)
		if parseErr != nil {
			return parseErr
		}
		response = append(response, record...)
		record = record[:0]
		if !isRecord {
			return nil
		}
		b, cipherErr := newCipher(v.Key)
		if cipherErr != nil {
			return cipherErr
		}
		ccm, cipherErr := kipher.NewCCM(b, len(v.Nonce), v.TagSize/8)
		if cipherErr != nil {
			return cipherErr
		}
		sealed := ccm.Seal(nil, v.Nonce, v.PlainText, v.Adata)
		response = append(response, "C = "+upperHex(sealed))
		return nil
	}
	for scanner.Scan() {
		raw := strings.TrimSuffix(scanner.Text(), "\r")
		if strings.TrimSpace(raw) == "" {
			if err := flush(); err != nil {
				return err
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
	return writeRecordResponse(output, response)
}

func TestCCMGe(
	t testing.TB,
	input string,
	newCipher func([]byte) (cipher.Block, error),
) {
	t.Helper()

	Test(t, input, func(v CCMGe) {
		b, err := newCipher(v.Key)
		require.NoError(t, err)

		cipherText, tag := v.CipherText[:len(v.PlainText)], v.CipherText[len(v.PlainText):]

		testAEAD(
			t,
			v.Nonce, v.PlainText, v.Adata,
			cipherText, tag,
			func() (cipher.AEAD, error) {
				return kipher.NewCCM(b, len(v.Nonce), v.TagSize/8)
			},
		)
	})
}

func ProcessCCMDv(
	input, output string,
	newCipher func([]byte) (cipher.Block, error),
) error {
	r, err := os.Open(input)
	if err != nil {
		return err
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)
	response := make([]string, 0, 1024)
	record := make([]string, 0, 8)
	flush := func() error {
		if len(record) == 0 {
			return nil
		}
		v, isRecord, parseErr := parseCCMRecord(record, true)
		if parseErr != nil {
			return parseErr
		}
		response = append(response, record...)
		record = record[:0]
		if !isRecord {
			return nil
		}
		b, cipherErr := newCipher(v.Key)
		if cipherErr != nil {
			return cipherErr
		}
		ccm, cipherErr := kipher.NewCCM(b, len(v.Nonce), v.TagSize/8)
		if cipherErr != nil {
			return cipherErr
		}
		plaintext, openErr := ccm.Open(nil, v.Nonce, v.CipherText, v.Adata)
		if openErr != nil {
			response = append(response, "INVALID")
		} else {
			response = append(response, "P = "+upperHex(plaintext))
		}
		return nil
	}
	for scanner.Scan() {
		raw := strings.TrimSuffix(scanner.Text(), "\r")
		if strings.TrimSpace(raw) == "" {
			if err := flush(); err != nil {
				return err
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
	return writeRecordResponse(output, response)
}

func parseCCMRecord(lines []string, decrypt bool) (CCMGe, bool, error) {
	var v CCMGe
	seen := make(map[string]bool, 7)
	isRecord := false
	for _, raw := range lines {
		entry := parseLine(raw)
		key := strings.ToLower(entry.key)
		var dst *[]byte
		switch key {
		case "count":
			isRecord = true
		case "k":
			dst, isRecord = &v.Key, true
		case "n":
			dst, isRecord = &v.Nonce, true
		case "a":
			dst, isRecord = &v.Adata, true
		case "p":
			dst, isRecord = &v.PlainText, true
		case "c":
			dst, isRecord = &v.CipherText, true
		case "tlen":
			isRecord = true
			var err error
			v.TagSize, err = strconv.Atoi(entry.val)
			if err != nil {
				return v, true, fmt.Errorf("cavp: invalid Tlen: %w", err)
			}
		}
		if dst != nil {
			var err error
			*dst, err = hex.DecodeString(strings.Join(strings.Fields(entry.val), ""))
			if err != nil {
				return v, true, fmt.Errorf("cavp: invalid %s: %w", entry.key, err)
			}
		}
		if key != "" {
			seen[key] = true
		}
	}
	if !isRecord {
		return v, false, nil
	}
	required := []string{"k", "n", "a", "tlen"}
	if decrypt {
		required = append(required, "c")
	} else {
		required = append(required, "p")
	}
	for _, key := range required {
		if !seen[key] {
			return v, true, fmt.Errorf("cavp: incomplete CCM record: missing %s", key)
		}
	}
	if v.TagSize <= 0 || v.TagSize%8 != 0 {
		return v, true, fmt.Errorf("cavp: invalid Tlen %d", v.TagSize)
	}
	return v, true, nil
}

func writeRecordResponse(output string, lines []string) error {
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	return os.WriteFile(output, []byte(strings.Join(lines, "\r\n")), 0o644)
}
