package cavp

import (
	"bufio"
	"bytes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/RyuaNerin/go-krypto/cmac"

	"github.com/stretchr/testify/require"
)

type (
	CMACGen struct {
		Count   int    `cavp:"COUNT"`
		Key     []byte `cavp:"K,hex"`
		Message []byte `cavp:"M,hex"`
		TagSize int    `cavp:"Tlen"`
		Tag     []byte `cavp:"T,hex,optional"` // Output
	}
	CMACVer struct {
		Count     int    `cavp:"COUNT"`
		Key       []byte `cavp:"K,hex"`
		Message   []byte `cavp:"M,hex"`
		TagSize   int    `cavp:"Tlen"`
		Tag       []byte `cavp:"T,hex"`
		IsInvalid bool   `cavp:"INVALID,optional"` // Output
		IsValid   bool   `cavp:"VALID,optional"`   // Output
	}
)

func ProcessCMACGen(
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
		key, message, tagBits, isRecord, parseErr := parseCMACRecord(record, false)
		if parseErr != nil {
			return parseErr
		}
		if !isRecord {
			response = append(response, record...)
			record = record[:0]
			return nil
		}
		b, cipherErr := newCipher(key)
		if cipherErr != nil {
			return cipherErr
		}
		h := cmac.New(b)
		n, writeErr := h.Write(message)
		if writeErr != nil {
			return writeErr
		}
		if n != len(message) {
			return io.ErrShortWrite
		}
		sum := h.Sum(nil)
		tagBytes := tagBits / 8
		if tagBytes > len(sum) {
			return fmt.Errorf("cavp: Tlen %d exceeds CMAC size", tagBits)
		}
		response = append(response, record...)
		response = append(response, "T = "+upperHex(sum[:tagBytes]))
		record = record[:0]
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
	for len(response) > 0 && strings.TrimSpace(response[len(response)-1]) == "" {
		response = response[:len(response)-1]
	}
	return os.WriteFile(output, []byte(strings.Join(response, "\r\n")), 0o644)
}

func TestCMACGen(
	t testing.TB,
	input string,
	newCipher func([]byte) (cipher.Block, error),
) {
	t.Helper()

	Test(t, input, func(v CMACGen) {
		b, err := newCipher(v.Key)
		require.NoError(t, err)

		testHashSMT(t, v.Message, v.Tag, func() hash.Hash { return cmac.New(b) })
	})
}

func ProcessCMACVer(
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
		key, message, tagBits, isRecord, parseErr := parseCMACRecord(record, true)
		if parseErr != nil {
			return parseErr
		}
		if !isRecord {
			response = append(response, record...)
			record = record[:0]
			return nil
		}
		var tag []byte
		for _, raw := range record {
			entry := parseLine(raw)
			if strings.EqualFold(entry.key, "T") {
				tag, err = hex.DecodeString(strings.Join(strings.Fields(entry.val), ""))
				if err != nil {
					return fmt.Errorf("cavp: invalid T: %w", err)
				}
			}
		}
		b, cipherErr := newCipher(key)
		if cipherErr != nil {
			return cipherErr
		}
		h := cmac.New(b)
		n, writeErr := h.Write(message)
		if writeErr != nil {
			return writeErr
		}
		if n != len(message) {
			return io.ErrShortWrite
		}
		sum := h.Sum(nil)
		tagBytes := tagBits / 8
		if tagBytes > len(sum) || len(tag) != tagBytes {
			return fmt.Errorf("cavp: invalid CMAC tag size: Tlen=%d, T=%d bytes", tagBits, len(tag))
		}
		result := "INVALID"
		if bytes.Equal(sum[:tagBytes], tag) {
			result = "VALID"
		}
		response = append(response, record...)
		response = append(response, result)
		record = record[:0]
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
	for len(response) > 0 && strings.TrimSpace(response[len(response)-1]) == "" {
		response = response[:len(response)-1]
	}
	return os.WriteFile(output, []byte(strings.Join(response, "\r\n")), 0o644)
}

func parseCMACRecord(lines []string, verification bool) (key, message []byte, tagBits int, isRecord bool, err error) {
	tagBits = -1
	haveKey, haveMessage, haveTag, haveTagBits := false, false, false, false
	for _, raw := range lines {
		entry := parseLine(raw)
		switch strings.ToLower(entry.key) {
		case "count":
			isRecord = true
		case "k":
			isRecord, haveKey = true, true
			key, err = hex.DecodeString(strings.Join(strings.Fields(entry.val), ""))
		case "m":
			isRecord, haveMessage = true, true
			message, err = hex.DecodeString(strings.Join(strings.Fields(entry.val), ""))
		case "tlen":
			isRecord, haveTagBits = true, true
			tagBits, err = strconv.Atoi(entry.val)
		case "t":
			isRecord, haveTag = true, true
		}
		if err != nil {
			return nil, nil, 0, true, fmt.Errorf("cavp: invalid %s: %w", entry.key, err)
		}
	}
	if !isRecord {
		return nil, nil, 0, false, nil
	}
	if !haveKey || !haveMessage || !haveTagBits || verification && !haveTag {
		return nil, nil, 0, true, fmt.Errorf("cavp: incomplete CMAC record")
	}
	if tagBits < 0 || tagBits%8 != 0 {
		return nil, nil, 0, true, fmt.Errorf("cavp: invalid Tlen %d", tagBits)
	}
	return key, message, tagBits, true, nil
}

func TestCMACVer(
	t testing.TB,
	input string,
	newCipher func([]byte) (cipher.Block, error),
) {
	t.Helper()

	Test(t, input, func(v CMACVer) {
		require.True(t, v.IsValid || v.IsInvalid, "both IsValid and IsInvalid are false")

		b, err := newCipher(v.Key)
		require.NoError(t, err)

		h := cmac.New(b)

		n, err := h.Write(v.Message)
		require.NoError(t, err)
		require.Equal(t, len(v.Message), n, "short write")

		sum := h.Sum(nil)[:v.TagSize/8]

		if v.IsValid {
			require.Equal(t, v.Tag, sum)
		} else {
			require.NotEqual(t, v.Tag, sum)
		}
	})
}
